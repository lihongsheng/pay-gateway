package adapter

import (
	"context"
	"encoding/json"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	alipayConfig "github.com/lihongsheng/payment-sdk/adapter/alipay/config"
	alipayEnum "github.com/lihongsheng/payment-sdk/adapter/alipay/enum"
	alipayModel "github.com/lihongsheng/payment-sdk/adapter/alipay/model"
	alipayUser "github.com/lihongsheng/payment-sdk/adapter/alipay/user"
	"github.com/lihongsheng/payment-sdk/config/proxy"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"net/url"
)

type AlipayUser struct {
	//redis redis.UniversalClient
}

func NewAlipayUser() AggregateUser {
	return &AlipayUser{}
}

func (a *AlipayUser) GetUserOpenID(ctx context.Context, authCode string, accountDetail *entity.PaymentAccount, app *model.Application, paymentMethod payment.Payment) (string, error) {
	var conf alipayConfig.Config
	err := json.Unmarshal([]byte(accountDetail.ChannelConfig), &conf)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效配置")
	}

	return GetAliUserOpenID(ctx, conf, authCode, app)
}

func GetAliUserOpenID(ctx context.Context, conf alipayConfig.Config, authCode string, app *model.Application) (string, error) {
	var p *proxy.Proxy
	if app.ProxyHost != "" {
		p = &proxy.Proxy{
			Host:     app.ProxyHost,
			Port:     int(app.ProxyPort),
			UserName: app.ProxyUser,
			Password: app.ProxyPwd,
		}
	}
	cl, err := alipayUser.NewUser(conf, p)
	if err != nil {
		return "", errors2.WrapError(errors2.ErrCodeInvalidParam, "无效配置", err)
	}
	r, err := cl.AuthToken(ctx, &alipayModel.UserAuthRequest{
		GrantType: alipayEnum.AUTH_TYPE_AUTHORIZATION_CODE,
		Code:      authCode,
	})
	if err != nil {
		return "", errors2.WrapError(errors2.ErrCodeInvalidParam, "无效配置", err)
	}
	return r.OpenId, nil
}

func GetAliRedirectUrl(ctx context.Context, conf alipayConfig.Config, callbackUrl string) string {
	oauthUrl, _ := url.Parse("https://openauth.alipay.com/oauth2/publicAppAuthorize.htm")
	query := url.Values{}
	query.Add("app_id", conf.AppID)
	query.Add("scope", "auth_base")
	query.Add("redirect_uri", callbackUrl)
	oauthUrl.RawQuery = query.Encode()
	return oauthUrl.String()
}

func (a *AlipayUser) RedirectUrl(ctx context.Context, callbackUrl string, accountDetail *entity.PaymentAccount, paymentMethod payment.Payment) (string, error) {
	var conf alipayConfig.Config
	err := json.Unmarshal([]byte(accountDetail.ChannelConfig), &conf)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效配置")
	}

	return GetAliRedirectUrl(ctx, conf, callbackUrl), nil
}
