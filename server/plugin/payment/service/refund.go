package service

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"

	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	enum2 "github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/infrastructure"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	paySdk "github.com/lihongsheng/payment-sdk"
	channel2 "github.com/lihongsheng/payment-sdk/enum/channel"
	refund2 "github.com/lihongsheng/payment-sdk/enum/refund"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RefundService interface {
	Refund(ctx context.Context, req *public.RefundRequest) (*public.RefundResponse, error)
	Query(ctx context.Context, req *public.RefundQueryRequest) (*public.RefundResponse, error)
	Callback(ctx context.Context, req *http.Request, appNo string, refundTradeNo string) (string, error)
	PublishCallback(ctx context.Context, req *http.Request, channel string, mchNo string, appNo string, refundTradeNo string) (response string, err error)
	// GenNotifyHandler 回调其他系统
	GenNotifyHandler(ctx context.Context, req event.RefundOrderStatusEvent) error
	NotifyHandler(ctx context.Context, req event.RefundNotifyRetryEvent) error
}

type refundService struct {
	paymentOrderRepo repo.PaymentOrderRepo
	refundRepo       repo.RefundRepo
	appRepo          repo.ApplicationRepo
	notifyRepo       repo.NotifyRepo
	event            infrastructure.Event
	domain           *domain.ServiceGroup
	redis            *redis.Client
}

func NewRefundService(
	paymentOrderRepo repo.PaymentOrderRepo,
	refundRepo repo.RefundRepo,
	appRepo repo.ApplicationRepo,
	notifyRepo repo.NotifyRepo,
	event infrastructure.Event,
	domain *domain.ServiceGroup,
	redis *redis.Client,
) RefundService {
	return &refundService{
		paymentOrderRepo: paymentOrderRepo,
		refundRepo:       refundRepo,
		appRepo:          appRepo,
		notifyRepo:       notifyRepo,
		event:            event,
		domain:           domain,
		redis:            redis,
	}
}

// DefaultRefund 包级单例
var DefaultRefund RefundService

func (s *refundService) NotifyHandler(ctx context.Context, req event.RefundNotifyRetryEvent) error {
	l := log.WithCtx(ctx)
	refund, err := s.refundRepo.GetFromCache(ctx, public.RefundQueryRequest{
		RefundNo: req.RefundNo,
		AppNo:    req.AppNo,
		MchNo:    req.MchNo,
	})
	if err != nil {
		return err
	}

	return s.notify(ctx, l, refund, req.NotifyCount, req.NotifyID)
}

func (s *refundService) GenNotifyHandler(ctx context.Context, req event.RefundOrderStatusEvent) error {
	if !(req.NewStatus == refund2.Status_Success) {
		return nil
	}
	key := "RefundNotify:" + req.EventId
	lock, err := s.redis.SetNX(ctx, key, time.Now().Unix(), enum2.PaymentLockExpire).Result()
	defer func() {
		s.redis.Del(ctx, key)
	}()
	if !lock && err == nil {
		return errors.NewError(errors.ErrPayPending, "已在处理，请稍后再试")
	}
	l := log.WithCtx(ctx)
	refund, err := s.refundRepo.GetFromCache(ctx, public.RefundQueryRequest{
		RefundNo: req.RefundNo,
		AppNo:    req.AppNo,
		MchNo:    req.MchNo,
	})
	if err != nil {
		return err
	}
	if refund.NotifyURL == "" {
		l.Warn("回调地址不存在", zap.Any("req", req))
		return nil
	}
	notify, err := s.notifyRepo.Get(ctx, req.MchNo, req.AppNo, enum2.NotifyType_Payment, req.RefundNo)
	if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if notify == nil {
		notify = s.buildNotifyRecord(req, refund)
		err := s.notifyRepo.Save(ctx, notify)
		if err != nil {
			return err
		}
	}

	return s.notify(ctx, l, refund, 1, notify.ID)
}

func (s *refundService) notify(ctx context.Context, l *zap.Logger, refund *model.RefundOrder, notifyCount int, notifyId int64) error {
	resp, err := utils.HttpJsonPost(ctx, refund.NotifyURL, s.buildSuccessNotifyRequest(refund), 30*time.Second)
	if err != nil {
		l.Error("退款回调处理失败", zap.Error(err), zap.String("refundTradeNo", refund.RefundTradeNo), zap.String("app_no", refund.AppNo), zap.String("mch_no", refund.MchNo))
	}
	if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
		if notifyCount < enum2.NotifyRetryMax {
			_ = s.notifyRepo.Retry(ctx, notifyId, time.Now().Add(time.Duration(notifyCount+1)*enum2.NotifyRetryStep))
			return err
		} else {
			_ = s.notifyRepo.UpdateStatus(ctx, notifyId, enum2.NotifyStatus_Fail, "已达最大重试次数")
			_ = s.refundRepo.ConfirmNotifyStatus(ctx, refund.ID, enum2.NotifyStatus_Fail)
		}
	}
	if resp != nil && resp.StatusCode == http.StatusOK {
		err = s.notifyRepo.UpdateStatus(ctx, notifyId, enum2.NotifyStatus_Success, fmt.Sprintf("httpStatus:%d", http.StatusOK))
		_ = s.refundRepo.ConfirmNotifyStatus(ctx, refund.ID, enum2.NotifyStatus_Success)
	}
	return err
}

func (s *refundService) buildSuccessNotifyRequest(refundOrder *model.RefundOrder) *public.NotifyRefundRequest {
	return &public.NotifyRefundRequest{
		AppNo:         refundOrder.AppNo,
		MchNo:         refundOrder.MchNo,
		OrderNo:       refundOrder.OrderNo,
		RefundNo:      refundOrder.RefundNo,
		RefundTradeNo: refundOrder.RefundTradeNo,
		OutTradeNo:    refundOrder.OutMchTradeNo,
		Status:        refund2.Status_Success,
	}
}

func (s *refundService) buildNotifyRecord(req event.RefundOrderStatusEvent, refund *model.RefundOrder) *model.NotifyRecord {
	t := time.Now()
	return &model.NotifyRecord{
		NotifyType:     int64(enum2.NotifyType_Refund),
		AppNo:          req.AppNo,
		MchNo:          req.MchNo,
		Resaon:         "",
		OutNo:          req.RefundNo,
		TradeNo:        req.RefundTradeNo,
		NotifyURL:      refund.NotifyURL,
		ResResult:      "",
		NotifyCount:    1,
		NotifyStatus:   int64(enum2.NotifyStatus_Init),
		LastNotifyTime: t,
		CreatedAt:      t,
		UpdatedAt:      t,
	}
}

func (s *refundService) PublishCallback(ctx context.Context, req *http.Request, channel string, mchNo string, appNo string, refundTradeNo string) (response string, err error) {
	l := log.WithCtx(ctx)
	eventInfo := event.NewRefundCallbackEvent(req, mchNo, appNo, channel, refundTradeNo)
	payDri, err := paySdk.GetRefundDriver(channel2.Channel(channel2.Channel_value[channel]))
	if err != nil {
		return "", err
	}
	err = s.event.Publish(ctx, eventInfo, global.Cfg.Payment.Topic.RefundCallback)
	if err != nil {
		l.Error("支付回调处理失败", zap.Error(err), zap.String("refundTradeNo", refundTradeNo), zap.String("app_no", appNo), zap.String("mch_no", mchNo))
		return "", err
	}
	return payDri.CallbackResponse(), nil
}

func (s *refundService) Callback(ctx context.Context, req *http.Request, appNo string, refundTradeNo string) (string, error) {
	l := log.WithCtx(ctx)
	appInfo, err := s.appRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		l.Error("应用不存在", zap.Error(err), zap.Any("app_no", appNo))
		return "", errors.WrapError(errors.ErrAppNotExist, "应用不存在", err)
	}
	body, _ := utils.GetRequestBody(req)
	l.Info("退款回调处理开始", zap.String("refundTradeNo", refundTradeNo), zap.String("body", string(body)), zap.String("mch_no", appInfo.MchNo))
	refundDetail, _, err := s.domain.RefundService.Callback(ctx, req, refundTradeNo, appInfo, nil, nil)
	if err != nil {
		l.Error("Callback", zap.Error(err), zap.Any("req", req))
		return "", err
	}
	return refundDetail.Response, nil
}

func (s *refundService) Query(ctx context.Context, req *public.RefundQueryRequest) (*public.RefundResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	l := log.WithCtx(ctx)
	appInfo, err := getAppInfoFromCtx(ctx, req.AppNo, s.appRepo)
	if err != nil {
		l.Error("应用不存在", zap.Error(err), zap.Any("req", req))
		return nil, errors.WrapError(errors.ErrAppNotExist, "应用不存在", err)
	}
	req.MchNo = appInfo.MchNo
	refundDetail, err := s.domain.RefundService.Query(ctx, req, appInfo)
	if err != nil {
		l.Error("Refund", zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	return buildRefundResponse(refundDetail), nil
}

func (s *refundService) Refund(ctx context.Context, req *public.RefundRequest) (*public.RefundResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	l := log.WithCtx(ctx)
	appInfo, err := getAppInfoFromCtx(ctx, req.AppNo, s.appRepo)
	if err != nil {
		l.Error("应用不存在", zap.Error(err), zap.Any("req", req))
		return nil, errors.WrapError(errors.ErrAppNotExist, "应用不存在", err)
	}
	req.MchNo = appInfo.MchNo
	refundDetail, err := s.domain.RefundService.Refund(ctx, req, appInfo)
	if err != nil {
		l.Error("Refund", zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	return buildRefundResponse(refundDetail), nil
}

func buildRefundResponse(refundDetail *model.RefundOrder) *public.RefundResponse {
	r := &public.RefundResponse{
		OrderNo:  refundDetail.OrderNo,
		RefundNo: refundDetail.RefundNo,
		Amount: public.Amount{
			Total:    refundDetail.RefundAmount,
			Currency: refundDetail.Currency,
		},
		Status:           refund2.Status(refundDetail.Status),
		RefundTradeNo:    refundDetail.RefundTradeNo,
		RefundOutTradeNo: refundDetail.OutMchTradeNo,
		OrderAmount: public.Amount{
			Total: refundDetail.PaymentAmount,
		},
		PassBackParams: refundDetail.ExtParam,
	}
	return r
}
