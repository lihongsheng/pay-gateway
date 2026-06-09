package event

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"github.com/lihongsheng/payment-sdk/enum/refund"
	"net/http"
	"time"
)

// 事件类型常量（统一管理，避免硬编码）
const (
	EventTypePaymentOrderStatus = "payment:status"
	EventTypeRefundStatus       = "refund:status"
	EventTypeCallbackPayment    = "callback:payment"
	EventTypeCallbackRefund     = "callback:refund"
	EventTypeRefundNotifyRetry  = "refund:notify:retry"
	EventTypePaymentNotifyRetry = "payment:notify:retry"
)

type Event interface {
	GetEventKey() string
}

type Base struct {
	EventType string `json:"event_type"`
	EventId   string `json:"event_id"`
}

func (e Base) GetEventKey() string {
	return e.EventId
}

type PaymentOrderStatusEvent struct {
	Base
	MchNo     string         `json:"mch_no"`
	OrderNo   string         `json:"order_no"`
	TradeNo   string         `json:"trade_no"`
	AppNo     string         `json:"app_no"`
	OldStatus payment.Status `json:"old_status"`
	NewStatus payment.Status `json:"new_status"`
	Amount    int64          `json:"amount"`
	ErrCode   int            `json:"err_code"`
	AccountNo string         `json:"account_no"`
	CreateAt  time.Time      `json:"create_at"`
	Retry     int64          `json:"retry"`
	ID        int64          `json:"id"`
}

func (e PaymentOrderStatusEvent) GetEventKey() string {
	return fmt.Sprintf("%s%s%s", e.MchNo, e.AppNo, e.OrderNo)
}

type RefundOrderStatusEvent struct {
	Base
	MchNo         string        `json:"mch_no"`
	OrderNo       string        `json:"order_no"`
	RefundNo      string        `json:"refund_no"`
	RefundTradeNo string        `json:"refund_trade_no"`
	AppNo         string        `json:"app_no"`
	OldStatus     refund.Status `json:"old_status"`
	NewStatus     refund.Status `json:"new_status"`
	AccountNo     string        `json:"account_no"`
	Amount        int64         `json:"amount"`
}

func (e RefundOrderStatusEvent) GetEventKey() string {
	return fmt.Sprintf("%s%s%s", e.MchNo, e.AppNo, e.RefundNo)
}

func GetRefundOrderStatusEvent(refundInfo *model.RefundOrder, newStatus refund.Status) *RefundOrderStatusEvent {
	return &RefundOrderStatusEvent{
		Base: Base{
			EventType: EventTypeRefundStatus,
			EventId:   fmt.Sprintf("%s_%d", refundInfo.RefundTradeNo, newStatus),
		},
		MchNo:         refundInfo.MchNo,
		OrderNo:       refundInfo.OrderNo,
		RefundNo:      refundInfo.RefundNo,
		RefundTradeNo: refundInfo.RefundTradeNo,
		AppNo:         refundInfo.AppNo,
		OldStatus:     refund.Status(refundInfo.Status),
		NewStatus:     newStatus,
		AccountNo:     refundInfo.PaymentAccountNo,
		Amount:        refundInfo.RefundAmount,
	}
}

type PaymentCallbackEvent struct {
	Base
	MchNo   string      `json:"mch_no"`
	AppNo   string      `json:"app_no"`
	OrderNo string      `json:"order_no"`
	Header  http.Header `json:"header"`
	Method  string      `json:"method"`
	Body    []byte      `json:"body"`
	Url     string      `json:"url"`
	Channel string      `json:"channel"`
	IsTest  bool        `json:"is_test"`
}

func NewPaymentCallbackEvent(req *http.Request, mchNo, appNo, Channel, orderNo string, isTest bool) *PaymentCallbackEvent {
	body, _ := utils.GetRequestBody(req)
	return &PaymentCallbackEvent{
		Base: Base{
			EventType: EventTypeCallbackPayment,
			EventId:   fmt.Sprintf("%s_%s", uuid.NewString(), Channel),
		},
		MchNo:   mchNo,
		AppNo:   appNo,
		OrderNo: orderNo,
		Header:  req.Header,
		Method:  req.Method,
		Body:    body,
		Url:     req.URL.String(),
	}
}

type RefundCallbackEvent struct {
	Base
	MchNo         string      `json:"mch_no"`
	AppNo         string      `json:"app_no"`
	RefundTradeNo string      `json:"refund_trade_no"`
	Header        http.Header `json:"header"`
	Method        string      `json:"method"`
	Body          []byte      `json:"body"`
	Url           string      `json:"url"`
	Channel       string      `json:"channel"`
}

func NewRefundCallbackEvent(req *http.Request, mchNo, appNo, Channel, RefundTradeNo string) *RefundCallbackEvent {
	body, _ := utils.GetRequestBody(req)
	return &RefundCallbackEvent{
		Base: Base{
			EventType: EventTypeCallbackRefund,
			EventId:   fmt.Sprintf("%s_%s", uuid.NewString(), Channel),
		},
		MchNo:         mchNo,
		AppNo:         appNo,
		RefundTradeNo: RefundTradeNo,
		Header:        req.Header,
		Method:        req.Method,
		Body:          body,
		Url:           req.URL.String(),
	}
}

// UserLoginLimitEvent 用户登录限制
type UserLoginLimitEvent struct {
	Base
	MchNo     string    `json:"mch_no"`
	AppNo     string    `json:"app_no"`
	Amount    int64     `json:"amount"`
	UserLimit bool      `json:"payment_limit"`
	ErrMsg    string    `json:"err_msg"`
	AccountNo string    `json:"account_no"`
	CreateAt  time.Time `json:"create_at"`
}

type PaymentNotifyRetryEvent struct {
	Base
	NotifyID    int64  `json:"notify_id"`
	MchNo       string `json:"mch_no"`
	AppNo       string `json:"app_no"`
	OrderNo     string `json:"order_no"`
	TradeNo     string `json:"trade_no"`
	NotifyCount int    `json:"notify_count"`
}

type RefundNotifyRetryEvent struct {
	Base
	NotifyID      int64  `json:"notify_id"`
	MchNo         string `json:"mch_no"`
	AppNo         string `json:"app_no"`
	RefundNo      string `json:"refund_no"`
	RefundTradeNo string `json:"refund_trade_no"`
	NotifyCount   int    `json:"notify_count"`
}
