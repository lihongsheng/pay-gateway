package logic

import (
	"context"

	"github.com/lihongsheng/pay-gateway/internal/svc"
	"github.com/lihongsheng/pay-gateway/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayLogic {
	return &PayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayLogic) Pay(in *payment.PayRequest) (*payment.PayResponse, error) {
	// todo: add your logic here and delete this line

	return &payment.PayResponse{}, nil
}
