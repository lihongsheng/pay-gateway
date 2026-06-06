package service

import (
	"context"
	errors2 "errors"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	config2 "github.com/lihongsheng/pay-gateway/plugin/payment/config"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/url"
	"time"
)

// CashierService 收银台支付
// 1. 创建支付订单-create; 只需订单信息
// 2. 基于支付订单创建成功的支付连接，引导用户进行支付； 调用方
// 3. 用户打开支付连接，完成支付 - payment；客户端渲染界面
// 4. 获取支付结果. - query
// 5. 获取支付结果成功，则返回支付成功页面，否则返回支付失败页面. - query
type CashierService interface {
	Create(ctx context.Context, req *public.CashierCreateOrder) (*public.CashierCreateOrderResponse, error)
	Payment(ctx context.Context, req *public.CashierPaymentOrder) (*public.PaymentOrderResponse, error)
	Query(ctx context.Context, req *public.CashierQueryOrder) (*public.PaymentOrderDetail, error)
}

type cashierService struct {
	svc    *svc.ServiceContext
	domain *domain.ServiceGroup
}

func NewCashierService(svc *svc.ServiceContext, domain *domain.ServiceGroup) CashierService {
	return &cashierService{
		svc:    svc,
		domain: domain,
	}
}

func (s *cashierService) Query(ctx context.Context, req *public.CashierQueryOrder) (*public.PaymentOrderDetail, error) {
	l := log.WithCtx(ctx)
	aggregateToken, err := getAggregateToken(req.Token)
	if err != nil {
		l.Error("cashierQuery", zap.Error(err), zap.Any("token", req))
		return nil, err
	}
	appInfo, err := s.svc.AppRepo.GetByAppNoFormCache(ctx, aggregateToken.AppNo)
	if err != nil {
		l.Error("cashierQuery", zap.Error(err), zap.String("aggregateToken.AppNo", aggregateToken.AppNo))
		return nil, errors.NewError(errors.ErrCodeInvalidParam, "token无效")
	}
	// 已有支付成功订单
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: req.OrderNo,
		TradeNo: "",
		AppNo:   appInfo.AppNo,
		MchNo:   appInfo.MchNo,
	})
	if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常扫后再试", err)
	}
	if payOrder == nil {
		return nil, errors.NewError(errors.ErrPayOrderNotExist, "订单不存在")
	}
	return buildPaymentOrderDetail(payOrder), nil
}

func (s *cashierService) getPaymentLockKey(orderNo, appNo string) string {
	return fmt.Sprintf(enum.CacheLockCashier, orderNo, appNo)
}

func (s *cashierService) Create(ctx context.Context, req *public.CashierCreateOrder) (*public.CashierCreateOrderResponse, error) {
	l := log.WithCtx(ctx)
	l.Info("CashierCreate", zap.Any("req", req))
	if err := req.Validate(); err != nil {
		return nil, err
	}
	appInfo, err := getAppInfoFromCtx(ctx, req.AppNo, s.svc)
	if err != nil {
		l.Error("CashierPayment get app info fail", zap.Error(err), zap.String("app_no", req.AppNo))
		return nil, err
	}
	lock, err := global.GVA_REDIS.SetNX(ctx, s.getPaymentLockKey(req.Order.OrderNo, req.AppNo), time.Now().Unix(), enum.PaymentLockExpire).Result()
	defer func() {
		global.GVA_REDIS.Del(ctx, s.getPaymentLockKey(req.Order.OrderNo, req.AppNo))
	}()
	if !lock && err == nil {
		return nil, errors.NewError(errors.ErrPayPending, "重复支付，请稍后再试")
	}
	req.MchNo = appInfo.MchNo
	// 4. 幂等校验：检查订单是否已存在
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: req.Order.OrderNo,
		AppNo:   req.AppNo,
		MchNo:   req.MchNo,
	})
	if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
		l.Error("CashierCreateFail", zap.Error(err), zap.Any("order", req.Order))
		return nil, errors.WrapError(errors.ErrCodeSysError, "查询订单失败", err)
	}
	token, err := genAggregateToken(dto.AggregateToken{
		AppNo: appInfo.AppNo,
	})
	if err != nil {
		return nil, errors.NewError(errors.ErrCodeInvalidParam, "token无效")
	}
	openUrl, err := s.genOpenUrl(req, appInfo, token)
	if err != nil {
		return nil, errors.WrapError(errors.ErrCodeSysError, "生成支付连接失败", err)
	}
	if payOrder == nil {
		payOrder, err = s.buildPaymentOrder(req, appInfo, token)
		if err != nil {
			l.Error("CashierCreateFail", zap.Error(err), zap.Any("order", req.Order))
			return nil, errors.WrapError(errors.ErrCodeSysError, "创建订单失败", err)
		}
		err = s.svc.PaymentOrderRepo.Save(ctx, payOrder)
		if err != nil {
			l.Error("CashierCreateFail", zap.Error(err), zap.Any("order", req.Order))
			return nil, errors.WrapError(errors.ErrCodeSysError, "创建订单失败", err)
		}
		eventInfo := payOrder.GetEvents()
		if eventInfo != nil {
			eventErr := s.svc.Event.Publish(ctx, eventInfo, config2.Config.Topic.PaymentStatus)
			if eventErr != nil {
				l.Error("发布订单状态失败", zap.Error(eventErr), zap.Any("event", eventInfo))
			}
		}
	}
	exT := time.Unix(payOrder.TimeExpire, 0)
	return &public.CashierCreateOrderResponse{
		OrderNo: payOrder.OrderNo,
		OpenURL: openUrl,
		Amount: public.Amount{
			Total:    payOrder.PaymentAmount,
			Currency: payOrder.Currency,
		},
		OrderExpireTime: exT,
		TradeNo:         payOrder.TradeNo,
		Status:          payment.Status(payOrder.Status),
	}, nil
}

func (s *cashierService) genOpenUrl(req *public.CashierCreateOrder, appInfo *model.Application, token string) (string, error) {
	u, err := getWebBaseUrl(appInfo)
	if err != nil {
		return "", err
	}
	u = u.JoinPath(enum.CashierPayment)
	query := url.Values{}
	query.Add("action_token", token)
	query.Add("action", enum.H5ActionHub)
	query.Add("order_no", req.Order.OrderNo)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (s *cashierService) buildPaymentOrder(req *public.CashierCreateOrder, app *model.Application, token string) (*entity.PaymentOrder, error) {
	result := &entity.PaymentOrder{
		PaymentAccountNo:     "",
		AppNo:                req.AppNo,
		MchNo:                req.MchNo,
		OrderSubject:         req.Order.Subject,
		PaymentMethod:        0,
		PaymentProduct:       0,
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
		UserOpenid:           "",
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
	if result.RedirectURL == "" {
		u, err := getNotifyUrl(app)
		if err != nil {
			return nil, err
		}
		u = u.JoinPath(enum.PaymentSuccessPath)
		query := u.Query()
		query.Add(enum.H5ActionToken, token)
		query.Add("order_no", req.Order.OrderNo)
		u.RawQuery = query.Encode()
		result.RedirectURL = u.String()
	}
	return result, nil
}

func (s *cashierService) Payment(ctx context.Context, req *public.CashierPaymentOrder) (*public.PaymentOrderResponse, error) {
	l := log.WithCtx(ctx)
	l.Info("CashierPayment", zap.Any("req", req))
	if err := req.Validate(); err != nil {
		return nil, err
	}
	aggregateToken, err := getAggregateToken(req.Token)
	if err != nil {
		l.Error("cashierPayment", zap.Error(err), zap.Any("token", req))
		return nil, errors.NewError(errors.ErrCodeInvalidParam, "token无效")
	}
	appInfo, err := s.svc.AppRepo.GetByAppNoFormCache(ctx, aggregateToken.AppNo)
	if err != nil {
		l.Error("cashierPayment", zap.Error(err), zap.String("aggregateToken.AppNo", aggregateToken.AppNo))
		return nil, errors.NewError(errors.ErrCodeInvalidParam, "token无效")
	}
	account, err := s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, req.AccountNo)
	if err != nil {
		l.Error("cashierPayment", zap.Error(err), zap.String("req.AccountNo", req.AccountNo))
		return nil, errors.NewError(errors.ErrCodeInvalidParam, "无效支付账号")
	}
	// 已有支付成功订单
	payOrder, err := s.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: req.OrderNo,
		TradeNo: "",
		AppNo:   appInfo.AppNo,
		MchNo:   appInfo.MchNo,
	})
	if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WrapError(errors.ErrCodeSysError, "系统异常扫后再试", err)
	}
	if payOrder.IsConfirmed() {
		return &public.PaymentOrderResponse{
			OrderNo:    payOrder.OrderNo,
			TradeNo:    payOrder.TradeNo,
			Status:     payment.Status(payOrder.Status),
			OutTradeNo: payOrder.OutMchTradeNo,
			Amount: public.Amount{
				Total:    payOrder.PaymentAmount,
				Currency: payOrder.Currency,
			},
			RedirectURL: payOrder.RedirectURL,
		}, nil
	}
	callbackUrl, err := getSelfCallbackUrl(appInfo, account, req.OrderNo, enum.PaymentNotify)
	if err != nil {
		return nil, err
	}
	result, _, err := s.domain.PaymentService.Payment(ctx, payOrder, s.buildPayment(appInfo, payOrder, req), appInfo, account, callbackUrl)
	if err != nil {
		l.Error("cashierPayment", zap.Error(err), zap.Any("req.AccountNo", req))
		return nil, err
	}
	return result, nil
}

func (s *cashierService) buildPayment(app *model.Application, payOrder *entity.PaymentOrder, req *public.CashierPaymentOrder) *public.PaymentOrder {
	r := &public.PaymentOrder{
		Order: public.Order{
			OrderNo: req.OrderNo,
			Amount: public.Amount{
				Total:    payOrder.PaymentAmount,
				Currency: payOrder.Currency,
			},
			Subject:  app.PaymentTitle,
			Desc:     app.PaymentTitle,
			CreateAt: time.Now(),
		},
		Payer: public.Payer{
			OpenID: req.OpenId,
		},
		RedirectUrl:    payOrder.RedirectURL,
		TimeExpire:     payOrder.TimeExpire,
		NotifyUrl:      payOrder.NotifyURL,
		PassBackParams: payOrder.ExtParam,
		AlipayExtra:    req.AlipayExtra,
		PaymentMethod:  req.PaymentMethod,
		PaymentProduct: req.PaymentProduct,
		SceneInfo:      req.SceneInfo,
		AccountNo:      req.AccountNo,
		AppNo:          app.AppNo,
		MchNo:          app.MchNo,
		RequestID:      req.RequestID,
	}
	return r
}
