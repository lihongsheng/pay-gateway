package log

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LogTraceId struct{}

func WithCtx(ctx context.Context) *zap.Logger {
	traceId := ctx.Value(LogTraceId{})
	if traceId == nil {
		return zap.L().With(zap.String("logTraceId", uuid.NewString()))
	}
	l := zap.L().With(zap.Any("logTraceId", traceId))
	return l
}

func NewCtx(ctx context.Context) (context.Context, *zap.Logger) {
	traceId := uuid.NewString()
	ctx = context.WithValue(ctx, LogTraceId{}, traceId)
	l := zap.L().With(zap.Any("logTraceId", traceId))
	return ctx, l
}
