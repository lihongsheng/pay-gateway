package enum

import "time"

type MchStatus int

const (
	Status_Down MchStatus = 2
	Status_Up   MchStatus = 1
)

func (m MchStatus) Disable() bool {
	return m == Status_Down
}

const (
	CustomerDomainStatus_Disable = 0
	CustomerDomainStatus_Enable  = 1
)

const (
	MultiChannel_Disable = 0
	MultiChannel_Enable  = 1
)

type RefundFrom int

const (
	RefundFrom_API   RefundFrom = 1
	RefundFrom_ADMIN RefundFrom = 2
)

var RefundFromDesc = map[RefundFrom]string{
	RefundFrom_API:   "API",
	RefundFrom_ADMIN: "后台",
}

type NotifyType int

const (
	NotifyRetryMax  = 5
	NotifyRetryStep = 5 * time.Second
	//  支付
	NotifyType_Payment NotifyType = 1
	// 退款
	NotifyType_Refund NotifyType = 2
	// 转账
	NotifyType_Transfer NotifyType = 3
)

type NotifyStatus int

const (
	NotifyStatus_Init    NotifyStatus = 1
	NotifyStatus_Success NotifyStatus = 2
	NotifyStatus_Fail    NotifyStatus = 3
)

const (
	MaxRetryPaymentLimit   = 3
	MaxRetryPaymentDefault = 1
)

const (
	StatisticsAll = "all"
	StatisticsKey = "TotalRequestOrder"
	StatisticsTag = "TotalCount"
)
