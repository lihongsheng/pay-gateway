package public

import (
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/payment-sdk/enum"
	"github.com/lihongsheng/payment-sdk/enum/payment"
)

type AggregateRedirectUrlRequest struct {
	Token         string          `form:"token" json:"token" binding:"required"`
	PaymentMethod payment.Payment `form:"payment_method" json:"payment_method"`
	// 支付产品 JSAPI
	PaymentProduct  payment.PaymentProduct `form:"payment_product" json:"payment_product"`
	Device          enum.Device            `form:"device" json:"device"`
	OrderNo         string                 `json:"order_no" form:"order_no"`
	FilterAccountNo string                 `json:"filter_account_no" form:"filter_account_no"`
	AccountNo       string                 `json:"account_no" form:"account_no"`
	Retry           bool                   `json:"retry" form:"retry"`
}

func (a AggregateRedirectUrlRequest) Validate() error {
	if a.PaymentMethod == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "未识别到支付方式")
	}
	if a.PaymentProduct == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "未识别到支付产品")
	}
	if a.Token == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "无效token")
	}
	return nil
}

type AggregateUserOpenIDRequest struct {
	Token         string          `form:"token" json:"token" binding:"required"`
	AuthCode      string          `form:"auth_code" json:"auth_code"`
	AccountNo     string          `form:"account_no" json:"account_no"`
	PaymentMethod payment.Payment `form:"payment_method" json:"payment_method"`
	// 支付产品 JSAPI
	PaymentProduct payment.PaymentProduct `form:"payment_product" json:"payment_product"`
	Device         enum.Device            `form:"device" json:"device"`
}

func (a AggregateUserOpenIDRequest) Validate() error {
	if a.PaymentMethod == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "未识别到支付方式")
	}
	if a.Token == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "无效token")
	}
	if a.AccountNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 account_no")
	}
	return nil
}

type AggregateApplication struct {
	AppNo        string `json:"app_no"`
	PayTitle     string `json:"pay_title"`
	PayIcon      string `json:"pay_icon"`
	MultiChannel bool   `json:"multi_channel"`
}

type AggregateOrder struct {
	Token   string `json:"token"`
	OrderNo string `json:"order_no"`
	Amount  int64  `json:"amount"`
	// 设备
	Device enum.Device `json:"device"`
	// 支付方式 Wxpay | Alipay
	PaymentMethod payment.Payment `json:"payment_method"`
	// 支付产品 JSAPI
	PaymentProduct payment.PaymentProduct `json:"payment_product"`
	OpenID         string                 `json:"open_id"`
	AccountNo      string                 `json:"account_no"`
	RequestID      string                 `json:"request_id"`
}

func (a *AggregateOrder) Validate() error {
	if a.Token == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "无效token")
	}
	if a.PaymentMethod == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "未识别到支付方式")
	}
	if a.PaymentProduct == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "未识别到支付产品")
	}
	if a.OpenID == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "未获取到用户")
	}
	if a.AccountNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "未获取到支付账户")
	}
	return nil
}
