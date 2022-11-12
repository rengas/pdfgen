package contexts

import (
	"context"
	"errors"
)

var (
	ErrNotInContext = errors.New("not in context")
)

type arguments string

var (
	userId   = arguments("userId")
	designId = arguments("designId")
)

func WithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userId, userId)
}

func UserIdFromContext(ctx context.Context) (string, error) {
	id, ok := ctx.Value(userId).(string)
	if !ok {
		return "", ErrNotInContext
	}
	return id, nil
}

func WithDesignId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userId, userId)
}

func DesignIdFromContext(ctx context.Context) (string, error) {
	id, ok := ctx.Value(designId).(string)
	if !ok {
		return "", ErrNotInContext
	}
	return id, nil
}
