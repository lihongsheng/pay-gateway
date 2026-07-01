package domain

import (
	"context"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	paySdk "github.com/lihongsheng/payment-sdk"
	enum2 "github.com/lihongsheng/payment-sdk/enum"
	"github.com/lihongsheng/payment-sdk/enum/channel"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type Router interface {
	GetAvailablePayment(ctx context.Context, req *public.QueryCheckoutPaymentMethod, appInfo *model.Application) ([]*entity.PaymentAccount, error)
	Router(ctx context.Context, apps []*entity.PaymentAccount, appInfo *model.Application, filterAccountNo []string) (available []*entity.PaymentAccount, last *entity.PaymentAccount, err error)
	Rand(req []*entity.PaymentAccount) (*entity.PaymentAccount, error)
}
type routerService struct {
	paymentAccountRepo repo.PaymentAccountRepo
	routerRepo         repo.RouterRepo
	RuleEngine         RuleEngine
}

func NewRouterService(paymentAccountRepo repo.PaymentAccountRepo, routerRepo repo.RouterRepo) Router {
	return &routerService{
		paymentAccountRepo: paymentAccountRepo,
		routerRepo:         routerRepo,
		RuleEngine:         NewRuleEngineService(),
	}
}
func (s *routerService) Rand(req []*entity.PaymentAccount) (*entity.PaymentAccount, error) {
	if len(req) == 0 {
		return nil, errors.NewError(errors.ErrPayChannelNotSupport, "未找到支持的支付渠道")
	}
	// 随机选择一个支付账户
	return req[rand.Intn(len(req))], nil
}
func (s *routerService) Router(ctx context.Context, apps []*entity.PaymentAccount, appInfo *model.Application, filterAccountNo []string) (available []*entity.PaymentAccount, last *entity.PaymentAccount, err error) {
	var result = make([]*entity.PaymentAccount, 0, len(apps))
	if len(apps) == 0 {
		return nil, nil, errors.NewError(errors.ErrPayChannelNotSupport, "请配置支付渠道")
	}
	log.WithCtx(ctx).Info("支付渠道路由", zap.Any("appInfo", appInfo))
	if appInfo.MultiChannel == enum.MultiChannel_Disable {
		for _, item := range apps {
			if item.Status.Disable() {
				continue
			} else {
				result = append(result, item)
			}
		}
		if len(result) == 0 {
			return nil, nil, errors.NewError(errors.ErrPayChannelNotSupport, "未找到可用支付渠道")
		}
		return result, apps[len(apps)-1], nil
	}
	r, err := s.multiChannelPayment(ctx, apps, appInfo, filterAccountNo)
	if err != nil {
		log.WithCtx(ctx).Error("支付渠道路由失败", zap.Error(err), zap.Any("apps", apps))
		return nil, apps[len(apps)-1], err
	}
	return r, apps[len(apps)-1], nil
}
func (s *routerService) multiChannelPayment(ctx context.Context, apps []*entity.PaymentAccount, appInfo *model.Application, filterAccountNo []string) ([]*entity.PaymentAccount, error) {
	// 可以有多个微信支付宝，需要查验某一个微信或者支付宝支付账户是否被 微信/支付宝 平台 封禁，封禁后，则不能再使用。只返回可以使用的。目前先按照 无多个微信和支付宝一样处理。
	// 需要当 获取不到微信或者支付宝用户openid , 支付返回非签名，网络等错误的时候，更新某一个微信或者支付宝账户为封禁状态，然后剔除。
	var result = make([]*entity.PaymentAccount, 0, len(apps))
	for _, item := range apps {
		if item.Status.Disable() {
			continue
		} else {
			result = append(result, item)
		}
	}
	statics, err := s.routerRepo.GetAppRouterStatisticFromCache(ctx, appInfo.AppNo, time.Now())
	log.WithCtx(ctx).Info("支付渠道路由统计", zap.Any("statics", statics))
	if err != nil {
		return result, nil
	}
	var routerConditions = make([]*dto.RuleConditions, 0, len(statics))
	if len(statics) > 0 {
		for _, item := range result {
			tmp := &dto.RuleConditions{
				AccountNO:       item.AccountNo,
				RequestSuccess:  0,
				UserLimit:       false,
				MaxLimit:        item.MaxLimit,
				FilterAccountNO: filterAccountNo,
			}
			for _, item2 := range statics {
				if item2.AccountNo == item.AccountNo {
					userLimit := false
					if item2.UserLimit > 0 {
						userLimit = true
					}
					tmp.UserLimit = userLimit
					tmp.RequestSuccess = item2.SuccessRequests
					paymentLimit := false
					if item2.PaymentLimit > 0 {
						paymentLimit = true
					}
					tmp.PaymentLimit = paymentLimit
					break
				}
			}
			routerConditions = append(routerConditions, tmp)
		}
		accountNos, err := s.RuleEngine.Evaluate(ctx, routerConditions)
		log.WithCtx(ctx).Info("支付渠道路由条件", zap.Any("routerConditions", routerConditions), zap.Any("accountNos", accountNos), zap.Error(err))
		if err != nil {
			return nil, err
		}
		var tmps = make([]*entity.PaymentAccount, 0, len(accountNos))
		for _, item := range result {
			for _, item2 := range accountNos {
				if item.AccountNo == item2 {
					tmps = append(tmps, item)
				}
			}
		}
		if len(tmps) > 0 {
			result = tmps
		}
	}
	if len(result) == 0 {
		return nil, errors.NewError(errors.ErrPayChannelNotSupport, "未找到支持的支付渠道")
	}
	return result, nil
}
func (s *routerService) GetAvailablePayment(ctx context.Context, req *public.QueryCheckoutPaymentMethod, appInfo *model.Application) ([]*entity.PaymentAccount, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	apps, err := s.paymentAccountRepo.GetByAppNoFormCache(ctx, req.AppNo, appInfo.AccountVersion)
	if err != nil {
		return nil, errors.WrapError(errors.ErrPayChannelNotSupport, "未找到支持的支付渠道", err)
	}
	var availableApps = make([]*entity.PaymentAccount, 0, len(apps))
	if req.PaymentMethod > 0 && req.PaymentProduct > 0 {
		availableApps = s.getPaymentAccount(ctx, apps, req.PaymentMethod, req.PaymentProduct)
	} else if req.Device != "" {
		availableApps, err = s.getPaymentAccountByDevice(ctx, apps, req.Device)
		if err != nil {
			return nil, err
		}
	}
	if len(availableApps) == 0 {
		return nil, errors.NewError(errors.ErrPayChannelNotSupport, "未找到支持的支付渠道")
	}
	return availableApps, err
}
func (s *routerService) getPaymentAccountByDevice(ctx context.Context, apps []*entity.PaymentAccount, device enum2.Device) ([]*entity.PaymentAccount, error) {
	var result = make([]*entity.PaymentAccount, 0, len(apps))
	for _, item := range apps {
		exists := false
		if item.Status.Disable() {
			continue
		}
		drivePay, err := paySdk.GetPaymentDriver(channel.Channel(channel.Channel_value[item.Channel]))
		if err != nil {
			return nil, errors.WrapError(errors.ErrPayChannelNotSupport, "未找到支持的支付渠道", err)
		}
		for _, paymentAccount := range item.PaymentMethod {
			if drivePay.IsSupportPayment(payment.PaymentProduct(payment.PaymentProduct_value[paymentAccount.Product]), device) {
				exists = true
				break
			}
		}
		if exists {
			result = append(result, item)
		}
	}
	return result, nil
}
func (s *routerService) getPaymentAccount(ctx context.Context, apps []*entity.PaymentAccount, paymentMethod payment.Payment, product payment.PaymentProduct) []*entity.PaymentAccount {
	var result = make([]*entity.PaymentAccount, 0, len(apps))
	for _, item := range apps {
		exists := false
		if item.Status.Disable() {
			continue
		}
		for _, paymentAccount := range item.PaymentMethod {
			if paymentAccount.Method == paymentMethod.String() && paymentAccount.Product == product.String() {
				exists = true
				break
			}
		}
		if exists {
			result = append(result, item)
		}
	}
	return result
}
