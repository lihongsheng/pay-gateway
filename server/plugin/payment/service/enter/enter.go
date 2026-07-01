package enter

import (
	"github.com/lihongsheng/pay-gateway/config"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/pay-gateway/service/system"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceApi struct {
	MchService              system.MchService
	AppService              service.ApplicationService
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

var ServiceApiApp *ServiceApi

func Init(svc *svc.ServiceContext, db *gorm.DB, domain *domain.ServiceGroup, redis *redis.Client, config config.Config) {
	ServiceApiApp = &ServiceApi{
		MchService:              system.DefaultMch,
		AppService:              service.NewApplicationService(svc.AppRepo, svc.MchRepo),
		PaymentAccount:          service.NewPaymentAccountService(svc.AppRepo, svc.MchRepo, svc.PaymentAccountRepo, db),
		PaymentService:          service.NewPaymentService(svc.PaymentOrderRepo, svc.PaymentAccountRepo, svc.AppRepo, svc.NotifyRepo, svc.Event, domain, svc.Config),
		RefundService:           service.NewRefundService(svc.PaymentOrderRepo, svc.RefundRepo, svc.AppRepo, svc.NotifyRepo, svc.Event, domain, redis, svc.Config),
		Aggregate:               service.NewAggregateService(svc.AppRepo, svc.PaymentAccountRepo, svc.AggregateAdapter, svc.Event, domain.PaymentService, domain.Router, config),
		TradeOrderService:       service.NewTradeOrderService(svc.PaymentOrderRepo, svc.MchRepo, svc.AppRepo, svc.PaymentAccountRepo),
		TradeRefundService:      service.NewTradeRefundService(svc.PaymentOrderRepo, svc.RefundRepo, svc.AppRepo, svc.MchRepo, svc.PaymentAccountRepo, domain.RefundService),
		RouterStatisticsService: service.NewRouterStatisticsService(svc.RouterRepo),
		TradeStatisticsService:  service.NewTradeStatisticsService(svc.MchRepo, svc.AppRepo, svc.PaymentAccountRepo, svc.PaymentOrderRepo, svc.TradeStatisticsRepo, svc.EventRecordRepo, svc.StatisticsRepo),
		CashierService:          service.NewCashierService(svc.AppRepo, svc.PaymentAccountRepo, svc.PaymentOrderRepo, svc.Event, domain.PaymentService, redis, svc.Config),
	}
}
