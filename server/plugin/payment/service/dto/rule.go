package dto

type RuleConditions struct {
	AccountNO string `json:"account_no"`
	// 用户登录是否已经被限制
	UserLimit bool `json:"user_limit"`
	// 最大收款次数
	MaxLimit int64 `json:"max_limit"`
	// 成功收款次数
	RequestSuccess int64 `json:"request_success"`
	// SuccessAmount 成功收款金额
	SuccessAmount int64 `json:"success_amount"`
	// MaxAmount 最大收款金额
	MaxAmount       int64    `json:"max_amount"`
	FilterAccountNO []string `json:"filter_account_no"`
	// 是否限制支付方式
	PaymentLimit bool `json:"payment_limit"`
}
