package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"github.com/lihongsheng/payment-sdk/enum/refund"
	"go.uber.org/zap"
)

type TradeRefundService interface {
	Refund(ctx context.Context, req *admin.Refund) (*public.RefundResponse, error)
	Detail(ctx context.Context, mchNO string, appNO string, refundTradeNo string) (*admin.RefundOrderDetail, error)
	Search(ctx context.Context, req *admin.RefundSearchRequest) ([]*admin.RefundOrderDetail, error)
	Count(ctx context.Context, req *admin.RefundSearchRequest) (int64, error)
	AvailableRefundAmount(ctx context.Context, mchNO string, appNO string, orderNo string) (int64, error)
}

type tradeRefundService struct {
	svc    *svc.ServiceContext
	domain *domain.ServiceGroup
}

func NewTradeRefundService(svc *svc.ServiceContext, domain *domain.ServiceGroup) TradeRefundService {
	return &tradeRefundService{
		svc:    svc,
		domain: domain,
	}
}

func (t *tradeRefundService) AvailableRefundAmount(ctx context.Context, mchNO string, appNO string, orderNo string) (int64, error) {
	order, err := t.svc.PaymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: orderNo,
		AppNo:   appNO,
		MchNo:   mchNO,
	})
	if err != nil {
		return 0, err
	}
	if !(order.Status == payment.Status_Success || order.Status == payment.Status_Refund) {
		return 0, errors.New("此订单不可退款")
	}
	amount, err := t.svc.RefundRepo.CountAmount(ctx, mchNO, appNO, orderNo, []refund.Status{refund.Status_Success})
	if err != nil {
		return 0, err
	}
	return order.PaymentAmount - amount, nil
}

func (t *tradeRefundService) Refund(ctx context.Context, req *admin.Refund) (*public.RefundResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	appInfo, err := t.svc.AppRepo.GetByAppNoFormCache(ctx, req.AppNo)
	if err != nil {
		return nil, err
	}
	if appInfo.MchNo != req.MchNo {
		return nil, fmt.Errorf("商户号不匹配")
	}
	refundTradeNo, _ := utils.GenTradeNo(req.MchNo)
	refundReq := public.RefundRequest{
		OrderNo:  req.OrderNo,
		RefundNo: refundTradeNo,
		Amount: public.Amount{
			Total: req.Amount,
		},
		Reason:         req.Reason,
		PassBackParams: "",
		NotifyUrl:      "",
		AppNo:          req.AppNo,
		MchNo:          req.MchNo,
		RefundFrom:     enum.RefundFrom_ADMIN,
	}
	resp, err := t.domain.RefundService.Refund(ctx, &refundReq, appInfo)
	if err != nil {
		log.WithCtx(ctx).Error("退款失败: ", zap.Error(err), zap.Any("refundReq", refundReq))
		return nil, err
	}
	if resp.Status == int64(refund.Status_Failed) || resp.Status == int64(refund.Status_Abnormal) {
		return nil, errors.New("退款失败: " + resp.ThirdMsg)
	}
	return buildRefundResponse(resp), nil
}

func (t *tradeRefundService) Detail(ctx context.Context, mchNO string, appNO string, refundTradeNo string) (*admin.RefundOrderDetail, error) {
	if mchNO == "" || refundTradeNo == "" {
		return nil, errors.New("参数错误")
	}
	m, err := t.svc.RefundRepo.GetRefund(ctx, mchNO, appNO, refundTradeNo)
	if err != nil {
		return nil, err
	}
	details, err := t.ModelToParamResponse(ctx, []*model.RefundOrder{m})
	return details[0], err
}

func (t *tradeRefundService) ModelToParamResponse(ctx context.Context, req []*model.RefundOrder) ([]*admin.RefundOrderDetail, error) {
	mch, app, paymentAccount, err := t.getLink(ctx, req)
	if err != nil {
		return nil, err
	}
	var res []*admin.RefundOrderDetail
	for _, v := range req {
		tmp := &admin.RefundOrderDetail{
			RefundNo:      v.RefundNo,
			RefundTradeNo: v.RefundTradeNo,
			MchNo:         v.MchNo,
			MchName:       mch[v.MchNo].MchName,
			AppNo:         v.AppNo,
			AppName:       app[v.AppNo].AppName,
			TradeNo:       v.TradeNo,
			OrderNo:       v.OrderNo,
			AccountNo:     v.PaymentAccountNo,
			AccountName:   paymentAccount[v.PaymentAccountNo].Name,
			OutMchTradeNo: v.OutMchTradeNo,
			RefundAmount:  v.RefundAmount,
			PaymentAmount: v.PaymentAmount,
			Status:        refund.Status(v.Status),
			NotifyUrl:     v.NotifyURL,
			ExtParams:     v.ExtParam,
			Reason:        v.Resaon,
			NotifyStatus:  int(v.NotifyStatus),
			CreatedAt:     v.CreatedAt,
			SuccessTime:   v.SuccessTime,
			ThirdMsg:      v.ThirdMsg,
			RefundFrom:    enum.RefundFromDesc[enum.RefundFrom(v.RefundFrom)],
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (t *tradeRefundService) getLink(ctx context.Context, req []*model.RefundOrder) (mch map[string]*model.Merchant, app map[string]*model.Application, paymentAccount map[string]*model.PaymentAccount, err error) {
	var mchNos = []string{}
	var existMchNos = map[string]string{}
	var appNos = []string{}
	var existAppNos = map[string]string{}
	var paymentAccountNos = []string{}
	var existPaymentAccountNos = map[string]string{}
	for _, v := range req {
		if _, ok := existMchNos[v.MchNo]; !ok {
			mchNos = append(mchNos, v.MchNo)
			existMchNos[v.MchNo] = v.MchNo
		}
		if _, ok := existAppNos[v.AppNo]; !ok {
			appNos = append(appNos, v.AppNo)
			existAppNos[v.AppNo] = v.AppNo
		}
		if _, ok := existPaymentAccountNos[v.PaymentAccountNo]; !ok {
			paymentAccountNos = append(paymentAccountNos, v.PaymentAccountNo)
			existPaymentAccountNos[v.PaymentAccountNo] = v.PaymentAccountNo
		}
	}
	mchs, err := t.svc.MchRepo.Search(ctx, admin.MchQueryRequest{MchNos: mchNos, Page: 1, PageSize: len(mchNos)})
	if err != nil {
		return nil, nil, nil, err
	}
	mch = map[string]*model.Merchant{}
	if len(mchs) == 0 {
		return nil, nil, nil, errors.New("商户不存在")
	}
	for _, v := range mchs {
		mch[v.MchNo] = v
	}
	apps, err := t.svc.AppRepo.Search(ctx, &admin.ApplicationQueryRequest{AppNos: appNos, Page: 1, PageSize: len(appNos)})
	if err != nil {
		return nil, nil, nil, err
	}
	app = map[string]*model.Application{}
	if len(apps) == 0 {
		return nil, nil, nil, errors.New("应用不存在")
	}
	for _, v := range apps {
		app[v.AppNo] = v
	}
	paymentAccounts, err := t.svc.PaymentAccountRepo.GetOnlyAccountNo(ctx, paymentAccountNos)
	paymentAccount = map[string]*model.PaymentAccount{}
	if err != nil {
		return nil, nil, nil, err
	}
	if len(paymentAccounts) == 0 {
		return nil, nil, nil, errors.New("支付账户不存在")
	}
	for _, v := range paymentAccounts {
		paymentAccount[v.AccountNo] = v
	}
	return mch, app, paymentAccount, nil
}

func (t *tradeRefundService) Search(ctx context.Context, req *admin.RefundSearchRequest) ([]*admin.RefundOrderDetail, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	m, err := t.svc.RefundRepo.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, nil
	}
	if len(m) == 0 {
		return nil, nil
	}
	return t.ModelToParamResponse(ctx, m)
}

func (t *tradeRefundService) Count(ctx context.Context, req *admin.RefundSearchRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}
	return t.svc.RefundRepo.Count(ctx, req)
}
