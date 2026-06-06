package cron

import (
	"context"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"go.uber.org/zap"
	"math"
	"time"
)

var Job = new(job)

type job struct{}

func (j *job) DeleteEventRecord(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "deleteEventRecordJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "deleteEventRecordJob")
	}()
	if !r {
		return nil
	}
	err := svc.ServiceContextApp.EventRecordRepo.Delete(NewCtx, time.Now().Add(-time.Hour*24*3))
	if err != nil {
		l.Error("删除事件记录失败", zap.Error(err))
	}
	return err
}

func (j *job) DeletePaymentLog(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "deletePaymentLogJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "deletePaymentLogJob")
	}()
	if !r {
		return nil
	}
	err := svc.ServiceContextApp.PaymentOrderRepo.DeletePaymentRequestLog(NewCtx, time.Now().Add(-time.Hour*24*7))
	if err != nil {
		l.Error("DeletePaymentLog", zap.Error(err))
	}
	return err
}

func (j *job) DeleteNotifyLog(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "deleteNotifyLogJob", fmt.Sprintf("%d", time.Now().Unix()), 2*time.Minute).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "deleteNotifyLogJob")
	}()
	if !r {
		return nil
	}
	err := svc.ServiceContextApp.NotifyRepo.Delete(NewCtx, time.Now().Add(-time.Hour*24))
	if err != nil {
		l.Error("DeleteNotifyLog", zap.Error(err))
	}
	return err
}

// CronPaymentPushNotifyHandler 定时处理回调
func (j *job) CronPaymentPushNotifyHandler(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "CronPaymentPushNotifyHandler", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "CronPaymentPushNotifyHandler")
	}()
	if !r {
		return nil
	}
	end := time.Now()
	start := end.Add(-time.Second * 15)
	page := 0
	count, err := svc.ServiceContextApp.NotifyRepo.CountRetry(NewCtx, start, end, enum.NotifyType_Refund)
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	maxPage := int(math.Ceil(float64(count) / 100))
	lastId := int64(0)
	for {
		if page > maxPage {
			break
		}
		page += 1
		records, err := svc.ServiceContextApp.NotifyRepo.GetRetry(NewCtx, start, end, lastId, 100, enum.NotifyType_Payment)
		if err != nil {
			l.Error("CronPaymentPushNotifyHandler", zap.Error(err), zap.Int64("lastID", lastId))
			return err
		}
		if len(records) == 0 {
			l.Info("CronPaymentPushNotifyHandler", zap.Int64("lastID", lastId))
			return nil
		}
		if lastId > 0 && lastId == records[len(records)-1].ID {
			l.Info("CronPaymentPushNotifyHandler", zap.Int64("lastID", lastId))
			return nil
		}
		for _, record := range records {
			// 扔到队列里，防止阻塞
			err := svc.ServiceContextApp.Event.Publish(ctx, &event.PaymentNotifyRetryEvent{
				Base: event.Base{
					EventType: event.EventTypePaymentNotifyRetry,
					EventId:   fmt.Sprintf("Notify:%d", record.ID),
				},
				NotifyID:    record.ID,
				AppNo:       record.AppNo,
				MchNo:       record.MchNo,
				OrderNo:     record.OutNo,
				TradeNo:     record.TradeNo,
				NotifyCount: int(record.NotifyCount),
			}, global.GVA_CONFIG.Kafka.Topic.PaymentNotifyRetry)
			if err != nil {
				l.Error("发送支付回调失败", zap.Error(err), zap.String("order_no", record.OutNo))
				return err
			}
			lastId = record.ID
		}

	}
	return nil
}

// CronRefundPushNotifyHandler 定时处理回调
func (j *job) CronRefundPushNotifyHandler(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "CronRefundPushNotifyHandler", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "CronRefundPushNotifyHandler")
	}()
	if !r {
		return nil
	}
	end := time.Now()
	start := end.Add(-time.Second * 15)
	page := 0
	count, err := svc.ServiceContextApp.NotifyRepo.CountRetry(NewCtx, start, end, enum.NotifyType_Refund)
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	maxPage := int(math.Ceil(float64(count) / 100))
	lastId := int64(0)
	for {
		if page > maxPage {
			break
		}
		page += 1
		records, err := svc.ServiceContextApp.NotifyRepo.GetRetry(NewCtx, start, end, lastId, 100, enum.NotifyType_Refund)
		if err != nil {
			l.Error("CronRefundPushNotifyHandler", zap.Error(err), zap.Int64("lastID", lastId))
			return err
		}
		if len(records) == 0 {
			l.Info("CronRefundPushNotifyHandler", zap.Int64("lastID", lastId))
			return nil
		}
		if lastId > 0 && lastId == records[len(records)-1].ID {
			l.Info("CronRefundPushNotifyHandler", zap.Int64("lastID", lastId))
			return nil
		}
		for _, record := range records {
			// 扔到队列里，防止阻塞
			err := svc.ServiceContextApp.Event.Publish(ctx, &event.RefundNotifyRetryEvent{
				Base: event.Base{
					EventType: event.EventTypePaymentNotifyRetry,
					EventId:   fmt.Sprintf("Notify:%d", record.ID),
				},
				NotifyID:      record.ID,
				AppNo:         record.AppNo,
				MchNo:         record.MchNo,
				RefundNo:      record.OutNo,
				RefundTradeNo: record.TradeNo,
				NotifyCount:   int(record.NotifyCount),
			}, global.GVA_CONFIG.Kafka.Topic.RefundNotifyRetry)
			if err != nil {
				l.Error("发送支付回调失败", zap.Error(err), zap.String("order_no", record.OutNo))
				return err
			}
			lastId = record.ID
		}
	}
	return nil
}

func (j *job) DeleteNotifyRecord(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "DeleteNotifyRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "DeleteNotifyRecord")
	}()
	if !r {
		return nil
	}
	err := svc.ServiceContextApp.NotifyRepo.Delete(NewCtx, time.Now().Add(-time.Hour*24))
	if err != nil {
		l.Error("DeleteNotifyRecord", zap.Error(err))
	}
	return err
}

func (j *job) PaymentExpireRecord(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "PaymentExpireRecordJob", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "PaymentExpireRecordJob")
	}()
	if !r {
		return nil
	}
	end := time.Now().Add(-enum.DefaultExpireTime)
	start := end.Add(-4 * enum.DefaultExpireTime)
	lastId := int64(0)
	for {
		records, err := svc.ServiceContextApp.PaymentOrderRepo.GetPaymentExpireRecord(NewCtx, start, end, 100, lastId)
		if err != nil {
			l.Error("PaymentExpireRecord", zap.Error(err), zap.Int64("lastID", lastId))
			return err
		}
		if len(records) == 0 {
			l.Info("PaymentExpireRecordComplete", zap.Int64("lastID", lastId))
			return nil
		}
		if lastId > 0 && lastId == records[len(records)-1].ID {
			l.Info("PaymentExpireRecordLast", zap.Int64("lastID", lastId))
			return nil
		}
		deleteIds := make([]int64, 0, len(records))
		for _, record := range records {
			payOrder, err := svc.ServiceContextApp.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
				OrderNo: record.OrderNo,
				TradeNo: "",
				AppNo:   record.AppNo,
				MchNo:   record.MchNo,
			})
			if err != nil {
				l.Error("更新订单状态失败", zap.Error(err), zap.Any("record", record))
				return err
			}
			if payOrder.IsTimeout() {
				err := svc.ServiceContextApp.PaymentOrderRepo.UpdateStatus(ctx, payOrder, payOrder.Status, map[string]interface{}{
					"status":      payment.Status_TimeOut,
					"status_from": "cronJob",
				})
				if err != nil {
					l.Error("更新订单状态失败", zap.Error(err), zap.Any("record", record))
					return err
				}
				deleteIds = append(deleteIds, record.ID)
			} else {
				deleteIds = append(deleteIds, record.ID)
			}
			fmt.Println(deleteIds)
			lastId = record.ID
		}
		if len(deleteIds) > 0 {
			_ = svc.ServiceContextApp.PaymentOrderRepo.DeletePaymentExpireRecordByIDs(ctx, deleteIds)
		}
	}
}

func (j *job) DeleteExpireRecord(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "DeleteExpireRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "DeleteExpireRecord")
	}()
	if !r {
		return nil
	}
	err := svc.ServiceContextApp.PaymentOrderRepo.DeletePaymentExpireRecord(NewCtx, time.Now().Add(-time.Minute*30))
	if err != nil {
		l.Error("DeleteExpireRecord", zap.Error(err))
	}
	return err
}

func (j *job) DeleteExpireRouterRecord(ctx context.Context) error {
	NewCtx, l := log.NewCtx(ctx)
	// 加锁防止多POD 同时运行
	r, _ := global.GVA_REDIS.SetNX(ctx, "DeleteExpireRouterRecord", fmt.Sprintf("%d", time.Now().Unix()), time.Hour).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, "DeleteExpireRouterRecord")
	}()
	if !r {
		return nil
	}
	err := svc.ServiceContextApp.RouterRepo.DeleteByDate(NewCtx, time.Now().Add(-time.Hour*24*30))
	if err != nil {
		l.Error("DeleteExpireRouterRecord", zap.Error(err))
	}
	return err
}
