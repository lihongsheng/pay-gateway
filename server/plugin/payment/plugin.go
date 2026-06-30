// Package payment 支付插件：支付网关核心功能
package payment

import (
	"github.com/lihongsheng/pay-gateway/enum"
	"github.com/lihongsheng/pay-gateway/model/system"
	"github.com/lihongsheng/pay-gateway/plugin"
	payApi "github.com/lihongsheng/pay-gateway/plugin/payment/api"
	payMiddleware "github.com/lihongsheng/pay-gateway/plugin/payment/middleware"
	payModel "github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type p struct{}

func (p) Name() string    { return "payment" }
func (p) Version() string { return "0.1.0" }

func (p) Models() []interface{} {
	return []interface{}{
		&payModel.Application{},
		&payModel.Merchant{},
		&payModel.PaymentAccount{},
		&payModel.PaymentMethod{},
		&payModel.PaymentOrder{},
		&payModel.PaymentOrderProduct{},
		&payModel.PaymentRequestLog{},
		&payModel.PaymentExpireRecord{},
		&payModel.RefundOrder{},
		&payModel.RefundOrderProduct{},
		&payModel.NotifyRecord{},
		&payModel.EventProcessRecord{},
		&payModel.TradeStatistic{},
		&payModel.Statistic{},
		&payModel.RouterAccountStatistic{},
	}
}

func (p) Menus() []system.SysMenu {
	return []system.SysMenu{
		// 平台端菜单
		{
			Type:       system.MenuTypeCatalog,
			Path:       "/plugin/payment",
			Name:       "PaymentManage",
			Component:  "Layout",
			Title:      "支付管理",
			Icon:       "money",
			Sort:       80,
			SystemType: enum.SystemTypePlatform,
			Children: []system.SysMenu{
				{
					Type:      system.MenuTypeMenu,
					Path:      "application",
					Name:      "PaymentApplication",
					Component: "plugin/payment/application/index",
					Title:     "应用管理",
					Icon:      "grid",
					Sort:      1,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/application","method":"GET"},{"path":"/api/v1/plugin/payment/application/search","method":"GET"},{"path":"/api/v1/plugin/payment/application","method":"POST"},{"path":"/api/v1/plugin/payment/application/status","method":"PUT"},{"path":"/api/v1/plugin/payment/application/qrcode","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增应用", Permission: "application:add"},
						{Type: system.MenuTypeButton, Name: "编辑应用", Permission: "application:edit"},
						{Type: system.MenuTypeButton, Name: "应用状态", Permission: "application:status"},
					},
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "payment-account",
					Name:      "PaymentAccount",
					Component: "plugin/payment/payment-account/index",
					Title:     "支付渠道",
					Icon:      "connection",
					Sort:      2,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/payment_account","method":"GET"},{"path":"/api/v1/plugin/payment/payment_account","method":"POST"},{"path":"/api/v1/plugin/payment/payment_account/list","method":"GET"},{"path":"/api/v1/plugin/payment/payment_account/channel/payment","method":"GET"},{"path":"/api/v1/plugin/payment/payment_account/channel/config","method":"GET"},{"path":"/api/v1/plugin/payment/payment_account/test/qrcode","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增渠道", Permission: "payment_account:add"},
						{Type: system.MenuTypeButton, Name: "编辑渠道", Permission: "payment_account:edit"},
					},
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "trade",
					Name:      "PaymentTrade",
					Component: "plugin/payment/trade/index",
					Title:     "订单管理",
					Icon:      "document",
					Sort:      3,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/trade","method":"GET"},{"path":"/api/v1/plugin/payment/trade/search","method":"GET"},{"path":"/api/v1/plugin/payment/refund","method":"POST"},{"path":"/api/v1/plugin/payment/refund/amount","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "退款", Permission: "trade:refund"},
					},
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "refund",
					Name:      "PaymentRefund",
					Component: "plugin/payment/refund/index",
					Title:     "退款管理",
					Icon:      "refresh-left",
					Sort:      4,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/refund","method":"GET"},{"path":"/api/v1/plugin/payment/refund/search","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "查看退款", Permission: "refund:view"},
					},
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "dashboard",
					Name:      "PaymentDashboard",
					Component: "plugin/payment/dashboard/index",
					Title:     "数据大盘",
					Icon:      "data-analysis",
					Sort:      5,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/dashboard/total_request","method":"GET"},{"path":"/api/v1/plugin/payment/dashboard/index","method":"GET"},{"path":"/api/v1/plugin/payment/dashboard/mch-index","method":"GET"},{"path":"/api/v1/plugin/payment/dashboard/mch-app-index","method":"GET"},{"path":"/api/v1/plugin/payment/dashboard/mch-app-account-all-index","method":"GET"},{"path":"/api/v1/plugin/payment/dashboard/search","method":"GET"}]`,
				},
			},
		},
		// 商户端菜单
		{
			Type:       system.MenuTypeCatalog,
			Path:       "/plugin/payment-mch",
			Name:       "PaymentManageMch",
			Component:  "Layout",
			Title:      "支付管理",
			Icon:       "money",
			Sort:       80,
			SystemType: enum.SystemTypeMch,
			Children: []system.SysMenu{
				{
					Type:      system.MenuTypeMenu,
					Path:      "application",
					Name:      "PaymentApplicationMch",
					Component: "plugin/payment/application/index-mch",
					Title:     "应用管理",
					Icon:      "grid",
					Sort:      1,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/application","method":"GET"},{"path":"/api/v1/plugin/payment/application/search","method":"GET"},{"path":"/api/v1/plugin/payment/application","method":"POST"}]`,
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "trade",
					Name:      "PaymentTradeMch",
					Component: "plugin/payment/trade/index-mch",
					Title:     "订单管理",
					Icon:      "document",
					Sort:      2,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/trade","method":"GET"},{"path":"/api/v1/plugin/payment/trade/search","method":"GET"}]`,
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "refund",
					Name:      "PaymentRefundMch",
					Component: "plugin/payment/refund/index-mch",
					Title:     "退款管理",
					Icon:      "refresh-left",
					Sort:      3,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/refund","method":"GET"},{"path":"/api/v1/plugin/payment/refund/search","method":"GET"}]`,
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "dashboard",
					Name:      "PaymentDashboardMch",
					Component: "plugin/payment/dashboard/index-mch",
					Title:     "数据大盘",
					Icon:      "data-analysis",
					Sort:      4,
					ApiRules:  `[{"path":"/api/v1/plugin/payment/dashboard/mch-index","method":"GET"},{"path":"/api/v1/plugin/payment/dashboard/mch-app-index","method":"GET"},{"path":"/api/v1/plugin/payment/dashboard/mch-app-account-all-index","method":"GET"}]`,
				},
			},
		},
	}
}

func (p) RegisterRoute(g *gin.Engine) {
	// 后台私有路由（需 JWT + Casbin 鉴权）
	// 注意：RegisterRoute 在 plugin 框架中已被放在 /api/v1/plugin/payment 前缀下
	// 但由于支付插件需要自定义路由结构，我们直接在 g 上注册
	// 公共路由（无需鉴权）
	public := g.Group("/public/v1")
	// 聚合支付
	public.GET("aggregate/redirect", payApi.AggregateRedirectUrl)
	public.GET("aggregate/user/auth", payApi.AggregateGetUserOpenID)
	public.GET("aggregate/application", payApi.AggregateGetApplication)
	public.POST("aggregate/payment", payApi.AggregatePayment)
	public.GET("aggregate/payment/query", payApi.AggregateQuery)
	public.POST("aggregate/test/payment", payApi.AggregateTestPayment)
	public.GET("aggregate/test/redirect", payApi.AggregateTestRedirectUrl)
	// 三方回调
	public.POST("notify/payment/:channel/:mchNo/:appNo/:orderNo", payApi.PublicPaymentCallback)
	public.POST("notify/test/payment/:channel/:mchNo/:appNo/:orderNo", payApi.PublicPaymentTestCallback)
	public.POST("notify/refund/:channel/:mchNo/:appNo/:tradeNo", payApi.PublicRefundCallback)
	// 收银台支付
	public.POST("cashier/payment", payApi.CashierPayment)
	public.GET("cashier/query", payApi.CashierQuery)
	public.GET("cashier/redirect", payApi.CashierRedirectUrl)
	// 需要验签的公共路由
	signPublic := public.Use(payMiddleware.SignMiddleware())
	signPublic.POST("cashier/create", payApi.CashierCreate)
	signPublic.POST("payment/pay", payApi.PublicPaymentPay)
	signPublic.GET("payment/query", payApi.PublicPaymentQuery)
	signPublic.POST("payment/close", payApi.PublicPaymentClose)
	signPublic.POST("refund", payApi.PublicRefundRefund)
	signPublic.GET("refund/query", payApi.PublicRefundQuery)
}

func (p) SeedTable(db *gorm.DB) error {
	return nil // 支付插件无需种子数据
}

func (p) Init() {

}

func init() { plugin.Register(p{}) }
