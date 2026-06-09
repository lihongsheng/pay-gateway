package enter

import (
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
)

type ServiceApi struct {
	MchService              service.MchService
	AppService              service.Application
	PaymentAccount          service.PaymentAccountService
	PaymentService          service.PaymentService
	RefundService           service.RefundService
	Aggregate               service.Aggregate
	TradeOrderService       service.TradeOrderService
	TradeRefundService      service.TradeRefundService
	RouterStatisticsService service.RouterStatisticsService
	TradeStatisticsService  service.TradeStatisticsService
	CashierService          service.CashierService
}

var ServiceApiApp = &ServiceApi{
	MchService:              service.NewMchService(svc.ServiceContextApp),
	AppService:              service.NewApplicationService(svc.ServiceContextApp),
	PaymentAccount:          service.NewPaymentAccountService(svc.ServiceContextApp),
	PaymentService:          service.NewPaymentService(svc.ServiceContextApp, domain.Service),
	RefundService:           service.NewRefundService(svc.ServiceContextApp, domain.Service),
	Aggregate:               service.NewAggregateService(svc.ServiceContextApp, domain.Service),
	TradeOrderService:       service.NewTradeOrderService(svc.ServiceContextApp),
	TradeRefundService:      service.NewTradeRefundService(svc.ServiceContextApp, domain.Service),
	RouterStatisticsService: service.NewRouterStatisticsService(svc.ServiceContextApp),
	TradeStatisticsService:  service.NewTradeStatisticsService(svc.ServiceContextApp),
	CashierService:          service.NewCashierService(svc.ServiceContextApp, domain.Service),
}
