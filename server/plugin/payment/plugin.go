// Package payment 支付插件：支付网关核心功能
package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/enum"
	"github.com/lihongsheng/pay-gateway/model/system"
	"github.com/lihongsheng/pay-gateway/plugin"
	payApi "github.com/lihongsheng/pay-gateway/plugin/payment/api"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	payMiddleware "github.com/lihongsheng/pay-gateway/plugin/payment/middleware"
	payModel "github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
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
					Component: "plugin/payment/view/application/index",
					Title:     "应用管理",
					Icon:      "grid",
					Sort:      1,
					ApiRules:  `[{"path":"/api/plugin/payment/application","method":"GET"},{"path":"/api/plugin/payment/application/search","method":"GET"},{"path":"/api/plugin/payment/application/qrcode","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增应用", Permission: "application:add", ApiRules: `[{"path":"/api/plugin/payment/application","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑应用", Permission: "application:edit", ApiRules: `[{"path":"/api/plugin/payment/application","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "应用状态", Permission: "application:status", ApiRules: `[{"path":"/api/plugin/payment/application/status","method":"POST"}]`},
					},
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "payment-account",
					Name:      "PaymentAccount",
					Component: "plugin/payment/view/payment-account/index",
					Title:     "支付渠道",
					Icon:      "connection",
					Sort:      2,
					ApiRules:  `[{"path":"/api/plugin/payment/payment_account","method":"GET"},{"path":"/api/plugin/payment/payment_account/list","method":"GET"},{"path":"/api/plugin/payment/payment_account/channel/payment","method":"GET"},{"path":"/api/plugin/payment/payment_account/channel/config","method":"GET"},{"path":"/api/plugin/payment/payment_account/test/qrcode","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增渠道", Permission: "payment_account:add", ApiRules: `[{"path":"/api/plugin/payment/payment_account","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑渠道", Permission: "payment_account:edit", ApiRules: `[{"path":"/api/plugin/payment/payment_account","method":"POST"}]`},
					},
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "trade",
					Name:      "PaymentTrade",
					Component: "plugin/payment/view/payment/index",
					Title:     "订单管理",
					Icon:      "document",
					Sort:      3,
					ApiRules:  `[{"path":"/api/plugin/payment/trade","method":"GET"},{"path":"/api/plugin/payment/trade/search","method":"GET"},{"path":"/api/plugin/payment/refund/amount","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "退款", Permission: "trade:refund", ApiRules: `[{"path":"/api/plugin/payment/refund","method":"POST"}]`},
					},
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "refund",
					Name:      "PaymentRefund",
					Component: "plugin/payment/view/refund/index",
					Title:     "退款管理",
					Icon:      "refresh-left",
					Sort:      4,
					ApiRules:  `[{"path":"/api/plugin/payment/refund","method":"GET"},{"path":"/api/plugin/payment/refund/search","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "查看退款", Permission: "refund:view"},
					},
				},
				{
					Type:      system.MenuTypeMenu,
					Path:      "dashboard",
					Name:      "PaymentDashboard",
					Component: "plugin/payment/view/dashboard/index",
					Title:     "数据大盘",
					Icon:      "data-analysis",
					Sort:      5,
					ApiRules:  `[{"path":"/api/plugin/payment/dashboard/total_request","method":"GET"},{"path":"/api/plugin/payment/dashboard/index","method":"GET"},{"path":"/api/plugin/payment/dashboard/mch-index","method":"GET"},{"path":"/api/plugin/payment/dashboard/mch-app-index","method":"GET"},{"path":"/api/plugin/payment/dashboard/mch-app-account-all-index","method":"GET"},{"path":"/api/plugin/payment/dashboard/search","method":"GET"}]`,
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
					Type:       system.MenuTypeMenu,
					Path:       "application",
					Name:       "PaymentApplicationMch",
					Component:  "plugin/payment/view/application/index-mch",
					Title:      "应用管理",
					Icon:       "grid",
					Sort:       1,
					SystemType: enum.SystemTypeMch,
					ApiRules:   `[{"path":"/api/plugin/payment/application","method":"GET"},{"path":"/api/plugin/payment/application/search","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增应用", Permission: "application:add", SystemType: enum.SystemTypeMch, ApiRules: `[{"path":"/api/plugin/payment/application","method":"POST"}]`},
					},
				},
				{
					Type:       system.MenuTypeMenu,
					Path:       "trade",
					Name:       "PaymentTradeMch",
					Component:  "plugin/payment/view/payment/index-mch",
					Title:      "订单管理",
					Icon:       "document",
					Sort:       2,
					SystemType: enum.SystemTypeMch,
					ApiRules:   `[{"path":"/api/plugin/payment/trade","method":"GET"},{"path":"/api/plugin/payment/trade/search","method":"GET"}]`,
				},
				{
					Type:       system.MenuTypeMenu,
					Path:       "refund",
					Name:       "PaymentRefundMch",
					Component:  "plugin/payment/view/refund/index-mch",
					Title:      "退款管理",
					Icon:       "refresh-left",
					Sort:       3,
					SystemType: enum.SystemTypeMch,
					ApiRules:   `[{"path":"/api/plugin/payment/refund","method":"GET"},{"path":"/api/plugin/payment/refund/search","method":"GET"}]`,
				},
				{
					Type:       system.MenuTypeMenu,
					Path:       "dashboard",
					Name:       "PaymentDashboardMch",
					Component:  "plugin/payment/view/dashboard/index-mch",
					Title:      "数据大盘",
					Icon:       "data-analysis",
					Sort:       4,
					SystemType: enum.SystemTypeMch,
					ApiRules:   `[{"path":"/api/plugin/payment/dashboard/mch-index","method":"GET"},{"path":"/api/plugin/payment/dashboard/mch-app-index","method":"GET"},{"path":"/api/plugin/payment/dashboard/mch-app-account-all-index","method":"GET"}]`,
				},
			},
		},
	}
}

// RegisterRoute g gin.Engine , privatePlugin 插件路由 /api/plugin/<name>
func (p) RegisterRoute(g *gin.Engine, privatePlugin *gin.RouterGroup) {
	// 后台相关
	{
		private := privatePlugin.Group("payment")
		private.GET("application", payApi.GetAppInfo)
		private.GET("application/search", payApi.SearchApplication)
		private.POST("application", payApi.SaveApplication)
		private.POST("application/status", payApi.ChangeAppStatus)
		private.GET("application/qrcode", payApi.GetQrCode)

		private.GET("payment_account", payApi.GetPaymentAccount)
		private.POST("payment_account", payApi.SavePaymentAccount)
		private.GET("payment_account/list", payApi.AccountList)
		private.GET("payment_account/channel/payment", payApi.GetChannelPaymentMethod)
		private.GET("payment_account/channel/config", payApi.GetChannelConfig)
		private.GET("payment_account/test/qrcode", payApi.GetTestQrCode)

		// 后台订单相关
		private.GET("trade", payApi.GetTradeOrder)
		private.GET("trade/search", payApi.SearchTradeOrder)
		private.POST("refund", payApi.TradeRefund)
		private.GET("refund", payApi.TradeRefundDetail)
		private.GET("refund/search", payApi.TradeRefundSearch)
		private.GET("refund/amount", payApi.GetAvailableRefundAmount)
		// 统计相关
		private.GET("dashboard/total_request", payApi.TotalRequest)
		private.GET("dashboard/index", payApi.DashboardIndex)
		private.GET("dashboard/mch-index", payApi.MchIndex)
		private.GET("dashboard/mch-app-index", payApi.MchAppIndex)
		private.GET("dashboard/mch-app-account-all-index", payApi.MchAppAccountAllIndex)
		private.GET("dashboard/search", payApi.SearchIndex)
	}

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
	public.POST("notify/payment/:channel/:mchNo/:appNo/:orderNo", payApi.PaymentCallback)
	public.POST("notify/test/payment/:channel/:mchNo/:appNo/:orderNo", payApi.PaymentTestCallback)
	public.POST("notify/refund/:channel/:mchNo/:appNo/:tradeNo", payApi.RefundCallback)
	// 收银台支付
	public.POST("cashier/payment", payApi.CashierPayment)
	public.GET("cashier/query", payApi.CashierQuery)
	public.GET("cashier/redirect", payApi.CashierRedirectUrl)
	// 需要验签的公共路由
	signPublic := public.Use(payMiddleware.SignMiddleware())
	signPublic.POST("cashier/create", payApi.CashierCreate)
	signPublic.POST("payment/pay", payApi.Pay)
	signPublic.GET("payment/query", payApi.QueryPayment)
	signPublic.POST("payment/close", payApi.ClosePayment)
	signPublic.POST("refund", payApi.Refund)
	signPublic.GET("refund/query", payApi.QueryRefund)
}

// InitServices 初始化支付插件服务层
func (p) InitServices(ctx plugin.InitContext) error {
	svcApp := svc.NewServiceContext(ctx.DB, ctx.Redis, nil, ctx.Config.Plugin.PaymentConfig)
	domainService := domain.NewServiceGroup(svcApp, ctx.Redis, ctx.Config)
	enter.Init(svcApp, ctx.DB, domainService, ctx.Redis, ctx.Config)
	return nil
}

func (p) SeedTable(db *gorm.DB) error {
	return nil // 支付插件无需种子数据
}

func init() { plugin.Register(p{}) }
