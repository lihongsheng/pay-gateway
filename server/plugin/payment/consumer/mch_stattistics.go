package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

var MchTotalStatistics = new(mchTotalStatistics)

type mchTotalStatistics struct{}

func (p *mchTotalStatistics) GroupName() string {
	return "mchTotalStatistics"
}
func (p *mchTotalStatistics) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("mchTotalStatistics", zap.Any("panic message", err))
		}
	}()
	l.Info("mchTotalStatistics", zap.String("msg", msg))
	var eventDetail event2.PaymentOrderStatusEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("mchTotalStatistics", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.TradeStatisticsService.MchAllOrderStatistics(ctx, &eventDetail, p.GroupName())
	if err != nil {
		l.Error("mchTotalStatistics", zap.Error(err))
	}
	return nil
}

func (p *mchTotalStatistics) NotifyClose() {

}
