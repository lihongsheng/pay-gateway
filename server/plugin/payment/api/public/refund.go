package public

import (
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/payment-sdk/enum/refund"
)

type RefundRequest struct {
	OrderNo        string          `json:"order_no"`
	RefundNo       string          `json:"refund_no"`
	Amount         Amount          `json:"amount"`
	Reason         string          `json:"reason"`
	PassBackParams string          `json:"pass_back_params"`
	NotifyUrl      string          `json:"notify_url"`
	AppNo          string          `json:"app_no"`
	MchNo          string          `json:"-"`
	RefundFrom     enum.RefundFrom `json:"-"`
}

func (r RefundRequest) Validate() error {
	if r.AppNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "app_no不能为空")
	}
	if r.OrderNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "order_no不能为空")
	}
	if r.RefundNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "refund_no不能为空")
	}
	if r.Amount.Total <= 0 {
		return errors.NewError(errors.ErrCodeInvalidParam, "amount.total不能小于0")
	}
	return nil
}

type RefundQueryRequest struct {
	RefundTradeNo string `json:"refund_trade_no" form:"refund_trade_no"`
	RefundNo      string `json:"refund_no" form:"refund_no"`
	AppNo         string `json:"app_no" form:"app_no"`
	MchNo         string `json:"-" form:"-"`
}

func (r RefundQueryRequest) Validate() error {
	if r.AppNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "app_no不能为空")
	}
	if r.RefundNo == "" && r.RefundTradeNo == "" {
		return errors.NewError(errors.ErrCodeInvalidParam, "trade_no或refund_no不能为空")
	}
	return nil
}

type RefundResponse struct {
	OrderNo          string        `json:"order_no"`
	RefundNo         string        `json:"refund_no"`
	Amount           Amount        `json:"amount"`
	Status           refund.Status `json:"refund_status"`
	RefundTradeNo    string        `json:"refund_trade_no"`
	RefundOutTradeNo string        `json:"refund_out_trade_no"`
	OrderAmount      Amount        `json:"order_amount"`
	PassBackParams   string        `json:"pass_back_params"`
}
