package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

var PaymentGenExpire = new(paymentGenExpire)

type paymentGenExpire struct {
}

func (p *paymentGenExpire) GroupName() string {
	return "paymentGenExpire"
}
func (p *paymentGenExpire) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("paymentGenExpire", zap.Any("panic message", err))
		}
	}()
	l.Info("paymentGenExpire", zap.String("msg", msg))
	var eventDetail event2.PaymentOrderStatusEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("paymentGenExpire", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.PaymentService.GenPaymentExpireRecord(ctx, eventDetail)
	if err != nil {
		l.Error("paymentGenExpire", zap.Error(err))
	}
	return nil
}

func (p *paymentGenExpire) NotifyClose() {

}
