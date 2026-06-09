package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

// RefundNotifyHandler 回调其他系统的处理
var RefundNotifyHandler = new(refundNotify)

type refundNotify struct{}

func (p *refundNotify) GroupName() string {
	return "refundNotify"
}
func (p *refundNotify) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("refundNotify", zap.Any("panic message", err))
		}
	}()
	l.Info("refundNotify", zap.String("msg", msg))
	var eventDetail event2.RefundNotifyRetryEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("refundNotify", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.RefundService.NotifyHandler(ctx, eventDetail)
	if err != nil {
		l.Error("refundNotify", zap.Error(err))
	}
	return nil
}

func (p *refundNotify) NotifyClose() {

}
