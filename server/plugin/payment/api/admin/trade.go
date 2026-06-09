package admin

import (
	"errors"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"github.com/lihongsheng/payment-sdk/enum/refund"
	"time"
)

type TradeSearchRequest struct {
	MchNo     string         `json:"mch_no" form:"mch_no"`         // 商户号
	AppNo     string         `json:"app_no" form:"app_no"`         // 应用编号
	TradeNo   string         `json:"trade_no" form:"trade_no"`     // 交易号
	OrderNo   string         `json:"order_no" form:"order_no"`     // 订单号
	StartTime time.Time      `json:"start_time" form:"start_time"` // 开始时间
	EndTime   time.Time      `json:"end_time" form:"end_time"`     // 结束时间
	Status    payment.Status `json:"status" form:"status"`         // 状态
	Page      int            `json:"page" form:"page"`             // 页码
	PageSize  int            `json:"page_size" form:"page_size"`   // 每页大小
}

func (t *TradeSearchRequest) Validate() error {
	if t.StartTime.IsZero() {
		return errors.New("开始时间不能为空")
	}
	if t.EndTime.IsZero() {
		return errors.New("结束时间不能为空")
	}
	// 开始时间和结束时不能大于两个月
	if t.EndTime.Sub(t.StartTime) > 60*24*time.Hour*2 {
		return errors.New("时间间隔不能超过两个月")
	}
	if t.Page <= 0 {
		t.Page = 1
	}
	if t.PageSize <= 0 {
		t.PageSize = 10
	}
	if t.MchNo == "" {
		return errors.New("商户号不能为空")
	}
	return nil
}

type TradeSearchResponse struct {
	MchNo          string         `json:"mch_no,omitempty"`          // 商户号
	MchName        string         `json:"mch_name,omitempty"`        // 商户名称
	AppNo          string         `json:"app_no,omitempty"`          // 应用编号
	AppName        string         `json:"app_name,omitempty"`        // 应用名称
	TradeNo        string         `json:"trade_no,omitempty"`        // 交易号
	OrderNo        string         `json:"order_no,omitempty"`        // 订单号
	AccountNo      string         `json:"account_no,omitempty"`      // 支付账户编号
	AccountName    string         `json:"account_name,omitempty"`    // 支付账户名称
	PaymentMethod  string         `json:"payment_method,omitempty"`  // 支付方式
	PaymentProduct string         `json:"payment_product,omitempty"` // 支付产品
	Amount         int64          `json:"amount,omitempty"`          // 金额
	Status         payment.Status `json:"status,omitempty"`          // 状态
	CreatedAt      time.Time      `json:"created_at"`                // 创建时间
	OrderTime      time.Time      `json:"order_time"`                // 下单时间
	Subject        string         `json:"subject,omitempty"`         // 商品标题
}

type TradeOrderDetail struct {
	MchNo          string         `json:"mch_no,omitempty"`
	MchName        string         `json:"mch_name,omitempty"`
	AppNo          string         `json:"app_no,omitempty"`
	AppName        string         `json:"app_name,omitempty"`
	TradeNo        string         `json:"trade_no,omitempty"`
	OrderNo        string         `json:"order_no,omitempty"`
	AccountNo      string         `json:"account_no,omitempty"`
	AccountName    string         `json:"account_name,omitempty"`
	PaymentMethod  string         `json:"payment_method,omitempty"`
	PaymentProduct string         `json:"payment_product,omitempty"`
	OutMchTradeNo  string         `json:"out_mch_trade_no,omitempty"`
	Amount         int64          `json:"amount,omitempty"`
	Status         payment.Status `json:"status,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	OrderTime      time.Time      `json:"order_time"`
	RedirectURL    string         `json:"redirect_url,omitempty"`
	NotifyUrl      string         `json:"notify_url,omitempty"`
	ExtParams      string         `json:"ext_params,omitempty"`
	ThirdMsg       string         `json:"third_msg,omitempty"`
	ThirdCode      string         `json:"third_code,omitempty"`
	Subject        string         `json:"subject,omitempty"`
	Desc           string         `json:"desc,omitempty"`
	NotifyStatus   int            `json:"notify_status"`
	UserOpenId     string         `json:"user_open_id,omitempty"`
}

type RefundSearchRequest struct {
	MchNo     string         `json:"mch_no" form:"mch_no"`
	AppNo     string         `json:"app_no" form:"app_no"`
	TradeNo   string         `json:"trade_no" form:"trade_no"`
	OrderNo   string         `json:"order_no" form:"order_no"`
	RefundNo  string         `json:"refund_no" form:"refund_no"`
	StartTime time.Time      `json:"start_time" form:"start_time"`
	EndTime   time.Time      `json:"end_time" form:"end_time"`
	Status    payment.Status `json:"status" form:"status"`
	Page      int            `json:"page" form:"page"`
	PageSize  int            `json:"page_size" form:"page_size"`
}

func (t *RefundSearchRequest) Validate() error {
	if t.StartTime.IsZero() {
		return errors.New("开始时间不能为空")
	}
	if t.EndTime.IsZero() {
		return errors.New("结束时间不能为空")
	}
	// 开始时间和结束时不能大于两个月
	if t.EndTime.Sub(t.StartTime) > 60*24*time.Hour*2 {
		return errors.New("时间间隔不能超过两个月")
	}
	if t.Page <= 0 {
		t.Page = 1
	}
	if t.PageSize <= 0 {
		t.PageSize = 10
	}
	if t.MchNo == "" {
		return errors.New("商户号不能为空")
	}
	return nil
}

type RefundOrderDetail struct {
	MchNo         string        `json:"mch_no,omitempty"`
	MchName       string        `json:"mch_name,omitempty"`
	AppNo         string        `json:"app_no,omitempty"`
	AppName       string        `json:"app_name,omitempty"`
	RefundNo      string        `json:"refund_no,omitempty"`
	RefundTradeNo string        `json:"refund_trade_no,omitempty"`
	TradeNo       string        `json:"trade_no,omitempty"`
	OrderNo       string        `json:"order_no,omitempty"`
	AccountNo     string        `json:"account_no,omitempty"`
	AccountName   string        `json:"account_name,omitempty"`
	OutMchTradeNo string        `json:"out_mch_trade_no,omitempty"`
	RefundAmount  int64         `json:"refund_amount,omitempty"`
	PaymentAmount int64         `json:"payment_amount,omitempty"`
	Status        refund.Status `json:"status,omitempty"`
	NotifyUrl     string        `json:"notify_url,omitempty"`
	ExtParams     string        `json:"ext_params,omitempty"`
	Reason        string        `json:"reason,omitempty"`
	NotifyStatus  int           `json:"notify_status"`
	CreatedAt     time.Time     `json:"created_at"`
	SuccessTime   time.Time     `json:"success_time"`
	ThirdMsg      string        `json:"third_msg"`
	RefundFrom    string        `json:"refund_from"`
}

type Refund struct {
	Amount  int64  `json:"amount" form:"amount"`
	Reason  string `json:"reason" form:"reason"`
	MchNo   string `json:"mch_no" form:"mch_no"`
	AppNo   string `json:"app_no" form:"app_no"`
	OrderNo string `json:"order_no" form:"order_no"`
}

func (t *Refund) Validate() error {
	if t.Amount <= 0 {
		return errors.New("金额不能小于0")
	}
	if t.Reason == "" {
		return errors.New("原因不能为空")
	}
	if t.OrderNo == "" {
		return errors.New("订单号不能为空")
	}
	if t.AppNo == "" {
		return errors.New("应用编号不能为空")
	}
	if t.MchNo == "" {
		return errors.New("商户编号不能为空")
	}
	return nil
}

type TradeStatistic struct {
	AppName           string `json:"app_name,omitempty"`     // 应用名称
	MchName           string `json:"mch_name,omitempty"`     // 商户名称
	AccountName       string `json:"account_name,omitempty"` // 支付账号名称
	AppNo             string `json:"app_no"`                 // 应用编号
	AccountNo         string `json:"account_no"`             // 支付账号编号
	MchNo             string `json:"mch_no"`                 // 商户编号
	StatisticDate     string `json:"statistic_date"`         // 统计日期
	TotalOrder        int64  `json:"total_order"`            // 发起支付的订单数
	SuccessOrder      int64  `json:"success_order"`          // 成功支付订单数
	TotalAmount       int64  `json:"total_amount"`           // 总金额 1元 = 100分
	SuccessAmount     int64  `json:"success_amount"`         // 成功金额
	RefundOrder       int64  `json:"refund_order"`           // 退款次数
	RefundAmount      int64  `json:"refund_amount"`          // 退款金额
	TotalRequestOrder int64  `json:"total_request_order"`    // 总请求订单数
}

type MchTradeStatisticRequest struct {
	MchNo     string    `json:"mch_no" form:"mch_no"`           // 商户编号
	AppNo     string    `json:"app_no"  form:"app_no"`          // 应用编号
	AccountNo string    `json:"account_no"   form:"account_no"` // 支付账号编号
	StartTime time.Time `json:"start_time" form:"start_time"`   // 开始时间
	EndTime   time.Time `json:"end_time" form:"end_time"`       // 结束时间
	Page      int       `json:"page" form:"page"`
	PageSize  int       `json:"page_size" form:"page_size"`
}
