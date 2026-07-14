package domain

import (
	"github.com/lihongsheng/pay-gateway/plugin/routing/dto"
	"github.com/lihongsheng/pay-gateway/plugin/routing/enum"
)

// RuleSchemaService 规则 schema 服务，供前端 QueryBuilder 使用
type RuleSchemaService struct{}

// NewRuleSchemaService 创建规则 schema 服务
func NewRuleSchemaService() *RuleSchemaService {
	return &RuleSchemaService{}
}

// GetSchema 返回完整的规则 schema 定义
func (s *RuleSchemaService) GetSchema() *dto.RuleSchema {
	return &dto.RuleSchema{
		Fields:               s.GetFieldDefinitions(),
		Operators:            s.GetOperatorDefinitions(),
		Conditions:           s.GetConditionDefinitions(),
		RuleStatusOptions:    s.GetRuleStatusOptions(),
		RuleAttributeOptions: s.GetRuleAttributeOptions(),
	}
}

// GetFieldDefinitions 条件字段定义
func (s *RuleSchemaService) GetFieldDefinitions() []dto.FieldDefinition {
	return []dto.FieldDefinition{
		{
			ID:        "merchant_category",
			Label:     "商户经验类别",
			Type:      string(enum.FieldTypeInteger),
			Input:     string(enum.InputTypeText),
			Operators: []string{"equal", "not_equal", "in", "not_in"},
		},
		{
			ID:        "merchant_account",
			Label:     "商户账号",
			Type:      string(enum.FieldTypeString),
			Input:     string(enum.InputTypeText),
			Operators: []string{"equal", "not_equal", "in", "not_in"},
		},
		{
			ID:        "currency",
			Label:     "支付币种",
			Type:      string(enum.FieldTypeString),
			Input:     string(enum.InputTypeSelect),
			Values:    currencies(),
			Operators: []string{"equal", "not_equal", "in", "not_in", "is_null", "is_not_null"},
		},
		{
			ID:        "merchant_region",
			Label:     "商户经营区域",
			Type:      string(enum.FieldTypeString),
			Input:     string(enum.InputTypeText),
			Operators: []string{"equal", "not_equal", "contains", "begins_with", "is_null", "is_not_null"},
		},
		{
			ID:         "account.[].day_order",
			Label:      "账户日交易数",
			Type:       string(enum.FieldTypeInteger),
			Input:      string(enum.InputTypeText),
			Operators:  []string{"equal", "not_equal", "less", "less_or_equal", "greater", "greater_or_equal"},
			Validation: &dto.FieldValidation{Min: ptrFloat64(0), Step: ptrFloat64(1)},
		},
		{
			ID:         "account.[].day_amount",
			Label:      "账户日金额",
			Type:       string(enum.FieldTypeDouble),
			Input:      string(enum.InputTypeText),
			Operators:  []string{"equal", "not_equal", "less", "less_or_equal", "greater", "greater_or_equal"},
			Validation: &dto.FieldValidation{Min: ptrFloat64(0), Step: ptrFloat64(0.01)},
		},
		{
			ID:         "order_amount",
			Label:      "订单金额",
			Type:       string(enum.FieldTypeDouble),
			Input:      string(enum.InputTypeText),
			Operators:  []string{"equal", "not_equal", "less", "less_or_equal", "greater", "greater_or_equal"},
			Validation: &dto.FieldValidation{Min: ptrFloat64(0), Step: ptrFloat64(0.01)},
		},
		{
			ID:        "app_id",
			Label:     "应用ID",
			Type:      string(enum.FieldTypeInteger),
			Input:     string(enum.InputTypeText),
			Operators: []string{"equal", "not_equal", "in", "not_in"},
		},
	}
}

// GetOperatorDefinitions 操作符定义
func (s *RuleSchemaService) GetOperatorDefinitions() []dto.OperatorDefinition {
	ops := enum.AllOperators()
	result := make([]dto.OperatorDefinition, 0, len(ops))
	for _, op := range ops {
		result = append(result, dto.OperatorDefinition{
			ID:     string(op),
			Label:  op.Label(),
			Accept: op.AcceptFieldTypes(),
		})
	}
	return result
}

// GetConditionDefinitions 条件逻辑定义 (AND/OR)
func (s *RuleSchemaService) GetConditionDefinitions() []dto.ConditionLogicDefinition {
	return []dto.ConditionLogicDefinition{
		{ID: "AND", Label: "且"},
		{ID: "OR", Label: "或"},
	}
}

// GetRuleStatusOptions 规则状态选项
func (s *RuleSchemaService) GetRuleStatusOptions() []dto.StatusOption {
	statuses := enum.AllRuleStatuses()
	result := make([]dto.StatusOption, 0, len(statuses))
	for _, status := range statuses {
		result = append(result, dto.StatusOption{
			Label: status.Label(),
			Value: int(status),
			Type:  status.TagType(),
		})
	}
	return result
}

// GetRuleAttributeOptions 规则属性选项
func (s *RuleSchemaService) GetRuleAttributeOptions() []dto.AttributeOption {
	attrs := enum.AllRuleAttributes()
	result := make([]dto.AttributeOption, 0, len(attrs))
	for _, attr := range attrs {
		result = append(result, dto.AttributeOption{
			Label: attr.Label(),
			Value: int(attr),
		})
	}
	return result
}

// GetFieldById 根据 field id 获取字段定义
func (s *RuleSchemaService) GetFieldById(id string) *dto.FieldDefinition {
	for _, field := range s.GetFieldDefinitions() {
		if field.ID == id {
			return &field
		}
	}
	return nil
}

// GetValidFieldIds 获取所有合法的 field id 列表
func (s *RuleSchemaService) GetValidFieldIds() []string {
	fields := s.GetFieldDefinitions()
	ids := make([]string, 0, len(fields))
	for _, field := range fields {
		ids = append(ids, field.ID)
	}
	return ids
}

// ─── 私有数据源 ───────────────────────────────────────────

func merchantCategories() []dto.OptionItem {
	return []dto.OptionItem{
		{Label: "服饰鞋包", Value: 1},
		{Label: "3C数码", Value: 2},
		{Label: "家居日用", Value: 3},
		{Label: "美妆个护", Value: 4},
		{Label: "运动户外", Value: 5},
		{Label: "图书文娱", Value: 6},
		{Label: "家电家装", Value: 7},
	}
}

func currencies() []dto.OptionItem {
	return []dto.OptionItem{
		{Label: "CNY", Value: "CNY"},
		{Label: "USD", Value: "USD"},
		{Label: "EUR", Value: "EUR"},
		{Label: "CAD", Value: "CAD"},
		{Label: "GBP", Value: "GBP"},
		{Label: "AUD", Value: "AUD"},
	}
}

func ptrFloat64(v float64) *float64 {
	return &v
}
