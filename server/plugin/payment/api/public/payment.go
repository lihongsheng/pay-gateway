package public

import (
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/payment-sdk/driver/dto"
	"github.com/lihongsheng/payment-sdk/enum"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"github.com/lihongsheng/payment-sdk/enum/refund"
	"time"
)

type Order struct {
	OrderNo string `json:"order_no"`
	// 实付金额
	Amount Amount `json:"amount"`
	// 订单商品
	Goods []Goods `json:"goods"`
	// 订单名称
	Subject string `json:"subject"`
	// 订单描述
	Desc string `json:"desc"`
	// 订单创建时间
	CreateAt time.Time `json:"create_at"`
}

type Amount struct {
	// 分单位，一元为100
	Total    int64  `json:"total"`
	Currency string `json:"currency"`
}

type Goods struct {
	// 商品名称
	Name string `json:"name"`
	// 商品SKU
	Sku string `json:"sku"`
	// 商品价格，单位分
	Price int64 `json:"price"`
	// 商品数量
	Quantity int `json:"quantity"`
	//
	Desc string `json:"desc"`
}

type PaymentOrder struct {
	Order Order `json:"order"`
	Payer Payer `json:"payer"`
	// 支付跳转地址，完成后前端页面跳转地址
	RedirectUrl string `json:"redirect_url"`
	// 订单超时时间，时间挫
	TimeExpire int64 `json:"time_expire"`
	// 支付回调地址
	NotifyUrl string `json:"notify_url"`
	// 透传参数 如果请求时传递了该参数，异步通知时将该参数原样返回。
	PassBackParams string              `json:"pass_back_params"`
	AlipayExtra    *AlipayPaymentExtra `json:"alipay_extra"`
	// 支付方式 Wxpay | Alipay
	PaymentMethod payment.Payment `json:"payment_method"`
	// 支付产品 JSAPI
	PaymentProduct payment.PaymentProduct `json:"payment_product"`
	// 支付场景信息
	SceneInfo *SceneInfo `json:"scene_info"`
	AccountNo string     `json:"account_no"`
	AppNo     string     `json:"app_no"`
	MchNo     string     `json:"-"`
	RequestID string     `json:"request_id"`
}

func (p *PaymentOrder) Validate() error {
	if p.PaymentMethod == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 payment_method")
	}
	if p.PaymentProduct == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 payment_product")
	}
	if p.Order.OrderNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 order_no")
	}
	if p.Order.Subject == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 订单标题")
	}
	if p.Order.Amount.Total <= 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 amount")
	}
	if p.Payer.OpenID == "" && p.Payer.UnionID == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 open_id")
	}
	if p.AppNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 app_no")
	}
	return nil
}

type PaymentOrderResponse struct {
	// 订单号
	OrderNo string `json:"order_no"`
	// 支付单号
	TradeNo string `json:"trade_no"`
	// 三方支付单号,如微信或者支付宝的交易号
	OutTradeNo string `json:"out_trade_no"`
	// 订单支付金额
	Amount      Amount         `json:"amount"`
	Status      payment.Status `json:"status"`
	Action      *dto.Action    `json:"action"`
	RedirectURL string         `json:"redirect_url"`
}

type QueryPaymentRequest struct {
	// 订单号
	OrderNo string `json:"order_no" form:"order_no"`
	TradeNo string `json:"trade_no" form:"trade_no"`
	AppNo   string `json:"app_no" form:"app_no"`
	MchNo   string `json:"-" form:"mch_no"`
}

func (q *QueryPaymentRequest) Validate() error {
	if q.OrderNo == "" && q.TradeNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 order_no 或者 trade_no")
	}
	if q.AppNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 app_no")
	}

	return nil
}

type Payer struct {
	// 用户在商户的标识
	UnionID string `json:"union_id"`
	// 用户在商户下某个应用的标识；微信，支付宝，富有都需要
	OpenID string `json:"open_id"`
	// 应用标识;如微信或者支付宝的应用ID；富友，拉卡拉需要透传 app_id 和 open_id
	AppID string `json:"app_id"`
}

type AlipayPaymentExtra struct {
	// 产品码。
	//商家和支付宝签约的产品码。
	//当面付场景下，如果签约的是当面付快捷版，则传 OFFLINE_PAYMENT;
	//【示例值】FACE_TO_FACE_PAYMENT
	// 使用JSAPI 接口，但是不适用JSAPI的产品码，透传
	ProductCode string `json:"product_code,omitempty"`
}

type SceneInfo struct {
	// 客户端IP
	ClientIp string `json:"client_ip"`
	// 设备ID
	DeviceID string `json:"device_id"`
	// 设备
	Device enum.Device `json:"device"`
	// 系统
	System enum.System `json:"system"`
	// 应用信息
	ApplicationInfo *ApplicationInfo `json:"application_info"`
}

type ApplicationInfo struct {
	AppName    string `json:"app_name"`
	Url        string `json:"url"`
	AppPackage string `json:"app_package"`
}

type PaymentOrderDetail struct {
	Order Order `json:"order"`
	// 支付跳转地址，完成后前端页面跳转地址
	RedirectUrl string `json:"redirect_url"`
	// 支付回调地址
	NotifyUrl string `json:"notify_url"`
	// 透传参数 如果请求时传递了该参数，异步通知时将该参数原样返回。
	PassBackParams string         `json:"pass_back_params"`
	Status         payment.Status `json:"status"`
	// 支付单号
	TradeNo string `json:"trade_no"`
	// 三方支付单号
	OutTradeNo string    `json:"out_trade_no"`
	Create     time.Time `json:"create"`
}

type CashierCreateOrder struct {
	Order Order `json:"order"`
	// 支付跳转地址，完成后前端页面跳转地址
	RedirectUrl string `json:"redirect_url"`
	// 订单超时时间，时间挫
	TimeExpire int64 `json:"time_expire"`
	// 支付回调地址
	NotifyUrl string `json:"notify_url"`
	// 透传参数 如果请求时传递了该参数，异步通知时将该参数原样返回。
	PassBackParams string `json:"pass_back_params"`
	AppNo          string `json:"app_no"`
	MchNo          string `json:"-"`
	RequestID      string `json:"request_id"`
}

func (p *CashierCreateOrder) Validate() error {
	if p.Order.OrderNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 order_no")
	}
	if p.Order.Subject == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 订单标题")
	}
	if p.Order.Amount.Total <= 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 amount")
	}
	if p.AppNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 app_no")
	}
	//if p.RedirectUrl == "" {
	//	return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 redirect_url")
	//}
	return nil
}

type CashierCreateOrderResponse struct {
	OrderNo         string         `json:"order_no"`          // 订单号
	OpenURL         string         `json:"open_url"`          // 跳转地址
	Amount          Amount         `json:"amount"`            // 支付金额
	TradeNo         string         `json:"trade_no"`          // 平台交易号
	OrderExpireTime time.Time      `json:"order_expire_time"` // 订单过期时间
	Status          payment.Status `json:"status"`
}

type CashierPaymentOrder struct {
	OrderNo     string              `json:"order_no"`
	OpenId      string              `json:"open_id"`
	AlipayExtra *AlipayPaymentExtra `json:"alipay_extra"`
	// 支付方式 Wxpay | Alipay
	PaymentMethod payment.Payment `json:"payment_method"`
	// 支付产品 JSAPI
	PaymentProduct payment.PaymentProduct `json:"payment_product"`
	// 支付场景信息
	SceneInfo *SceneInfo `json:"scene_info"`
	AccountNo string     `json:"account_no"`
	RequestID string     `json:"request_id"`
	Token     string     `json:"token"`
}

func (p *CashierPaymentOrder) Validate() error {
	if p.PaymentMethod == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 payment_method")
	}
	if p.PaymentProduct == 0 {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 payment_product")
	}
	if p.OrderNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 order_no")
	}
	if p.OpenId == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 open_id")
	}
	if p.AccountNo == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 account_no")
	}
	if p.Token == "" {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "请传入 token")
	}
	return nil
}

type CashierQueryOrder struct {
	OrderNo string `json:"order_no" form:"order_no"`
	Token   string `json:"token" form:"token"`
}

type NotifyPaymentRequest struct {
	AppNo      string         `json:"app_no" form:"app_no"`
	MchNo      string         `json:"mch_no" form:"mch_no"`
	OrderNo    string         `json:"order_no" form:"order_no"`
	TradeNo    string         `json:"trade_no" form:"trade_no"`
	OutTradeNo string         `json:"out_trade_no" form:"out_trade_no"`
	Status     payment.Status `json:"status" form:"status"`
}

type NotifyRefundRequest struct {
	AppNo         string        `json:"app_no" form:"app_no"`
	MchNo         string        `json:"mch_no" form:"mch_no"`
	OrderNo       string        `json:"order_no" form:"order_no"`
	RefundNo      string        `json:"refund_no" form:"refund_no"`
	RefundTradeNo string        `json:"refund_trade_no" form:"refund_trade_no"`
	OutTradeNo    string        `json:"out_trade_no" form:"out_trade_no"`
	Status        refund.Status `json:"status" form:"status"`
}
