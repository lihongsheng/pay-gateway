package domain

import "github.com/lihongsheng/pay-gateway/plugin/payment/svc"

var Service = NewServiceGroup(svc.ServiceContextApp)

type ServiceGroup struct {
	PaymentService
	Router
	RefundService
}

func NewServiceGroup(svc *svc.ServiceContext) *ServiceGroup {
	return &ServiceGroup{
		PaymentService: NewPaymentService(svc),
		Router:         NewRouterService(svc),
		RefundService:  NewRefundService(svc),
	}
}
