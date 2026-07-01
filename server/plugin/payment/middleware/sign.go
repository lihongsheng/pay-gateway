package middleware

import (
	"fmt"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"

	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/contextkey"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
)

// SignMiddleware 验签中间件：获取应用信息+验签，注入到标准 Context
func SignMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appNo := c.GetHeader(enum.MetaAppNo)
		if appNo == "" {
			c.AbortWithStatusJSON(400, gin.H{"code": 400, "msg": "应用编号为空"})
			return
		}
		appInfo, err := svc.ServiceContextApp.AppRepo.GetByAppNoFormCache(c.Request.Context(), appNo)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"code": errors.ErrAppNotExist, "msg": fmt.Sprintf("获取应用信息失败：%v", err)})
			return
		}
		if err := utils.SignValidate(appInfo.Secret, c.Request); err != nil {
			c.AbortWithStatusJSON(403, err)
			return
		}
		newCtx := contextkey.WithAppInfo(c.Request.Context(), appInfo)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
