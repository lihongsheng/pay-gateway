package dto

type Amount struct {
	// 分单位，一元为100
	Total    int64
	Currency string
}

type Payer struct {
	// 用户在商户的标识
	Unionid string
	// 用户在商户下某个应用的标识
	OpenID string
}

type Goods struct {
	// 商品名称
	Name string
	// 商品SKU
	Sku string
	// 商品价格，单位分
	Price int64
	// 商品数量
	Quantity int
	//
	Desc string
}

type Order struct {
	OrderNo string
	// 订单金额
	Amount Amount
	// 实付金额
	PayAmount Amount
	// 订单商品
	Goods []Goods
	// 订单名称
	Name string
	// 订单描述
	Desc string
	//
}

type ApplicationInfo struct {
	AppID          string
	Url            string
	IOSPackage     string
	AndroidPackage string
}

type H5 struct {
	ApplicationInfo ApplicationInfo
}

type SceneInfo struct {
	// 客户端IP
	ClientIp string
	// 设备ID
	DeviceID string
	// 门店ID
	StoreID string
}

type PayOrder struct {
	Order  Order
	Amount Amount
	Payer  Payer
	// 支付跳转地址
	RedirectUrl string
	// 订单超时时间
	TimeExpire int64
	// 支付回调地址
	NotifyUrl string
	// 透传参数 如果请求时传递了该参数，异步通知时将该参数原样返回。
	PassbackParams string
}

type PayResponse struct {
	// 订单号
	OrderNo string
	// 交易单号
	TradeNo string
	// 实付金额
	PayAmount Amount
	// Pending | Success | Fail
	Status int
	// H5 | JSAPI | NATIVE | APP | 扫码
	PaymentMethod string
}

type SettleInfo struct {
	// 是否分账
	ProfitSharing bool
}

type Action struct {
	// Redirect | Qrcode | Prepay
	Action     string
	Method     string
	Url        string
	Parameters map[string]string
}

type Query struct {
	// 订单号
	OrderNo string
	// 交易单号
	TradeNo string
}

type Detail struct {
	// 订单优惠金额
	DiscountAmount Amount
}
