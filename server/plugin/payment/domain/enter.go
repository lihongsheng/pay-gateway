package domain

import (
	"github.com/lihongsheng/pay-gateway/plugin/payment/infrastructure"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
)

type ServiceGroup struct {
	PaymentService
	Router
	RefundService
}

func NewServiceGroup(
	paymentOrderRepo repo.PaymentOrderRepo,
	refundRepo repo.RefundRepo,
	paymentAccountRepo repo.PaymentAccountRepo,
	routerRepo repo.RouterRepo,
	event infrastructure.Event,
) *ServiceGroup {
	return &ServiceGroup{
		PaymentService: NewPaymentService(paymentOrderRepo, paymentAccountRepo, event),
		Router:         NewRouterService(paymentAccountRepo, routerRepo),
		RefundService:  NewRefundService(refundRepo, paymentOrderRepo, paymentAccountRepo, event),
	}
}

// DefaultServiceGroup 包级单例
var Service *ServiceGroup
