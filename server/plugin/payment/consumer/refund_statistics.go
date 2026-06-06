package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

var RefundStatistics = new(refundStatistics)

type refundStatistics struct{}

func (p *refundStatistics) GroupName() string {
	return "refundStatistics"
}
func (p *refundStatistics) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("refundStatistics", zap.Any("panic message", err))
		}
	}()
	l.Info("refundStatistics", zap.String("msg", msg))
	var eventDetail event2.RefundOrderStatusEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("refundStatistics", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.TradeStatisticsService.HandleRefundEvent(ctx, &eventDetail, p.GroupName())
	if err != nil {
		l.Error("refundStatistics", zap.Error(err))
	}
	return nil
}

func (p *refundStatistics) NotifyClose() {

}
