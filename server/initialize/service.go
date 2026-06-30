// Package initialize service / repo 装配
package initialize

import (
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/log"
	repoExampleNote "github.com/lihongsheng/pay-gateway/plugin/example/repo"
	serviceExampleNote "github.com/lihongsheng/pay-gateway/plugin/example/service"
	payRepo "github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	paySvc "github.com/lihongsheng/pay-gateway/plugin/payment/service"
	repoSys "github.com/lihongsheng/pay-gateway/repo/system"
	serviceBase "github.com/lihongsheng/pay-gateway/service/base"
	serviceInstall "github.com/lihongsheng/pay-gateway/service/install"
	serviceSys "github.com/lihongsheng/pay-gateway/service/system"
	casbinUtil "github.com/lihongsheng/pay-gateway/utils/casbin"
)

// InitInstallService 仅装配 install service（与 DB 是否就绪无关）
func InitInstallService() {
	if serviceInstall.Default != nil {
		return
	}
	svc := serviceInstall.NewService()
	if setter, ok := svc.(interface{ SetLogger(log.Logger) }); ok {
		setter.SetLogger(log.Global())
	}
	svc.OnReady(InitDBServices)
	serviceInstall.Default = svc
}

// InitDBServices 装配所有依赖 global.DB 的 service 单例。
func InitDBServices() {
	if global.DB == nil {
		return
	}
	casbinPort := casbinUtil.NewPort()
	// ---------- system repo ----------
	userRepo := repoSys.NewUserRepo(global.DB)
	roleRepo := repoSys.NewRoleRepo(global.DB)
	menuRepo := repoSys.NewMenuRepo(global.DB)
	mchRepo := repoSys.NewMchRepo(global.DB)
	noteRepo := repoExampleNote.NewNoteRepo(global.DB)
	// ---------- system service ----------
	serviceSys.DefaultUser = serviceSys.NewUserService(userRepo)
	serviceSys.DefaultRole = serviceSys.NewRoleService(roleRepo, menuRepo, casbinPort)
	serviceSys.DefaultMenu = serviceSys.NewMenuService(menuRepo, userRepo)
	serviceSys.DefaultMch = serviceSys.NewMchService(mchRepo)
	// ---------- base service ----------
	serviceBase.Default = serviceBase.NewService(userRepo)
	// ---------- example plugin service ----------
	serviceExampleNote.DefaultNote = serviceExampleNote.NewNoteService(noteRepo)

	// ---------- payment plugin repo ----------
	payMchRepo := payRepo.NewMchRepo(global.DB)
	payAppRepo := payRepo.NewApplicationRepo(global.DB, global.Redis)
	payAccountRepo := payRepo.NewPaymentAccountRepo(global.DB, global.Redis)
	payOrderRepo := payRepo.NewPaymentOrderRepo(global.DB, global.Redis)
	payRefundRepo := payRepo.NewRefundRepo(global.DB, global.Redis)
	payRouterRepo := payRepo.NewRouterRepo(global.DB, global.Redis)
	payTradeStatsRepo := payRepo.NewTradeStaticsRepo(global.DB)
	payNotifyRepo := payRepo.NewNotifyRepo(global.DB)
	payEventRepo := payRepo.NewEventRecordRepo(global.DB)
	payStatsRepo := payRepo.NewStatistics(global.DB)

	// ---------- payment plugin service ----------
	_ = payMchRepo
	_ = payAppRepo
	_ = payAccountRepo
	_ = payOrderRepo
	_ = payRefundRepo
	_ = payRouterRepo
	_ = payTradeStatsRepo
	_ = payNotifyRepo
	_ = payEventRepo
	_ = payStatsRepo
	_ = paySvc.DefaultMch // Ensure package is imported
	// TODO: Wire payment service singletons once service layer is refactored
	// paySvc.DefaultMch = paySvc.NewMchService(payMchRepo, payAppRepo)
	// paySvc.DefaultApplication = paySvc.NewApplicationService(payAppRepo, payMchRepo)
	// ... etc

	log.Info("db services initialized")
}
