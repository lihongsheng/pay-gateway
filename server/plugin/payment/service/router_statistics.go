package service

import (
	"context"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"go.uber.org/zap"
)

type RouterStatisticsService interface {
	HandlePaymentEvent(ctx context.Context, event *event2.PaymentOrderStatusEvent) error
	HandleUserLimitEvent(ctx context.Context, event *event2.UserLoginLimitEvent) error
}

type routerStatisticsService struct {
	svc *svc.ServiceContext
}

func NewRouterStatisticsService(svc *svc.ServiceContext) RouterStatisticsService {
	return &routerStatisticsService{
		svc: svc,
	}
}

// HandleUserLimitEvent 处理用户登录限制事件
// 微信可能会限制获取用户openid导致支付失败
func (s *routerStatisticsService) HandleUserLimitEvent(ctx context.Context, event *event2.UserLoginLimitEvent) error {
	if !event.UserLimit {
		return nil
	}
	statics := &model.RouterAccountStatistic{
		AppNo:         event.AppNo,
		AccountNo:     event.AccountNo,
		StatisticDate: event.CreateAt,
		UserLimit:     1,
	}
	err := s.svc.RouterRepo.Save(ctx, statics)
	if err != nil {
		cacheErr := s.svc.RouterRepo.UpdateCache(ctx, event.AppNo, event.CreateAt)
		if cacheErr != nil {
			log.WithCtx(ctx).Error("更新路由统计缓存失败", zap.Error(cacheErr), zap.String("appNO", event.AppNo), zap.String("date", event.CreateAt.Format("2006-01-02")))
		}
	}
	return err
}

// HandlePaymentEvent 消费支付事件数据，统计路由数据，用于自动切换支付账户。
// @param ctx
// @param event
// @return error
func (s *routerStatisticsService) HandlePaymentEvent(ctx context.Context, event *event2.PaymentOrderStatusEvent) error {
	if !(event.NewStatus == payment.Status_Pending || event.NewStatus == payment.Status_Success ||
		event.NewStatus == payment.Status_Failed || event.NewStatus == payment.Status_TempFailed) {
		return nil
	}
	statics := s.buildStatistic(event)
	err := s.svc.RouterRepo.Save(ctx, statics)
	if err != nil {
		cacheErr := s.svc.RouterRepo.UpdateCache(ctx, event.AppNo, event.CreateAt)
		if cacheErr != nil {
			log.WithCtx(ctx).Error("更新路由统计缓存失败", zap.Error(cacheErr), zap.String("appNO", event.AppNo), zap.String("date", event.CreateAt.Format("2006-01-02")))
		}
	}
	return err
}

func (s *routerStatisticsService) buildStatistic(event *event2.PaymentOrderStatusEvent) *model.RouterAccountStatistic {
	statics := &model.RouterAccountStatistic{
		AppNo:         event.AppNo,
		AccountNo:     event.AccountNo,
		StatisticDate: event.CreateAt,
	}
	switch event.NewStatus {
	case payment.Status_Pending:
		statics.TotalRequests = 1
	case payment.Status_Success:
		statics.SuccessRequests = 1
		statics.SuccessAmount = event.Amount
	case payment.Status_Failed, payment.Status_TempFailed:
		statics.TotalRequests = 1
		statics.FailureRequests = 1
	}
	if event.ErrCode == errors.ErrPaymentLimit {
		statics.PaymentLimit = 1
	}
	return statics
}
