package entity

import (
	"fmt"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"time"
)

type PaymentOrder struct {
	ID                   int64                  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键ID" json:"id"`                                                   // 主键ID
	PaymentAccountNo     string                 `gorm:"column:payment_account_no;type:varchar(19);not null;comment:payment_account_no" json:"payment_account_no"`                     // payment_account_no
	AppNo                string                 `gorm:"column:app_no;type:varchar(100);not null;comment:应用编号" json:"app_no"`                                                          // 应用编号
	MchNo                string                 `gorm:"column:mch_no;type:varchar(100);primaryKey;comment:商户编号" json:"mch_no"`                                                        // 商户编号
	OrderSubject         string                 `gorm:"column:order_subject;type:varchar(100);not null;comment:订单标题" json:"order_subject"`                                            // 订单标题
	PaymentMethod        payment.Payment        `gorm:"column:payment_method;type:varchar(60);not null;comment:支付方式：Wechat 微信 | Alipay 支付宝 | 云闪付" json:"payment_method"`              // 支付方式：Wechat 微信 | Alipay 支付宝 | 云闪付
	PaymentProduct       payment.PaymentProduct `gorm:"column:payment_product;type:varchar(60);not null;comment:支付产品：H5 | JSAPI | LITE | APP | Qrcode | Card" json:"payment_product"` // 支付产品：H5 | JSAPI | LITE | APP | Qrcode | Card
	OrderNo              string                 `gorm:"column:order_no;type:varchar(64);not null;comment:订单号" json:"order_no"`                                                        // 订单号
	OutMchTradeNo        string                 `gorm:"column:out_mch_trade_no;type:varchar(100);not null;comment:渠道订单号：微信|支付宝" json:"out_mch_trade_no"`                              // 渠道订单号：微信|支付宝
	TradeNo              string                 `gorm:"column:trade_no;type:varchar(100);not null;comment:支付系统交易号" json:"trade_no"`                                                   // 支付系统交易号
	Device               string                 `gorm:"column:device;type:varchar(20);not null;comment:设备" json:"device"`                                                             // 设备
	System               string                 `gorm:"column:system;type:varchar(20);not null;comment:操作系统" json:"system"`                                                           // 操作系统
	OrderAmount          int64                  `gorm:"column:order_amount;type:bigint;not null;comment:订单金额分单位" json:"order_amount"`                                                 // 订单金额分单位
	PaymentAmount        int64                  `gorm:"column:payment_amount;type:bigint;not null;comment:订单实付金额分单位" json:"payment_amount"`                                           // 订单实付金额分单位
	UserPaymentAmount    int64                  `gorm:"column:user_payment_amount;type:bigint;not null;comment:用户实付金额分单位" json:"user_payment_amount"`                                 // 用户实付金额分单位
	Currency             string                 `gorm:"column:currency;type:varchar(10);not null;comment:币种" json:"currency"`                                                         // 币种
	Status               payment.Status         `gorm:"column:status;type:tinyint;not null;comment:支付状态" json:"status"`                                                               // 支付状态
	StatusFrom           string                 `gorm:"column:status_from;type:varchar(100);not null;comment:状态修改来源API|WEBHOOK||ADMIN" json:"status_from"`                            // 状态修改来源API|WEBHOOK||ADMIN
	TimeExpire           int64                  `gorm:"column:time_expire;type:bigint;not null;comment:支付过期时间" json:"time_expire"`                                                    // 支付过期时间
	OrderCreate          int64                  `gorm:"column:order_create;type:bigint;not null;comment:订单创建时间，时间戳" json:"order_create"`                                              // 订单创建时间，时间戳
	ExtParam             string                 `gorm:"column:ext_param;type:varchar(255);not null;comment:支付扩展参数" json:"ext_param"`                                                  // 支付扩展参数
	RedirectURL          string                 `gorm:"column:redirect_url;type:varchar(200);not null;comment:支付完成跳转url" json:"redirect_url"`                                         // 支付完成跳转url
	NotifyURL            string                 `gorm:"column:notify_url;type:varchar(200);not null;comment:异步通知url" json:"notify_url"`                                               // 异步通知url
	OrderDesc            string                 `gorm:"column:order_desc;type:varchar(255);not null;comment:订单描述" json:"order_desc"`                                                  // 订单描述
	ThirdCode            string                 `gorm:"column:third_code;type:varchar(100);not null;comment:三方接口原返回码code" json:"third_code"`                                          // 三方接口原返回码code
	ThirdMsg             string                 `gorm:"column:third_msg;type:varchar(100);not null;comment:三方接口返回描述" json:"third_msg"`                                                // 三方接口返回描述
	UserOpenid           string                 `gorm:"column:user_openid;type:varchar(100);not null;comment:三方用户标识" json:"user_openid"`                                              // 三方用户标识
	NotifyStatus         int64                  `gorm:"column:notify_status;type:tinyint;not null;comment:1 已经回调" json:"notify_status"`                                               // 1 已经回调
	CreatedAt            time.Time              `gorm:"column:created_at;type:datetime;primaryKey;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                          // 创建时间
	PaymentOrderProducts []Product              `json:"products"`
	Retry                int64                  `gorm:"column:retry;type:int;not null;comment:重试次数" json:"retry"` // 重试次数
	event                *event.PaymentOrderStatusEvent
}

type Product struct {
	ID          int64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	ProductName string `gorm:"column:product_name;type:varchar(100);not null;comment:订单标题" json:"product_name"` // 订单标题
	AppNo       string `gorm:"column:app_no;type:varchar(100);not null;comment:应用编号" json:"app_no"`             // 应用编号
	MchNo       string `gorm:"column:mch_no;type:varchar(100);primaryKey;comment:商户编号" json:"mch_no"`           // 商户编号
	Quantity    int64  `gorm:"column:quantity;type:int;not null;comment:购买数量" json:"quantity"`                  // 购买数量
	Price       int64  `gorm:"column:price;type:int;not null;comment:单价" json:"price"`                          // 单价
	Sku         string `gorm:"column:sku;type:varchar(100);not null;comment:商品sku" json:"sku"`                  // 商品sku
	ProductDesc string `gorm:"column:product_desc;type:varchar(255);not null;comment:商品描述" json:"product_desc"` // 商品描述
	URL         string `gorm:"column:url;type:varchar(255);not null;comment:商品链接" json:"url"`                   // 商品链接
	OrderNo     string `gorm:"column:order_no;type:varchar(64);not null;comment:订单号" json:"order_no"`           // 订单号
	TradeNo     string `gorm:"column:trade_no;type:varchar(100);not null;comment:支付单号" json:"trade_no"`         // 支付单号
}

//func (p *PaymentOrder) Reset() {
//	p.Status = payment.Status_Created
//	if p.TimeExpire < time.Now().Unix() {
//		p.TimeExpire = time.Now().Add(enum.DefaultExpireTime).Unix()
//	}
//}

func (p *PaymentOrder) IsTimeout() bool {
	// 核心规则：待支付状态 + 当前时间超过超时时间
	if p.IsConfirmed() {
		return false
	}
	return time.Now().Unix() > p.TimeExpire
}

func (p *PaymentOrder) Validate() error {
	if p.TradeNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "请传入 trade_no")
	}
	if p.AppNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "请传入 app_no")
	}
	if p.MchNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "请传入 mch_no")
	}
	if p.PaymentMethod == payment.Payment_Payment_UNKNOWN {
		return errors.NewError(errors.ErrCodeInvalidParam, "请传入 payment_method")
	}
	if p.PaymentProduct == payment.PaymentProduct_PaymentMethod_UNKNOWN {
		return errors.NewError(errors.ErrCodeInvalidParam, "请传入 payment_product")
	}
	if p.OrderNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "请传入 order_no")
	}
	if p.OrderSubject == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "请传入 order_subject")
	}
	if p.PaymentAmount <= 0 {
		return errors.NewError(errors.ErrCodeInvalidParam, "请传入 payment_amount")
	}
	return nil
}

func (p *PaymentOrder) SetTradeNo() {
	tradeNo, t := utils.GenTradeNo(p.MchNo)
	p.TradeNo = tradeNo
	p.CreatedAt = utils.GetGenIDTimestamp(t)
}

func (p *PaymentOrder) IsCreated() bool {
	return p.Status == payment.Status_Created
}

func (p *PaymentOrder) IsPending() bool {
	return p.Status == payment.Status_Pending || p.Status == payment.Status_WaitConfirm || p.Status == payment.Status_WaitPay || p.Status == payment.Status_TempFailed
}

func (p *PaymentOrder) IsConfirmed() bool {
	return !(p.Status == payment.Status_Created || p.Status == payment.Status_Pending || p.Status == payment.Status_WaitConfirm || p.Status == payment.Status_WaitPay || p.Status == payment.Status_TempFailed)
}

func (p *PaymentOrder) IsRetry() bool {
	if p.Status == payment.Status_Success || p.Status == payment.Status_Refund {
		return false
	}
	return true
}

func (p *PaymentOrder) SetStatus(status payment.Status) {
	old := p.Status
	p.Status = status
	//if !p.IsCreated() && !p.IsPending() {
	p.event = &event.PaymentOrderStatusEvent{
		Base: event.Base{
			EventType: event.EventTypePaymentOrderStatus,
			EventId:   fmt.Sprintf("%s_%d", p.TradeNo, p.Status),
		},
		MchNo:     p.MchNo,
		OrderNo:   p.OrderNo,
		TradeNo:   p.TradeNo,
		AppNo:     p.AppNo,
		OldStatus: old,
		NewStatus: p.Status,
		Amount:    p.PaymentAmount,
		AccountNo: p.PaymentAccountNo,
		CreateAt:  p.CreatedAt,
		ID:        p.ID,
	}
	//}
}

func (p *PaymentOrder) SetErrStatus(status payment.Status, errCode int) {
	old := p.Status
	p.Status = status
	//if !p.IsCreated() && !p.IsPending() {
	p.event = &event.PaymentOrderStatusEvent{
		Base: event.Base{
			EventType: event.EventTypePaymentOrderStatus,
			EventId:   fmt.Sprintf("%s_%d", p.TradeNo, p.Status),
		},
		MchNo:     p.MchNo,
		OrderNo:   p.OrderNo,
		TradeNo:   p.TradeNo,
		AppNo:     p.AppNo,
		OldStatus: old,
		NewStatus: p.Status,
		Amount:    p.PaymentAmount,
		ErrCode:   errCode,
		AccountNo: p.PaymentAccountNo,
		CreateAt:  p.CreatedAt,
		Retry:     p.Retry,
		ID:        p.ID,
	}
	//}
}

func (p *PaymentOrder) GetEvents() *event.PaymentOrderStatusEvent {
	if p.event != nil {
		return p.event
	}
	if p.ID == 0 && p.Status == payment.Status_Created {
		return &event.PaymentOrderStatusEvent{
			Base: event.Base{
				EventType: event.EventTypePaymentOrderStatus,
				EventId:   fmt.Sprintf("%s_%d", p.TradeNo, p.Status),
			},
			MchNo:     p.MchNo,
			OrderNo:   p.OrderNo,
			TradeNo:   p.TradeNo,
			AppNo:     p.AppNo,
			OldStatus: 0,
			NewStatus: payment.Status_Created,
			Amount:    p.PaymentAmount,
			AccountNo: p.PaymentAccountNo,
			CreateAt:  p.CreatedAt,
			ID:        p.ID,
		}
	}
	return nil
}
