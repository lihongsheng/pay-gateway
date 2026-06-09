package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/initialize"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/dao"
)

var Plugin = new(plugin)

type plugin struct{}

func (p *plugin) Register(group *gin.Engine) {
	// 读取配置文件
	initialize.Viper()
	// 注册api
	initialize.Api()
	//initialize.Menu(ctx)
	// 注册路由
	initialize.Router(group)
	dao.SetDefault(global.GVA_PAY_DB)
}
