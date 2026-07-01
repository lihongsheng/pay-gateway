package domain

import (
	"github.com/lihongsheng/pay-gateway/config"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/redis/go-redis/v9"
)

type ServiceGroup struct {
	PaymentService
	Router
	RefundService
}

func NewServiceGroup(
	svc *svc.ServiceContext,
	redis *redis.Client,
	cfg config.Config,
) *ServiceGroup {
	Service = &ServiceGroup{
		PaymentService: NewPaymentService(svc, redis, cfg),
		Router:         NewRouterService(svc.PaymentAccountRepo, svc.RouterRepo),
		RefundService:  NewRefundService(svc),
	}
	return Service
}

// DefaultServiceGroup 包级单例
var Service *ServiceGroup
