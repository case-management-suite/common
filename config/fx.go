package config

import "go.uber.org/fx"

func FxConfig(appConfig AppConfig) fx.Option {
	return fx.Options(fx.Supply(appConfig))
}
