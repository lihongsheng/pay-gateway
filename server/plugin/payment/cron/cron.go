package cron

import (
	"context"
	"fmt"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"time"

	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"go.uber.org/zap"
)

var NewCtx = context.Background()

type CronService struct {
	svc *svc.ServiceContext
}

func NewCronService(svc *svc.ServiceContext) *CronService {
	return &CronService{
		svc: svc,
	}
}

// DeleteEventRecord 删除事件处理记录
func (c *CronService) DeleteEventRecord(ctx context.Context) {
	r, _ := c.svc.Redis.SetNX(NewCtx, "deleteEventRecordJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "deleteEventRecordJob")
	err := c.svc.EventRecordRepo.Delete(NewCtx, time.Now().Add(-time.Hour*24*3))
	if err != nil {
		log.WithCtx(NewCtx).Error("DeleteEventRecord", zap.Error(err))
	}
}

// DeletePaymentRequestLog 删除支付请求日志
func (c *CronService) DeletePaymentRequestLog(ctx context.Context) {
	r, _ := c.svc.Redis.SetNX(NewCtx, "deletePaymentLogJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "deletePaymentLogJob")
	err := c.svc.PaymentOrderRepo.DeletePaymentRequestLog(NewCtx, time.Now().Add(-time.Hour*24*7))
	if err != nil {
		log.WithCtx(NewCtx).Error("DeletePaymentRequestLog", zap.Error(err))
	}
}

// DeleteNotifyRecord 删除通知记录
func (c *CronService) DeleteNotifyRecord(ctx context.Context) {
	r, _ := c.svc.Redis.SetNX(NewCtx, "deleteNotifyLogJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "deleteNotifyLogJob")
	err := c.svc.NotifyRepo.Delete(NewCtx, time.Now().Add(-time.Hour*24))
	if err != nil {
		log.WithCtx(NewCtx).Error("DeleteNotifyRecord", zap.Error(err))
	}
}

// CronPaymentPushNotifyHandler 支付通知重试
func (c *CronService) CronPaymentPushNotifyHandler( (ctx context.Context) {
	r, _ := c.svc.Redis.SetNX(NewCtx, "CronPaymentPushNotifyHandler", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "CronPaymentPushNotifyHandler")
	l := log.WithCtx(NewCtx)
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)
	count, err := c.svc.NotifyRepo.CountRetry(NewCtx, start, end, enum.NotifyType_Payment)
	if err != nil || count == 0 {
		return
	}
	var lastId int64
	for {
		records, err := c.svc.NotifyRepo.GetRetry(NewCtx, start, end, lastId, 100, enum.NotifyType_Payment)
		if err != nil || len(records) == 0 {
			break
		}
		for _, record := range records {
			err := c.svc.Event.Publish(ctx, &event.PaymentNotifyRetryEvent{
				NotifyID:    record.ID,
				AppNo:       record.AppNo,
				MchNo:       record.MchNo,
				OrderNo:     record.OutNo,
				TradeNo:     record.TradeNo,
				NotifyCount: int(record.NotifyCount),
			}, c.svc.Config.Topic.PaymentNotifyRetry)
			if err != nil {
				l.Error("发送支付回调失败", zap.Error(err), zap.String("order_no", record.OutNo))
			}
		}
		lastId = records[len(records)-1].ID
		if len(records) < 100 {
			break
		}
	}
	return nil
}

// CronRefundPushNotifyHandler 退款通知重试
func (c *CronService) CronRefundPushNotifyHandler(ctx context.Context) {
	r, _ := c.svc.Redis.SetNX(NewCtx, "CronRefundPushNotifyHandler", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "CronRefundPushNotifyHandler")
	l := log.WithCtx(NewCtx)
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)
	count, err := c.svc.NotifyRepo.CountRetry(NewCtx, start, end, enum.NotifyType_Refund)
	if err != nil || count == 0 {
		return
	}
	var lastId int64
	for {
		records, err := c.svc.NotifyRepo.GetRetry(NewCtx, start, end, lastId, 100, enum.NotifyType_Refund)
		if err != nil || len(records) == 0 {
			break
		}
		for _, record := range records {
			err := c.svc.Event.Publish(ctx, &event.RefundNotifyRetryEvent{
				NotifyID:    record.ID,
				AppNo:       record.AppNo,
				MchNo:       record.MchNo,
				RefundNo:    record.OutNo,
				TradeNo:     record.TradeNo,
				NotifyCount: int(record.NotifyCount),
			}, c.svc.Config.Topic.RefundNotifyRetry)
			if err != nil {
				l.Error("发送退款回调失败", zap.Error(err), zap.String("refund_no", record.OutNo))
			}
		}
		lastId = records[len(records)-1].ID
		if len(records) < 100 {
			break
		}
	}
}

// DeleteExpireNotifyRecord 删除过期通知
func (c *CronService) DeleteExpireNotifyRecord(ctx context.Context) {
	r, _ := c.svc.Redis.SetNX(NewCtx, "DeleteNotifyRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "DeleteNotifyRecord")
	err := c.svc.NotifyRepo.Delete(NewCtx, time.Now().Add(-time.Hour*24))
	if err != nil {
		log.WithCtx(NewCtx).Error("DeleteExpireNotifyRecord", zap.Error(err))
	}
}

// PaymentExpireRecordJob 处理支付过期订单
func (c *CronService) PaymentExpireRecordJob(ctx context.Context) {
	r, _ := c.svc.Redis.SetNX(NewCtx, "PaymentExpireRecordJob", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "PaymentExpireRecordJob")
	l := log.WithCtx(NewCtx)
	now := time.Now()
	start := now.Add(-time.Hour)
	end := now
	var lastId int64
	for {
		records, err := c.svc.PaymentOrderRepo.GetPaymentExpireRecord(NewCtx, start, end, 100, lastId)
		if err != nil || len(records) == 0 {
			break
		}
		var deleteIds []int64
		for _, record := range records {
			payOrder, err := c.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
				OrderNo: record.OrderNo,
				AppNo:   record.AppNo,
				MchNo:   record.MchNo,
			})
			if err != nil {
				continue
			}
			if payOrder.IsTimeout() {
				err := c.svc.PaymentOrderRepo.UpdateStatus(ctx, payOrder, payOrder.Status, map[string]interface{}{
					"status":      "timeout",
					"status_from": "cronExpire",
				})
				if err != nil {
					l.Error("更新订单为关闭状态失败", zap.Error(err), zap.String("order_no", record.OrderNo))
				}
			}
			deleteIds = append(deleteIds, record.ID)
		}
		if len(deleteIds) > 0 {
			_ = c.svc.PaymentOrderRepo.DeletePaymentExpireRecordByIDs(ctx, deleteIds)
		}
		lastId = records[len(records)-1].ID
		if len(records) < 100 {
			break
		}
	}
}

// DeleteExpireRecord 删除过期支付记录
func (c *CronService) DeleteExpireRecord() {
	r, _ := c.svc.Redis.SetNX(NewCtx, "DeleteExpireRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "DeleteExpireRecord")
	err := c.svc.PaymentOrderRepo.DeletePaymentExpireRecord(NewCtx, time.Now().Add(-time.Minute*30))
	if err != nil {
		log.WithCtx(NewCtx).Error("DeleteExpireRecord", zap.Error(err))
	}
}

// DeleteExpireRouterRecord 删除过期路由记录
func (c *CronService) DeleteExpireRouterRecord(ctx context.Context) {
	r, _ := c.svc.Redis.SetNX(NewCtx, "DeleteExpireRouterRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return
	}
	defer c.svc.Redis.Del(NewCtx, "DeleteExpireRouterRecord")
	err := c.svc.RouterRepo.DeleteByDate(NewCtx, time.Now().Add(-time.Hour*24*30))
	if err != nil {
		log.WithCtx(NewCtx).Error("DeleteExpireRouterRecord", zap.Error(err))
	}
}
