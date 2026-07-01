package service

import (
	"context"
	"errors"
	system2 "github.com/lihongsheng/pay-gateway/model/system"
	"github.com/lihongsheng/pay-gateway/repo/system"
	"time"

	dtoSys "github.com/lihongsheng/pay-gateway/dto/system"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/payment-sdk/enum/payment"
)

type TradeOrderService interface {
	Search(ctx context.Context, req *dto.TradeSearchRequest) ([]*dto.TradeSearchResponse, error)
	Count(ctx context.Context, req *dto.TradeSearchRequest) (int64, error)
	Detail(ctx context.Context, mchNO string, appNO string, orderNo string) (*dto.TradeOrderDetail, error)
	Log(ctx context.Context, machID int64, orderNo string, user *dto.User) ([]*model.PaymentRequestLog, error)
}

type tradeOrderService struct {
	paymentOrderRepo   repo.PaymentOrderRepo
	mchRepo            system.MchRepo
	appRepo            repo.ApplicationRepo
	paymentAccountRepo repo.PaymentAccountRepo
}

func NewTradeOrderService(
	paymentOrderRepo repo.PaymentOrderRepo,
	mchRepo system.MchRepo,
	appRepo repo.ApplicationRepo,
	paymentAccountRepo repo.PaymentAccountRepo,
) TradeOrderService {
	return &tradeOrderService{
		paymentOrderRepo:   paymentOrderRepo,
		mchRepo:            mchRepo,
		appRepo:            appRepo,
		paymentAccountRepo: paymentAccountRepo,
	}
}

func (t *tradeOrderService) Search(ctx context.Context, req *dto.TradeSearchRequest) ([]*dto.TradeSearchResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	models, err := t.paymentOrderRepo.SearchModel(ctx, &dto.SearchRequest{
		OrderNo:   req.OrderNo,
		TradeNo:   req.TradeNo,
		AppNo:     req.AppNo,
		MchNo:     req.MchNo,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    req.Status,
		Page:      req.Page,
		PageSize:  req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	}
	return t.EntityToParamResponse(ctx, models)
}

func (t *tradeOrderService) EntityToParamResponse(ctx context.Context, req []*model.PaymentOrder) ([]*dto.TradeSearchResponse, error) {
	mch, app, paymentAccount, err := t.getLink(ctx, req)
	if err != nil {
		return nil, err
	}
	var res []*dto.TradeSearchResponse
	for _, v := range req {
		accountName := ""
		if paymentAccount[v.PaymentAccountNo] != nil {
			accountName = paymentAccount[v.PaymentAccountNo].Name
		}
		tmp := &dto.TradeSearchResponse{
			MchNo:          v.MchNo,
			MchName:        mch[v.MchNo].MchName,
			AppNo:          v.AppNo,
			AppName:        app[v.AppNo].AppName,
			TradeNo:        v.TradeNo,
			OrderNo:        v.OrderNo,
			AccountNo:      v.PaymentAccountNo,
			AccountName:    accountName,
			PaymentMethod:  v.PaymentMethod,
			PaymentProduct: v.PaymentProduct,
			Amount:         v.PaymentAmount,
			Status:         payment.Status(v.Status),
			CreatedAt:      v.CreatedAt,
			OrderTime:      time.Unix(v.OrderCreate, 0),
			Subject:        v.OrderSubject,
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (t *tradeOrderService) getLink(ctx context.Context, req []*model.PaymentOrder) (mch map[string]*system2.Merchant, app map[string]*model.Application, paymentAccount map[string]*model.PaymentAccount, err error) {
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
	mchs, err := t.mchRepo.Search(ctx, dtoSys.MchQueryRequest{MchNos: mchNos, Page: 1, PageSize: len(mchNos)})
	if err != nil {
		return nil, nil, nil, err
	}
	mch = map[string]*system2.Merchant{}
	if len(mchs) == 0 {
		return nil, nil, nil, errors.New("商户不存在")
	}
	for _, v := range mchs {
		mch[v.MchNo] = v
	}
	apps, err := t.appRepo.Search(ctx, &dto.ApplicationQueryRequest{AppNos: appNos, Page: 1, PageSize: len(appNos)})
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
	paymentAccounts, err := t.paymentAccountRepo.GetOnlyAccountNo(ctx, paymentAccountNos)
	paymentAccount = map[string]*model.PaymentAccount{}
	if err != nil {
		return nil, nil, nil, err
	}
	for _, v := range paymentAccounts {
		paymentAccount[v.AccountNo] = v
	}
	return mch, app, paymentAccount, nil
}

func (t *tradeOrderService) Count(ctx context.Context, req *dto.TradeSearchRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}
	c, err := t.paymentOrderRepo.Count(ctx, &dto.SearchRequest{
		OrderNo:   req.OrderNo,
		TradeNo:   req.TradeNo,
		AppNo:     req.AppNo,
		MchNo:     req.MchNo,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    req.Status,
		Page:      req.Page,
		PageSize:  req.PageSize,
	})
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (t *tradeOrderService) Detail(ctx context.Context, mchNO string, appNO string, orderNo string) (*dto.TradeOrderDetail, error) {
	if mchNO == "" || orderNo == "" || appNO == "" {
		return nil, errors.New("参数错误")
	}
	m, err := t.paymentOrderRepo.GetModel(ctx, public.QueryPaymentRequest{
		OrderNo: orderNo,
		MchNo:   mchNO,
		AppNo:   appNO,
	}, nil)
	if err != nil {
		return nil, err
	}
	detail, err := t.EntityToParamResponseDetail(ctx, []*model.PaymentOrder{m})
	if err != nil {
		return nil, err
	}
	return detail[0], nil
}

func (t *tradeOrderService) EntityToParamResponseDetail(ctx context.Context, req []*model.PaymentOrder) ([]*dto.TradeOrderDetail, error) {
	mch, app, paymentAccount, err := t.getLink(ctx, req)
	if err != nil {
		return nil, err
	}
	var res []*dto.TradeOrderDetail
	for _, v := range req {
		tmp := &dto.TradeOrderDetail{
			MchNo:          v.MchNo,
			MchName:        mch[v.MchNo].MchName,
			AppNo:          v.AppNo,
			AppName:        app[v.AppNo].AppName,
			TradeNo:        v.TradeNo,
			OrderNo:        v.OrderNo,
			AccountNo:      v.PaymentAccountNo,
			AccountName:    paymentAccount[v.PaymentAccountNo].Name,
			PaymentMethod:  v.PaymentMethod,
			PaymentProduct: v.PaymentProduct,
			OutMchTradeNo:  v.OutMchTradeNo,
			Amount:         v.PaymentAmount,
			Status:         payment.Status(v.Status),
			CreatedAt:      v.CreatedAt,
			OrderTime:      time.Unix(v.OrderCreate, 0),
			RedirectURL:    v.RedirectURL,
			NotifyUrl:      v.NotifyURL,
			ExtParams:      v.ExtParam,
			ThirdMsg:       v.ThirdMsg,
			ThirdCode:      v.ThirdCode,
			Subject:        v.OrderSubject,
			Desc:           v.OrderDesc,
			NotifyStatus:   int(v.NotifyStatus),
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (t *tradeOrderService) Log(ctx context.Context, machID int64, orderNo string, user *dto.User) ([]*model.PaymentRequestLog, error) {
	// 未实现
	return nil, errors.New("未实现")
}
