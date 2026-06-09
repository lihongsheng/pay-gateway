package admin

import (
	"errors"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"time"
)

type PaymentAccountCreateRequest struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"  binding:"required"`    // 名字标识
	Remark        string          `json:"remark"`                      // 描述
	AppNo         string          `json:"app_no"  binding:"required"`  // 应用编号
	MchNo         string          `json:"mch_no"`                      // 商户编号
	Channel       string          `json:"channel"  binding:"required"` // 支付渠道：Wechat 微信 | Alipay 支付宝 | Lakala 拉卡拉 | Fuiou 富有
	Status        int64           `json:"status"  binding:"required"`  // 1 qi 2 下线
	Extend        string          `json:"extend"`
	PaymentMethod []PaymentMethod `json:"payment_method"  binding:"required"`
	ChannelConfig string          `json:"channel_config"  binding:"required"`
	MaxLimit      int64           `json:"max_limit"` // 最高请求次数
}

func (c PaymentAccountCreateRequest) Validate() error {
	if c.Name == "" {
		return errors.New("名称不能为空")
	}
	if c.AppNo == "" {
		return errors.New("应用编号不能为空")
	}
	if c.Channel == "" {
		return errors.New("支付渠道不能为空")
	}
	if c.ChannelConfig == "" {
		return errors.New("支付渠道配置不能为空")
	}
	if len(c.PaymentMethod) == 0 {
		return errors.New("支付产品不能为空")
	}
	return nil
}

type PaymentMethod struct {
	Method  string `json:"method"`
	Product string `json:"product"`
	Extend  string `json:"extend"`
}

type PaymentAccountQueryRequest struct {
	Name    string `json:"name" form:"name"`
	AppNo   string `json:"app_no" form:"app_no"`
	MchNo   string `json:"mch_no" form:"mch_no"`
	Channel string `json:"channel" form:"channel"`
	Status  int64  `json:"status" form:"status"`
}

type SearchRequest struct {
	// 订单号
	OrderNo string `json:"order_no" form:"order_no"`
	// 支付单号
	TradeNo string `json:"trade_no" form:"trade_no"`
	// 应用标识
	AppNo string `json:"app_no" form:"app_no"`
	// 商户标识
	MchNo     string         `json:"mch_no" form:"mch_no"`
	StartTime time.Time      `json:"start_time" form:"start_time"`
	EndTime   time.Time      `json:"end_time" form:"end_time"`
	Status    payment.Status `form:"status"`
	Page      int            `json:"page" form:"page"`
	PageSize  int            `json:"page_size" form:"page_size"`
}

func (q SearchRequest) Validate() error {
	if q.MchNo == "" {
		return errors.New("商户不能为空")
	}
	if q.StartTime.IsZero() {
		return errors.New("开始时间不能为空")
	}
	if q.EndTime.IsZero() {
		return errors.New("结束时间不能为空")
	}
	return nil
}
