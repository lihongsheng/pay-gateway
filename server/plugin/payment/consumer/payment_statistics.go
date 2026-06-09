package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

var TradeStatistics = new(tradeStatistics)

type tradeStatistics struct{}

func (p *tradeStatistics) GroupName() string {
	return "tradeStatistics"
}
func (p *tradeStatistics) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("tradeStatistics", zap.Any("panic message", err))
		}
	}()
	l.Info("tradeStatistics", zap.String("msg", msg))
	var eventDetail event2.PaymentOrderStatusEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("tradeStatistics", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.TradeStatisticsService.HandlePaymentEvent(ctx, &eventDetail, p.GroupName())
	if err != nil {
		l.Error("tradeStatistics", zap.Error(err))
	}
	return nil
}

func (p *tradeStatistics) NotifyClose() {

}
