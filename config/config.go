package config

import (
	"os"

	"github.com/case-management-suite/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	LogConfig          LogConfig
	CasesService       CasesServiceConfig
	CasesStorage       DatabaseConfig
	GraphQLConfig      GraphQLConfig
	RulesServiceConfig RulesServiceConfig
}

type DatabaseType byte

const (
	Sqlite DatabaseType = iota
	Postgres
)

type DatabaseConfig struct {
	/*
		Connection string.

		- For Sqlite: Use path to the *.db file
		- For Postgres: A connection string like follows
			"host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	*/
	Address      string
	CreateSQL    string
	DatabaseType DatabaseType
	LogSQL       bool
}

type Channel = string

type QueueConnectionConfig struct {
	Address                  string
	CaseActionsChannel       Channel
	CaseNotificationsChannel Channel
	PurgeOnStart             bool
	SendRetries              int
	LogLevel                 zerolog.Level
}

type CasesServiceConfig struct {
	Host string
	Port int16
}

type GraphQLConfig struct {
	Port int
}

type RulesServiceConfig struct {
	QueueConfig QueueConnectionConfig
	Logger      zerolog.Logger
	QueueType   QueueType
}

func NewLocalAppConfig() AppConfig {
	return AppConfig{
		LogConfig: LogConfig{
			Logger:        common.BuildDefaultLogger(),
			CasesService:  zerolog.ErrorLevel,
			CasesDB:       zerolog.ErrorLevel,
			RulesServce:   zerolog.ErrorLevel,
			RulesServer:   zerolog.DebugLevel,
			QueueService:  zerolog.ErrorLevel,
			WorkScheduler: zerolog.ErrorLevel,
		},
		CasesService: CasesServiceConfig{
			Host: "localhost",
			Port: 7777,
		},
		CasesStorage: DatabaseConfig{
			CreateSQL:    "create table if not exists cases(case_id INTEGER PRIMARY KEY, status TEXT);",
			Address:      "./test_cases.db",
			DatabaseType: Sqlite,
			LogSQL:       false,
		},
		GraphQLConfig: GraphQLConfig{
			Port: 8080,
		},
		RulesServiceConfig: RulesServiceConfig{
			QueueType: RabbitMQ,
			QueueConfig: QueueConnectionConfig{
				Address:                  "amqp://guest:guest@0.0.0.0:5672/",
				CaseActionsChannel:       "case_action_channel",
				CaseNotificationsChannel: "case_notifications_channel",
				SendRetries:              5,
				LogLevel:                 zerolog.ErrorLevel,
			},
		},
	}
}

func NewLocalTestAppConfig() AppConfig {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return AppConfig{
		LogConfig: LogConfig{
			Logger:        common.BuildDefaultLogger(),
			CasesService:  zerolog.DebugLevel,
			CasesDB:       zerolog.DebugLevel,
			RulesServce:   zerolog.DebugLevel,
			RulesServer:   zerolog.DebugLevel,
			QueueService:  zerolog.DebugLevel,
			WorkScheduler: zerolog.DebugLevel,
		},
		CasesService: CasesServiceConfig{
			Host: "localhost",
			Port: 7777,
		},
		CasesStorage: DatabaseConfig{
			CreateSQL:    "create table if not exists cases(case_id INTEGER PRIMARY KEY, status TEXT);",
			Address:      "./test_cases.db",
			DatabaseType: Sqlite,
			LogSQL:       false,
		},
		GraphQLConfig: GraphQLConfig{
			Port: 8080,
		},
		RulesServiceConfig: RulesServiceConfig{
			QueueType: GoChannels,
			Logger:    log.With().Str("server", "RulesServiceServer").Logger(),
			QueueConfig: QueueConnectionConfig{
				Address:                  "amqp://guest:guest@0.0.0.0:5672/",
				CaseActionsChannel:       "test_case_actions_channel",
				CaseNotificationsChannel: "test_case_events_channel",
				PurgeOnStart:             true,
				SendRetries:              2,
				LogLevel:                 zerolog.DebugLevel,
			},
		},
	}
}
