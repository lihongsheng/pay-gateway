package model

import (
	"time"
)

const TableNamePayRoutingRuleCondition = "pay_routing_rule_condition"

// PayRoutingRuleCondition 路由规则条件
type PayRoutingRuleCondition struct {
	ID           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	RuleID       int64     `gorm:"column:rule_id;type:bigint;not null;index;comment:规则ID" json:"rule_id"`                     // 规则ID
	FieldKey     string    `gorm:"column:field_key;type:varchar(64);not null;comment:条件左值" json:"field_key"`                 // 条件左值
	Operator     string    `gorm:"column:operator;type:varchar(32);not null;comment:操作符" json:"operator"`                    // 操作符
	ValueType    string    `gorm:"column:value_type;type:varchar(32);not null;default:input;comment:值类型" json:"value_type"`   // 值类型
	FieldValue   string    `gorm:"column:field_value;type:varchar(500);comment:条件右值" json:"field_value"`                     // 条件右值
	GroupKey     string    `gorm:"column:group_key;type:varchar(32);not null;default:G1;comment:条件组" json:"group_key"`       // 条件组
	GroupLogic   string    `gorm:"column:group_logic;type:varchar(16);not null;default:AND;comment:组内逻辑" json:"group_logic"` // 组内逻辑
	SortOrder    int       `gorm:"column:sort_order;type:int unsigned;not null;default:0;comment:排序" json:"sort_order"`      // 排序
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName PayRoutingRuleCondition's table name
func (*PayRoutingRuleCondition) TableName() string {
	return TableNamePayRoutingRuleCondition
}
