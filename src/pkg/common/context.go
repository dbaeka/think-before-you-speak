package common

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/lithammer/shortuuid/v3"
	"think-before-you-speak/pkg/log"
)

type ctxKey int

const (
	correlationIDKey ctxKey = iota
	sseServerKey     ctxKey = iota
)

func EchoToGoContext(c echo.Context) context.Context {
	return c.Request().Context()
}

func SetCorrelationID(ctx context.Context, correlationID string) context.Context {
	return setContext(ctx, correlationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
	v := getContext[string](ctx, correlationIDKey)

	if v == nil {
		l := log.Get()
		l.Warn().Msg("Correlation ID not found in context")

		// add "gen_" prefix to distinguish generated correlation IDs from correlation IDs passed by the client
		// it's useful to detect if correlation ID was not passed properly

		return "gen_" + shortuuid.New()
	}

	return *v
}

func setContext(ctx context.Context, key ctxKey, value interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, key, value)
}

func getContext[T any](ctx context.Context, key ctxKey) *T {
	if ctx == nil {
		return nil
	}

	v, ok := ctx.Value(key).(T)
	if ok {
		return &v
	}
	return nil
}
