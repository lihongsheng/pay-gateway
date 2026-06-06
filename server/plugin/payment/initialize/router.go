package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/middleware"
	middle "github.com/lihongsheng/pay-gateway/plugin/payment/middleware"
	"github.com/lihongsheng/pay-gateway/plugin/payment/router"
)

func Router(engine *gin.Engine) {
	public := engine.Group("/public")
	public.Use(middle.LogTraceID())
	private := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("private/v1")
	private.Use(middleware.JWTAuth()).Use(middleware.OperationRecord()).Use(middleware.CasbinHandler()).Use(middle.LogTraceID())
	router.InitRoute(private, public)
}
