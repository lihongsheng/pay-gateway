package wxpay

import (
	"context"
	"errors"
	"fmt"

	"github.com/lihongsheng/pay-gateway/config"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type Api struct {
	c      config.Config
	client *core.Client
}

func InitClient(c config.Config) (*Api, error) {
	w := &Api{c: c}
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(c.CertPath)
	if err != nil {
		return nil, errors.New("wxpay load merchant private key error")
	}
	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(c.MchID, c.CertificateSerialNumber, mchPrivateKey, c.Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("new wechat pay client err:%s", err.Error()))
	}
	w.client = client
	return w, nil
}
