package logging

import (
	"context"
)

func ContextWithFields(ctx context.Context, fields ...Field) context.Context {
	if ctx == nil {
		return ctx
	}

	values := make([]interface{}, 0)
	for _, field := range fields {
		values = append(values, field.Label)
		values = append(values, field.Value)
	}

	if val, ok := ctx.Value(Logger{}).([]interface{}); ok && val != nil {
		return context.WithValue(ctx, Logger{}, append(val, values...))
	}

	return context.WithValue(ctx, Logger{}, values)
}
