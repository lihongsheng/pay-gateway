package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	wechatConfig "github.com/lihongsheng/payment-sdk/adapter/wxpay/config"
	"github.com/lihongsheng/payment-sdk/config/proxy"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"github.com/silenceper/wechat/v2"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	wechatUtil "github.com/silenceper/wechat/v2/util"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AggregateUser interface {
	RedirectUrl(ctx context.Context, callbackUrl string, accountDetail *entity.PaymentAccount, paymentMethod payment.Payment) (string, error)
	GetUserOpenID(ctx context.Context, authCode string, accountDetail *entity.PaymentAccount, app *model.Application, paymentMethod payment.Payment) (string, error)
}

type WechatUser struct {
	Cache     *CacheWechat
	Transport *Transport
}
type proxyCtx struct{}

type CacheWechat struct {
}
type Transport struct{}

var DefaultTransport = &Transport{}

func init() {
	wechatUtil.DefaultHTTPClient.Transport = DefaultTransport
}
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	proxyCtx := req.Context().Value(proxyCtx{})
	if p, ok := proxyCtx.(http.RoundTripper); ok {
		return p.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}
func NewWechatUser() AggregateUser {
	return &WechatUser{
		Cache: &CacheWechat{},
	}
}

func (a *CacheWechat) Get(key string) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return global.Redis.Get(ctx, key).String()
}

func (a *CacheWechat) Set(key string, val interface{}, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return global.Redis.Set(ctx, key, val, timeout).Err()
}

func (a *CacheWechat) IsExist(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return global.Redis.Exists(ctx, key).Val() > 0
}

func (a *CacheWechat) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return global.Redis.Del(ctx, key).Err()
}

func (a *WechatUser) GetUserOpenID(ctx context.Context, authCode string, accountDetail *entity.PaymentAccount, app *model.Application, paymentMethod payment.Payment) (string, error) {
	// 测试后删除
	if (global.Cfg.Payment.IsTest() || global.Cfg.Payment.IsTest()) && accountDetail.AccountNo == "P6b5df34ad2800" {
		return "", errors2.NewError(errors2.ErrUserLimit, "限制登录")
	}
	var conf wechatConfig.Config
	err := json.Unmarshal([]byte(accountDetail.ChannelConfig), &conf)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效配置")
	}
	return GetWechatOpenID(ctx, a.Cache, conf, authCode, app)
}

func GetWechatOpenID(ctx context.Context, cache *CacheWechat, conf wechatConfig.Config, authCode string, app *model.Application) (string, error) {
	token, aes := utils.GenerateTokenAndEncodingKey(conf.Merchant.AppID)
	c := &offConfig.Config{
		AppID:          conf.Merchant.AppID,
		AppSecret:      conf.Merchant.AppSecret,
		Token:          token,
		EncodingAESKey: aes,
		Cache:          cache,
		UseStableAK:    true,
	}
	wc := wechat.NewWechat().GetOfficialAccount(c)
	if app.ProxyHost != "" {
		proxyInfo := proxy.Proxy{
			Host:     app.ProxyHost,
			Port:     int(app.ProxyPort),
			UserName: app.ProxyUser,
			Password: app.ProxyPwd,
		}
		ctx = buildProxyURL(ctx, proxyInfo)
	}
	r, err := wc.GetOauth().GetUserAccessTokenContext(ctx, authCode)
	if err != nil {
		if strings.Contains(err.Error(), "user limited") {
			return "", errors2.NewError(errors2.ErrUserLimit, "限制登录")
		}
		return "", errors2.WrapError(errors2.ErrCodeInvalidRequest, "无效获取openid", err)
	}
	return r.OpenID, nil
}

func (a *WechatUser) RedirectUrl(ctx context.Context, callbackUrl string, accountDetail *entity.PaymentAccount, paymentMethod payment.Payment) (string, error) {
	var conf wechatConfig.Config
	err := json.Unmarshal([]byte(accountDetail.ChannelConfig), &conf)
	if err != nil {
		return "", errors2.NewError(errors2.ErrCodeInvalidParam, "无效配置")
	}
	return GetWechatRedirectUrl(conf, callbackUrl), nil
}

func GetWechatRedirectUrl(conf wechatConfig.Config, callbackUrl string) string {
	oauthUrl, _ := url.Parse("https://open.weixin.qq.com/connect/oauth2/authorize")
	query := url.Values{}
	query.Add("appid", conf.Merchant.AppID)
	query.Add("scope", "snsapi_base")
	query.Add("redirect_uri", callbackUrl)
	query.Add("response_type", "code")
	oauthUrl.RawQuery = query.Encode()
	return oauthUrl.String() + "#wechat_redirect"
}

func buildProxyURL(ctx context.Context, proxy proxy.Proxy) context.Context {
	r := &http.Transport{
		Proxy: func(req *http.Request) (u *url.URL, err error) {
			u, err = url.Parse(fmt.Sprintf("http://%s:%d", proxy.Host, proxy.Port))
			if err != nil {
				return nil, err
			}
			if proxy.UserName != "" && proxy.Password != "" {
				u.User = url.UserPassword(proxy.UserName, proxy.Password)
			}
			if proxy.UserName != "" {
				u.User = url.User(proxy.UserName)
			}
			return u, nil
		},
	}
	return context.WithValue(ctx, proxyCtx{}, r)
}
