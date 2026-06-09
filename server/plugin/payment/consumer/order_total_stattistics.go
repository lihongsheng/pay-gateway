package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

var OrderTotalStatistics = new(orderTotalStatistics)

type orderTotalStatistics struct{}

func (p *orderTotalStatistics) GroupName() string {
	return "orderTotalStatistics"
}
func (p *orderTotalStatistics) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("orderTotalStatistics", zap.Any("panic message", err))
		}
	}()
	l.Info("orderTotalStatistics", zap.String("msg", msg))
	var eventDetail event2.PaymentOrderStatusEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("orderTotalStatistics", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.TradeStatisticsService.AllOrderStatistics(ctx, &eventDetail, p.GroupName())
	if err != nil {
		l.Error("orderTotalStatistics", zap.Error(err))
	}
	return nil
}

func (p *orderTotalStatistics) NotifyClose() {

}
