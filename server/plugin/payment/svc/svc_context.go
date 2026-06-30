package svc

import (
	"github.com/IBM/sarama"
	"github.com/lihongsheng/pay-gateway/plugin/payment/infrastructure"
	adapter2 "github.com/lihongsheng/pay-gateway/plugin/payment/infrastructure/adapter"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/pay-gateway/repo/system"
	"github.com/lihongsheng/payment-sdk/enum/channel"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceContext struct {
	MchRepo             system.MchRepo
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

func NewServiceContext(db *gorm.DB, redis *redis.Client, producer sarama.SyncProducer) *ServiceContext {
	return &ServiceContext{
		MchRepo:             system.NewMchRepo(db),
		AppRepo:             repo.NewApplicationRepo(db, redis),
		PaymentAccountRepo:  repo.NewPaymentAccountRepo(db, redis),
		PaymentOrderRepo:    repo.NewPaymentOrderRepo(db, redis),
		RefundRepo:          repo.NewRefundRepo(db, redis),
		RouterRepo:          repo.NewRouterRepo(db, redis),
		TradeStatisticsRepo: repo.NewTradeStaticsRepo(db),
		EventRecordRepo:     repo.NewEventRecordRepo(db),
		NotifyRepo:          repo.NewNotifyRepo(db),
		StatisticsRepo:      repo.NewStatistics(db),
		AggregateAdapter: map[channel.Channel]adapter2.AggregateUser{
			channel.Channel_Alipay: adapter2.NewAlipayUser(),
			channel.Channel_Wechat: adapter2.NewWechatUser(),
			channel.Channel_Fuiou:  adapter2.NewFuiouUser(),
		},
		Event: infrastructure.NewEvent(producer),
	}
}
