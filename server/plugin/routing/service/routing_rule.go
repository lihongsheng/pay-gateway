package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lihongsheng/pay-gateway/plugin/routing/dto"
	"github.com/lihongsheng/pay-gateway/plugin/routing/domain"
	"github.com/lihongsheng/pay-gateway/plugin/routing/enum"
	"github.com/lihongsheng/pay-gateway/plugin/routing/repo"
	"github.com/lihongsheng/pay-gateway/plugin/routing/repo/model"
)

// RoutingRuleService 路由规则服务接口
type RoutingRuleService interface {
	// 分页查询路由规则
	Page(ctx context.Context, req *dto.RoutingRuleQueryRequest) (*dto.RoutingRulePageResult, error)
	// 获取路由规则详情
	GetByID(ctx context.Context, id int64) (*dto.RoutingRuleResponse, error)
	// 新增路由规则
	Create(ctx context.Context, req *dto.RoutingRuleCreateRequest) (string, error)
	// 更新路由规则
	Update(ctx context.Context, req *dto.RoutingRuleCreateRequest) (string, error)
	// 切换规则状态
	ToggleStatus(ctx context.Context, req *dto.RoutingRuleStatusRequest) error
	// 获取可用单商户收单账号
	ListAvailableSingleAccounts(ctx context.Context, ruleID int64) ([]dto.AvailableAccount, error)
}

type routingRuleService struct {
	ruleRepo      repo.RoutingRuleRepo
	ruleEngine    *domain.RuleEngine
	schemaService *domain.RuleSchemaService
}

// NewRoutingRuleService 创建路由规则服务
func NewRoutingRuleService(ruleRepo repo.RoutingRuleRepo, ruleEngine *domain.RuleEngine, schemaService *domain.RuleSchemaService) RoutingRuleService {
	return &routingRuleService{
		ruleRepo:      ruleRepo,
		ruleEngine:    ruleEngine,
		schemaService: schemaService,
	}
}

// Page 分页查询路由规则
func (s *routingRuleService) Page(ctx context.Context, req *dto.RoutingRuleQueryRequest) (*dto.RoutingRulePageResult, error) {
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}

	rules, total, err := s.ruleRepo.Page(ctx, req)
	if err != nil {
		return nil, err
	}

	// 获取统计摘要
	summary, err := s.ruleRepo.Summary(ctx, req)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	rows := make([]dto.RoutingRuleResponse, 0, len(rules))
	for _, r := range rules {
		resp := repo.ModelToResponse(r)
		rows = append(rows, *resp)
	}

	return &dto.RoutingRulePageResult{
		Rows:    rows,
		Total:   total,
		Summary: *summary,
	}, nil
}

// GetByID 获取路由规则详情
func (s *routingRuleService) GetByID(ctx context.Context, id int64) (*dto.RoutingRuleResponse, error) {
	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return repo.ModelToResponse(rule), nil
}

// Create 新增路由规则
func (s *routingRuleService) Create(ctx context.Context, req *dto.RoutingRuleCreateRequest) (string, error) {
	if err := s.validateRule(req); err != nil {
		return "", err
	}

	// 生成规则编号
	ruleNo, err := s.ruleRepo.NextRuleNo(ctx)
	if err != nil {
		return "", err
	}

	// 序列化规则 JSON
	ruleJSON, err := json.Marshal(req.Rules)
	if err != nil {
		return "", errors.New("规则JSON序列化失败")
	}

	// 序列化动作（存入 rule_attribute_value）
	attrVal := s.buildRuleAttributeValue(req)
	attrValJSON, err := json.Marshal(attrVal)
	if err != nil {
		return "", errors.New("规则属性值序列化失败")
	}

	rule := &model.PayRoutingRule{
		RuleNo:             ruleNo,
		MchNo:              req.MchNo,
		AppNo:              req.AppNo,
		RuleName:           req.RuleName,
		RuleDesc:           req.RuleDesc,
		RuleStatus:         enum.RuleStatusInactive, // 新建默认未激活
		RuleAttribute:      req.RuleAttribute,
		ConditionLogic:     "AND",
		RuleJSON:           string(ruleJSON),
		RuleAttributeValue: string(attrValJSON),
		Priority:           req.Priority,
	}

	if err := s.ruleRepo.Create(ctx, rule); err != nil {
		return "", err
	}

	return "创建成功", nil
}

// Update 更新路由规则
func (s *routingRuleService) Update(ctx context.Context, req *dto.RoutingRuleCreateRequest) (string, error) {
	if req.ID <= 0 {
		return "", errors.New("规则ID不能为空")
	}

	existing, err := s.ruleRepo.GetByID(ctx, req.ID)
	if err != nil {
		return "", errors.New("规则不存在")
	}

	// 激活状态不允许编辑
	if existing.RuleStatus == enum.RuleStatusActive {
		return "", errors.New("激活中的规则不允许编辑")
	}

	if err := s.validateRule(req); err != nil {
		return "", err
	}

	// 序列化规则 JSON
	ruleJSON, err := json.Marshal(req.Rules)
	if err != nil {
		return "", errors.New("规则JSON序列化失败")
	}

	// 序列化动作
	attrVal := s.buildRuleAttributeValue(req)
	attrValJSON, err := json.Marshal(attrVal)
	if err != nil {
		return "", errors.New("规则属性值序列化失败")
	}

	existing.MchNo = req.MchNo
	existing.AppNo = req.AppNo
	existing.RuleName = req.RuleName
	existing.RuleDesc = req.RuleDesc
	existing.RuleAttribute = req.RuleAttribute
	existing.RuleJSON = string(ruleJSON)
	existing.RuleAttributeValue = string(attrValJSON)
	existing.Priority = req.Priority

	if err := s.ruleRepo.Update(ctx, existing); err != nil {
		return "", err
	}

	return "更新成功", nil
}

// ToggleStatus 切换规则状态
func (s *routingRuleService) ToggleStatus(ctx context.Context, req *dto.RoutingRuleStatusRequest) error {
	existing, err := s.ruleRepo.GetByID(ctx, req.ID)
	if err != nil {
		return errors.New("规则不存在")
	}

	targetStatus := enum.RuleStatusFromString(req.RuleStatus)

	// 如果已经是目标状态，直接返回
	if existing.RuleStatus == targetStatus {
		return nil
	}

	return s.ruleRepo.UpdateStatus(ctx, req.ID, targetStatus)
}

// ListAvailableSingleAccounts 获取可用单商户收单账号
func (s *routingRuleService) ListAvailableSingleAccounts(ctx context.Context, ruleID int64) ([]dto.AvailableAccount, error) {
	// TODO: 从 payment_account 表查询可用的 PayPal 收单账号
	// 当前返回空列表，待与 payment 插件集成后实现
	return []dto.AvailableAccount{}, nil
}

// validateRule 校验规则
func (s *routingRuleService) validateRule(req *dto.RoutingRuleCreateRequest) error {
	if req.MchNo == "" {
		return errors.New("请选择商户")
	}
	if req.AppNo == "" {
		return errors.New("请选择应用")
	}
	if req.RuleName == "" {
		return errors.New("规则名称不能为空")
	}
	if req.RuleAttribute == "" {
		return errors.New("规则属性不能为空")
	}
	switch req.RuleAttribute {
	case "single_merchant", "multi_merchant":
		// valid
	default:
		return fmt.Errorf("不支持的规则属性: %s", req.RuleAttribute)
	}

	// 校验规则 JSON 结构
	if req.Rules != nil {
		if err := s.ruleEngine.Validate(req.Rules, 0); err != nil {
			return fmt.Errorf("规则校验失败: %s", err.Error())
		}
	}

	// 单商户规则必须指定账号
	if req.RuleAttribute == "single_merchant" {
		if len(req.Actions) == 0 || req.Actions[0].PaypalAccountNo == "" {
			return errors.New("单商户规则必须指定收单账号")
		}
	}

	return nil
}

// buildRuleAttributeValue 从请求构建规则属性值
func (s *routingRuleService) buildRuleAttributeValue(req *dto.RoutingRuleCreateRequest) enum.RuleAttributeValue {
	attrVal := enum.RuleAttributeValue{}

	switch req.RuleAttribute {
	case "single_merchant":
		for _, action := range req.Actions {
			if action.PaypalAccountNo != "" {
				attrVal.Account = append(attrVal.Account, action.PaypalAccountNo)
			}
		}
	case "multi_merchant":
		// 多商户不需要指定账号
	}

	return attrVal
}
