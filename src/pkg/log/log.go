package log

import (
	"context"
	"github.com/rs/zerolog"
)

func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func ToContext(ctx context.Context, fields map[string]interface{}) context.Context {
	l := FromContext(ctx)
	return l.With().Fields(fields).Logger().WithContext(ctx)
}
