package enum

import "time"

const MetaSign = "X-Sign"
const MetaTimestamp = "X-Timestamp"
const MetaNonce = "X-Nonce"
const MetaAppNo = "X-App-NO"
const MetaSignMethod = "X-Method"

const TestPath = "/payment/test.html"
const PaymentBaseIndexPath = "/payment/web/public/payment.html"
const PaymentSuccessPath = "/payment/web/public/success.html"
const PaymentTestPath = "/payment/web/public/test.html"
const CashierPayment = "/payment/web/public/cashier.html"

// const AuthPaymentIndexPath = PaymentBaseIndexPath + "#/oauth/%s"
const H5Action = "action"
const H5ActionHub = "hub"
const H5ActionAuth = "oauth"
const H5ActionToken = "action_token"

// 支付回调

const PaymentNotify = "/public/v1/notify/payment/%s/%s/%s/%s"
const PaymentTestNotify = "/public/v1/notify/test/payment/%s/%s/%s/%s"
const RefundNotify = "/public/v1/notify/refund/%s/%s/%s/%s"

const (
	CacheLockPayment = "paymentLock:%s:%s"
	CacheLockCashier = "cashierLock:%s:%s"
	CacheLockRefund  = "refundLock:%s:%s"
)

const (
	PaymentLockExpire = time.Minute      // 支付锁过期时间
	DefaultExpireTime = 10 * time.Minute // 默认订单过期时间
	SaveOrderTimeout  = time.Minute      // 保存订单超时时间
)
