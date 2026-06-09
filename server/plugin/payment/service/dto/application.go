package dto

import "github.com/lihongsheng/pay-gateway/plugin/payment/enum"

type ApplicationCacheInfo struct {
	ID               int64                    `json:"id"`
	AppName          string                   `json:"app_name,omitempty"`           // 应用名字
	Secret           string                   `json:"secret,omitempty"`             // 应用秘密
	AppNo            string                   `json:"app_no,omitempty"`             // 应用编号
	MchNo            string                   `json:"mch_no,omitempty"`             // 商户编号
	PaymentTitle     string                   `json:"payment_title,omitempty"`      // 支付的tiltle
	PayIcon          string                   `json:"pay_icon,omitempty"`           // 支付的icon
	Status           enum.MchStatus           `json:"status,omitempty"`             // 2 停用 1 正常
	IsCustomerDomain int64                    `json:"is_customer_domain,omitempty"` // 是否自定义域名方式
	CustomerDomain   string                   `json:"customer_domain,omitempty"`    // 自定义域名
	ProxyHost        string                   `json:"proxy_host,omitempty"`
	ProxyPort        int64                    `json:"proxy_port,omitempty"` // 端口
	ProxyUser        string                   `json:"proxy_user,omitempty"` // 用户名
	ProxyPwd         string                   `json:"proxy_pwd,omitempty"`  // 密码
	MchInfo          *ApplicationMchCacheInfo `json:"mch_info,omitempty"`
}

type ApplicationMchCacheInfo struct {
	ID      int64          `json:"id"`
	MchName string         `json:"mch_name,omitempty"` // 公司名字
	MchNo   string         `json:"mch_no,omitempty"`   // 编号
	Status  enum.MchStatus `json:"status,omitempty"`   // 2 停用 1 正常
}
