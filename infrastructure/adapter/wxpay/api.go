package wxpay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/lihongsheng/pay-gateway/infrastructure/config"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var apies = sync.Map{}

type Api struct {
	C      config.Config
	Client *core.Client
}

func InitClient(c config.Config) (*Api, error) {
	//privateKey := ""
	//if c.CertPrivateKey != "" {
	//	privateKey = tools.Md5(c.CertPrivateKey)
	//}
	//key := fmt.Sprintf("%s%s%s%s%s%s", c.PaymentName, c.MchID, c.CertificateSerialNumber, c.APIKey, c.AppID, privateKey)
	//if v, ok := apies.Load(key); ok {
	//	return v.(*Api), nil
	//}
	w := &Api{C: c}
	// 使用 utils 提供的函数从私钥字符串中加载商户私钥
	mchPrivateKey, err := utils.LoadPrivateKey(c.CertPrivateKey)
	publicKeyStr := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArpfqIBmK2gTiu25eJ4Rh
9FxyowIbZuhbO0YfzdAek0T76a1IGI/ziOJZBz6HSOnwBOYnbdpz3zd5Qjes9ak8
dmHXmVK/D951rf0FWVI13RSckINEna0m0VmjuWDryGU1FyiNWwWPxM/g2dKWvvkQ
nxC8BuuMwl+ov2W6zDC51YO+DQ9ziaukAqEvRr7jbj4owvYX1yLeOVGDlYA3HiNR
a7wuIbFzPpStfl3jtWZQ2WSmPoIdkKtHNSl0lDa0anPm6/Euxw4OT0mXLsqN4st4
PtIPIZamBZ4jdZqzMfig/3omzY1iSWfd9OiylUWGi62z9pbSG8iTJ1sOqHQgACZF
LwIDAQAB
-----END PUBLIC KEY-----`
	publicKey, err := utils.LoadPublicKey(publicKeyStr)
	if err != nil {
		return nil, errors.New("wxpay load merchant private key errors")
	}
	ctx := context.Background()
	opts := []core.ClientOption{}
	if c.Proxy.Host != "" {
		opts = append(opts, proxy(c))
	}
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	//opts = append(opts, option.WithWechatPayAutoAuthCipher(c.MchID, c.CertificateSerialNumber, mchPrivateKey, c.APIKey))
	opts = append(opts, option.WithWechatPayPublicKeyAuthCipher(c.MchID, c.CertificateSerialNumber, mchPrivateKey, "sdfsdfs", publicKey))
	//opts = append(opts, option.WithMerchantCredential(c.MchID, c.CertificateSerialNumber, mchPrivateKey))
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("new wechat pay client err:%s", err.Error()))
	}
	w.Client = client
	//apies.Store(key, w)
	return w, nil
}

type WithProxyOption struct {
	C config.Config
}

func (w *WithProxyOption) Apply(settings *core.DialSettings) error {
	settings.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (u *url.URL, err error) {
				u, err = url.Parse(fmt.Sprintf("http://%s:%d", w.C.Proxy.Host, w.C.Proxy.Port))
				if err != nil {
					return nil, err
				}
				if w.C.Proxy.UserName != "" && w.C.Proxy.Password != "" {
					u.User = url.UserPassword(w.C.Proxy.UserName, w.C.Proxy.Password)
				}
				if w.C.Proxy.UserName != "" {
					u.User = url.User(w.C.Proxy.UserName)
				}
				return u, nil
			},
		},
	}
	return nil
}
func proxy(c config.Config) core.ClientOption {
	return &WithProxyOption{C: c}
}
