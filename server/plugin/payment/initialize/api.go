package initialize

import (
	model "github.com/lihongsheng/pay-gateway/model/system"
	"github.com/lihongsheng/pay-gateway/plugin/plugin-tool/utils"
)

func Api() {
	entities := []model.SysApi{
		{
			Path:        "/private/v1/merchant/status",
			Description: "改变商户状态",
			ApiGroup:    "商户管理",
			Method:      "POST",
		},
		{
			Path:        "/private/v1/merchant",
			Description: "保存商户",
			ApiGroup:    "商户管理",
			Method:      "POST",
		},
		{
			Path:        "/private/v1/merchant/search",
			Description: "商户搜索",
			ApiGroup:    "商户管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/merchant",
			Description: "商户详情",
			ApiGroup:    "商户管理",
			Method:      "GET",
		},
		// application
		{
			Path:        "/private/v1/application/status",
			Description: "改变应用状态",
			ApiGroup:    "应用管理",
			Method:      "POST",
		},
		{
			Path:        "/private/v1/application",
			Description: "保存应用",
			ApiGroup:    "应用管理",
			Method:      "POST",
		},
		{
			Path:        "/private/v1/application/search",
			Description: "应用搜索",
			ApiGroup:    "应用管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/application",
			Description: "应用详情",
			ApiGroup:    "应用管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/application/qrcode",
			Description: "应用二维码",
			ApiGroup:    "应用管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/payment_account/test/qrcode",
			Description: "支付渠道测试二维码",
			ApiGroup:    "支付渠道管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/payment_account",
			Description: "支付渠道配置",
			ApiGroup:    "支付渠道管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/payment_account",
			Description: "支付渠道保存",
			ApiGroup:    "支付渠道管理",
			Method:      "POST",
		},
		{
			Path:        "/private/v1/payment_account/list",
			Description: "应用支付渠道列表",
			ApiGroup:    "支付渠道管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/payment_account/channel/payment",
			Description: "支付渠道支持的产品",
			ApiGroup:    "支付渠道管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/payment_account/channel/config",
			Description: "支付渠道配置详情",
			ApiGroup:    "支付渠道管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/trade",
			Description: "支付单详情",
			ApiGroup:    "订单管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/trade/search",
			Description: "支付搜索",
			ApiGroup:    "订单管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/refund",
			Description: "退款详情",
			ApiGroup:    "退款管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/refund/search",
			Description: "退款列表查询",
			ApiGroup:    "退款管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/refund",
			Description: "退款",
			ApiGroup:    "退款管理",
			Method:      "POST",
		},
		{
			Path:        "/private/v1/refund/amount",
			Description: "可退款金额",
			ApiGroup:    "退款管理",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/dashboard/index",
			Description: "首页数据统计",
			ApiGroup:    "数据大盘",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/dashboard/mch-index",
			Description: "商户首页数据统计",
			ApiGroup:    "数据大盘",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/dashboard/mch-index",
			Description: "商户首页数据统计",
			ApiGroup:    "数据大盘",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/dashboard/mch-app-index",
			Description: "应用数据统计",
			ApiGroup:    "数据大盘",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/dashboard/mch-app-account-all-index",
			Description: "应用支付账户数据统计",
			ApiGroup:    "数据大盘",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/dashboard/search",
			Description: "搜索数据统计",
			ApiGroup:    "数据大盘",
			Method:      "GET",
		},
		{
			Path:        "/private/v1/dashboard/total_request",
			Description: "总请求订单数",
			ApiGroup:    "数据大盘",
			Method:      "GET",
		},
	}
	utils.RegisterApis(entities...)
}
