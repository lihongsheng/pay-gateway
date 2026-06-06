package service

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	config2 "github.com/lihongsheng/pay-gateway/plugin/payment/config"
	"github.com/lihongsheng/pay-gateway/plugin/payment/contextkey"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	enum2 "github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	paySdk "github.com/lihongsheng/payment-sdk"
	"github.com/lihongsheng/payment-sdk/enum"
	channel2 "github.com/lihongsheng/payment-sdk/enum/channel"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// PaymentService 支付服务

type PaymentService interface {
	Payment(ctx context.Context, req *public.PaymentOrder, requestID string) (*public.PaymentOrderResponse, error)
	Query(ctx context.Context, req *public.QueryPaymentRequest) (*public.PaymentOrderDetail, error)
	Close(ctx context.Context, req *public.QueryPaymentRequest) error
	Callback(ctx context.Context, req *http.Request, appNo string, orderNo string, isTest bool) (response string, err error)
	PublishCallback(ctx context.Context, req *http.Request, channel string, mchNo string, appNo string, orderNo string, isTest bool) (response string, err error)
	// GenNotifyHandler 回调其他系统
	GenNotifyHandler(ctx context.Context, req event.PaymentOrderStatusEvent) error
	//	CronPushNotifyHandler(ctx context.Context, start time.Time, end time.Time) error
	NotifyHandler(ctx context.Context, req event.PaymentNotifyRetryEvent) error
	GenPaymentExpireRecord(ctx context.Context, req event.PaymentOrderStatusEvent) error
}
type paymentService struct {
	svc    *svc.ServiceContext
	domain *domain.ServiceGroup
}

func NewPaymentService(svc *svc.ServiceContext, domain *domain.ServiceGroup) PaymentService {
	return &paymentService{
		svc:    svc,
		domain: domain,
	}
}

func (s *paymentService) GenPaymentExpireRecord(ctx context.Context, req event.PaymentOrderStatusEvent) error {
	if !(req.NewStatus == payment.Status_Pending || req.NewStatus == payment.Status_TempFailed) {
		return nil
	}
	l := log.WithCtx(ctx)
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: req.OrderNo,
		TradeNo: "",
		AppNo:   req.AppNo,
		MchNo:   req.MchNo,
	})
	if err != nil {
		return err
	}
	if payOrder.IsTimeout() {
		err := s.svc.PaymentOrderRepo.UpdateStatus(ctx, payOrder, payOrder.Status, map[string]interface{}{
			"status":      payment.Status_TimeOut,
			"status_from": "consumerExpire",
		})
		if err != nil {
			l.Error("更新订单状态失败", zap.Error(err), zap.String("order_no", req.OrderNo))
		}
		return err
	}
	te := time.Unix(payOrder.TimeExpire, 0)
	err = s.svc.PaymentOrderRepo.SavePaymentExpireRecord(ctx, &model.PaymentExpireRecord{
		AppNo:      payOrder.AppNo,
		MchNo:      payOrder.MchNo,
		OrderNo:    payOrder.OrderNo,
		TradeNo:    payOrder.TradeNo,
		ExpireTime: te,
	})
	if err != nil {
		l.Error("保存订单过期记录失败", zap.Error(err), zap.String("order_no", req.OrderNo))
	}

	return nil
}

func (s *paymentService) NotifyHandler(ctx context.Context, req event.PaymentNotifyRetryEvent) error {
	l := log.WithCtx(ctx)
	notify, err := s.svc.NotifyRepo.GetById(ctx, req.NotifyID)
	if err != nil {
		return err
	}
	if notify.NotifyStatus == int64(enum2.NotifyStatus_Success) || notify.NotifyStatus == int64(enum2.NotifyStatus_Fail) {
		return nil
	}
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: req.OrderNo,
		TradeNo: "",
		AppNo:   req.AppNo,
		MchNo:   req.MchNo,
	})
	if err != nil {
		return err
	}

	return s.notify(ctx, l, payOrder, req.NotifyCount, req.NotifyID)
}

// CronPushNotifyHandler 定时处理回调
//func (s *paymentService) CronPushNotifyHandler(ctx context.Context, start time.Time, end time.Time) error {
//	records, err := s.svc.NotifyRepo.GetRetry(ctx, start, end, 100)
//	if err != nil {
//		return err
//	}
//	if len(records) == 0 {
//		return nil
//	}
//	l := log.WithCtx(ctx)
//	for _, record := range records {
//		// 扔到队列里，防止阻塞
//		err := s.svc.Event.Publish(ctx, &event.PaymentNotifyRetryEvent{
//			Base: event.Base{
//				EventType: event.EventTypePaymentNotifyRetry,
//				EventId:   fmt.Sprintf("Notify:%d", record.ID),
//			},
//			NotifyID:    record.ID,
//			AppNo:       record.AppNo,
//			MchNo:       record.MchNo,
//			OrderNo:     record.OutNo,
//			TradeNo:     record.TradeNo,
//			NotifyCount: int(record.NotifyCount),
//		}, config2.Config.Topic.PaymentNotifyRetry)
//		if err != nil {
//			l.Error("发送支付回调失败", zap.Error(err), zap.String("order_no", record.OutNo))
//			return err
//		}
//	}
//	return nil
//}

// GenNotifyHandler 回调其他系统
func (s *paymentService) GenNotifyHandler(ctx context.Context, req event.PaymentOrderStatusEvent) error {
	if !(req.NewStatus == payment.Status_Success) {
		return nil
	}
	key := "Notify:" + req.EventId
	lock, err := global.GVA_REDIS.SetNX(ctx, key, time.Now().Unix(), enum2.PaymentLockExpire).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, key)
	}()
	if !lock && err == nil {
		return errors.NewError(errors.ErrPayPending, "已在处理，请稍后再试")
	}
	l := log.WithCtx(ctx)
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: req.OrderNo,
		TradeNo: "",
		AppNo:   req.AppNo,
		MchNo:   req.MchNo,
	})
	if err != nil {
		return err
	}
	if payOrder.NotifyURL == "" {
		l.Warn("回调地址不存在", zap.String("order_no", req.OrderNo))
		return nil
	}
	notify, err := s.svc.NotifyRepo.Get(ctx, req.MchNo, req.AppNo, enum2.NotifyType_Payment, req.OrderNo)
	if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if notify == nil {
		notify = s.buildNotifyRecord(req, payOrder)
		err := s.svc.NotifyRepo.Save(ctx, notify)
		if err != nil {
			return err
		}
	}

	return s.notify(ctx, l, payOrder, 1, notify.ID)
}

func (s *paymentService) notify(ctx context.Context, l *zap.Logger, payOrder *entity.PaymentOrder, notifyCount int, notifyId int64) error {
	resp, err := utils.HttpJsonPost(ctx, payOrder.NotifyURL, s.buildSuccessNotifyRequest(payOrder), 30*time.Second)
	if err != nil {
		l.Error("支付回调处理失败", zap.Error(err), zap.String("OrderNo", payOrder.OrderNo), zap.String("app_no", payOrder.AppNo), zap.String("mch_no", payOrder.MchNo))
	}

	if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
		if notifyCount < enum2.NotifyRetryMax {
			_ = s.svc.NotifyRepo.Retry(ctx, notifyId, time.Now().Add(time.Duration(notifyCount+1)*enum2.NotifyRetryStep))
			return err
		} else {
			_ = s.svc.NotifyRepo.UpdateStatus(ctx, notifyId, enum2.NotifyStatus_Fail, "已达最大重试次数")
			_ = s.svc.PaymentOrderRepo.ConfirmNotifyStatus(ctx, payOrder.ID, enum2.NotifyStatus_Fail)
		}
	}
	if resp != nil && resp.StatusCode == http.StatusOK {
		err = s.svc.NotifyRepo.UpdateStatus(ctx, notifyId, enum2.NotifyStatus_Success, fmt.Sprintf("httpStatus:%d", http.StatusOK))
		_ = s.svc.PaymentOrderRepo.ConfirmNotifyStatus(ctx, payOrder.ID, enum2.NotifyStatus_Success)
	}
	return err
}

func (s *paymentService) buildSuccessNotifyRequest(payOrder *entity.PaymentOrder) *public.NotifyPaymentRequest {
	return &public.NotifyPaymentRequest{
		AppNo:      payOrder.AppNo,
		MchNo:      payOrder.MchNo,
		OrderNo:    payOrder.OrderNo,
		TradeNo:    payOrder.TradeNo,
		OutTradeNo: payOrder.OutMchTradeNo,
		Status:     payment.Status_Success,
	}
}

func (s *paymentService) buildNotifyRecord(req event.PaymentOrderStatusEvent, payOrder *entity.PaymentOrder) *model.NotifyRecord {
	t := time.Now()
	return &model.NotifyRecord{
		NotifyType:     int64(enum2.NotifyType_Payment),
		AppNo:          req.AppNo,
		MchNo:          req.MchNo,
		Resaon:         "",
		OutNo:          req.OrderNo,
		TradeNo:        req.TradeNo,
		NotifyURL:      payOrder.NotifyURL,
		ResResult:      "",
		NotifyCount:    1,
		NotifyStatus:   int64(enum2.NotifyStatus_Init),
		LastNotifyTime: t,
		CreatedAt:      t,
		UpdatedAt:      t,
	}
}

func (s *paymentService) PublishCallback(ctx context.Context, req *http.Request, channel string, mchNo string, appNo string, orderNo string, isTest bool) (response string, err error) {
	l := log.WithCtx(ctx)
	eventInfo := event.NewPaymentCallbackEvent(req, mchNo, appNo, channel, orderNo, isTest)
	payDri, err := paySdk.GetPaymentDriver(channel2.Channel(channel2.Channel_value[channel]))
	if err != nil {
		return "", err
	}
	err = s.svc.Event.Publish(ctx, eventInfo, config2.Config.Topic.PaymentCallback)
	if err != nil {
		l.Error("支付回调处理失败", zap.Error(err), zap.String("trade_no", orderNo), zap.String("app_no", appNo), zap.String("mch_no", mchNo))
		return "", err
	}
	return payDri.CallbackResponse(), nil
}

func (s *paymentService) Callback(ctx context.Context, req *http.Request, appNo string, orderNo string, isTest bool) (response string, err error) {
	l := log.WithCtx(ctx)
	appInfo, err := s.svc.AppRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		return "", errors.WrapError(errors.ErrAppNotExist, "应用不存在", err)
	}
	body, _ := utils.GetRequestBody(req)
	l.Info("回调处理开始", zap.String("orderNo", orderNo), zap.String("body", string(body)), zap.String("mch_no", appInfo.MchNo))
	if err != nil {
		l.Error("回调处理失败", zap.Error(err), zap.String("trade_no", orderNo), zap.String("app_no", appNo), zap.String("mch_no", appInfo.MchNo))
		return "", errors.WrapError(errors.ErrAppNotExist, "应用不存在", err)
	}
	detail, payOrder, account, err := s.domain.PaymentService.Callback(ctx, req, appInfo, orderNo)
	if err != nil {
		l.Error("回调处理失败", zap.Error(err), zap.String("orderNo", orderNo), zap.String("app_no", appNo), zap.String("mch_no", appInfo.MchNo))
		return "", err
	}
	l.Info("回调处理成功", zap.Any("detail", detail), zap.String("orderNo", orderNo))
	if isTest && detail != nil && detail.Status == payment.Status_Success {
		// 测试环境回调处理
		gErr := s.svc.PaymentAccountRepo.SetPayValidateSuccess(ctx, account.AccountNo)
		if gErr != nil {
			l.Error("回调处理失败PaymentAccountRepo", zap.Error(err), zap.String("orderNo", orderNo), zap.String("app_no", appNo), zap.String("mch_no", appInfo.MchNo))
		}
	}
	response = detail.Response
	// 支付宝使用的支付时候的回调地址
	if detail.EventAction == enum.Event_REFUND && detail.EventRefund != nil {
		_, _, err := s.domain.RefundService.Callback(ctx, req, detail.EventRefund.RefundNo, appInfo, payOrder, account)
		if err != nil {
			l.Error("回调处理失败", zap.Error(err), zap.String("orderNo", orderNo), zap.String("app_no", appNo), zap.String("mch_no", appInfo.MchNo))
			return "", err
		}
	}
	return response, nil
}

func (s *paymentService) Close(ctx context.Context, req *public.QueryPaymentRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	l := log.WithCtx(ctx)
	appInfo, err := getAppInfoFromCtx(ctx, req.AppNo, s.svc)
	if err != nil {
		return err
	}
	req.MchNo = appInfo.MchNo
	if err != nil {
		return errors.WrapError(errors.ErrAppNotExist, "应用不存在", err)
	}
	err = s.domain.PaymentService.Close(ctx, *req, appInfo)
	if err != nil {
		l.Error("Close", zap.Error(err), zap.Any("req", req))
		return err
	}
	return nil
}

func (s *paymentService) Query(ctx context.Context, req *public.QueryPaymentRequest) (*public.PaymentOrderDetail, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	l := log.WithCtx(ctx)
	appInfo, err := getAppInfoFromCtx(ctx, req.AppNo, s.svc)
	if err != nil {
		return nil, err
	}
	req.MchNo = appInfo.MchNo
	order, err := s.domain.PaymentService.Query(ctx, req, appInfo)
	if err != nil {
		l.Error("Query", zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	return buildPaymentOrderDetail(order), nil
}

func buildPaymentOrderDetail(req *entity.PaymentOrder) *public.PaymentOrderDetail {
	r := &public.PaymentOrderDetail{
		Order: public.Order{
			OrderNo: req.OrderNo,
			Amount: public.Amount{
				Total:    req.PaymentAmount,
				Currency: req.Currency,
			},
			Goods:    make([]public.Goods, 0, len(req.PaymentOrderProducts)),
			Subject:  req.OrderSubject,
			Desc:     req.OrderDesc,
			CreateAt: time.Unix(req.OrderCreate, 0),
		},
		RedirectUrl:    req.RedirectURL,
		NotifyUrl:      req.NotifyURL,
		PassBackParams: req.ExtParam,
		Status:         req.Status,
		TradeNo:        req.TradeNo,
		OutTradeNo:     req.OutMchTradeNo,
		Create:         req.CreatedAt,
	}
	for _, product := range req.PaymentOrderProducts {
		r.Order.Goods = append(r.Order.Goods, public.Goods{
			Name:     product.ProductName,
			Sku:      product.Sku,
			Price:    product.Price,
			Quantity: int(product.Quantity),
			Desc:     product.ProductDesc,
		})
	}
	return r
}

func (s *paymentService) Payment(ctx context.Context, req *public.PaymentOrder, requestID string) (*public.PaymentOrderResponse, error) {
	l := log.WithCtx(ctx)
	if err := req.Validate(); err != nil {
		return nil, err
	}
	appInfo, err := getAppInfoFromCtx(ctx, req.AppNo, s.svc)
	if err != nil {
		return nil, errors.WrapError(errors.ErrAppNotExist, "应用不存在", err)
	}
	var account *entity.PaymentAccount
	req.MchNo = appInfo.MchNo
	if req.RequestID == "" {
		req.RequestID = requestID
	}
	if req.AccountNo != "" {
		account, err = s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, req.AccountNo)
		if err != nil {
			return nil, errors.WrapError(errors.ErrPayChannelNotSupport, "未找到支持的支付渠道", err)
		}
	} else {
		// 微信或者支付宝 open 和支付渠道 accountNo 对应
		if appInfo.MultiChannel == enum2.MultiChannel_Enable {
			return nil, errors.NewError(errors.ErrPayChannelNotSupport, "多渠道模式下请指定支付渠道")
		}
		accounts, err := s.domain.Router.GetAvailablePayment(ctx, &public.QueryCheckoutPaymentMethod{
			PaymentMethod:  req.PaymentMethod,
			PaymentProduct: req.PaymentProduct,
			AppNo:          req.AppNo,
		}, appInfo)
		if err != nil {
			return nil, errors.WrapError(errors.ErrPayChannelNotSupport, "未找到支持的支付渠道", err)
		}
		account = accounts[0]
	}
	callbackUrl, err := getSelfCallbackUrl(appInfo, account, req.Order.OrderNo, enum2.PaymentNotify)
	if err != nil {
		return nil, err
	}
	result, _, err := s.domain.PaymentService.Payment(ctx, nil, req, appInfo, account, callbackUrl)
	if err != nil {
		l.Error("GetUserOpenID", zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	return result, nil
}

// getAppInfoFromCtx
func getAppInfoFromCtx(ctx context.Context, appNo string, svcCtx *svc.ServiceContext) (*model.Application, error) {
	// 1. 优先从标准 Context 提取（已注入则直接复用）
	appInfo := contextkey.GetAppInfo(ctx)
	if appInfo != nil {
		return appInfo, nil
	}
	if appNo != appInfo.AppNo {
		return nil, errors.NewError(errors.ErrAppNotExist, "应用不存在")
	}
	return svcCtx.AppRepo.GetByAppNoFormCache(ctx, appNo)
}
