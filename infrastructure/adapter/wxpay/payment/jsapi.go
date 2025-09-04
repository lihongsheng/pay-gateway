package payment

import (
	"context"
	"time"

	"github.com/lihongsheng/pay-gateway/enum"
	"github.com/lihongsheng/pay-gateway/enum/action"
	"github.com/lihongsheng/pay-gateway/infrastructure/adapter/wxpay"
	"github.com/lihongsheng/pay-gateway/infrastructure/config"
	"github.com/lihongsheng/pay-gateway/infrastructure/driver"
	"github.com/lihongsheng/pay-gateway/infrastructure/driver/dto"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

type Jsapi struct {
	*wxpay.Api
	client jsapi.JsapiApiService
}

func NewJsApi(conf config.Config) (driver.Pay, error) {
	api, err := wxpay.InitClient(conf)
	if err != nil {
		return nil, err
	}
	svc := jsapi.JsapiApiService{Client: api.Client}
	return &Jsapi{
		Api:    api,
		client: svc,
	}, nil
}

func (j *Jsapi) Pay(ctx context.Context, req *dto.PayOrder) (*dto.PayResponse, error) {
	resp, result, err := j.client.Prepay(ctx, j.buildPayParmams(req))
	if err != nil {
		return nil, ErrorHandler(ctx, result, err, "")
	}
	if resp.PrepayId == nil || *resp.PrepayId == "" {
		return nil, ErrorHandler(ctx, result, err, "not return PrepayId")
	}
	return &dto.PayResponse{
		PaymentMethod: enum.PaymentMethod_JSAPI.String(),
		Action: dto.Action{
			Action: action.Action_Prepay.String(),
			Parameters: map[string]string{
				"prepay_id": *resp.PrepayId,
			},
			Url: "",
		},
		OrderNo:   req.Order.OrderNo,
		PayAmount: dto.Amount{},
	}, nil
}
func (j *Jsapi) buildPayParmams(req *dto.PayOrder) jsapi.PrepayRequest {
	var t *time.Time
	if req.TimeExpire > 0 {
		t = core.Time(time.Unix(req.TimeExpire, 0))
	}
	amount := &jsapi.Amount{
		Total: core.Int64(req.Order.PayAmount.Total),
	}
	if req.Order.PayAmount.Currency != "" {
		amount.Currency = core.String(req.Order.PayAmount.Currency)
	}
	resp := jsapi.PrepayRequest{
		Appid:       core.String(j.C.AppID),
		Mchid:       core.String(j.C.MchID),
		OutTradeNo:  core.String(req.Order.OrderNo),
		TimeExpire:  t,
		Attach:      core.String(req.PassbackParams),
		NotifyUrl:   core.String(req.NotifyUrl),
		Description: core.String(req.Order.Subject),
		Amount:      amount,
		Payer: &jsapi.Payer{
			Openid: core.String(req.Payer.OpenID),
		},
	}
	if req.SettleInfo != nil {
		resp.SettleInfo = &jsapi.SettleInfo{
			ProfitSharing: core.Bool(req.SettleInfo.ProfitSharing),
		}
	}
	if req.SceneInfo != nil {
		resp.SceneInfo = &jsapi.SceneInfo{
			StoreInfo: &jsapi.StoreInfo{
				Id: core.String(req.SceneInfo.StoreID),
			},
			PayerClientIp: core.String(req.SceneInfo.ClientIp),
			DeviceId:      core.String(req.SceneInfo.DeviceID),
		}
	}
	return resp
}
func (j *Jsapi) Query(ctx context.Context, req dto.Query) (*dto.PayDetail, error) {

	return nil, nil
}

func (j *Jsapi) Close() {

}
