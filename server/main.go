package main

import (
	"github.com/lihongsheng/pay-gateway/core"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/initialize"
	"github.com/lihongsheng/pay-gateway/queue/initalize"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// 这部分 @Tag 设置用于排序, 需要排序的接口请按照下面的格式添加
// swag init 对 @Tag 只会从入口文件解析, 默认 main.go
// 也可通过 --generalInfo flag 指定其他文件
// @Tag.Name        Base
// @Tag.Name        SysUser
// @Tag.Description 用户

// @title                       Gin-Vue-Admin Swagger API接口文档
// @version                     v2.8.7
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	// 初始化系统
	job := initializeSystem()
	// 运行服务器
	core.RunServer()
	if job != nil {
		defer job()
	}
}

// initializeSystem 初始化系统所有组件
// 提取为单独函数以便于系统重载时调用
func initializeSystem() func() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	global.GVA_PAY_DB = initialize.GormMysqlByConfig(global.GVA_CONFIG.PaymentMysql)
	producer, err := initialize.InitKafka()
	if err != nil {
		panic(err)
	}
	global.GVA_KAFKA_PRODUCER = producer
	initialize.Redis()
	initialize.SetupHandlers() // 注册全局函数
	// 运行消费者和定时任务
	consumer := initalize.Init()
	//cronSer := cronInit.Init()
	return func() {
		//	_ = cronSer.Stop()
		_ = consumer.Stop()
	}
}
