package common

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ServerUtils struct {
	Logger zerolog.Logger
}

func NewTestServerUtils() ServerUtils {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().CallerWithSkipFrameCount(4).Logger()
	return ServerUtils{Logger: logger}
}
