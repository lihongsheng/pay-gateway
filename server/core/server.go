package core

import (
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/initialize"
	"github.com/lihongsheng/pay-gateway/service/system"
	"go.uber.org/zap"
	"time"
)

func RunServer() {
	if global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
		if global.GVA_CONFIG.System.UseMultipoint {
			initialize.RedisList()
		}
	}

	if global.GVA_CONFIG.System.UseMongo {
		err := initialize.Mongo.Initialization()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}
	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)

	fmt.Printf(`
	欢迎使用
	当前版本:%s
`, global.Version)
	// address, address, global.GVA_CONFIG.MCP.SSEPath, address, global.GVA_CONFIG.MCP.MessagePath
	initServer(address, Router, 10*time.Minute, 10*time.Minute)
}
