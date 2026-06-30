package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/infrastructure"
	adapter2 "github.com/lihongsheng/pay-gateway/plugin/payment/infrastructure/adapter"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"github.com/lihongsheng/payment-sdk/enum/channel"
	"go.uber.org/zap"
)

// 聚合扫码支付

type Aggregate interface {
	GetApplication(ctx context.Context, appNo string) (*model.Application, error)
	GetUserOpenID(ctx context.Context, req public.AggregateUserOpenIDRequest) (string, error)
	GetRedirectUrl(ctx context.Context, req *public.AggregateRedirectUrlRequest, basePath string) (string, error)
	Payment(ctx context.Context, req *public.AggregateOrder, notify string) (*public.PaymentOrderResponse, error)
	Query(ctx context.Context, token string, orderNo string) (*public.PaymentOrderDetail, error)
	GenQrCode(ctx context.Context, appNo string) (string, error)
	GenTestQrCode(ctx context.Context, appNo string, accNo string) (qrLink string, orderNo string, err error)
}

type aggregateService struct {
	appRepo            repo.ApplicationRepo
	paymentAccountRepo repo.PaymentAccountRepo
	aggregateAdapter   map[channel.Channel]adapter2.AggregateUser
	eventPublisher     infrastructure.Event
	domainPayment      domain.PaymentService
	domainRouter       domain.Router
}

func NewAggregateService(
	appRepo repo.ApplicationRepo,
	paymentAccountRepo repo.PaymentAccountRepo,
	aggregateAdapter map[channel.Channel]adapter2.AggregateUser,
	eventPublisher infrastructure.Event,
	domainPayment domain.PaymentService,
	domainRouter domain.Router,
) Aggregate {
	return &aggregateService{
		appRepo:            appRepo,
		paymentAccountRepo: paymentAccountRepo,
		aggregateAdapter:   aggregateAdapter,
		eventPublisher:     eventPublisher,
		domainPayment:      domainPayment,
		domainRouter:       domainRouter,
	}
}

// DefaultAggregate 包级单例
var DefaultAggregate Aggregate

// GenTestQrCode 聚合支付测试二维码
func (a *aggregateService) GenTestQrCode(ctx context.Context, appNo string, accNo string) (qrLink string, orderNo string, err error) {
	l := log.WithCtx(ctx)
	appInfo, err := a.appRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		l.Error("AggregatePayment", zap.Error(err), zap.String("aggregateToken.AppNo", appNo))
		return "", "", errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	token, err := genAggregateToken(dto.AggregateToken{
		AppNo: appInfo.AppNo,
	})
	if err != nil {
		return "", "", errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	orderNo = fmt.Sprintf("PT%d%s%s", time.Now().UnixMilli(), utils.GenRandNumToStr(), utils.HashCrc(token, 4))
	u, err := getWebBaseUrl(appInfo)
	if err != nil {
		return "", "", errors2.WrapError(errors2.ErrCodeInvalidParam, "web域名格式错误:", err)
	}
	u = u.JoinPath(enum.PaymentTestPath)
	query := url.Values{}
	query.Add("action_token", token)
	query.Add("action", enum.H5ActionHub)
	query.Add("account_no", accNo)
	query.Add("order_no", orderNo)
	u.RawQuery = query.Encode()
	return u.String(), orderNo, nil
}

// GenQrCode 聚合支付二维码
func (a *aggregateService) GenQrCode(ctx context.Context, appNo string) (string, error) {
	l := log.WithCtx(ctx)
	appInfo, err := a.appRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		l.Error("AggregatePayment", zap.Error(err), zap.String("aggregateToken.AppNo", appNo))
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	u, err := getWebBaseUrl(appInfo)
	if err != nil {
		return "", errors2.WrapError(errors2.ErrCodeInvalidParam, "web域名格式错误:", err)
	}
	u = u.JoinPath(enum.PaymentBaseIndexPath)
	query := url.Values{}
	token, err := genAggregateToken(dto.AggregateToken{
		AppNo: appInfo.AppNo,
	})
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	query.Add("action_token", token)
	query.Add("action", enum.H5ActionHub)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (a *aggregateService) GetApplication(ctx context.Context, token string) (*model.Application, error) {
	l := log.WithCtx(ctx)
	aggregateToken, err := getAggregateToken(token)
	if err != nil {
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	app, err := a.appRepo.GetByAppNoFormCache(ctx, aggregateToken.AppNo)
	if err != nil {
		l.Error("GetApplication", zap.Error(err), zap.String("accountNo", aggregateToken.AppNo))
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "无效应用")
	}
	return app, nil
}

func (a *aggregateService) GetUserOpenID(ctx context.Context, req public.AggregateUserOpenIDRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}
	l := log.WithCtx(ctx)
	aggregateToken, err := getAggregateToken(req.Token)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	accountDetail, err := a.paymentAccountRepo.GetByAccountNoFormCache(ctx, req.AccountNo)
	if err != nil {
		l.Error("GetApplication", zap.Error(err), zap.String("accountNo", req.AccountNo))
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效支付账号")
	}
	if aggregateToken.AppNo != accountDetail.AppNo {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效应用")
	}
	app, err := a.appRepo.GetByAppNoFormCache(ctx, accountDetail.AppNo)
	if err != nil {
		l.Error("GetApplication", zap.Error(err), zap.String("app_no", accountDetail.AppNo), zap.Any("req", req))
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效应用")
	}
	agg, ok := a.aggregateAdapter[channel.Channel(channel.Channel_value[accountDetail.Channel])]
	if ok {
		openid, err := agg.GetUserOpenID(ctx, req.AuthCode, accountDetail, app, req.PaymentMethod)
		if err != nil {
			if errors2.IsUserLimitError(err) {
				e := a.eventPublisher.Publish(ctx, &event.UserLoginLimitEvent{
					AppNo:     app.AppNo,
					UserLimit: true,
					AccountNo: req.AccountNo,
					CreateAt:  time.Now(),
				}, global.Cfg.Payment.Topic.UserLimit)
				if e != nil {
					l.Error("GetUserOpenIDPublish", zap.Error(e), zap.String("app_no", app.AppNo), zap.Any("req", req))
				}
			}
			l.Error("GetUserOpenID", zap.Error(err), zap.Any("req", req))
			return "", err
		}
		return openid, nil
	}
	return "", errors2.NewError(errors2.ErrCodeInvalidParam, "不支持的支付渠道")
}

func (a *aggregateService) GetRedirectUrl(ctx context.Context, req *public.AggregateRedirectUrlRequest, basePath string) (string, error) {
	l := log.WithCtx(ctx)
	l.Info("GetRedirectUrl", zap.Any("req", req))
	if err := req.Validate(); err != nil {
		return "", err
	}
	if global.Cfg.Payment.IsTest() && req.FilterAccountNo == "" && req.OrderNo == "T1772159917" {
		req.FilterAccountNo = "P4966c7e946c00,P6b5e44d8e7c00"
	}
	if global.Cfg.Payment.IsTest() && req.FilterAccountNo == "" && req.OrderNo == "T17721599111" {
		req.FilterAccountNo = "P4966c7e946c00"
	}
	aggregateToken, err := getAggregateToken(req.Token)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	appInfo, err := a.appRepo.GetByAppNoFormCache(ctx, aggregateToken.AppNo)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	account, err := a.getAccount(ctx, req, appInfo, l)
	if err != nil {
		return "", err
	}
	u, err := a.buildRedirectUrl(ctx, req, appInfo, account, basePath)
	if err != nil {
		l.Error("GetRedirectUrl", zap.Error(err), zap.String("app_no", appInfo.AppNo))
		return "", err
	}
	return u, nil
}

func (a *aggregateService) getAccount(ctx context.Context, req *public.AggregateRedirectUrlRequest, app *model.Application, l *zap.Logger) (*entity.PaymentAccount, error) {
	if req.AccountNo != "" {
		account, err := a.paymentAccountRepo.GetByAccountNoFormCache(ctx, req.AccountNo)
		if err != nil {
			return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "无效支付账号")
		}
		if account.AppNo != app.AppNo {
			return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "无效应用")
		}
		return account, nil
	}
	available, err := a.domainRouter.GetAvailablePayment(ctx, &public.QueryCheckoutPaymentMethod{
		PaymentMethod:  req.PaymentMethod,
		PaymentProduct: req.PaymentProduct,
		AppNo:          app.AppNo,
	}, app)
	if err != nil {
		l.Error("GetRedirectUrl", zap.Error(err), zap.String("app_no", app.AppNo))
		return nil, errors2.WrapError(errors2.ErrPayChannelNotSupport, "未找到支持的支付渠道", err)
	}
	var filterAccountNo []string
	if req.FilterAccountNo != "" {
		filterAccountNo = strings.Split(req.FilterAccountNo, ",")
	}
	var lastAvailable *entity.PaymentAccount
	available, lastAvailable, err = a.domainRouter.Router(ctx, available, app, filterAccountNo)
	if err != nil {
		l.Error("GetRedirectUrlRouter", zap.Error(err), zap.String("app_no", app.AppNo), zap.Any("req", req))
		// 兜底，保障使用最早的一个。GetByAccountNoFormCache 是倒序。
		if lastAvailable == nil {
			return nil, err
		}
		available = []*entity.PaymentAccount{available[len(available)-1]}
	}
	account, _ := a.domainRouter.Rand(available)
	l.Info("GetRedirectUrl", zap.String("app_no", app.AppNo), zap.Any("account", account), zap.Any("available", available))
	return account, nil
}

func (a *aggregateService) buildRedirectUrl(ctx context.Context, req *public.AggregateRedirectUrlRequest, app *model.Application, accountDetail *entity.PaymentAccount, basePath string) (string, error) {
	callbackUrl := ""
	u, err := getWebBaseUrl(app)
	if err != nil {
		return "", errors2.WrapError(errors2.ErrCodeInvalidParam, "web域名格式错误:", err)
	}
	// u = u.JoinPath(enum.PaymentBaseIndexPath)
	u = u.JoinPath(basePath)
	query := u.Query()
	query.Add(enum.H5Action, enum.H5ActionAuth)
	query.Add(enum.H5ActionToken, req.Token)
	query.Add("account_no", accountDetail.AccountNo)
	if req.OrderNo != "" {
		query.Add("order_no", req.OrderNo)
	}
	if len(req.FilterAccountNo) > 0 {
		query.Add("filter_account_no", req.FilterAccountNo)
	}
	u.RawQuery = query.Encode()
	callbackUrl = u.String()
	agg, ok := a.aggregateAdapter[channel.Channel(channel.Channel_value[accountDetail.Channel])]
	if ok {
		return agg.RedirectUrl(ctx, callbackUrl, accountDetail, req.PaymentMethod)
	}
	return "", errors2.NewError(errors2.ErrCodeInvalidParam, "不支持的支付渠道")
}

func (a *aggregateService) Payment(ctx context.Context, req *public.AggregateOrder, notify string) (*public.PaymentOrderResponse, error) {
	l := log.WithCtx(ctx)
	l.Info("AggregatePayment", zap.Any("req", req))
	if err := req.Validate(); err != nil {
		return nil, err
	}
	aggregateToken, err := getAggregateToken(req.Token)
	if err != nil {
		l.Error("AggregatePayment", zap.Error(err), zap.Any("token", req))
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	appInfo, err := a.appRepo.GetByAppNoFormCache(ctx, aggregateToken.AppNo)
	if err != nil {
		l.Error("AggregatePayment", zap.Error(err), zap.String("aggregateToken.AppNo", aggregateToken.AppNo))
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	account, err := a.paymentAccountRepo.GetByAccountNoFormCache(ctx, req.AccountNo)
	if err != nil {
		l.Error("AggregatePayment", zap.Error(err), zap.String("req.AccountNo", req.AccountNo))
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "无效支付账号")
	}
	order, err := a.buildPayment(appInfo, req)
	if err != nil {
		l.Error("AggregatePayment", zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	callbackUrl, err := getSelfCallbackUrl(appInfo, account, order.Order.OrderNo, notify)
	r, _, err := a.domainPayment.Payment(ctx, nil, order, appInfo, account, callbackUrl)
	if err != nil {
		l.Error("AggregatePayment", zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	return r, nil
}

func (a *aggregateService) buildPayment(app *model.Application, req *public.AggregateOrder) (*public.PaymentOrder, error) {
	r := &public.PaymentOrder{
		Order: public.Order{
			OrderNo: req.OrderNo,
			Amount: public.Amount{
				Total:    req.Amount,
				Currency: "CNY",
			},
			Subject:  app.PaymentTitle,
			Desc:     app.PaymentTitle,
			CreateAt: time.Now(),
		},
		Payer: public.Payer{
			OpenID: req.OpenID,
		},
		RedirectUrl:    "",
		TimeExpire:     time.Now().Add(enum.DefaultExpireTime).Unix(),
		NotifyUrl:      "",
		PassBackParams: "",
		AlipayExtra:    nil,
		PaymentMethod:  req.PaymentMethod,
		PaymentProduct: req.PaymentProduct,
		SceneInfo:      nil,
		AccountNo:      req.AccountNo,
		AppNo:          app.AppNo,
		MchNo:          app.MchNo,
		RequestID:      req.RequestID,
	}
	u, err := getNotifyUrl(app)
	if err != nil {
		return nil, err
	}
	u = u.JoinPath(enum.PaymentSuccessPath)
	query := u.Query()
	query.Add(enum.H5ActionToken, req.Token)
	query.Add("order_no", req.OrderNo)
	u.RawQuery = query.Encode()
	r.RedirectUrl = u.String()
	return r, nil
}

func (a *aggregateService) Query(ctx context.Context, token string, orderNo string) (*public.PaymentOrderDetail, error) {
	l := log.WithCtx(ctx)
	aggregateToken, err := getAggregateToken(token)
	if err != nil {
		l.Error("AggregatePaymentQuery", zap.Error(err), zap.Any("token", token), zap.String("orderNo", orderNo))
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	appInfo, err := a.appRepo.GetByAppNoFormCache(ctx, aggregateToken.AppNo)
	if err != nil {
		l.Error("AggregatePaymentQuery", zap.Error(err), zap.String("aggregateToken.AppNo", aggregateToken.AppNo), zap.String("orderNo", orderNo))
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	order, err := a.domainPayment.Query(ctx, &public.QueryPaymentRequest{
		OrderNo: orderNo,
		AppNo:   appInfo.AppNo,
		TradeNo: "",
		MchNo:   appInfo.MchNo,
	}, appInfo)
	if err != nil {
		l.Error("AggregatePaymentQuery", zap.Error(err), zap.Any("req", orderNo))
		return nil, err
	}
	return buildPaymentOrderDetail(order), nil
}

func getAggregateToken(token string) (*dto.AggregateToken, error) {
	var aggregateToken *dto.AggregateToken
	str, err := utils.AesGCMDecryptHex(token, []byte(global.Cfg.Payment.Salt))
	if err != nil {
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	err = json.Unmarshal([]byte(str), &aggregateToken)
	if err != nil || aggregateToken.AppNo == "" {
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	return aggregateToken, nil
}

func genAggregateToken(req dto.AggregateToken) (string, error) {
	by, _ := json.Marshal(req)
	str, err := utils.AesGCMEncryptHex(by, []byte(global.Cfg.Payment.Salt))
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "token无效")
	}
	return str, nil
}

func getNotifyUrl(appInfo *model.Application) (*url.URL, error) {
	var parse *url.URL
	var err error
	if appInfo.IsCustomerDomain == 1 && appInfo.CustomerDomain != "" {
		parse, err = url.Parse(appInfo.CustomerDomain)
		if err != nil {
			return nil, errors2.WrapError(errors2.ErrCodeInvalidParam, "构建支付回调地址失败", err)
		}
		parse = parse.JoinPath(global.Cfg.Payment.ProxyNotifyPrefix)
	} else {
		parse, err = url.Parse(global.Cfg.Payment.ApiHost)
		if err != nil {
			return nil, errors2.WrapError(errors2.ErrCodeInvalidParam, "构建支付回调地址失败", err)
		}
	}
	return parse, nil
}

func getSelfCallbackUrl(appInfo *model.Application, account *entity.PaymentAccount, orderNO string, notify string) (string, error) {
	parse, err := getNotifyUrl(appInfo)
	if err != nil {
		return "", err
	}
	if notify == "" {
		notify = enum.PaymentNotify
	}
	path := fmt.Sprintf(notify, account.Channel, appInfo.MchNo, appInfo.AppNo, orderNO)
	parse = parse.JoinPath(path)
	return parse.String(), nil
}

func getWebBaseUrl(app *model.Application) (*url.URL, error) {
	if global.Cfg.Payment.WebHost == "" {
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "未配置web域名")
	}
	u, err := url.Parse(global.Cfg.Payment.WebHost)
	if err != nil {
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, fmt.Sprintf("web域名格式错误: %s", err.Error()))
	}
	if app.IsCustomerDomain == enum.CustomerDomainStatus_Enable {
		u, err = url.Parse(app.CustomerDomain)
		if err != nil {
			return nil, errors2.WrapError(errors2.ErrCodeInvalidParam, "web域名格式错误:", err)
		}
		u = u.JoinPath(global.Cfg.Payment.ProxyNotifyPrefix)
	}
	return u, nil
}
