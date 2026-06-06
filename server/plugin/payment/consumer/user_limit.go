package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

var UserLimitConsumer = new(userLimitConsumer)

type userLimitConsumer struct{}

func (p *userLimitConsumer) GroupName() string {
	return "userLimit"
}
func (p *userLimitConsumer) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("userLimitConsumer", zap.Any("panic message", err))
		}
	}()
	l.Info("userLimitConsumer", zap.String("msg", msg))
	var eventDetail event2.UserLoginLimitEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("userLimitConsumer", zap.Error(err))
		return nil
	}
	err = enter.ServiceApiApp.RouterStatisticsService.HandleUserLimitEvent(ctx, &eventDetail)
	if err != nil {
		l.Error("userLimitConsumer", zap.Error(err))
	}
	return nil
}

func (p *userLimitConsumer) NotifyClose() {

}
