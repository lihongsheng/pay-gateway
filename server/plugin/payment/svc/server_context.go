package svc

import (
	"github.com/lihongsheng/pay-gateway/plugin/payment/infrastructure"
	adapter2 "github.com/lihongsheng/pay-gateway/plugin/payment/infrastructure/adapter"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/payment-sdk/enum/channel"
)

type ServiceContext struct {
	MchRepo             repo.MchRepo
	AppRepo             repo.ApplicationRepo
	PaymentAccountRepo  repo.PaymentAccountRepo
	PaymentOrderRepo    repo.PaymentOrderRepo
	RefundRepo          repo.RefundRepo
	RouterRepo          repo.RouterRepo
	AggregateAdapter    map[channel.Channel]adapter2.AggregateUser
	Event               infrastructure.Event
	TradeStatisticsRepo repo.TradeStaticsRepo
	NotifyRepo          repo.NotifyRepo
	EventRecordRepo     repo.EventRecordRepo
	StatisticsRepo      repo.Statistics
}

var ServiceContextApp = &ServiceContext{
	MchRepo:             repo.NewMchRepo(),
	AppRepo:             repo.NewApplicationRepo(),
	PaymentAccountRepo:  repo.NewPaymentAccountRepo(),
	PaymentOrderRepo:    repo.NewPaymentOrderRepo(),
	RefundRepo:          repo.NewRefundRepo(),
	RouterRepo:          repo.NewRouterRepo(),
	TradeStatisticsRepo: repo.NewTradeStaticsRepo(),
	EventRecordRepo:     repo.NewEventRecordRepo(),
	NotifyRepo:          repo.NewNotifyRepo(),
	StatisticsRepo:      repo.NewStatistics(),
	AggregateAdapter: map[channel.Channel]adapter2.AggregateUser{
		channel.Channel_Alipay: adapter2.NewAlipayUser(),
		channel.Channel_Wechat: adapter2.NewWechatUser(),
		channel.Channel_Fuiou:  adapter2.NewFuiouUser(),
	},
	Event: infrastructure.NewEvent(),
}
