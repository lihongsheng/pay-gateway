package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
)

// SignMiddleware 验签中间件：获取应用信息+验签，注入到标准 Context
func LogTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 4. 核心：将应用信息注入到标准 Context，并替换 Gin 请求的 Context
		newCtx := context.WithValue(c.Request.Context(), log.LogTraceId{}, uuid.NewString())
		c.Request = c.Request.WithContext(newCtx) // 替换后，后续所有流程的 ctx 都包含应用信息
		c.Next()                                  // 继续执行后续处理逻辑
	}
}
