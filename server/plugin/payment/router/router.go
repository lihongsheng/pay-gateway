package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api"
	middle "github.com/lihongsheng/pay-gateway/plugin/payment/middleware"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
)

func InitRoute(private *gin.RouterGroup, public *gin.RouterGroup) {
	private.GET("merchant", api.ApiGroupApp.MchApi.Get)
	private.GET("merchant/search", api.ApiGroupApp.MchApi.Search)
	private.POST("merchant", api.ApiGroupApp.MchApi.Save)
	private.POST("merchant/status", api.ApiGroupApp.MchApi.ChangeStatus)

	private.GET("application", api.ApiGroupApp.AppApi.GetAppInfo)
	private.GET("application/search", api.ApiGroupApp.AppApi.Search)
	private.POST("application", api.ApiGroupApp.AppApi.Save)
	private.POST("application/status", api.ApiGroupApp.AppApi.ChangeStatus)
	private.GET("application/qrcode", api.ApiGroupApp.AppApi.GetQrCode)

	private.GET("payment_account", api.ApiGroupApp.AccountApi.Get)
	private.POST("payment_account", api.ApiGroupApp.AccountApi.Save)
	private.GET("payment_account/list", api.ApiGroupApp.AccountApi.AccountList)
	private.GET("payment_account/channel/payment", api.ApiGroupApp.AccountApi.GetChannelPaymentMethod)
	private.GET("payment_account/channel/config", api.ApiGroupApp.AccountApi.GetChannelConfig)
	private.GET("payment_account/test/qrcode", api.ApiGroupApp.AccountApi.GetTestQrCode)

	// 后台订单相关
	private.GET("trade", api.ApiGroupApp.TradeApi.GetTradeOrderApi)
	private.GET("trade/search", api.ApiGroupApp.TradeApi.SearchTradeOrderApi)
	private.POST("refund", api.ApiGroupApp.TradeApi.TradeRefundApi)
	private.GET("refund", api.ApiGroupApp.TradeApi.TradeRefundDetailApi)
	private.GET("refund/search", api.ApiGroupApp.TradeApi.TradeRefundSearchApi)
	private.GET("refund/amount", api.ApiGroupApp.TradeApi.GetAvailableRefundAmount)
	// 统计相关
	private.GET("dashboard/total_request", api.ApiGroupApp.DashboardApi.TotalRequest)
	private.GET("dashboard/index", api.ApiGroupApp.DashboardApi.Index)
	private.GET("dashboard/mch-index", api.ApiGroupApp.DashboardApi.MchIndex)
	private.GET("dashboard/mch-app-index", api.ApiGroupApp.DashboardApi.MchAppIndex)
	private.GET("dashboard/mch-app-account-all-index", api.ApiGroupApp.DashboardApi.MchAppAccountAllIndex)
	private.GET("dashboard/search", api.ApiGroupApp.DashboardApi.SearchIndex)

	// 聚合支付，无需header传签名,基于加密token 验证
	public.GET("v1/aggregate/redirect", api.ApiGroupApp.AggregateApi.RedirectUrl)
	public.GET("v1/aggregate/user/auth", api.ApiGroupApp.AggregateApi.GetUserOpenID)
	public.GET("v1/aggregate/application", api.ApiGroupApp.AggregateApi.GetApplication)
	public.POST("v1/aggregate/payment", api.ApiGroupApp.AggregateApi.Payment)
	public.GET("v1/aggregate/payment/query", api.ApiGroupApp.AggregateApi.Query)
	public.POST("v1/aggregate/test/payment", api.ApiGroupApp.AggregateApi.TestPayment)
	public.GET("v1/aggregate/test/redirect", api.ApiGroupApp.AggregateApi.TestRedirectUrl)

	// 三方回调无需header传签名
	public.POST("v1/notify/payment/:channel/:mchNo/:appNo/:orderNo", api.ApiGroupApp.PublicPaymentApi.Callback)
	public.POST("v1/notify/test/payment/:channel/:mchNo/:appNo/:orderNo", api.ApiGroupApp.PublicPaymentApi.TestCallback)
	// 退款回调
	public.POST("v1/notify/refund/:channel/:mchNo/:appNo/:tradeNo", api.ApiGroupApp.PublicRefundApi.Callback)

	// 收银台支付
	public.POST("v1/cashier/payment", api.ApiGroupApp.CashierApi.Payment)
	public.GET("v1/cashier/query", api.ApiGroupApp.CashierApi.Query)
	public.GET("v1/cashier/redirect", api.ApiGroupApp.CashierApi.RedirectUrl)
	// 需要验签名
	signPublic := public.Use(middle.SignMiddleware(svc.ServiceContextApp))
	signPublic.POST("v1/cashier/create", api.ApiGroupApp.CashierApi.Create)
	// 接口支付
	signPublic.POST("v1/payment/pay", api.ApiGroupApp.PublicPaymentApi.Pay)
	signPublic.GET("v1/payment/query", api.ApiGroupApp.PublicPaymentApi.Query)
	signPublic.POST("v1/payment/close", api.ApiGroupApp.PublicPaymentApi.Close)
	// 接口退款
	signPublic.POST("v1/refund", api.ApiGroupApp.PublicRefundApi.Refund)
	signPublic.GET("v1/refund", api.ApiGroupApp.PublicRefundApi.Query)
}
