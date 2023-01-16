package service

import (
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/logger"
	"github.com/case-management-suite/common/server"
)

type ServiceUtils struct {
	IsSet  bool
	Logger logger.Logger
}

func (*ServiceUtils) mustImplementServiceable() {}

func (su *ServiceUtils) clone(other ServiceUtils) {
	su.Logger = other.Logger
	su.IsSet = other.IsSet
}

func NewTestServiceUtils() ServiceUtils {
	return ServiceUtils{IsSet: true, Logger: logger.NewTestLogger()}
}

func NewServiceUtils(serviceName string, appConfig config.AppConfig) ServiceUtils {
	return ServiceUtils{IsSet: true, Logger: logger.NewServiceLogger(serviceName, appConfig.Env)}
}

func NewServiceUtilsFromServerUtils(utls server.ServerUtils) ServiceUtils {
	return ServiceUtils{IsSet: true, Logger: utls.Logger}
}

type Serviceable interface {
	clone(other ServiceUtils)
	mustImplementServiceable()
}

type Service[T Serviceable] struct {
	Value T
}

func NewService[T Serviceable](appConfig config.AppConfig, serviceName string, serviceable T) Service[T] {
	serviceable.clone(NewServiceUtils(serviceName, appConfig))
	return Service[T]{Value: serviceable}
}

// type FxServiceInput struct {
// 	fx.In
// 	AppConfig config.AppConfig
// }

// type FxServiceOutput[T Serviceable] struct {
// 	fx.Out
// 	Service Service[T]
// }

// type FxServiceFactory func()

// func NewFxService[T Serviceable](serviceName, opts fx.Option) fx.Option {

// }
