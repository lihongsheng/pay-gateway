package enum

type RuleAttribute int

const (
	// RuleAttributeAll ALl
	RuleAttributeAll RuleAttribute = iota
	// 指定账号
	RuleAttributeSinglePayAccount
	// 指定支付渠道 Wechat, Alipay, Lakala
	RuleAttributeSinglePayChannel
	// 指定支付方式 Wechat, Alipay
	RuleAttributeSinglePayMethod
)

func AllRuleAttributes() []RuleAttribute {
	return []RuleAttribute{RuleAttributeAll, RuleAttributeSinglePayAccount, RuleAttributeSinglePayChannel, RuleAttributeSinglePayMethod}
}

func (a RuleAttribute) Label() string {
	switch a {
	case RuleAttributeAll:
		return "全部可用账户"
	case RuleAttributeSinglePayAccount:
		return "指定账号"
	case RuleAttributeSinglePayChannel:
		return "指定支付渠道"
	case RuleAttributeSinglePayMethod:
		return "指定支付方式"
	default:
		return ""
	}
}

type RuleAttributeValue struct {
	Account []string `json:"account,omitempty"`
	Channel []string `json:"channel,omitempty"`
	Method  []string `json:"method,omitempty"`
}
