package admin

import (
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/payment-sdk/driver/iface"
	"time"
)

type PaymentAccountDetail struct {
	ID                  int64                 `json:"id"`
	Name                string                `json:"name"  binding:"required"`    // 名字标识
	Remark              string                `json:"remark"`                      // 描述
	AppNo               string                `json:"app_no"  binding:"required"`  // 应用编号
	MchNo               string                `json:"mch_no"  binding:"required"`  // 商户编号
	AccountNo           string                `json:"account_no"`                  // 账户编号
	Channel             string                `json:"channel"  binding:"required"` // 支付渠道：Wechat 微信 | Alipay 支付宝 | Lakala 拉卡拉 | Fuiou 富有
	ChannelName         string                `json:"channel_name"`
	Status              enum.MchStatus        `json:"status"  binding:"required"`         // 1 启用 2 停用
	Extend              string                `json:"extend"`                             // 扩展配置
	ChannelConfig       string                `json:"channel_config"  binding:"required"` // 渠道配置
	PaymentMethodConfig []PaymentMethodConfig `json:"payment_method_config"  binding:"required"`
	ChannelOption       *iface.ChannelOption  `json:"channel_option"`
	CreatedAt           time.Time             `json:"created_at"`
	ValidateStatus      int64                 `json:"validate_status"`
	MaxLimit            int64                 `json:"max_limit"` // 日最大支付次数
}

type PaymentMethodConfig struct {
	// 支付方式
	Method string `json:"method"`
	// 支付方式描述
	Label string `json:"label"`
	// 支持的支付产品
	Product []PaymentProductConfig `json:"product"`
}

type PaymentProductConfig struct {
	// 产品标识
	Product string `json:"product"`
	// 产品描述
	Label string `json:"label"`
	// 是否启用
	Used bool `json:"used"`
}
