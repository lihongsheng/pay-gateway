package initalize

import (
	"github.com/lihongsheng/pay-gateway/global"
	_ "github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/consumer"
	"github.com/lihongsheng/pay-gateway/queue/kafka"
)

func GetKafkaConsumer() []kafka.ConsumerConfig {
	consumerConfigs := []kafka.ConsumerConfig{
		// 处理微信支付宝回调
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.PaymentCallback.GroupName(),
			Handler: consumer.PaymentCallback,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.PaymentCallback,
		},
		// 处理微信支付宝回调
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.RefundCallback.GroupName(),
			Handler: consumer.RefundCallback,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.RefundCallback,
		},
		// 统计路由数据
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.RouterStatisticsConsumer.GroupName(),
			Handler: consumer.RouterStatisticsConsumer,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.PaymentStatus,
		},
		// 统计支付数据
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.TradeStatistics.GroupName(),
			Handler: consumer.TradeStatistics,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.PaymentStatus,
		},
		// 统计支付订单总数
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.OrderTotalStatistics.GroupName(),
			Handler: consumer.OrderTotalStatistics,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.PaymentStatus,
		},
		// 统计商户支付订单总数
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.MchTotalStatistics.GroupName(),
			Handler: consumer.MchTotalStatistics,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.PaymentStatus,
		},
		// 统计退款数据
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.RefundStatistics.GroupName(),
			Handler: consumer.RefundStatistics,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.RefundStatus,
		},
		// 微信支付账户登录被限制
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.UserLimitConsumer.GroupName(),
			Handler: consumer.UserLimitConsumer,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.UserLimit,
		},
		// 生成支付三方通知记录
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.PaymentGenNotifyHandler.GroupName(),
			Handler: consumer.PaymentGenNotifyHandler,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.PaymentStatus,
		},
		// 生成退款三方通知记录
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.RefundGenNotifyHandler.GroupName(),
			Handler: consumer.RefundGenNotifyHandler,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.RefundStatus,
		},
		// 支付通知重试
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.PaymentNotifyHandler.GroupName(),
			Handler: consumer.PaymentNotifyHandler,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.PaymentNotifyRetry,
		},
		// 退款通知重试
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.RefundNotifyHandler.GroupName(),
			Handler: consumer.RefundNotifyHandler,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.RefundNotifyRetry,
		},
		// 生成超时支付订单记录
		kafka.ConsumerConfig{
			Group:   global.Cfg.Env.String() + consumer.PaymentGenExpire.GroupName(),
			Handler: consumer.PaymentGenExpire,
			Topic:   global.Cfg.Plugin.PaymentConfig.Topic.PaymentStatus,
		},
	}
	return consumerConfigs
}
