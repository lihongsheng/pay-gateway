package logic

import (
	"context"

	"github.com/lihongsheng/pay-gateway/internal/svc"
	"github.com/lihongsheng/pay-gateway/refund"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogic {
	return &QueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryLogic) Query(in *refund.RefundQuery) (*refund.RefundResponse, error) {
	// todo: add your logic here and delete this line

	return &refund.RefundResponse{}, nil
}
