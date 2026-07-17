// Package routing 路由规则插件：支付路由策略管理
package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/enum"
	"github.com/lihongsheng/pay-gateway/model/system"
	"github.com/lihongsheng/pay-gateway/plugin"
	routingApi "github.com/lihongsheng/pay-gateway/plugin/routing/api"
	"github.com/lihongsheng/pay-gateway/plugin/routing/domain"
	routingModel "github.com/lihongsheng/pay-gateway/plugin/routing/repo/model"
	routingRepo "github.com/lihongsheng/pay-gateway/plugin/routing/repo"
	"github.com/lihongsheng/pay-gateway/plugin/routing/service"
	"gorm.io/gorm"
)

type p struct{}

func (p) Name() string    { return "routing" }
func (p) Version() string { return "0.1.0" }

func (p) Models() []interface{} {
	return []interface{}{
		&routingModel.PayRoutingRule{},
		&routingModel.PayRoutingRuleCondition{},
		&routingModel.PayRoutingRuleLog{},
		&routingModel.PayRoutingRuleOrderLog{},
	}
}

func (p) Menus() []system.SysMenu {
	return []system.SysMenu{
		{
			Type:       system.MenuTypeCatalog,
			Path:       "/plugin/routing",
			Name:       "RoutingManage",
			Component:  "Layout",
			Title:      "路由管理",
			Icon:       "guide",
			Sort:       90,
			SystemType: enum.SystemTypePlatform,
			Children: []system.SysMenu{
				{
					Type:      system.MenuTypeMenu,
					Path:      "routing-rule",
					Name:      "RoutingRule",
					Component: "plugin/routing/index",
					Title:     "路由策略",
					Icon:      "switch",
					Sort:      1,
					ApiRules:  `[{"path":"/api/plugin/routing/schema","method":"GET"},{"path":"/api/plugin/routing/rule/search","method":"GET"},{"path":"/api/plugin/routing/rule","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增规则", Permission: "routing_rule:add", ApiRules: `[{"path":"/api/plugin/routing/rule","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑规则", Permission: "routing_rule:edit", ApiRules: `[{"path":"/api/plugin/routing/rule","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "切换状态", Permission: "routing_rule:status", ApiRules: `[{"path":"/api/plugin/routing/rule/status","method":"POST"}]`},
					},
				},
			},
		},
	}
}

// RegisterRoute 注册路由
func (p) RegisterRoute(g *gin.Engine, privatePlugin *gin.RouterGroup) {
	private := privatePlugin.Group("routing")

	// Schema
	private.GET("schema", routingApi.GetRoutingSchema)

	// 路由规则 CRUD
	private.GET("rule/search", routingApi.PageRoutingRules)
	private.GET("rule", routingApi.GetRoutingRule)
	private.POST("rule", routingApi.SaveRoutingRule)
	private.POST("rule/status", routingApi.ToggleRoutingRuleStatus)

	// 可用收单账号
	private.GET("rule/available-accounts", routingApi.ListAvailableSingleAccounts)
}

// InitServices 初始化路由插件服务层
func (p) InitServices(ctx plugin.InitContext) error {
	// 创建 domain 服务
	schemaService := domain.NewRuleSchemaService()
	ruleEngine := domain.NewRuleEngine(schemaService)

	// 创建 repo 层
	ruleRepo := routingRepo.NewRoutingRuleRepo(ctx.DB)

	// 创建 service 层
	ruleService := service.NewRoutingRuleService(ruleRepo, ruleEngine, schemaService)

	// 初始化 API 层
	routingApi.InitAPI(ruleService, schemaService)

	return nil
}

func (p) SeedTable(db *gorm.DB) error {
	return nil // 路由插件无需种子数据
}

func init() { plugin.Register(p{}) }
