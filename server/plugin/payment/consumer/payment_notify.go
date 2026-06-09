package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

// PaymentNotifyHandler 回调其他系统
var PaymentNotifyHandler = new(paymentNotify)

type paymentNotify struct {
}

func (p *paymentNotify) GroupName() string {
	return "paymentNotify"
}
func (p *paymentNotify) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("paymentNotify", zap.Any("panic message", err))
		}
	}()
	l.Info("paymentNotify", zap.String("msg", msg))
	var eventDetail event2.PaymentNotifyRetryEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("paymentNotify", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.PaymentService.NotifyHandler(ctx, eventDetail)
	if err != nil {
		l.Error("paymentNotify", zap.Error(err))
	}
	return nil
}

func (p *paymentNotify) NotifyClose() {

}
