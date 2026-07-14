package model

import (
	"github.com/lihongsheng/pay-gateway/plugin/routing/enum"
	"time"
)

const TableNamePayRoutingRule = "pay_routing_rule"

// PayRoutingRule 路由规则
type PayRoutingRule struct {
	ID                 int64                     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	RuleNo             string                    `gorm:"column:rule_no;type:varchar(32);not null;uniqueIndex;comment:规则编号" json:"rule_no"`                         // 规则编号
	RuleName           string                    `gorm:"column:rule_name;type:varchar(128);not null;comment:规则名称" json:"rule_name"`                                // 规则名称
	RuleDesc           string                    `gorm:"column:rule_desc;type:text;comment:规则描述" json:"rule_desc"`                                                 // 规则描述
	RuleStatus         enum.RuleStatus           `gorm:"column:rule_status;type:int;not null;default:1;comment:1激活 2未激活" json:"rule_status"`                      // 规则状态
	RuleAttribute      string                    `gorm:"column:rule_attribute;type:varchar(32);not null;comment:single_merchant/multi_merchant" json:"rule_attribute"` // 规则属性
	ConditionLogic     string                    `gorm:"column:condition_logic;type:varchar(16);not null;default:AND;comment:条件逻辑" json:"condition_logic"`         // 条件逻辑
	RuleJSON           string                    `gorm:"column:rule_json;type:json;comment:嵌套规则JSON（QueryBuilder格式）" json:"rule_json"`                           // 嵌套规则JSON
	RuleAttributeValue string                    `gorm:"column:rule_attribute_value;type:json;default:{};" json:"rule_attribute_value"`                                // 规则属性
	Priority           int                       `gorm:"column:priority;type:int unsigned;not null;default:100" json:"priority"`                                       // 优先级
	CreateBy           int                       `gorm:"column:create_by;type:int unsigned;not null;default:0;comment:创建人" json:"create_by"`                        // 创建人
	UpdateBy           int                       `gorm:"column:update_by;type:int unsigned;not null;default:0;comment:更新人" json:"update_by"`                        // 更新人
	Remark             string                    `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`                                                   // 备注
	MchNo              string                    `gorm:"column:mch_no;type:varchar(100);mchAppCreateIndex;comment:商户编号" json:"mch_no"`                             // 商户编号
	AppNo              string                    `gorm:"column:app_no;type:varchar(100);mchAppCreateIndex;not null;comment:应用编号" json:"app_no"`                    // 应用编号
	CreatedAt          time.Time                 `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`        // 创建时间
	UpdatedAt          time.Time                 `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`        // 更新时间
	DeletedAt          *time.Time                `gorm:"column:deleted_at;type:datetime;index;comment:删除时间" json:"deleted_at"`                                     // 删除时间
	Conditions         []PayRoutingRuleCondition `gorm:"foreignKey:RuleID;references:ID;constraint:false" json:"conditions,omitempty"`
}

// TableName PayRoutingRule's table name
func (*PayRoutingRule) TableName() string {
	return TableNamePayRoutingRule
}
