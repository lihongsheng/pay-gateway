package dto

// RoutingRuleQueryRequest 路由规则分页查询请求
type RoutingRuleQueryRequest struct {
	RuleName      string `json:"ruleName" form:"ruleName"`             // 规则名称（模糊搜索）
	RuleStatus    string `json:"ruleStatus" form:"ruleStatus"`         // 规则状态 active/inactive
	RuleAttribute string `json:"ruleAttribute" form:"ruleAttribute"`   // 规则属性 single_merchant/multi_merchant
	MchNo         string `json:"mchNo" form:"mchNo"`                   // 商户编号
	AppNo         string `json:"appNo" form:"appNo"`                   // 应用编号
	BeginTime     string `json:"beginTime" form:"beginTime"`           // 开始时间
	EndTime       string `json:"endTime" form:"endTime"`               // 结束时间
	Page          int    `json:"page" form:"page"`                     // 页码
	PageSize      int    `json:"pageSize" form:"pageSize"`             // 每页条数
	OrderBy       string `json:"orderBy" form:"orderBy"`               // 排序字段
	OrderDir      string `json:"orderDir" form:"orderDir"`             // 排序方向 asc/desc
}

// RoutingRuleCreateRequest 新增/编辑路由规则请求
type RoutingRuleCreateRequest struct {
	ID            int64          `json:"id"`                                    // 编辑时传入
	MchNo         string         `json:"mchNo"`                                 // 商户编号
	AppNo         string         `json:"appNo"`                                 // 应用编号
	RuleName      string         `json:"ruleName" binding:"required,min=1"`     // 规则名称
	RuleDesc      string         `json:"ruleDesc"`                              // 规则描述
	RuleAttribute string         `json:"ruleAttribute" binding:"required"`      // 规则属性 single_merchant/multi_merchant
	Priority      int            `json:"priority"`                              // 优先级
	Rules         map[string]any `json:"rules"`                                 // QueryBuilder 规则 JSON
	Actions       []RuleAction   `json:"actions"`                               // 执行动作列表
}

// RuleAction 执行动作
type RuleAction struct {
	ActionType       string `json:"actionType"`        // 动作类型 select_paypal_account
	AccountAttribute string `json:"accountAttribute"`  // 账号属性
	PaypalAccountNo  string `json:"paypalAccountNo"`   // PayPal 账号（单商户时指定）
}

// RoutingRuleStatusRequest 切换规则状态请求
type RoutingRuleStatusRequest struct {
	ID         int64  `json:"id" binding:"required,gt=0"`
	RuleStatus string `json:"ruleStatus" binding:"required,oneof=active inactive"` // active/inactive
}

// RoutingRuleSummary 路由规则统计摘要
type RoutingRuleSummary struct {
	Total          int64 `json:"total"`          // 规则总数
	Active         int64 `json:"active"`         // 激活中
	Inactive       int64 `json:"inactive"`       // 未激活
	SingleMerchant int64 `json:"singleMerchant"` // 单商户规则
	MultiMerchant  int64 `json:"multiMerchant"`  // 多商户规则
}

// RoutingRuleResponse 路由规则响应（前端展示用）
type RoutingRuleResponse struct {
	ID            int64          `json:"id"`
	RuleNo        string         `json:"ruleNo"`
	MchNo         string         `json:"mchNo"`
	AppNo         string         `json:"appNo"`
	RuleName      string         `json:"ruleName"`
	RuleDesc      string         `json:"ruleDesc"`
	RuleStatus    string         `json:"ruleStatus"`    // active/inactive
	RuleAttribute string         `json:"ruleAttribute"` // single_merchant/multi_merchant
	Priority      int            `json:"priority"`
	Rules         map[string]any `json:"rules,omitempty"`
	Actions       []RuleAction   `json:"actions,omitempty"`
	CreateTime    string         `json:"createTime"`
	UpdateTime    string         `json:"updateTime"`
}

// RoutingRulePageResult 路由规则分页结果
type RoutingRulePageResult struct {
	Rows    []RoutingRuleResponse `json:"rows"`
	Total   int64                 `json:"total"`
	Summary RoutingRuleSummary    `json:"summary"`
}

// AvailableAccount 可用收单账号
type AvailableAccount struct {
	AccountNo    string `json:"accountNo"`
	AccountEmail string `json:"accountEmail"`
	AccountName  string `json:"accountName"`
}
