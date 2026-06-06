package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

// PaymentGenNotifyHandler 回调其他系统
var PaymentGenNotifyHandler = new(paymentGenNotify)

type paymentGenNotify struct {
}

func (p *paymentGenNotify) GroupName() string {
	return "paymentGenNotify"
}
func (p *paymentGenNotify) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("paymentGenNotify", zap.Any("panic message", err))
		}
	}()
	l.Info("paymentGenNotify", zap.String("msg", msg))
	var eventDetail event2.PaymentOrderStatusEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("paymentGenNotify", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.PaymentService.GenNotifyHandler(ctx, eventDetail)
	if err != nil {
		l.Error("paymentGenNotify", zap.Error(err))
	}
	return nil
}

func (p *paymentGenNotify) NotifyClose() {

}
