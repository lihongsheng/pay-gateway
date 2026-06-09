package errors

// ===================== 全局通用错误码 =====================
const (
	ErrCodeInvalidParam       = 10001 // 参数错误
	ErrCodeEmptyParam         = 10002 // 必填参数为空
	ErrCodeInvalidRequest     = 10003 // 非法请求
	ErrCodeSignValidateFailed = 10004
)

// ===================== 支付模块专属错误码(重点) =====================
const (
	ErrPayOrderNotExist     = 20001 // 支付订单不存在
	ErrPayOrderStatusErr    = 20002 // 订单状态异常(已支付/已关闭)
	ErrPayAmountErr         = 20003 // 支付金额错误
	ErrPayChannelNotConfig  = 20004 // 支付渠道未配置(微信/支付宝未开通)
	ErrPayChannelNotSupport = 20005 // 当前渠道不支持该支付方式
	ErrPayCreateFail        = 20006 // 创建支付订单失败
	ErrPayChannelNotExist   = 20007
	ErrPayPending           = 20008
	ErrPayOrderNotPending   = 20009
	ErrAppNotExist          = 20010
	ErrPayThirdParty        = 20101 // 支付接口异常
	ErrPayNotifySignErr     = 20301 // 支付回调签名验证失败
	ErrPayFail              = 20401

	// ===================== 退款码对应默认描述信息 =====================
	ErrRefundOrderNotExist    = 30001
	ErrRefundOrderStatusErr   = 30002
	ErrRefundAmountErr        = 30003
	ErrRefundChannelNotConfig = 30004
	ErrRefundCreateFail       = 30005
	ErrRefundThirdParty       = 30006
	ErrRefundNotifySignErr    = 30008
	ErrRefundPending          = 30009
	// ===================== 三方支付渠道错误(重点) =====================
	ErrUserLimit    = 41101 // 用户登录限制
	ErrPaymentLimit = 41202 // 渠道支付限制f
	ErrRetryFail    = 41203
)

// ===================== 系统级错误码 =====================
const (
	ErrCodeSysError = 50000 // 系统内部错误
)

// ===================== 错误码对应默认描述信息 =====================
var ErrMsgMap = map[int]string{
	ErrCodeInvalidParam:       "请求参数格式错误，请核对后重试",
	ErrCodeEmptyParam:         "必填参数不能为空",
	ErrCodeInvalidRequest:     "非法请求，禁止访问",
	ErrCodeSignValidateFailed: "签名验证失败",

	ErrPayOrderNotExist:     "支付订单不存在",
	ErrPayOrderStatusErr:    "订单状态异常，无法支付",
	ErrPayAmountErr:         "支付金额错误，请核对订单金额",
	ErrPayChannelNotConfig:  "该支付渠道未配置，请联系管理员开通",
	ErrPayChannelNotSupport: "当前终端不支持该支付方式",
	ErrPayCreateFail:        "创建支付订单失败，请稍后重试",
	ErrPayThirdParty:        "支付接口调用失败",
	ErrPayChannelNotExist:   "支付渠道不存在",
	ErrPayNotifySignErr:     "支付回调签名验证失败，非法请求",
	ErrPayPending:           "支付订单处理中，请稍后重试",
	ErrPayOrderNotPending:   "支付订单已完成",
	ErrAppNotExist:          "应用不存在",
	ErrPayFail:              "支付失败",

	ErrRefundOrderNotExist:    "退款订单不存在",
	ErrRefundOrderStatusErr:   "订单状态异常，无法退款",
	ErrRefundAmountErr:        "退款金额错误",
	ErrRefundChannelNotConfig: "该支付渠道未配置，请联系管理员开通",
	ErrRefundCreateFail:       "创建退款订单失败，请稍后重试",
	ErrRefundThirdParty:       "退款接口调用失败",
	ErrRefundNotifySignErr:    "退款回调签名验证失败，非法",
	ErrRefundPending:          "退款订单处理中，请稍后重试",
	ErrUserLimit:              "用户登录限制",
	ErrPaymentLimit:           "支付被限制",
	ErrRetryFail:              "支付失败,正在重试中",

	ErrCodeSysError: "系统内部错误，请联系管理员",
}

// GetErrMsg 根据错误码获取默认描述信息
func GetErrMsg(code int) string {
	if msg, ok := ErrMsgMap[code]; ok {
		return msg
	}
	return "未知错误"
}
