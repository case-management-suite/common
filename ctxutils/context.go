package ctxutils

import "context"

type ContextDecoration struct {
	Name string
}

func DecorateInitialContext(ctx context.Context, env string) context.Context {
	ctx = WithEnvContext(ctx, env)
	ctx = WithLoggerConfig(ctx)
	return ctx
}

func DecorateContext(ctx context.Context, deco ContextDecoration) context.Context {
	ctx = WithServiceName(ctx, deco.Name)
	ctx = WithLoggerConfig(ctx)
	return ctx
}
