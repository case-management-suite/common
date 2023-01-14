package ctxutils

import "context"

type ctxServiceNameKeyType string

const serviceNameKey = ctxServiceNameKeyType("service_name")

func GetServiceNames(ctx context.Context) *[]string {
	names, ok := ctx.Value(serviceNameKey).([]string)
	if ok {
		return &names
	}
	return nil
}

func WithServiceName(ctx context.Context, name string) context.Context {
	namesAdd := GetServiceNames(ctx)
	if namesAdd != nil {
		names := *namesAdd
		names = append(names, name)
		return context.WithValue(ctx, serviceNameKey, names)
	}
	return context.WithValue(ctx, serviceNameKey, []string{})
}
