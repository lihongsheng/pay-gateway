package dto

import (
	"errors"
	"github.com/lihongsheng/payment-sdk/driver/iface"
	"time"

	"github.com/lihongsheng/payment-sdk/enum/payment"
)

// ---------- 支付账户 DTO ----------

type PaymentAccountCreateRequest struct {
	ID            int64                `json:"id"`
	Name          string               `json:"name"  binding:"required"`
	Remark        string               `json:"remark"`
	AppNo         string               `json:"app_no"  binding:"required"`
	MchNo         string               `json:"mch_no"`
	Channel       string               `json:"channel"  binding:"required"`
	Status        int64                `json:"status"  binding:"required"`
	Extend        string               `json:"extend"`
	PaymentMethod []PaymentMethod      `json:"payment_method"  binding:"required"`
	ChannelOption *iface.ChannelOption `json:"channel_option"`
	ChannelConfig string               `json:"channel_config"  binding:"required"`
	MaxLimit      int64                `json:"max_limit"`
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
	OrderNo   string         `json:"order_no" form:"order_no"`
	TradeNo   string         `json:"trade_no" form:"trade_no"`
	AppNo     string         `json:"app_no" form:"app_no"`
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

// ---------- 支付账户响应 DTO ----------

type PaymentAccountDetail struct {
	ID                  int64                 `json:"id"`
	Name                string                `json:"name"`
	Remark              string                `json:"remark"`
	AppNo               string                `json:"app_no"`
	MchNo               string                `json:"mch_no"`
	AccountNo           string                `json:"account_no"`
	Channel             string                `json:"channel"`
	ChannelName         string                `json:"channel_name"`
	Status              int                   `json:"status"`
	Extend              string                `json:"extend"`
	ChannelConfig       string                `json:"channel_config"`
	PaymentMethodConfig []PaymentMethodConfig `json:"payment_method_config"`
	ChannelOption       *iface.ChannelOption  `json:"channel_option"`
	CreatedAt           time.Time             `json:"created_at"`
	ValidateStatus      int64                 `json:"validate_status"`
	MaxLimit            int64                 `json:"max_limit"`
}

type PaymentMethodConfig struct {
	Method  string                 `json:"method"`
	Label   string                 `json:"label"`
	Product []PaymentProductConfig `json:"product"`
}

type PaymentProductConfig struct {
	Product string `json:"product"`
	Label   string `json:"label"`
	Used    bool   `json:"used"`
}
