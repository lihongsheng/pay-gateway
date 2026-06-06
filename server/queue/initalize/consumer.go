package initalize

import (
	"context"
	"github.com/lihongsheng/pay-gateway/global"
	_ "github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/consumer"
	"github.com/lihongsheng/pay-gateway/queue"
	"github.com/lihongsheng/pay-gateway/queue/kafka"
)

func getConsumer() []kafka.ConsumerConfig {
	consumerConfigs := []kafka.ConsumerConfig{
		// 处理微信支付宝回调
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.PaymentCallback.GroupName(),
			Handler: consumer.PaymentCallback,
			Topic:   global.GVA_CONFIG.Kafka.Topic.PaymentCallback,
		},
		// 处理微信支付宝回调
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.RefundCallback.GroupName(),
			Handler: consumer.RefundCallback,
			Topic:   global.GVA_CONFIG.Kafka.Topic.RefundCallback,
		},
		// 统计路由数据
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.RouterStatisticsConsumer.GroupName(),
			Handler: consumer.RouterStatisticsConsumer,
			Topic:   global.GVA_CONFIG.Kafka.Topic.PaymentStatus,
		},
		// 统计支付数据
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.TradeStatistics.GroupName(),
			Handler: consumer.TradeStatistics,
			Topic:   global.GVA_CONFIG.Kafka.Topic.PaymentStatus,
		},
		// 统计支付订单总数
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.OrderTotalStatistics.GroupName(),
			Handler: consumer.OrderTotalStatistics,
			Topic:   global.GVA_CONFIG.Kafka.Topic.PaymentStatus,
		},
		// 统计商户支付订单总数
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.MchTotalStatistics.GroupName(),
			Handler: consumer.MchTotalStatistics,
			Topic:   global.GVA_CONFIG.Kafka.Topic.PaymentStatus,
		},
		// 统计退款数据
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.RefundStatistics.GroupName(),
			Handler: consumer.RefundStatistics,
			Topic:   global.GVA_CONFIG.Kafka.Topic.RefundStatus,
		},
		// 微信支付账户登录被限制
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.UserLimitConsumer.GroupName(),
			Handler: consumer.UserLimitConsumer,
			Topic:   global.GVA_CONFIG.Kafka.Topic.UserLimit,
		},
		// 生成支付三方通知记录
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.PaymentGenNotifyHandler.GroupName(),
			Handler: consumer.PaymentGenNotifyHandler,
			Topic:   global.GVA_CONFIG.Kafka.Topic.PaymentStatus,
		},
		// 生成退款三方通知记录
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.RefundGenNotifyHandler.GroupName(),
			Handler: consumer.RefundGenNotifyHandler,
			Topic:   global.GVA_CONFIG.Kafka.Topic.RefundStatus,
		},
		// 支付通知重试
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.PaymentNotifyHandler.GroupName(),
			Handler: consumer.PaymentNotifyHandler,
			Topic:   global.GVA_CONFIG.Kafka.Topic.PaymentNotifyRetry,
		},
		// 退款通知重试
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.RefundNotifyHandler.GroupName(),
			Handler: consumer.RefundNotifyHandler,
			Topic:   global.GVA_CONFIG.Kafka.Topic.RefundNotifyRetry,
		},
		// 生成超时支付订单记录
		kafka.ConsumerConfig{
			Group:   global.GVA_CONFIG.Env.String() + consumer.PaymentGenExpire.GroupName(),
			Handler: consumer.PaymentGenExpire,
			Topic:   global.GVA_CONFIG.Kafka.Topic.PaymentStatus,
		},
	}
	return consumerConfigs
}

func Init() queue.Server {
	c := kafka.NewConsumer(getConsumer())
	go func() {
		err := c.Start(context.Background())
		if err != nil {
			panic(err)
		}
	}()
	return c
}
