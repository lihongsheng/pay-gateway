package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lihongsheng/pay-gateway/plugin/routing/dto"
	"github.com/lihongsheng/pay-gateway/plugin/routing/enum"
	"github.com/lihongsheng/pay-gateway/plugin/routing/repo/model"
	"gorm.io/gorm"
)

// RoutingRuleRepo 路由规则仓储接口
type RoutingRuleRepo interface {
	// 分页查询
	Page(ctx context.Context, req *dto.RoutingRuleQueryRequest) ([]*model.PayRoutingRule, int64, error)
	// 统计摘要
	Summary(ctx context.Context, req *dto.RoutingRuleQueryRequest) (*dto.RoutingRuleSummary, error)
	// 根据ID获取
	GetByID(ctx context.Context, id int64) (*model.PayRoutingRule, error)
	// 创建
	Create(ctx context.Context, rule *model.PayRoutingRule) error
	// 更新
	Update(ctx context.Context, rule *model.PayRoutingRule) error
	// 切换状态
	UpdateStatus(ctx context.Context, id int64, status enum.RuleStatus) error
	// 生成规则编号
	NextRuleNo(ctx context.Context) (string, error)
}

type routingRuleRepoImpl struct {
	db *gorm.DB
}

// NewRoutingRuleRepo 创建路由规则仓储
func NewRoutingRuleRepo(db *gorm.DB) RoutingRuleRepo {
	return &routingRuleRepoImpl{db: db}
}

// Page 分页查询路由规则
func (r *routingRuleRepoImpl) Page(ctx context.Context, req *dto.RoutingRuleQueryRequest) ([]*model.PayRoutingRule, int64, error) {
	var rules []*model.PayRoutingRule
	var total int64

	query := r.buildQuery(ctx, req)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}

	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 默认按优先级降序、创建时间降序
	orderBy := "priority DESC, created_at DESC"
	if req.OrderBy != "" {
		dir := "ASC"
		if req.OrderDir == "desc" || req.OrderDir == "DESC" {
			dir = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s, priority DESC, created_at DESC", req.OrderBy, dir)
	}
	query = query.Order(orderBy)

	if err := query.Find(&rules).Error; err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// Summary 统计摘要
func (r *routingRuleRepoImpl) Summary(ctx context.Context, req *dto.RoutingRuleQueryRequest) (*dto.RoutingRuleSummary, error) {
	baseQuery := r.db.WithContext(ctx).Model(&model.PayRoutingRule{}).Where("deleted_at IS NULL")

	var summary dto.RoutingRuleSummary

	// 总数
	if err := baseQuery.Count(&summary.Total).Error; err != nil {
		return nil, err
	}

	// 激活中
	if err := r.db.WithContext(ctx).Model(&model.PayRoutingRule{}).
		Where("rule_status = ? AND deleted_at IS NULL", enum.RuleStatusActive).
		Count(&summary.Active).Error; err != nil {
		return nil, err
	}

	// 未激活
	summary.Inactive = summary.Total - summary.Active

	// 单商户规则
	if err := r.db.WithContext(ctx).Model(&model.PayRoutingRule{}).
		Where("rule_attribute = ? AND deleted_at IS NULL", "single_merchant").
		Count(&summary.SingleMerchant).Error; err != nil {
		return nil, err
	}

	// 多商户规则
	if err := r.db.WithContext(ctx).Model(&model.PayRoutingRule{}).
		Where("rule_attribute = ? AND deleted_at IS NULL", "multi_merchant").
		Count(&summary.MultiMerchant).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

// GetByID 根据ID获取路由规则
func (r *routingRuleRepoImpl) GetByID(ctx context.Context, id int64) (*model.PayRoutingRule, error) {
	var rule model.PayRoutingRule
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Create 创建路由规则
func (r *routingRuleRepoImpl) Create(ctx context.Context, rule *model.PayRoutingRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// Update 更新路由规则
func (r *routingRuleRepoImpl) Update(ctx context.Context, rule *model.PayRoutingRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

// UpdateStatus 切换规则状态
func (r *routingRuleRepoImpl) UpdateStatus(ctx context.Context, id int64, status enum.RuleStatus) error {
	return r.db.WithContext(ctx).Model(&model.PayRoutingRule{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("rule_status", status).Error
}

// NextRuleNo 生成规则编号
func (r *routingRuleRepoImpl) NextRuleNo(ctx context.Context) (string, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.PayRoutingRule{}).Count(&count).Error; err != nil {
		return "", err
	}
	now := time.Now()
	return fmt.Sprintf("RR%s%04d", now.Format("20060102"), count+1), nil
}

// buildQuery 构建查询条件
func (r *routingRuleRepoImpl) buildQuery(ctx context.Context, req *dto.RoutingRuleQueryRequest) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&model.PayRoutingRule{}).Where("deleted_at IS NULL")

	if req.RuleName != "" {
		query = query.Where("rule_name LIKE ?", "%"+req.RuleName+"%")
	}
	if req.RuleStatus != "" {
		status := enum.RuleStatusFromString(req.RuleStatus)
		query = query.Where("rule_status = ?", status)
	}
	if req.RuleAttribute != "" {
		query = query.Where("rule_attribute = ?", req.RuleAttribute)
	}
	if req.MchNo != "" {
		query = query.Where("mch_no = ?", req.MchNo)
	}
	if req.AppNo != "" {
		query = query.Where("app_no = ?", req.AppNo)
	}
	if req.BeginTime != "" {
		query = query.Where("created_at >= ?", req.BeginTime+" 00:00:00")
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime+" 23:59:59")
	}

	return query
}

// ─── 辅助：model → response 转换 ───────────────────────────────────────────

// ModelToResponse 将数据库模型转换为前端响应格式
func ModelToResponse(m *model.PayRoutingRule) *dto.RoutingRuleResponse {
	resp := &dto.RoutingRuleResponse{
		ID:            m.ID,
		RuleNo:        m.RuleNo,
		MchNo:         m.MchNo,
		AppNo:         m.AppNo,
		RuleName:      m.RuleName,
		RuleDesc:      m.RuleDesc,
		RuleStatus:    m.RuleStatus.String(),
		RuleAttribute: m.RuleAttribute,
		Priority:      m.Priority,
		CreateTime:    m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateTime:    m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 解析 RuleJSON
	if m.RuleJSON != "" && m.RuleJSON != "{}" {
		var rules map[string]any
		if err := json.Unmarshal([]byte(m.RuleJSON), &rules); err == nil {
			resp.Rules = rules
		}
	}

	// 解析 RuleAttributeValue 为 actions
	if m.RuleAttributeValue != "" && m.RuleAttributeValue != "{}" {
		var attrVal enum.RuleAttributeValue
		if err := json.Unmarshal([]byte(m.RuleAttributeValue), &attrVal); err == nil {
			var actions []dto.RuleAction
			actionType := "select_paypal_account"
			if m.RuleAttribute == "single_merchant" && len(attrVal.Account) > 0 {
				for _, acc := range attrVal.Account {
					actions = append(actions, dto.RuleAction{
						ActionType:       actionType,
						AccountAttribute: m.RuleAttribute,
						PaypalAccountNo:  acc,
					})
				}
			} else {
				actions = append(actions, dto.RuleAction{
					ActionType:       actionType,
					AccountAttribute: m.RuleAttribute,
					PaypalAccountNo:  "",
				})
			}
			resp.Actions = actions
		}
	}

	return resp
}
