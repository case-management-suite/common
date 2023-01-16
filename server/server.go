package server

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/ctxutils"
	"github.com/case-management-suite/common/logger"
	"go.uber.org/fx"
)

type ServerConnectionType string

const (
	HttpServerType     ServerConnectionType = ServerConnectionType("HttpServer")
	GRPCServerType     ServerConnectionType = ServerConnectionType("gRPCServer")
	ProcessServerType  ServerConnectionType = ServerConnectionType("ProcessServer")
	GroupOfServersType ServerConnectionType = ServerConnectionType("GroupOfServers")
)

type ServerConfig struct {
	ServerName string
	Host       string
	Port       int
	Type       ServerConnectionType
}

type Serveable interface {
	GetName() string
	Start(context.Context) error
	Stop(context.Context) error
	GetServerConfig() *ServerConfig
}

type Server[T Serveable] struct {
	Env           config.EnvType
	Logger        logger.Logger
	ServerMetrics ServerMetrics
	Server        T
}

func (s *Server[T]) logServerInfo(msg string) {
	serverName := s.Server.GetName()
	serverInfo := s.Server.GetServerConfig()
	l := s.Logger.Info().Str("server_name", serverName)
	if serverInfo != nil {
		stype := serverInfo.Type
		l = l.Str("server_type", string(stype))
		switch stype {
		case HttpServerType, GRPCServerType:
			l = l.Str("host", serverInfo.Host).Int("port", serverInfo.Port)
		}
	}
	l.Msg(msg)
}

func (s *Server[T]) Start(ctx context.Context) error {

	serverName := s.Server.GetName()
	s.logServerInfo("Starting server...")
	ctx = ctxutils.DecorateContext(ctx, ctxutils.ContextDecoration{Name: serverName})
	err := s.Server.Start(ctx)
	if err != nil {
		s.logServerInfo("Failed to start")
	} else {
		s.logServerInfo("Started")
	}
	return err
}

func (s *Server[T]) Stop(ctx context.Context) error {
	s.logServerInfo("Stopping server...")
	err := s.Server.Stop(ctx)
	if err != nil {
		s.logServerInfo("Failed to stop")
	} else {
		s.logServerInfo("Stopped")
	}
	return err
}

func (s *Server[T]) fxServer(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})
}

func (s *Server[T]) GetFxOption() fx.Option {
	return fx.Options(
		fx.Invoke(s.fxServer),
	)
}

type factoryFn[T Serveable] func(ServerUtils) T

func NewServer[T Serveable](factory factoryFn[T], appConfig config.AppConfig) Server[T] {
	l := logger.NewLogger(appConfig.Env)
	params := ServerUtils{Logger: l}
	return Server[T]{Server: factory(params), Logger: l}
}

type ServerUtils struct {
	Logger logger.Logger
}

func (ServerUtils) RunAllAsync(ctx context.Context, timeout time.Duration, fns ...func(context.Context) error) []chan error {
	return RunAllAsync(ctx, timeout, fns...)
}

func (ServerUtils) RunAsync(fn func(context.Context) error, ctx context.Context, timeout time.Duration) chan error {
	errchan := make(chan error, 1)
	go func() {
		for {
			select {
			case errchan <- fn(ctx):
				return
			case <-time.After(timeout):
				errchan <- errors.New("timeout")
			}
		}
	}()
	return errchan
}

func NewTestServerUtils() ServerUtils {
	logger := logger.NewTestLogger()
	return ServerUtils{Logger: logger}
}

type ServerMetrics struct{}

func StartServer[T Serveable](srv Server[T]) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.Start(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer srv.Stop(ctx)
}

func RunAllAsync(ctx context.Context, timeout time.Duration, fns ...func(context.Context) error) []chan error {
	errchans := make([]chan error, len(fns))
	for i, fn := range fns {
		errchans[i] = make(chan error, 1)
		go func(i int, fn func(context.Context) error) {
			for {
				select {
				case errchans[i] <- fn(ctx):
					return
				case <-time.After(timeout):
					errchans[i] <- errors.New("timeout")
				}
			}
		}(i, fn)
	}

	return errchans
}
