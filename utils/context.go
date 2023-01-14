package utils

import (
	"context"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ContextData string
type ContextLogger string
type ExecutionCtx string

const (
	ExecutionCtxKey ExecutionCtx = ExecutionCtx("key")
	Prod            ExecutionCtx = ExecutionCtx("prod")
	Test            ExecutionCtx = ExecutionCtx("test")
)

const (
	ServiceName ContextData   = ContextData("ServiceName")
	CtxLogger   ContextLogger = ContextLogger("ContextLogger")
)

type ContextDataValue struct {
	Service string
}

func WithTestContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ExecutionCtxKey, Test)
}

func WithProdContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ExecutionCtxKey, Prod)
}

func IsTestContext(ctx context.Context) bool {
	switch t := ctx.Value(ExecutionCtxKey).(type) {
	case ExecutionCtx:
		return t == Test
	default:
		return false
	}
}

func IsProdContext(ctx context.Context) bool {
	switch t := ctx.Value(ExecutionCtxKey).(type) {
	case ExecutionCtx:
		return t == Test
	default:
		return false
	}
}

func DecorateContext(ctx context.Context, serviceName string) context.Context {
	defaultLogger := log.Output(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: false,
	}).With().Caller().Logger()
	zerolog.DefaultContextLogger = &defaultLogger

	ctx = context.WithValue(ctx, ServiceName, ContextDataValue{serviceName})

	if IsTestContext(ctx) {
		ctx = zerolog.Ctx(ctx).With().Str("execution", "test").Logger().WithContext(ctx)
	} else if IsTestContext(ctx) {
		ctx = zerolog.Ctx(ctx).With().Str("execution", "prod").Logger().WithContext(ctx)
	}
	return ctx
}

// func LoggerFromContext(ctx context.Context) zerolog.Logger {
// 	l := ctx.Value(CtxLogger)
// 	switch v := l.(type) {
// 	case zerolog.Logger:
// 		return v
// 		case
// 	}
// }
