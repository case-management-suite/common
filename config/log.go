package config

import "github.com/rs/zerolog"

type LogConfig struct {
	QueueService  zerolog.Level
	WorkScheduler zerolog.Level
	RulesServce   zerolog.Level
	RulesServer   zerolog.Level
	CasesService  zerolog.Level
	CasesServer   zerolog.Level
	CasesDB       zerolog.Level
	Logger        zerolog.Logger
}
