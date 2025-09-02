package payment

import (
	"context"

	"github.com/lihongsheng/pay-gateway/adapter/wxpay"
	"github.com/lihongsheng/pay-gateway/config"
	"github.com/lihongsheng/pay-gateway/driver"
	"github.com/lihongsheng/pay-gateway/driver/dto"
)

type Jsapi struct {
	*wxpay.Api
}

func NewJsApi(conf config.Config) (driver.Pay, error) {
	api, err := wxpay.InitClient(conf)
	if err != nil {
		return nil, err
	}
	return &Jsapi{
		api,
	}
}

func (j *Jsapi) Pay(ctx context.Context, req *dto.PayOrder) (*dto.PayResponse, error) {

}
