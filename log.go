package common

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func BuildDefaultLogger() zerolog.Logger {
	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
}
