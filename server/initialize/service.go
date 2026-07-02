// Package initialize service / repo 装配
package initialize

import (
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/log"
	repoExampleNote "github.com/lihongsheng/pay-gateway/plugin/example/repo"
	serviceExampleNote "github.com/lihongsheng/pay-gateway/plugin/example/service"
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
	log.Info("db services initialized")
}
