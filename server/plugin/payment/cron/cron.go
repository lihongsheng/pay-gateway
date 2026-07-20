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

type Handler struct {
	svc *svc.ServiceContext
}

func NewJob() *Handler {
	return &Handler{
		svc: svc.ServiceContextApp,
	}
}

// DeleteEventRecord 删除事件处理记录
func (c *Handler) DeleteEventRecord(ctx context.Context) error {
	r, err := c.svc.Redis.SetNX(ctx, "deleteEventRecordJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "deleteEventRecordJob")
	err = c.svc.EventRecordRepo.Delete(ctx, time.Now().Add(-time.Hour*24*3))
	if err != nil {
		log.WithCtx(ctx).Error("DeleteEventRecord", zap.Error(err))
		return err
	}
	return nil
}

// DeletePaymentRequestLog 删除支付请求日志
func (c *Handler) DeletePaymentRequestLog(ctx context.Context) error {
	r, _ := c.svc.Redis.SetNX(ctx, "deletePaymentLogJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "deletePaymentLogJob")
	err := c.svc.PaymentOrderRepo.DeletePaymentRequestLog(ctx, time.Now().Add(-time.Hour*24*7))
	if err != nil {
		log.WithCtx(ctx).Error("DeletePaymentRequestLog", zap.Error(err))
		return err
	}
	return nil
}

// DeleteNotifyRecord 删除通知记录
func (c *Handler) DeleteNotifyRecord(ctx context.Context) error {
	r, _ := c.svc.Redis.SetNX(ctx, "deleteNotifyLogJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "deleteNotifyLogJob")
	err := c.svc.NotifyRepo.Delete(ctx, time.Now().Add(-time.Hour*24))
	if err != nil {
		log.WithCtx(ctx).Error("DeleteNotifyRecord", zap.Error(err))
		return err
	}
	return nil
}

// CronPaymentPushNotifyHandler 支付通知重试
func (c *Handler) CronPaymentPushNotifyHandler(ctx context.Context) error {
	r, _ := c.svc.Redis.SetNX(ctx, "CronPaymentPushNotifyHandler", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "CronPaymentPushNotifyHandler")
	l := log.WithCtx(ctx)
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)
	count, err := c.svc.NotifyRepo.CountRetry(ctx, start, end, enum.NotifyType_Payment)
	if err != nil || count == 0 {
		return nil
	}
	var lastId int64
	for {
		records, err := c.svc.NotifyRepo.GetRetry(ctx, start, end, lastId, 100, enum.NotifyType_Payment)
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
func (c *Handler) CronRefundPushNotifyHandler(ctx context.Context) error {
	r, _ := c.svc.Redis.SetNX(ctx, "CronRefundPushNotifyHandler", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "CronRefundPushNotifyHandler")
	l := log.WithCtx(ctx)
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)
	count, err := c.svc.NotifyRepo.CountRetry(ctx, start, end, enum.NotifyType_Refund)
	if err != nil || count == 0 {
		return nil
	}
	var lastId int64
	for {
		records, err := c.svc.NotifyRepo.GetRetry(ctx, start, end, lastId, 100, enum.NotifyType_Refund)
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
	return nil
}

// DeleteExpireNotifyRecord 删除过期通知
func (c *Handler) DeleteExpireNotifyRecord(ctx context.Context) error {
	r, _ := c.svc.Redis.SetNX(ctx, "DeleteNotifyRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "DeleteNotifyRecord")
	err := c.svc.NotifyRepo.Delete(ctx, time.Now().Add(-time.Hour*24))
	if err != nil {
		log.WithCtx(ctx).Error("DeleteExpireNotifyRecord", zap.Error(err))
		return err
	}
	return nil

}

// PaymentExpireRecordJob 处理支付过期订单
func (c *Handler) PaymentExpireRecordJob(ctx context.Context) error {
	r, _ := c.svc.Redis.SetNX(ctx, "PaymentExpireRecordJob", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "PaymentExpireRecordJob")
	l := log.WithCtx(ctx)
	now := time.Now()
	start := now.Add(-time.Hour)
	end := now
	var lastId int64
	for {
		records, err := c.svc.PaymentOrderRepo.GetPaymentExpireRecord(ctx, start, end, 100, lastId)
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
	return nil
}

// DeleteExpireRecord 删除过期支付记录
func (c *Handler) DeleteExpireRecord(ctx context.Context) error {
	r, _ := c.svc.Redis.SetNX(ctx, "DeleteExpireRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "DeleteExpireRecord")
	err := c.svc.PaymentOrderRepo.DeletePaymentExpireRecord(ctx, time.Now().Add(-time.Minute*30))
	if err != nil {
		log.WithCtx(ctx).Error("DeleteExpireRecord", zap.Error(err))
		return err
	}
	return nil

}

// DeleteExpireRouterRecord 删除过期路由记录
func (c *Handler) DeleteExpireRouterRecord(ctx context.Context) error {
	r, _ := c.svc.Redis.SetNX(ctx, "DeleteExpireRouterRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	if !r {
		return nil
	}
	defer c.svc.Redis.Del(ctx, "DeleteExpireRouterRecord")
	err := c.svc.RouterRepo.DeleteByDate(ctx, time.Now().Add(-time.Hour*24*30))
	if err != nil {
		log.WithCtx(ctx).Error("DeleteExpireRouterRecord", zap.Error(err))
		return err
	}
	return nil
}
