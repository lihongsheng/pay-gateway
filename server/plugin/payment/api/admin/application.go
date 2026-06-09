package admin

import (
	"fmt"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
)

type ApplicationCreateRequest struct {
	ID      int64  `json:"id"`
	AppName string `json:"app_name" binding:"required,min=2"` // 应用名字
	//	Secret           string         `json:"secret" binding:"required,len=32"`             // 应用秘密
	MchNo            string         `json:"mch_no"  binding:"omitempty,min=2"`       // 公司编号
	Desc             string         `json:"desc"`                                    // 应用描述
	PaymentTitle     string         `json:"payment_title"`                           // 支付的tiltle
	PayIcon          string         `json:"pay_icon"`                                // 支付的icon
	Status           enum.MchStatus `json:"status" binding:"required"`               // 2 停用 1 正常
	IsCustomerDomain int64          `json:"is_customer_domain"`                      // 是否自定义域名方式 ,0 非 1自定义，当页面选择1时下边是必须输入
	CustomerDomain   string         `json:"customer_domain" binding:"omitempty,url"` // 自定义域名
	ProxyHost        string         `json:"proxy_host"`
	ProxyPort        int64          `json:"proxy_port"`    // 端口
	ProxyUser        string         `json:"proxy_user"`    // 用户名
	ProxyPwd         string         `json:"proxy_pwd"`     // 密码
	MultiChannel     int64          `json:"multi_channel"` // 是否支持单一channel配置多个
}

func (a ApplicationCreateRequest) Validate() error {
	if a.IsCustomerDomain == 1 {
		if a.CustomerDomain == "" {
			return fmt.Errorf("代理模式下域名必须设置")
		}
		if a.ProxyHost == "" {
			return fmt.Errorf("代理模式下proxy_host 必须设置")
		}
		if a.ProxyPort == 0 {
			return fmt.Errorf("代理模式下proxy_port 必须设置")
		}
	}
	if a.MchNo == "" {
		return fmt.Errorf("商户编号不能为空")
	}
	if a.AppName == "" {
		return fmt.Errorf("应用名称不能为空")
	}
	return nil
}

type ApplicationQueryRequest struct {
	AppName  string         `json:"app_name" form:"app_name"`
	Status   enum.MchStatus `json:"status" form:"status"`                           // 2 停用 1 正常
	MchNo    string         `json:"mch_no" form:"mch_no" binding:"omitempty,min=2"` // 公司编号
	AppNo    string         `json:"app_no" form:"app_no" binding:"omitempty,min=2"`
	AppNos   []string       `json:"app_nos" form:"app_nos"`
	PageSize int            `json:"page_size" form:"page_size"` // 限制数量
	Page     int            `json:"page" form:"page"`           // 页大小
}

type ApplicationStatusRequest struct {
	ID     int64          `json:"id" form:"id" binding:"required,gt=0"`
	Status enum.MchStatus `json:"status" form:"status" binding:"required,gt=0"` // 0 停用 1 正常
}
