package public

import (
	"errors"
	"github.com/lihongsheng/payment-sdk/enum"
	"github.com/lihongsheng/payment-sdk/enum/payment"
)

// TestQrcode 后台生成测试订单
type TestQrcode struct {
	AppNo          string                 `json:"app_no" form:"app_no"`
	OrderNo        string                 `json:"order_no" form:"order_no"`
	Amount         int64                  `json:"amount" form:"amount"`
	OrderTitle     string                 `json:"order_title" form:"order_title"`
	PaymentMethod  payment.Payment        `json:"payment_method" form:"payment_method"`
	PaymentProduct payment.PaymentProduct `json:"payment_product" form:"payment_product"`
	AccountNo      string                 `json:"account_no" form:"account_no"`
}

// QueryCheckoutPaymentMethod 查询支持的支付方式
type QueryCheckoutPaymentMethod struct {
	PaymentMethod  payment.Payment        `json:"payment_method"`
	PaymentProduct payment.PaymentProduct `json:"payment_product"`
	// 设备
	Device enum.Device `json:"device"`
	// 系统
	System enum.System `json:"system"`
	AppNo  string      `json:"app_no"`
}

func (q QueryCheckoutPaymentMethod) Validate() error {
	if q.AppNo == "" {
		return errors.New("app_no不能为空")
	}
	if (q.PaymentMethod == 0 || q.PaymentProduct == 0) && q.Device == "" {
		return errors.New("请传入 device 或者 payment_method 和 payment_product")
	}
	return nil
}

type QueryCheckoutPaymentMethodResponse struct {
	PaymentMethod payment.Payment `json:"payment_method"`
	AccountNo     string          `json:"account_no"`
	Icon          string          `json:"icon"`
}
