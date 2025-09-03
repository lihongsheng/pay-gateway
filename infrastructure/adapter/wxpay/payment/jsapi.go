package payment

import (
	"context"
	"errors"
	"time"

	errors2 "github.com/lihongsheng/pay-gateway/errors"

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
	resp, _, err := j.client.Prepay(ctx, j.buildPayParmams(req))
	if err != nil {
		var apiErr *core.APIError
		if errors.As(err, &apiErr) {
			return nil, errors2.ErrorSystemError("ddd").WithCause(err)
		}
		return nil, errors2.ErrorSystemError("wxpay is error").WithCause(err)
	}
	if resp.PrepayId == nil {
		return nil, errors2.ErrorSystemError("not return PrepayId")
	}
	return &dto.PayResponse{
		Action: dto.Action{
			Action:     "wxpay",
			Method:     "POST",
			Parameters: map[string]string{"appId": j.C.AppID, "timeStamp": time.Now().Format("2006-01-02 15:04:05"), "nonceStr": j.C.MchID, "package": "prepay_id=" + *resp.PrepayId, "signType": "HMAC-SHA256"},
			Url:        "https://api.mch.weixin.qq.com/pay/unifiedorder",
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
