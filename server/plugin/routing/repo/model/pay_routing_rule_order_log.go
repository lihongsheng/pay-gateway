package model

import (
	"time"
)

const TableNamePayRoutingRuleOrderLog = "pay_routing_rule_order_log"

// PayRoutingRuleOrderLog 订单使用的路由规则记录
type PayRoutingRuleOrderLog struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	RuleID    int64     `gorm:"column:rule_id;type:bigint;not null;index;comment:规则ID" json:"rule_id"`                         // 规则ID
	OrderNo   string    `gorm:"column:order_no;type:varchar(64);not null;comment:订单号" json:"order_no"`                         // 订单号
	RuleJSON  string    `gorm:"column:rule_json;type:json;comment:嵌套规则JSON（QueryBuilder格式）" json:"rule_json"`                 // 嵌套规则JSON
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
}

// TableName PayRoutingRuleOrderLog's table name
func (*PayRoutingRuleOrderLog) TableName() string {
	return TableNamePayRoutingRuleOrderLog
}
