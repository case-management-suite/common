package ctxutils

import "context"

type executionCtxType string

const executionCtxKey = executionCtxType("env")

func WithEnvContext(ctx context.Context, env string) context.Context {
	return context.WithValue(ctx, executionCtxKey, env)
}

func GetEnvType(ctx context.Context) string {
	env, ok := ctx.Value(executionCtxKey).(string)
	if ok {
		return env
	}
	return "unkown"
}
