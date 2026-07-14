package model

import (
	"time"
)

const TableNamePayRoutingRuleLog = "pay_routing_rule_log"

// PayRoutingRuleLog 路由规则修改历史记录
type PayRoutingRuleLog struct {
	ID            int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	RuleID        int64      `gorm:"column:rule_id;type:bigint;not null;index;comment:规则ID" json:"rule_id"`                     // 规则ID
	OperationType string     `gorm:"column:operation_type;type:varchar(64);not null;comment:操作类型" json:"operation_type"`        // 操作类型
	BeforeValue   string     `gorm:"column:before_value;type:json;comment:修改前值" json:"before_value"`                           // 修改前值
	AfterValue    string     `gorm:"column:after_value;type:json;comment:修改后值" json:"after_value"`                             // 修改后值
	Operator      string     `gorm:"column:operator;type:varchar(64);comment:操作人" json:"operator"`                             // 操作人
	Remark        string     `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`                                 // 备注
	CreatedAt     time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
}

// TableName PayRoutingRuleLog's table name
func (*PayRoutingRuleLog) TableName() string {
	return TableNamePayRoutingRuleLog
}
