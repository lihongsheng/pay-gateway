package dto

import "github.com/lihongsheng/payment-sdk/enum/payment"

type RouterPaymentMethod struct {
	PaymentMethod payment.Payment `json:"payment_method"`
	AccountNo     string          `json:"account_no"`
	Icon          string          `json:"icon"`
}

type AggregateToken struct {
	AppNo     string `json:"a,omitempty"`
	CreatedAt int64  `json:"c,omitempty"`
	//	AccountNo string `json:"ao"`
}
