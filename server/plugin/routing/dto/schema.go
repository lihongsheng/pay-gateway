package dto

// FieldDefinition 条件字段定义
type FieldDefinition struct {
	ID         string           `json:"id"`
	Label      string           `json:"label"`
	Type       string           `json:"type"`
	Input      string           `json:"input"`
	Values     []OptionItem     `json:"values,omitempty"`
	Operators  []string         `json:"operators"`
	Validation *FieldValidation `json:"validation,omitempty"`
}

// OptionItem 选项键值对
type OptionItem struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

// FieldValidation 字段校验规则
type FieldValidation struct {
	Min  *float64 `json:"min,omitempty"`
	Step *float64 `json:"step,omitempty"`
}

// OperatorDefinition 操作符定义
type OperatorDefinition struct {
	ID     string   `json:"id"`
	Label  string   `json:"label"`
	Accept []string `json:"accept"`
}

// ConditionLogicDefinition 条件逻辑定义 (AND/OR)
type ConditionLogicDefinition struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// StatusOption 规则状态选项
type StatusOption struct {
	Label string `json:"label"`
	Value int    `json:"value"`
	Type  string `json:"type"`
}

// AttributeOption 规则属性选项
type AttributeOption struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

// RuleSchema 完整的规则 schema 定义
type RuleSchema struct {
	Fields               []FieldDefinition          `json:"fields"`
	Operators            []OperatorDefinition       `json:"operators"`
	Conditions           []ConditionLogicDefinition `json:"conditions"`
	RuleStatusOptions    []StatusOption             `json:"ruleStatusOptions"`
	RuleAttributeOptions []AttributeOption          `json:"ruleAttributeOptions"`
}

type RuleContext struct {
	Data     map[string]any       `json:"data,omitempty"`
	Accounts []RuleAccountContext `json:"accounts,omitempty"`
}

type RuleAccountContext struct {
	AccountNo string  `json:"account_no,omitempty"`
	DayAmount float64 `json:"day_amount,omitempty"`
	DayOrder  int64   `json:"day_order,omitempty"`
}
