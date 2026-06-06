package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/plugin/announcement"
	"github.com/lihongsheng/pay-gateway/plugin/payment"
	"github.com/lihongsheng/pay-gateway/utils/plugin/v2"
)

func PluginInitV2(group *gin.Engine, plugins ...plugin.Plugin) {
	for i := 0; i < len(plugins); i++ {
		plugins[i].Register(group)
	}
}
func bizPluginV2(engine *gin.Engine) {
	PluginInitV2(engine, announcement.Plugin)
	PluginInitV2(engine, payment.Plugin)
}
