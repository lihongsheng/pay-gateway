package driver

type Refund interface {
	Refund(orderId int64, amount int64) error
	Query()
}
