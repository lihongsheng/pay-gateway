package consumer

// 消费支付事件数据，统计路由数据，用于自动切换支付账户。

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

var RouterStatisticsConsumer = new(routerStatisticsConsumer)

type routerStatisticsConsumer struct{}

func (p *routerStatisticsConsumer) GroupName() string {
	return "routerStatistics"
}
func (p *routerStatisticsConsumer) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("routerStatisticsConsumer", zap.Any("panic message", err))
		}
	}()
	l.Info("routerStatisticsConsumer", zap.String("msg", msg))
	var eventDetail event2.PaymentOrderStatusEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("routerStatisticsConsumer", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.RouterStatisticsService.HandlePaymentEvent(ctx, &eventDetail)
	if err != nil {
		l.Error("routerStatisticsConsumer", zap.Error(err))
	}
	return nil
}

func (p *routerStatisticsConsumer) NotifyClose() {

}
