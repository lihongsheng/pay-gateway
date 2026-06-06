package log

import (
	"context"
	"github.com/google/uuid"
	"github.com/lihongsheng/pay-gateway/global"
	"go.uber.org/zap"
)

type LogTraceId struct{}

func WithCtx(ctx context.Context) *zap.Logger {
	traceId := ctx.Value(LogTraceId{})
	if traceId == nil {
		return global.GVA_LOG.With(zap.String("logTraceId", uuid.NewString()))
	}
	l := global.GVA_LOG.With(zap.Any("logTraceId", traceId))
	return l
}

func NewCtx(ctx context.Context) (context.Context, *zap.Logger) {
	traceId := uuid.NewString()
	ctx = context.WithValue(ctx, LogTraceId{}, traceId)
	l := global.GVA_LOG.With(zap.Any("logTraceId", traceId))
	return ctx, l
}
