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

var RefundCallback = &refundCallback{}

type refundCallback struct{}

func (p *refundCallback) GroupName() string {
	return "refundCallback"
}

func (r *refundCallback) Message(ctx context.Context, msg string) error {
	ctx, l := log.NewCtx(ctx)
	defer func() {
		if err := recover(); err != nil {
			l.Error("refundCallback", zap.Any("panic message", err))
		}
	}()
	l.Info("refundCallback", zap.String("msg", msg))
	var eventDetail event2.RefundCallbackEvent
	err := json.Unmarshal([]byte(msg), &eventDetail)
	if err != nil {
		l.Error("refundCallback", zap.Error(err))
	}
	u, _ := url.Parse(eventDetail.Url)
	var req = &http.Request{
		Header: eventDetail.Header,
		Body:   ioutil.NopCloser(bytes.NewBuffer(eventDetail.Body)),
		Method: eventDetail.Method,
		URL:    u,
	}
	_, err = enter.ServiceApiApp.RefundService.Callback(ctx, req, eventDetail.AppNo, eventDetail.RefundTradeNo)
	if err != nil {
		l.Error("paymentCallbackErr", zap.Error(err))
	}
	return nil
}

func (d *refundCallback) NotifyClose() {

}
