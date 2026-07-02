package adapter

import (
	"context"
	"encoding/json"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	alipayConfig "github.com/lihongsheng/payment-sdk/adapter/alipay/config"
	fuiouConfig "github.com/lihongsheng/payment-sdk/adapter/fuiou/config"
	wechatConfig "github.com/lihongsheng/payment-sdk/adapter/wxpay/config"
	"github.com/lihongsheng/payment-sdk/enum/payment"
)

type FuiouUser struct {
	alipayUser *AlipayUser
	wechatUser *WechatUser
}

func NewFuiouUser() AggregateUser {
	return &FuiouUser{
		alipayUser: &AlipayUser{},
		wechatUser: &WechatUser{},
	}
}

func (a *FuiouUser) GetUserOpenID(ctx context.Context, authCode string, accountDetail *entity.PaymentAccount, app *model.Application, paymentMethod payment.Payment) (string, error) {
	var fuiouConf fuiouConfig.Config
	err := json.Unmarshal([]byte(accountDetail.ChannelConfig), &fuiouConf)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效配置")
	}

	switch paymentMethod {
	case payment.Payment_Wechat:
		return a.GetWechatOpenID(ctx, &fuiouConf, authCode, accountDetail, app)
	case payment.Payment_Alipay:
		return a.GetAliOpenID(ctx, &fuiouConf, authCode, accountDetail, app)
	}
	return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效的支付方式:"+paymentMethod.String())
}

func (a *FuiouUser) GetWechatOpenID(ctx context.Context, fuiouConf *fuiouConfig.Config, authCode string, accountDetail *entity.PaymentAccount, app *model.Application) (string, error) {
	conf := wechatConfig.Config{
		Merchant: wechatConfig.Merchant{
			AppID:     fuiouConf.Wechat.AppID,
			AppSecret: fuiouConf.Wechat.AppSecret,
		},
	}
	if conf.Merchant.AppID == "" || conf.Merchant.AppSecret == "" {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效的微信配置")
	}
	return GetWechatOpenID(ctx, a.wechatUser.Cache, conf, authCode, app)
}

func (a *FuiouUser) GetAliOpenID(ctx context.Context, fuiouConf *fuiouConfig.Config, authCode string, accountDetail *entity.PaymentAccount, app *model.Application) (string, error) {
	conf := alipayConfig.Config{
		Merchant: alipayConfig.Merchant{
			AppID: fuiouConf.Alipay.AppID,
		},
		Cert: alipayConfig.Cert{
			RsaPrivate: fuiouConf.Alipay.RsaPrivate,
			RsaRootCrt: fuiouConf.Alipay.RsaRootCrt,
		},
	}

	if conf.Merchant.AppID == "" || conf.Cert.RsaPrivate == "" || conf.Cert.RsaRootCrt == "" {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效的支付宝配置")
	}
	return GetAliUserOpenID(ctx, conf, authCode, app)
}

func (a *FuiouUser) RedirectUrl(ctx context.Context, callbackUrl string, accountDetail *entity.PaymentAccount, paymentMethod payment.Payment) (string, error) {
	var fuiouConf fuiouConfig.Config
	err := json.Unmarshal([]byte(accountDetail.ChannelConfig), &fuiouConf)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效配置")
	}

	switch paymentMethod {
	case payment.Payment_Wechat:
		return a.GetWechatRedirectUrl(ctx, &fuiouConf, callbackUrl)
	case payment.Payment_Alipay:
		return a.GetAliRedirectUrl(ctx, &fuiouConf, callbackUrl)
	}
	return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效的支付方式:"+paymentMethod.String())
}

func (a *FuiouUser) GetWechatRedirectUrl(ctx context.Context, fuiouConf *fuiouConfig.Config, callbackUrl string) (string, error) {
	conf := wechatConfig.Config{
		Merchant: wechatConfig.Merchant{
			AppID:     fuiouConf.Wechat.AppID,
			AppSecret: fuiouConf.Wechat.AppSecret,
		},
	}
	if conf.Merchant.AppID == "" || conf.Merchant.AppSecret == "" {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效的微信配置")
	}
	return GetWechatRedirectUrl(conf, callbackUrl), nil
}

func (a *FuiouUser) GetAliRedirectUrl(ctx context.Context, fuiouConf *fuiouConfig.Config, callbackUrl string) (string, error) {
	conf := alipayConfig.Config{
		Merchant: alipayConfig.Merchant{
			AppID: fuiouConf.Alipay.AppID,
		},
		Cert: alipayConfig.Cert{
			RsaPrivate: fuiouConf.Alipay.RsaPrivate,
			RsaRootCrt: fuiouConf.Alipay.RsaRootCrt,
		},
	}

	if conf.Merchant.AppID == "" || conf.Cert.RsaPrivate == "" || conf.Cert.RsaRootCrt == "" {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效的支付宝配置")
	}
	return GetAliRedirectUrl(ctx, conf, callbackUrl), nil
}
