package domain

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	config2 "github.com/lihongsheng/pay-gateway/plugin/payment/config"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	paySdk "github.com/lihongsheng/payment-sdk"
	"github.com/lihongsheng/payment-sdk/config"
	"github.com/lihongsheng/payment-sdk/config/proxy"
	"github.com/lihongsheng/payment-sdk/driver/dto"
	"github.com/lihongsheng/payment-sdk/driver/iface"
	"github.com/lihongsheng/payment-sdk/enum/channel"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"github.com/lihongsheng/payment-sdk/enum/refund"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type RefundService interface {
	Refund(ctx context.Context, req *public.RefundRequest, appInfo *model.Application) (*model.RefundOrder, error)
	Query(ctx context.Context, req *public.RefundQueryRequest, appInfo *model.Application) (*model.RefundOrder, error)
	Callback(ctx context.Context, req *http.Request, refundTradeNo string, appInfo *model.Application, payOrder *entity.PaymentOrder, account *entity.PaymentAccount) (*dto.CallbackRefundDetail, *model.RefundOrder, error)
}

type refundService struct {
	svc *svc.ServiceContext
}

func NewRefundService(svc *svc.ServiceContext) RefundService {
	return &refundService{
		svc: svc,
	}
}

func (s *refundService) Callback(ctx context.Context, req *http.Request, refundTradeNo string, appInfo *model.Application, payOrder *entity.PaymentOrder, account *entity.PaymentAccount) (*dto.CallbackRefundDetail, *model.RefundOrder, error) {
	m, err := s.svc.RefundRepo.GetFromCache(ctx, public.RefundQueryRequest{
		RefundTradeNo: refundTradeNo,
		RefundNo:      "",
		AppNo:         appInfo.AppNo,
		MchNo:         appInfo.MchNo,
	})
	if err != nil {
		if !errors2.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
		}
		return nil, nil, errors.NewError(errors.ErrPayOrderNotExist, "订单不存在")
	}
	if payOrder == nil {
		payOrder, err = s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
			OrderNo: m.OrderNo,
			TradeNo: "",
			AppNo:   m.AppNo,
			MchNo:   m.MchNo,
		})
		if err != nil {
			if !errors2.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
			}
			return nil, nil, errors.NewError(errors.ErrPayOrderNotExist, "订单不存在")
		}
	}
	if account == nil {
		account, err = s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, payOrder.PaymentAccountNo)
		if err != nil {
			return nil, nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
		}
	}
	payDri, err := s.buildPayDri(payOrder.PaymentMethod, payOrder.PaymentProduct, appInfo, account)
	if err != nil {
		return nil, nil, err
	}
	resp, err := payDri.Callback(ctx, req)
	if err != nil {
		return nil, nil, errors.WrapError(errors.ErrPayChannelNotSupport, "回调处理失败", err)
	}
	l := log.WithCtx(ctx)
	if int64(resp.Status) != m.Status {
		eventInfo := event.GetRefundOrderStatusEvent(m, resp.Status)
		m.Status = int64(refund.Status(resp.Status))
		if resp.TradeRefundNo != "" {
			m.OutMchTradeNo = resp.TradeRefundNo
		}
		if !resp.SuccessTime.IsZero() {
			m.SuccessTime = resp.SuccessTime
		} else if resp.Status == refund.Status_Success {
			m.SuccessTime = time.Now()
		}
		m.StatusFrom = "callback"
		m.RefundAmount = resp.Amount.Total
		if err := s.svc.RefundRepo.Save(ctx, m); err != nil {
			return nil, nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
		}
		if eventInfo != nil {
			l.Info("发布退款状态", zap.Error(err), zap.Any("event", eventInfo))
			if err := s.svc.Event.Publish(ctx, eventInfo, config2.Config.Topic.RefundStatus); err != nil {
				l.Error("发布退款状态失败", zap.Error(err), zap.Any("event", eventInfo))
			}
		}
	}
	if m.Status == int64(refund.Status(resp.Status)) && payOrder.Status == payment.Status_Success {
		updatePaymentErr := s.svc.PaymentOrderRepo.UpdateStatus(ctx, payOrder, payOrder.Status, map[string]interface{}{
			"status":      payment.Status_Refund,
			"status_from": "refundCallback",
		})
		if updatePaymentErr != nil {
			l.Error("更新支付订单状态失败", zap.Error(updatePaymentErr))
		}
	}
	return resp, m, nil
}

func (s *refundService) Refund(ctx context.Context, req *public.RefundRequest, appInfo *model.Application) (m *model.RefundOrder, err error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	l := log.WithCtx(ctx)
	lock, err := global.GVA_REDIS.SetNX(ctx, s.getRefundLockKey(req.OrderNo, req.AppNo), time.Now().Unix(), enum.PaymentLockExpire).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, s.getRefundLockKey(req.OrderNo, req.AppNo))
	}()
	if !lock && err == nil {
		return nil, errors.NewError(errors.ErrRefundPending, "有退款处理中, 请稍后再试")
	}
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: req.OrderNo,
		TradeNo: "",
		AppNo:   req.AppNo,
		MchNo:   req.MchNo,
	})
	if err != nil {
		if !errors2.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
		}
		return nil, errors.NewError(errors.ErrPayOrderNotExist, "订单不存在")
	}
	if payOrder.IsPending() {
		return nil, errors.NewError(errors.ErrPayOrderNotPending, "订单未完成支付")
	}
	if req.Amount.Total > payOrder.PaymentAmount {
		return nil, errors.NewError(errors.ErrRefundAmountErr, "退款金额大于支付金额")
	}
	account, err := s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, payOrder.PaymentAccountNo)
	if err != nil {
		return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
	}
	m, err = s.svc.RefundRepo.GetFromCache(ctx, public.RefundQueryRequest{
		RefundTradeNo: "",
		RefundNo:      req.RefundNo,
		AppNo:         req.AppNo,
		MchNo:         req.MchNo,
	})
	if err != nil {
		if !errors2.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
		}
	}
	if m == nil {
		m = s.buildRefundModel(req, payOrder)
	} else {
		if m.RefundAmount != req.Amount.Total {
			return nil, errors.NewError(errors.ErrRefundAmountErr, "退款金额不一致")
		}
		m.Resaon = req.Reason
		m.UpdatedAt = time.Now()
	}
	if m.Status == int64(refund.Status_Success) {
		return m, nil
	}
	refundAmount, err := s.svc.RefundRepo.CountAmount(ctx, req.MchNo, req.AppNo, req.OrderNo, []refund.Status{refund.Status_Success})
	if err != nil {
		return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
	}
	if refundAmount+req.Amount.Total > payOrder.PaymentAmount {
		return nil, errors.NewError(errors.ErrRefundAmountErr, "退款总金额大于支付金额")
	}
	// 6. 异步保存订单（优化：独立ctx防止主ctx超时
	defer func() {
		c, cancel := context.WithTimeout(context.Background(), enum.SaveOrderTimeout)
		if err := s.svc.RefundRepo.Save(c, m); err != nil {
			l.Error("创建退款订单失败", zap.Error(err), zap.Any("refundOrder", m))
		}
		cancel()
	}()
	refundResult, err := s.invokeChannelPayment(ctx, m, payOrder, appInfo, account)
	oldStatus := m.Status
	if err != nil {
		m.Status = int64(refund.Status_Failed)
		m.ThirdMsg = err.Error()
		l.Error("退款失败", zap.Error(err), zap.Any("refundOrder", m))
	}
	if refundResult != nil {
		m.Status = int64(refundResult.Status)
		//m.RefundAmount = refundResult.Amount.Total
		if refundResult.TradeRefundNo != "" {
			m.OutMchTradeNo = refundResult.TradeRefundNo
		}
	}
	if oldStatus != m.Status {
		eventInfo := event.GetRefundOrderStatusEvent(m, refund.Status(m.Status))
		if eventInfo != nil {
			if err := s.svc.Event.Publish(ctx, eventInfo, config2.Config.Topic.RefundStatus); err != nil {
				log.WithCtx(ctx).Error("发布订单状态失败", zap.Error(err), zap.Any("event", eventInfo))
			}
		}
	}
	return m, nil
}

func (s *refundService) invokeChannelPayment(ctx context.Context, refund *model.RefundOrder, payOrder *entity.PaymentOrder, appInfo *model.Application, account *entity.PaymentAccount) (*dto.RefundDetail, error) {
	payDri, err := s.buildPayDri(payOrder.PaymentMethod, payOrder.PaymentProduct, appInfo, account)
	if err != nil {
		return nil, err
	}
	payReq, err := s.buildChannelPayOrder(refund, payOrder, appInfo, account)
	if err != nil {
		return nil, err
	}
	return payDri.Refund(ctx, payReq)
}

func (s *refundService) buildChannelPayOrder(refund *model.RefundOrder, payOrder *entity.PaymentOrder, appInfo *model.Application, account *entity.PaymentAccount) (*dto.RefundRequest, error) {
	path := fmt.Sprintf(enum.RefundNotify, account.Channel, refund.MchNo, refund.AppNo, refund.RefundTradeNo)
	parse, err := getNotifyUrl(appInfo)
	if err != nil {
		return nil, err
	}
	parse = parse.JoinPath(path)
	request := &dto.RefundRequest{
		RefundNo:  refund.RefundTradeNo,
		OrderNo:   payOrder.TradeNo,
		Reason:    refund.Resaon,
		NotifyUrl: "",
		Amount: dto.Amount{
			Total:    refund.RefundAmount,
			Currency: payOrder.Currency,
		},
		Goods: nil,
		OrderAmount: dto.Amount{
			Total:    payOrder.PaymentAmount,
			Currency: payOrder.Currency,
		},
	}
	request.NotifyUrl = parse.String()

	return request, nil
}

func (s *refundService) buildPayDri(paymentMethod payment.Payment, paymentProduct payment.PaymentProduct, appInfo *model.Application, account *entity.PaymentAccount) (iface.Refund, error) {
	opts := []config.Option{
		config.WithPayment(paymentMethod),
		config.WithPaymentProduct(paymentProduct),
		config.WithConfig(account.ChannelConfig),
	}
	if appInfo.ProxyHost != "" && appInfo.ProxyPort > 0 {
		opts = append(opts, config.WithProxy(&proxy.Proxy{
			Host:     appInfo.ProxyHost,
			Port:     int(appInfo.ProxyPort),
			UserName: appInfo.ProxyUser,
			Password: appInfo.ProxyPwd,
		}))
	}
	payDri, err := paySdk.Refund(channel.Channel(channel.Channel_value[account.Channel]), opts...)
	if err != nil {
		return nil, errors.WrapError(errors.ErrPayChannelNotSupport, fmt.Sprintf("构建[%s]渠道驱动失败:", account.Channel), err)
	}
	return payDri, nil
}

func (s *refundService) buildRefundModel(req *public.RefundRequest, payOrder *entity.PaymentOrder) *model.RefundOrder {
	no, t := utils.GenTradeNo(req.MchNo)
	r := &model.RefundOrder{
		ID:               0,
		PaymentAccountNo: payOrder.PaymentAccountNo,
		AppNo:            req.AppNo,
		MchNo:            req.MchNo,
		Resaon:           req.Reason,
		OrderNo:          payOrder.OrderNo,
		TradeNo:          payOrder.TradeNo,
		RefundTradeNo:    no,
		RefundNo:         req.RefundNo,
		Device:           "",
		RefundAmount:     req.Amount.Total,
		PaymentAmount:    payOrder.PaymentAmount,
		Currency:         req.Amount.Currency,
		Status:           int64(refund.Status_Created),
		StatusFrom:       "api",
		NotifyURL:        req.NotifyUrl,
		NotifyStatus:     0,
		OutMchTradeNo:    "",
		ExtParam:         req.PassBackParams,
		CreatedAt:        utils.GetGenIDTimestamp(t),
		RefundFrom:       int64(req.RefundFrom),
	}
	return r
}
func (s *refundService) Query(ctx context.Context, req *public.RefundQueryRequest, appInfo *model.Application) (*model.RefundOrder, error) {
	m, err := s.svc.RefundRepo.GetFromCache(ctx, public.RefundQueryRequest{
		RefundTradeNo: req.RefundTradeNo,
		RefundNo:      req.RefundNo,
		AppNo:         appInfo.AppNo,
		MchNo:         appInfo.MchNo,
	})
	if err != nil {
		if !errors2.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
		}
		return nil, errors.NewError(errors.ErrPayOrderNotExist, "订单不存在")
	}
	if !(m.Status == int64(refund.Status_Created) || m.Status == int64(refund.Status_Pending)) {
		return m, nil
	}
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: m.OrderNo,
		TradeNo: m.TradeNo,
		AppNo:   m.AppNo,
		MchNo:   m.MchNo,
	})
	if err != nil {
		if !errors2.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
		}
		return nil, errors.NewError(errors.ErrPayOrderNotExist, "订单不存在")
	}
	account, err := s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, payOrder.PaymentAccountNo)
	if err != nil {
		return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
	}
	payDri, err := s.buildPayDri(payOrder.PaymentMethod, payOrder.PaymentProduct, appInfo, account)
	if err != nil {
		return nil, err
	}
	resp, err := payDri.Query(ctx, dto.RefundQuery{
		RefundNo: m.RefundNo,
		OrderNo:  "",
		TradeNo:  m.RefundTradeNo,
	})
	if err != nil {
		return nil, errors.WrapError(errors.ErrCodeSysError, "查询失败，请稍后再试", err)
	}
	if int64(resp.Status) != int64(m.Status) {
		eventInfo := event.GetRefundOrderStatusEvent(m, resp.Status)
		m.Status = int64(resp.Status)
		m.OutMchTradeNo = resp.TradeRefundNo
		err := s.svc.RefundRepo.Save(ctx, m)
		if err != nil {
			return nil, errors.WrapError(errors.ErrCodeSysError, "更新订单状态失败", err)
		}
		if eventInfo != nil {
			if err := s.svc.Event.Publish(ctx, eventInfo, config2.Config.Topic.RefundStatus); err != nil {
				log.WithCtx(ctx).Error("发布订单状态失败", zap.Error(err), zap.Any("event", eventInfo))
			}
		}
	}
	return m, nil
}

func (s *refundService) getRefundLockKey(orderNo, appNo string) string {
	return fmt.Sprintf(enum.CacheLockRefund, orderNo, appNo)
}
