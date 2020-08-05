package ctxutil

import (
	"context"
	"errors"

	"github.com/cubny/cart/internal/auth"
)

type ctxKeyType int

const (
	ctxAuthoriseAccess ctxKeyType = iota
)

// SetUserAuthAccessKey to the provided context.
func SetUserAuthAccessKey(ctx context.Context, accessKey *auth.AccessKey) (context.Context, error) {
	switch {
	case ctx == nil:
		return nil, errors.New("context is required")
	case accessKey == nil:
		return nil, errors.New("access key is required")
	case accessKey.AccessKey == "":
		return nil, errors.New("access key is required")
	}

	return context.WithValue(ctx, ctxAuthoriseAccess, *accessKey), nil
}

// GetUserAuthAccessKey retrieved from the provided context.
func GetUserAuthAccessKey(ctx context.Context) (auth.AccessKey, error) {
	userAuthoriseAccess, ok := ctx.Value(ctxAuthoriseAccess).(auth.AccessKey)
	if !ok {
		return auth.AccessKey{}, errors.New("access key is not set on the context")
	}

	if userAuthoriseAccess.UserID == 0 {
		return auth.AccessKey{}, errors.New("access key is not set on the context")
	}

	return userAuthoriseAccess, nil
}
