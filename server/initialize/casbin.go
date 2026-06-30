package initialize

import (
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/log"
	"github.com/lihongsheng/pay-gateway/utils/casbin"
)

// SetupCasbin 在 global.DB 就绪后初始化 enforcer + 自动建 casbin_rule 表
func SetupCasbin() {
	if global.DB == nil {
		return
	}
	if err := casbin.Setup(global.DB); err != nil {
		log.Error("casbin setup: " + err.Error())
		return
	}
	log.Info("casbin enforcer ready")
}
