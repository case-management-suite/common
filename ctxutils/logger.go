package ctxutils

import (
	"context"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func WithLoggerConfig(ctx context.Context) context.Context {
	defaultLogger := log.Output(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: false,
	}).With().Caller().Logger()
	zerolog.DefaultContextLogger = &defaultLogger

	env := GetEnvType(ctx)
	names := GetServiceNames(ctx)

	ctx = zerolog.Ctx(ctx).With().Interface("service_names", names).Str("env", env).Logger().WithContext(ctx)
	return ctx
}
