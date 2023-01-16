package logger

import (
	"os"

	"github.com/case-management-suite/common/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerLevel = zerolog.Level

var LogDebug = zerolog.LevelDebugValue
var LogError = zerolog.LevelErrorValue

type Logger struct {
	zerolog.Logger
}

func BuildDefaultLogger() Logger {
	return Logger{Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().CallerWithSkipFrameCount(4).Logger()}
}

func NewLogger(env config.EnvType) Logger {
	switch env {
	case config.Env.Local, config.Env.Test:
		return Logger{Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().CallerWithSkipFrameCount(4).Logger()}
	case config.Env.Prod:
		return Logger{Logger: log.With().CallerWithSkipFrameCount(4).Logger()}
	default:
		log.Warn().Interface("env", env).Msg("The request environment is not recognized")
		return Logger{Logger: log.Logger}
	}
}

func NewServiceLogger(name string, env config.EnvType) Logger {
	return Logger{Logger: NewLogger(env).With().Str("service", name).Logger()}
}

func NewTestLogger() Logger {
	return Logger{Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().CallerWithSkipFrameCount(4).Logger()}
}
