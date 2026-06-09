package consumer

// 处理微信或者支付宝的回调
import (
	"bytes"
	"context"
	"encoding/json"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
)

var PaymentCallback = new(paymentCallback)

type paymentCallback struct{}

func (p *paymentCallback) GroupName() string {
	return "paymentCallback"
}
func (p *paymentCallback) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("paymentCallback", zap.Any("panic message", err))
		}
	}()
	l.Info("paymentCallback", zap.String("msg", msg))
	var eventDetail event2.PaymentCallbackEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("paymentCallbackErr", zap.Error(err))
	}
	u, _ := url.Parse(eventDetail.Url)
	var req = &http.Request{
		Header: eventDetail.Header,
		Body:   ioutil.NopCloser(bytes.NewBuffer(eventDetail.Body)),
		Method: eventDetail.Method,
		URL:    u,
	}
	_, err = enter.ServiceApiApp.PaymentService.Callback(ctx, req, eventDetail.AppNo, eventDetail.OrderNo, eventDetail.IsTest)
	if err != nil {
		l.Error("paymentCallbackErr", zap.Error(err))
	}
	return nil
}

func (p *paymentCallback) NotifyClose() {

}
