package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/enum"
	"github.com/lihongsheng/pay-gateway/plugin/routing/dto"
	"github.com/lihongsheng/pay-gateway/plugin/routing/domain"
	"github.com/lihongsheng/pay-gateway/plugin/routing/service"
	serviceSys "github.com/lihongsheng/pay-gateway/service/system"
	"github.com/lihongsheng/pay-gateway/utils/jwt"
	"github.com/lihongsheng/pay-gateway/utils/response"
	"go.uber.org/zap"
)

// ─── 全局服务引用，由 plugin.go InitServices 初始化 ─────────────────────────

var (
	ruleService    service.RoutingRuleService
	schemaService  *domain.RuleSchemaService
)

// getRoutingUserInfo 从 JWT 获取当前用户信息，返回 (systemType, mchNo)
// 对于商户用户，将 MchID（数据库主键）转换为 mch_no（业务编号）
func getRoutingUserInfo(c *gin.Context) (systemType enum.SystemType, mchNo string) {
	user, err := jwt.GetUser(c.Request.Context())
	if err != nil {
		return
	}
	systemType = user.SystemType
	if user.MchID > 0 {
		// 通过 MchService 将 MchID 转为 MchNo
		mch, err := serviceSys.DefaultMch.Get(c.Request.Context(), user.MchID)
		if err == nil && mch != nil {
			mchNo = mch.MchNo
		} else {
			// 降级：使用数字 ID 作为兜底
			mchNo = fmt.Sprintf("%d", user.MchID)
		}
	}
	return
}

// InitAPI 初始化 API 层依赖
func InitAPI(svc service.RoutingRuleService, schema *domain.RuleSchemaService) {
	ruleService = svc
	schemaService = schema
}

// ─── API 处理函数 ─────────────────────────────────────────────────────────

// GetRoutingSchema
// @Summary   获取路由规则 Schema
// @Tags      RoutingRule
// @Security  ApiKeyAuth
// @Produce   json
// @Success   200  {object}  response.Body{data=dto.RuleSchema}
// @Router    /api/plugin/routing/schema [get]
func GetRoutingSchema(c *gin.Context) {
	result := schemaService.GetSchema()
	response.OK(c, result)
}

// PageRoutingRules
// @Summary   分页查询路由规则
// @Tags      RoutingRule
// @Security  ApiKeyAuth
// @Produce   json
// @Param     data  query  dto.RoutingRuleQueryRequest  true  "查询参数"
// @Success   200  {object}  response.Body{data=dto.RoutingRulePageResult}
// @Router    /api/plugin/routing/rule/search [get]
func PageRoutingRules(c *gin.Context) {
	var req dto.RoutingRuleQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 商户用户只能查询自己的规则
	systemType, mchNo := getRoutingUserInfo(c)
	if systemType == enum.SystemTypeMch && mchNo != "" {
		req.MchNo = mchNo
	}

	result, err := ruleService.Page(c.Request.Context(), &req)
	if err != nil {
		zap.L().Error("查询路由规则失败", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}

	response.OK(c, result)
}

// GetRoutingRule
// @Summary   获取路由规则详情
// @Tags      RoutingRule
// @Security  ApiKeyAuth
// @Produce   json
// @Param     id  query  int  true  "规则ID"
// @Success   200  {object}  response.Body{data=dto.RoutingRuleResponse}
// @Router    /api/plugin/routing/rule [get]
func GetRoutingRule(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.Fail(c, "参数错误：id不能为空")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, "参数错误：id无效")
		return
	}

	result, err := ruleService.GetByID(c.Request.Context(), id)
	if err != nil {
		zap.L().Error("获取路由规则详情失败", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}

	response.OK(c, result)
}

// SaveRoutingRule
// @Summary   新增/更新路由规则（有id为更新，无id为新增）
// @Tags      RoutingRule
// @Security  ApiKeyAuth
// @Accept    json
// @Produce   json
// @Param     data  body  dto.RoutingRuleCreateRequest  true  "规则数据"
// @Success   200  {object}  response.Body{msg=string}
// @Router    /api/plugin/routing/rule [post]
func SaveRoutingRule(c *gin.Context) {
	var req dto.RoutingRuleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 商户用户自动填充 mchNo，不允许修改
	systemType, mchNo := getRoutingUserInfo(c)
	if systemType == enum.SystemTypeMch && mchNo != "" {
		req.MchNo = mchNo
	}

	var msg string
	var err error
	if req.ID > 0 {
		msg, err = ruleService.Update(c.Request.Context(), &req)
	} else {
		msg, err = ruleService.Create(c.Request.Context(), &req)
	}

	if err != nil {
		zap.L().Error("保存路由规则失败", zap.Error(err))
		response.Fail(c, err.Error())
		return
	}

	response.OKMsg(c, msg)
}

// ToggleRoutingRuleStatus
// @Summary   切换路由规则状态
// @Tags      RoutingRule
// @Security  ApiKeyAuth
// @Accept    json
// @Produce   json
// @Param     data  body  dto.RoutingRuleStatusRequest  true  "状态切换请求"
// @Success   200  {object}  response.Body{msg=string}
// @Router    /api/plugin/routing/rule/status [post]
func ToggleRoutingRuleStatus(c *gin.Context) {
	var req dto.RoutingRuleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := ruleService.ToggleStatus(c.Request.Context(), &req); err != nil {
		zap.L().Error("切换路由规则状态失败", zap.Error(err))
		response.Fail(c, err.Error())
		return
	}

	response.OKMsg(c, "状态更新成功")
}

// ListAvailableSingleAccounts
// @Summary   获取可用单商户收单账号
// @Tags      RoutingRule
// @Security  ApiKeyAuth
// @Produce   json
// @Param     ruleId  query  int  false  "规则ID（排除已绑定的账号）"
// @Success   200  {object}  response.Body{data=[]dto.AvailableAccount}
// @Router    /api/plugin/routing/rule/available-accounts [get]
func ListAvailableSingleAccounts(c *gin.Context) {
	var ruleID int64
	if ruleIdStr := c.Query("ruleId"); ruleIdStr != "" {
		var err error
		ruleID, err = strconv.ParseInt(ruleIdStr, 10, 64)
		if err != nil {
			ruleID = 0
		}
	}

	result, err := ruleService.ListAvailableSingleAccounts(c.Request.Context(), ruleID)
	if err != nil {
		zap.L().Error("获取可用收单账号失败", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}

	response.OK(c, result)
}
