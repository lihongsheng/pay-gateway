package domain

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/lihongsheng/pay-gateway/config"
	config3 "github.com/lihongsheng/pay-gateway/plugin/payment/config"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	paySdk "github.com/lihongsheng/payment-sdk"
	config2 "github.com/lihongsheng/payment-sdk/config"
	"github.com/lihongsheng/payment-sdk/config/proxy"
	"github.com/lihongsheng/payment-sdk/driver/dto"
	"github.com/lihongsheng/payment-sdk/driver/iface"
	"github.com/lihongsheng/payment-sdk/enum/channel"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	errors3 "github.com/lihongsheng/payment-sdk/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"time"
)

type PaymentService interface {
	Payment(ctx context.Context, payOrder *entity.PaymentOrder, req *public.PaymentOrder, appInfo *model.Application, account *entity.PaymentAccount, selfCallback string) (*public.PaymentOrderResponse, *entity.PaymentOrder, error)
	Query(ctx context.Context, req *public.QueryPaymentRequest, appInfo *model.Application) (*entity.PaymentOrder, error)
	Close(ctx context.Context, req public.QueryPaymentRequest, appInfo *model.Application) error
	Callback(ctx context.Context, req *http.Request, appInfo *model.Application, orderNO string) (detail *dto.CallbackPayDetail, payOrder *entity.PaymentOrder, account *entity.PaymentAccount, err error)
}
type paymentService struct {
	svc   *svc.ServiceContext
	redis *redis.Client
	cfg   config.Config
}

func NewPaymentService(svc *svc.ServiceContext, redis *redis.Client, cfg config.Config) PaymentService {
	return &paymentService{
		svc:   svc,
		redis: redis,
		cfg:   cfg,
	}
}
func (s *paymentService) Callback(ctx context.Context, req *http.Request, appInfo *model.Application, orderNO string) (detail *dto.CallbackPayDetail, payOrder *entity.PaymentOrder, account *entity.PaymentAccount, err error) {
	payOrder, err = s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: orderNO,
		TradeNo: "",
		AppNo:   appInfo.AppNo,
		MchNo:   appInfo.MchNo,
	})
	if err != nil {
		return nil, nil, nil, errors.WrapCode(errors.ErrPayOrderNotExist, err)
	}
	account, err = s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, payOrder.PaymentAccountNo)
	if err != nil {
		return nil, nil, nil, errors.WrapCode(errors.ErrPayChannelNotSupport, err)
	}
	payDri, err := s.buildPayDri(payOrder.PaymentMethod, payOrder.PaymentProduct, appInfo, account)
	detail, err = payDri.Callback(ctx, req)
	if err != nil {
		return nil, nil, nil, errors.WrapError(errors.ErrPayChannelNotSupport, "回调处理失败", err)
	}
	if payOrder.IsPending() {
		updates := make(map[string]interface{})
		updates["out_mch_trade_no"] = detail.TradeNo
		updates["status"] = detail.Status
		updates["status_from"] = "callback"
		err = s.svc.PaymentOrderRepo.UpdateStatus(ctx, payOrder, payOrder.Status, updates)
		if err != nil {
			return nil, nil, nil, errors.WrapError(errors.ErrCodeSysError, "更新订单状态失败", err)
		}
		payOrder.SetStatus(detail.Status)
		payOrder.OutMchTradeNo = detail.TradeNo
		eventInfo := payOrder.GetEvents()
		if eventInfo != nil {
			err := s.svc.Event.Publish(ctx, eventInfo, s.svc.Config.Topic.PaymentStatus)
			if err != nil {
				log.WithCtx(ctx).Error("发布订单状态失败", zap.Error(err), zap.Any("event", eventInfo))
			}
		}
	}
	return detail, payOrder, account, nil
}
func (s *paymentService) Close(ctx context.Context, req public.QueryPaymentRequest, appInfo *model.Application) error {
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, req)
	if err != nil {
		if !errors2.Is(err, gorm.ErrRecordNotFound) {
			return errors.WrapError(errors.ErrCodeSysError, "系统异常，请稍后再试", err)
		}
		return errors.NewError(errors.ErrPayOrderNotExist, "订单不存在")
	}
	if !payOrder.IsPending() {
		return errors.NewError(errors.ErrPayOrderNotPending, "支付中订单才可关闭")
	}
	account, err := s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, payOrder.PaymentAccountNo)
	if err != nil {
		return errors.WrapError(errors.ErrPayChannelNotSupport, "获取支付账户失败", err)
	}
	payDri, err := s.buildPayDri(payOrder.PaymentMethod, payOrder.PaymentProduct, appInfo, account)
	if err != nil {
		return errors.WrapCode(errors.ErrPayChannelNotSupport, err)
	}
	err = payDri.Close(ctx, dto.CloseQuery{
		OrderNo: payOrder.TradeNo,
		TradeNo: "",
	})
	if err != nil {
		if !errors3.IsNoSupport(err) {
			return errors.WrapError(errors.ErrPayChannelNotSupport, "关闭订单失败", err)
		}
	}
	err = s.svc.PaymentOrderRepo.UpdateStatus(ctx, payOrder, payOrder.Status, map[string]interface{}{
		"status":      payment.Status_Close,
		"status_from": "apiClose",
	})
	if err != nil {
		zap.L().Error("更新订单为关闭状态失败",
			zap.String("tradeNo", payOrder.TradeNo),
			zap.String("mchNo", payOrder.MchNo),
			zap.Error(err))
	}
	return nil
}
func (s *paymentService) Query(ctx context.Context, req *public.QueryPaymentRequest, appInfo *model.Application) (*entity.PaymentOrder, error) {
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, *req)
	if err != nil {
		return nil, errors.WrapCode(errors.ErrPayOrderNotExist, err)
	}
	if !payOrder.IsPending() {
		return payOrder, nil
	}
	account, err := s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, payOrder.PaymentAccountNo)
	if err != nil {
		return nil, errors.WrapError(errors.ErrPayChannelNotSupport, "获取支付账户失败", err)
	}
	payDri, err := s.buildPayDri(payOrder.PaymentMethod, payOrder.PaymentProduct, appInfo, account)
	if err != nil {
		return nil, errors.WrapCode(errors.ErrPayChannelNotSupport, err)
	}
	resp, err := payDri.Query(ctx, dto.Query{
		OrderNo: payOrder.TradeNo,
		TradeNo: "",
	})
	if err != nil {
		return nil, errors.WrapCode(errors.ErrPayChannelNotSupport, err)
	}
	updates := make(map[string]interface{})
	if resp.Status != payOrder.Status {
		updates["status"] = resp.Status
	}
	oldStatus := payOrder.Status
	payOrder.SetStatus(resp.Status)
	if payOrder.OutMchTradeNo == "" {
		payOrder.OutMchTradeNo = resp.TradeNo
		updates["out_mch_trade_no"] = payOrder.OutMchTradeNo
	}
	if len(updates) > 0 {
		l := log.WithCtx(ctx)
		err := s.svc.PaymentOrderRepo.UpdateStatus(ctx, payOrder, oldStatus, updates)
		if err != nil {
			l.Error("更新支付订单失败", zap.Error(err), zap.Any("updates", updates), zap.String("RefundTradeNo", payOrder.TradeNo), zap.String("mchNo", payOrder.MchNo))
		}
	}
	eventInfo := payOrder.GetEvents()
	if eventInfo != nil {
		err := s.svc.Event.Publish(ctx, eventInfo, s.svc.Config.Topic.PaymentStatus)
		if err != nil {
			log.WithCtx(ctx).Error("发布订单状态失败", zap.Error(err), zap.Any("event", eventInfo))
		}
	}
	return payOrder, nil
}
func (s *paymentService) Payment(ctx context.Context, payOrder *entity.PaymentOrder, req *public.PaymentOrder, appInfo *model.Application, account *entity.PaymentAccount, selfCallback string) (*public.PaymentOrderResponse, *entity.PaymentOrder, error) {
	lock, err := s.redis.SetNX(ctx, s.getPaymentLockKey(req.Order.OrderNo, req.AppNo), time.Now().Unix(), enum.PaymentLockExpire).Result()
	defer func() {
		s.redis.Del(ctx, s.getPaymentLockKey(req.Order.OrderNo, req.AppNo))
	}()
	if !lock && err == nil {
		return nil, nil, errors.NewError(errors.ErrPayPending, "重复支付，请稍后再试")
	}
	if payOrder == nil {
		// 已有支付成功订单
		payOrder, err = s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
			OrderNo: req.Order.OrderNo,
			TradeNo: "",
			AppNo:   req.AppNo,
			MchNo:   req.MchNo,
		})
		if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.WrapError(errors.ErrCodeSysError, "系统异常扫后再试", err)
		}
	}
	var m *entity.PaymentOrder
	if payOrder != nil {
		if !payOrder.IsRetry() {
			return &public.PaymentOrderResponse{
				OrderNo: payOrder.OrderNo,
				TradeNo: payOrder.TradeNo,
				Amount: public.Amount{
					Total:    payOrder.PaymentAmount,
					Currency: "CNY",
				},
				Status:      payment.Status(payOrder.Status),
				Action:      nil,
				RedirectURL: req.RedirectUrl,
			}, payOrder, nil
		}
		if payOrder.PaymentAmount != req.Order.Amount.Total {
			return nil, nil, errors.NewError(errors.ErrPayAmountErr, fmt.Sprintf("订单金额不匹配，本地金额：%d，请求金额：%d", payOrder.PaymentAmount, req.Order.Amount.Total))
		}
		m = payOrder
		m.PaymentProduct = req.PaymentProduct
		m.PaymentMethod = req.PaymentMethod
		m.PaymentAccountNo = account.AccountNo
		m.UserOpenid = req.Payer.OpenID
	} else {
		m = buildModel(req, appInfo, account)
	}
	tradeNo := m.TradeNo
	if err := m.Validate(); err != nil {
		return nil, nil, err
	}
	l := log.WithCtx(ctx)
	// 6. 异步保存订单（优化：独立ctx防止主ctx超时
	defer func() {
		c, cancel := context.WithTimeout(context.Background(), enum.SaveOrderTimeout)
		if err := s.svc.PaymentOrderRepo.Save(c, m); err != nil {
			l.Error("创建支付订单失败", zap.Error(err), zap.Any("payOrder", m))
		}
		cancel()
	}()
	payRes, err := s.invokeChannelPayment(ctx, req, appInfo, account, tradeNo, selfCallback)
	l.Info("支付结果", zap.Any("payRes", payRes), zap.Any("err", err))
	if err != nil {
		m.ThirdMsg = err.Error()
		m.Retry++
		if m.Retry <= enum.MaxRetryPaymentLimit {
			if errors3.IsPaymentLimited(err) {
				m.SetErrStatus(payment.Status_TempFailed, errors.ErrPaymentLimit)
			} else {
				m.SetErrStatus(payment.Status_TempFailed, errors.ErrPayFail)
			}
		} else {
			m.SetErrStatus(payment.Status_Failed, errors.ErrPayFail)
		}
	} else {
		m.SetStatus(payRes.Status)
	}
	eventInfo := m.GetEvents()
	if eventInfo != nil {
		eventErr := s.svc.Event.Publish(ctx, eventInfo, s.svc.Config.Topic.PaymentStatus)
		if eventErr != nil {
			l.Error("发布订单状态失败", zap.Error(eventErr), zap.Any("event", eventInfo))
		}
	}
	if err != nil {
		if m.Retry > enum.MaxRetryPaymentLimit {
			return nil, m, errors.WrapError(errors.ErrPayFail, "支付失败", err)
		}
		return nil, m, errors.WrapError(errors.ErrRetryFail, "支付失败，正在重试", err)
	}
	if payRes != nil && payRes.TradeNo != "" {
		m.OutMchTradeNo = payRes.TradeNo
	}
	return &public.PaymentOrderResponse{
		OrderNo:    m.OrderNo,
		TradeNo:    m.TradeNo,
		OutTradeNo: payRes.TradeNo,
		Amount: public.Amount{
			Total:    m.PaymentAmount,
			Currency: req.Order.Amount.Currency,
		},
		Status:      m.Status,
		Action:      &payRes.Action,
		RedirectURL: m.RedirectURL,
	}, m, nil
}
func (s *paymentService) invokeChannelPayment(ctx context.Context, req *public.PaymentOrder, appInfo *model.Application, account *entity.PaymentAccount, tradeNo string, selfCallback string) (*dto.PayResponse, error) {
	// 测试后删除
	if s.cfg.Env.IsTest() && account.AccountNo == "P6b5e44d8e7c00" {
		return nil, errors3.ErrorPaymentLimited("此商家的收款功能已被限制", nil)
	}
	payDri, err := s.buildPayDri(req.PaymentMethod, req.PaymentProduct, appInfo, account)
	if err != nil {
		return nil, err
	}
	payReq, err := s.buildChannelPayOrder(req, appInfo, account, tradeNo, selfCallback)
	if err != nil {
		return nil, err
	}
	log.WithCtx(ctx).Info("支付请求", zap.Any("payReq", payReq))
	return payDri.Pay(ctx, payReq)
}
func (s *paymentService) buildPayDri(paymentMethod payment.Payment, paymentProduct payment.PaymentProduct, appInfo *model.Application, account *entity.PaymentAccount) (iface.Pay, error) {
	opts := []config2.Option{
		config2.WithPayment(paymentMethod),
		config2.WithPaymentProduct(paymentProduct),
		config2.WithConfig(account.ChannelConfig),
	}
	if appInfo.ProxyHost != "" && appInfo.ProxyPort > 0 {
		opts = append(opts, config2.WithProxy(&proxy.Proxy{
			Host:     appInfo.ProxyHost,
			Port:     int(appInfo.ProxyPort),
			UserName: appInfo.ProxyUser,
			Password: appInfo.ProxyPwd,
		}))
	}
	payDri, err := paySdk.Payment(channel.Channel(channel.Channel_value[account.Channel]), opts...)
	if err != nil {
		return nil, errors.WrapError(errors.ErrPayChannelNotSupport, fmt.Sprintf("构建[%s]渠道驱动失败:", account.Channel), err)
	}
	return payDri, nil
}
func buildModel(req *public.PaymentOrder, appInfo *model.Application, account *entity.PaymentAccount) *entity.PaymentOrder {
	result := &entity.PaymentOrder{
		PaymentAccountNo:     account.AccountNo,
		AppNo:                req.AppNo,
		MchNo:                req.MchNo,
		OrderSubject:         req.Order.Subject,
		PaymentMethod:        req.PaymentMethod,
		PaymentProduct:       req.PaymentProduct,
		OrderNo:              req.Order.OrderNo,
		OutMchTradeNo:        "",
		Device:               "",
		System:               "",
		OrderAmount:          req.Order.Amount.Total,
		PaymentAmount:        req.Order.Amount.Total,
		UserPaymentAmount:    0,
		Currency:             req.Order.Amount.Currency,
		Status:               payment.Status_Created,
		StatusFrom:           "api",
		TimeExpire:           req.TimeExpire,
		OrderCreate:          req.Order.CreateAt.Unix(),
		ExtParam:             req.PassBackParams,
		RedirectURL:          req.RedirectUrl,
		NotifyURL:            req.NotifyUrl,
		OrderDesc:            req.Order.Desc,
		ThirdCode:            "",
		ThirdMsg:             "",
		CreatedAt:            time.Now(),
		PaymentOrderProducts: make([]entity.Product, 0, len(req.Order.Goods)),
		UserOpenid:           req.Payer.OpenID,
		Retry:                enum.MaxRetryPaymentDefault,
	}
	if result.TimeExpire == 0 {
		result.TimeExpire = time.Now().Add(enum.DefaultExpireTime).Unix()
	}
	result.SetTradeNo()
	if len(req.Order.Goods) > 0 {
		for _, item := range req.Order.Goods {
			result.PaymentOrderProducts = append(result.PaymentOrderProducts, entity.Product{
				TradeNo:     result.TradeNo,
				ProductName: item.Name,
				AppNo:       req.AppNo,
				MchNo:       req.MchNo,
				Quantity:    int64(item.Quantity),
				Price:       int64(item.Price),
				Sku:         item.Sku,
				ProductDesc: item.Desc,
				URL:         "",
				OrderNo:     req.Order.OrderNo,
			})
		}
	}
	return result
}
func (s *paymentService) buildChannelPayOrder(req *public.PaymentOrder, appInfo *model.Application, account *entity.PaymentAccount, tradeNo string, selfCallback string) (*dto.PayOrder, error) {
	request := &dto.PayOrder{
		Order: dto.Order{
			OrderNo: tradeNo,
			Amount: dto.Amount{
				Total:    req.Order.Amount.Total,
				Currency: req.Order.Amount.Currency,
			},
			PayAmount: dto.Amount{
				Total:    req.Order.Amount.Total,
				Currency: req.Order.Amount.Currency,
			},
			Goods:    make([]dto.Goods, 0, len(req.Order.Goods)),
			Subject:  req.Order.Subject,
			Desc:     req.Order.Desc,
			CreateAt: req.Order.CreateAt,
		},
		Payer: dto.Payer{
			OpenID:  req.Payer.OpenID,
			UnionID: req.Payer.UnionID,
			AppID:   req.Payer.AppID,
		},
		RedirectUrl:    req.RedirectUrl,
		TimeExpire:     req.TimeExpire,
		NotifyUrl:      "",
		PassBackParams: req.PassBackParams,
		SettleInfo:     nil,
		SceneInfo:      nil,
		RiskFund:       nil,
		AlipayExtra:    nil,
	}
	request.NotifyUrl = selfCallback
	if req.Order.Goods != nil {
		for _, v := range req.Order.Goods {
			request.Order.Goods = append(request.Order.Goods, dto.Goods{
				Name:     v.Name,
				Sku:      v.Sku,
				Price:    v.Price,
				Quantity: v.Quantity,
				Desc:     v.Desc,
			})
		}
	}
	if req.SceneInfo != nil {
		request.SceneInfo = &dto.SceneInfo{
			ClientIp: req.SceneInfo.ClientIp,
			DeviceID: req.SceneInfo.DeviceID,
			Device:   req.SceneInfo.Device,
			System:   req.SceneInfo.System,
		}
		if req.SceneInfo.ApplicationInfo != nil {
			request.SceneInfo.ApplicationInfo = dto.ApplicationInfo{
				AppName:    req.SceneInfo.ApplicationInfo.AppName,
				Url:        req.SceneInfo.ApplicationInfo.Url,
				AppPackage: req.SceneInfo.ApplicationInfo.AppPackage,
			}
		}
	}
	return request, nil
}
func (s *paymentService) getPaymentLockKey(orderNo, appNo string) string {
	return fmt.Sprintf(enum.CacheLockPayment, orderNo, appNo)
}
func getNotifyUrl(appInfo *model.Application, cfg config3.Config) (*url.URL, error) {
	var parse *url.URL
	var err error
	if appInfo.IsCustomerDomain == 1 && appInfo.CustomerDomain != "" {
		parse, err = url.Parse(appInfo.CustomerDomain)
		if err != nil {
			return nil, errors.WrapError(errors.ErrCodeInvalidParam, "构建支付回调地址失败", err)
		}
		parse = parse.JoinPath(cfg.ProxyNotifyPrefix)
	} else {
		parse, err = url.Parse(cfg.ApiHost)
		if err != nil {
			return nil, errors.WrapError(errors.ErrCodeInvalidParam, "构建支付回调地址失败", err)
		}
	}
	return parse, nil
}
