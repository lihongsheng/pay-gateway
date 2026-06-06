// middleware/sign_middleware.go
package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/contextkey"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
)

// SignMiddleware 验签中间件：获取应用信息+验签，注入到标准 Context
func SignMiddleware(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求提取应用标识（如 Header/参数，按你的业务调整）
		appNo := c.GetHeader(enum.MetaAppNo)
		if appNo == "" {
			c.AbortWithStatusJSON(400, gin.H{"code": 400, "msg": "应用编号为空"})
			return
		}
		// 2. 首次获取应用信息（可加缓存优化，见下文）
		appInfo, err := svcCtx.AppRepo.GetByAppNoFormCache(c.Request.Context(), appNo)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"code": errors.ErrAppNotExist, "msg": fmt.Sprintf("获取应用信息失败：%v", err)})
			return
		}
		// 3. 验签逻辑（使用 appInfo 的秘钥）
		if err := utils.SignValidate(appInfo.Secret, c.Request); err != nil {
			c.AbortWithStatusJSON(403, err)
			return
		}
		// 4. 核心：将应用信息注入到标准 Context，并替换 Gin 请求的 Context
		newCtx := contextkey.WithAppInfo(c.Request.Context(), appInfo)
		c.Request = c.Request.WithContext(newCtx) // 替换后，后续所有流程的 ctx 都包含应用信息
		c.Next()                                  // 继续执行后续处理逻辑
	}
}
