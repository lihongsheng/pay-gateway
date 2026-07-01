package service

import (
	"context"
	"errors"
	system2 "github.com/lihongsheng/pay-gateway/model/system"
	"github.com/lihongsheng/pay-gateway/repo/system"
	"time"

	dtoSys "github.com/lihongsheng/pay-gateway/dto/system"
	"github.com/lihongsheng/pay-gateway/global"
	event2 "github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"github.com/lihongsheng/payment-sdk/enum/refund"
	"gorm.io/gorm"
)

type TradeStatisticsService interface {
	HandlePaymentEvent(ctx context.Context, event *event2.PaymentOrderStatusEvent, tag string) error
	HandleRefundEvent(ctx context.Context, event *event2.RefundOrderStatusEvent, tag string) error
	CountGroup(ctx context.Context, start, end time.Time) ([]*dto.TradeStatistic, error)
	CountGroupMch(ctx context.Context, mchNo string, start, end time.Time) ([]*dto.TradeStatistic, error)
	CountGroupMchApp(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*dto.TradeStatistic, error)
	CountGroupMchAppAccount(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*dto.TradeStatistic, error)
	SearchDashboard(ctx context.Context, req *dto.MchTradeStatisticRequest) ([]*dto.TradeStatistic, error)
	Count(ctx context.Context, req *dto.MchTradeStatisticRequest) (int64, error)
	AllOrderStatistics(ctx context.Context, event *event2.PaymentOrderStatusEvent, tag string) error
	MchAllOrderStatistics(ctx context.Context, event *event2.PaymentOrderStatusEvent, tag string) error
	GetAllRequestOrder(ctx context.Context, mchNo string, start, end time.Time) (int64, error)
}
type tradeStatisticsService struct {
	mchRepo             system.MchRepo
	appRepo             repo.ApplicationRepo
	paymentAccountRepo  repo.PaymentAccountRepo
	paymentOrderRepo    repo.PaymentOrderRepo
	tradeStatisticsRepo repo.TradeStaticsRepo
	eventRecordRepo     repo.EventRecordRepo
	statisticsRepo      repo.Statistics
}

func NewTradeStatisticsService(
	mchRepo system.MchRepo,
	appRepo repo.ApplicationRepo,
	paymentAccountRepo repo.PaymentAccountRepo,
	paymentOrderRepo repo.PaymentOrderRepo,
	tradeStatisticsRepo repo.TradeStaticsRepo,
	eventRecordRepo repo.EventRecordRepo,
	statisticsRepo repo.Statistics,
) TradeStatisticsService {
	return &tradeStatisticsService{
		mchRepo:             mchRepo,
		appRepo:             appRepo,
		paymentAccountRepo:  paymentAccountRepo,
		paymentOrderRepo:    paymentOrderRepo,
		tradeStatisticsRepo: tradeStatisticsRepo,
		eventRecordRepo:     eventRecordRepo,
		statisticsRepo:      statisticsRepo,
	}
}

func (t *tradeStatisticsService) GetAllRequestOrder(ctx context.Context, mchNo string, start, end time.Time) (int64, error) {
	mch := enum.StatisticsAll
	if mchNo != "" {
		mch = mchNo
	}
	totalRequest, err := t.statisticsRepo.CountGroupMch(ctx, mch, enum.StatisticsAll, enum.StatisticsTag, enum.StatisticsKey, start, end)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if totalRequest == nil {
		return 0, nil
	}
	return totalRequest.Statistic, nil
}

func (t *tradeStatisticsService) AllOrderStatistics(ctx context.Context, event *event2.PaymentOrderStatusEvent, tag string) error {
	if !(event.NewStatus == payment.Status_Created || event.NewStatus == payment.Status_TempFailed || event.NewStatus == payment.Status_Failed || event.NewStatus == payment.Status_Pending) {
		return nil
	}
	statics := t.buildAllOrderStatistic(event)
	if statics == nil {
		return nil
	}
	exits, _ := t.eventRecordRepo.Get(ctx, event.EventId, event.EventType, tag)
	if exits != nil {
		return nil
	}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		er := &model.EventProcessRecord{
			EventID:     event.EventId,
			EventType:   event.EventType,
			ProcessedAt: time.Now().Unix(),
			CreatedAt:   time.Now(),
			Consumer:    tag,
		}
		// 唯一键防重复统计
		if err := t.eventRecordRepo.Save(ctx, er, tx); err != nil {
			return err
		}
		return t.statisticsRepo.Save(ctx, statics, tx)
	})
	return err
}

func (t *tradeStatisticsService) buildAllOrderStatistic(event *event2.PaymentOrderStatusEvent) *model.Statistic {
	statics := &model.Statistic{
		AppNo:         enum.StatisticsAll,
		StatisticDate: event.CreateAt,
		MchNo:         enum.StatisticsAll,
		StatisticKey:  enum.StatisticsKey,
		StatisticTag:  enum.StatisticsTag,
	}
	handler := false
	switch event.NewStatus {
	case payment.Status_Created, payment.Status_Pending, payment.Status_TempFailed, payment.Status_Failed:
		if event.ID == 0 {
			statics.Statistic = 1
			handler = true
		}
	}
	if handler {
		return statics
	}
	return nil
}
func (t *tradeStatisticsService) MchAllOrderStatistics(ctx context.Context, event *event2.PaymentOrderStatusEvent, tag string) error {
	if !(event.NewStatus == payment.Status_Created || event.NewStatus == payment.Status_TempFailed || event.NewStatus == payment.Status_Failed || event.NewStatus == payment.Status_Pending) {
		return nil
	}
	statics := t.buildMchOrderStatistic(event)
	if statics == nil {
		return nil
	}
	exits, _ := t.eventRecordRepo.Get(ctx, event.EventId, event.EventType, tag)
	if exits != nil {
		return nil
	}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		er := &model.EventProcessRecord{
			EventID:     event.EventId,
			EventType:   event.EventType,
			ProcessedAt: time.Now().Unix(),
			CreatedAt:   time.Now(),
			Consumer:    tag,
		}
		// 唯一键防重复统计
		if err := t.eventRecordRepo.Save(ctx, er, tx); err != nil {
			return err
		}
		return t.statisticsRepo.Save(ctx, statics, tx)
	})
	return err
}

func (t *tradeStatisticsService) buildMchOrderStatistic(event *event2.PaymentOrderStatusEvent) *model.Statistic {
	statics := &model.Statistic{
		AppNo:         enum.StatisticsAll,
		StatisticDate: event.CreateAt,
		MchNo:         event.MchNo,
		StatisticKey:  enum.StatisticsKey,
		StatisticTag:  enum.StatisticsTag,
	}
	handler := false
	switch event.NewStatus {
	case payment.Status_Created, payment.Status_Pending, payment.Status_TempFailed, payment.Status_Failed:
		if event.ID == 0 {
			statics.Statistic = 1
			handler = true
		}
	}
	if handler {
		return statics
	}
	return nil
}

func (t *tradeStatisticsService) Count(ctx context.Context, req *dto.MchTradeStatisticRequest) (int64, error) {
	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		return 0, errors.New("start or end time can not be empty")
	}

	return t.tradeStatisticsRepo.Count(ctx, req)
}

func (t *tradeStatisticsService) SearchDashboard(ctx context.Context, req *dto.MchTradeStatisticRequest) ([]*dto.TradeStatistic, error) {
	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		return nil, errors.New("start or end time can not be empty")
	}

	models, err := t.tradeStatisticsRepo.SearchDashboard(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return []*dto.TradeStatistic{}, nil
	}
	return t.buildTradeStatistic(ctx, models)
}
func (t *tradeStatisticsService) CountGroup(ctx context.Context, start, end time.Time) ([]*dto.TradeStatistic, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start or end time can not be empty")
	}

	models, err := t.tradeStatisticsRepo.CountGroup(ctx, start, end)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return []*dto.TradeStatistic{}, nil
	}

	return t.buildTradeStatistic(ctx, models)
}
func (t *tradeStatisticsService) CountGroupMch(ctx context.Context, mchNo string, start, end time.Time) ([]*dto.TradeStatistic, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start or end time can not be empty")
	}
	if mchNo == "" {
		return nil, errors.New("mchNo can not be empty")
	}
	models, err := t.tradeStatisticsRepo.CountGroupMch(ctx, mchNo, start, end)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return []*dto.TradeStatistic{}, nil
	}
	return t.buildTradeStatistic(ctx, models)
}
func (t *tradeStatisticsService) CountGroupMchApp(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*dto.TradeStatistic, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start or end time can not be empty")
	}
	if mchNo == "" {
		return nil, errors.New("mchNo can not be empty")
	}
	if appNo == "" {
		return nil, errors.New("appNo can not be empty")
	}
	models, err := t.tradeStatisticsRepo.CountGroupMchApp(ctx, mchNo, appNo, start, end)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return []*dto.TradeStatistic{}, nil
	}
	return t.buildTradeStatistic(ctx, models)
}
func (t *tradeStatisticsService) CountGroupMchAppAccount(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*dto.TradeStatistic, error) {
	models, err := t.tradeStatisticsRepo.CountGroupMchAppAccount(ctx, mchNo, appNo, start, end)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return []*dto.TradeStatistic{}, nil
	}
	return t.buildTradeStatistic(ctx, models)
}

func (t *tradeStatisticsService) buildTradeStatistic(ctx context.Context, models []*model.TradeStatistic) ([]*dto.TradeStatistic, error) {
	var mchNos = []string{}
	var mchNosExists = map[string]struct{}{}
	var appNos = []string{}
	var accountNos = []string{}
	for _, v := range models {
		if v.MchNo != "" {
			if _, ok := mchNosExists[v.MchNo]; !ok {
				mchNos = append(mchNos, v.MchNo)
				mchNosExists[v.MchNo] = struct{}{}
			}
		}
		if v.AppNo != "" {
			if _, ok := mchNosExists[v.AppNo]; !ok {
				appNos = append(appNos, v.AppNo)
				mchNosExists[v.AppNo] = struct{}{}
			}
		}
		if v.AccountNo != "" {
			if _, ok := mchNosExists[v.AccountNo]; !ok {
				accountNos = append(accountNos, v.AccountNo)
				mchNosExists[v.AccountNo] = struct{}{}
			}
		}
	}
	var mchMap = map[string]*system2.Merchant{}
	var mchAppMap = map[string]*model.Application{}
	var paymentAccountMap = map[string]*model.PaymentAccount{}
	if len(mchNos) > 0 {
		mchInfos, err := t.mchRepo.Search(ctx, dtoSys.MchQueryRequest{MchNos: mchNos, Page: 1, PageSize: len(mchNos)})
		if err != nil {
			return nil, err
		}
		for _, v := range mchInfos {
			mchMap[v.MchNo] = v
		}
	}
	if len(appNos) > 0 {
		mchApps, err := t.appRepo.Search(ctx, &dto.ApplicationQueryRequest{AppNos: appNos, Page: 1, PageSize: len(appNos)})
		if err != nil {
			return nil, err
		}
		for _, v := range mchApps {
			mchAppMap[v.AppNo] = v
		}
	}
	if len(accountNos) > 0 {
		paymentAccounts, err := t.paymentAccountRepo.GetOnlyAccountNo(ctx, accountNos)
		if err != nil {
			return nil, err
		}
		for _, v := range paymentAccounts {
			paymentAccountMap[v.AccountNo] = v
		}
	}
	var tradeStatistics = make([]*dto.TradeStatistic, 0, len(models))
	for _, v := range models {
		mchName := ""
		if mchMap[v.MchNo] != nil {
			mchName = mchMap[v.MchNo].MchName
		}
		appName := ""
		if mchAppMap[v.AppNo] != nil {
			appName = mchAppMap[v.AppNo].AppName
		}
		accountName := ""
		if paymentAccountMap[v.AccountNo] != nil {
			accountName = paymentAccountMap[v.AccountNo].Name
		}
		tradeStatistics = append(tradeStatistics, &dto.TradeStatistic{
			AppName:       appName,
			MchName:       mchName,
			AccountName:   accountName,
			AppNo:         v.AppNo,
			AccountNo:     v.AccountNo,
			MchNo:         v.MchNo,
			StatisticDate: v.StatisticDate.Format("2006-01-02"),
			TotalOrder:    v.TotalOrder,
			SuccessOrder:  v.SuccessOrder,
			TotalAmount:   v.TotalAmount,
			SuccessAmount: v.SuccessAmount,
			RefundOrder:   v.RefundOrder,
			RefundAmount:  v.RefundAmount,
		})
	}
	return tradeStatistics, nil
}

func (t *tradeStatisticsService) HandlePaymentEvent(ctx context.Context, event *event2.PaymentOrderStatusEvent, tag string) error {
	if !(event.NewStatus == payment.Status_Success || event.NewStatus == payment.Status_Failed || event.NewStatus == payment.Status_Pending) {
		return nil
	}
	statics := t.buildPaymentTradeStatistic(event)
	if statics == nil {
		return nil
	}
	exits, _ := t.eventRecordRepo.Get(ctx, event.EventId, event.EventType, tag)
	if exits != nil {
		return nil
	}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		er := &model.EventProcessRecord{
			EventID:     event.EventId,
			EventType:   event.EventType,
			ProcessedAt: time.Now().Unix(),
			CreatedAt:   time.Now(),
			Consumer:    tag,
		}
		// 唯一键防重复统计
		if err := t.eventRecordRepo.Save(ctx, er, tx); err != nil {
			return err
		}
		return t.tradeStatisticsRepo.Save(ctx, statics, tx)
	})
	return err
}

func (t *tradeStatisticsService) buildPaymentTradeStatistic(event *event2.PaymentOrderStatusEvent) *model.TradeStatistic {
	statics := &model.TradeStatistic{
		AppNo:         event.AppNo,
		AccountNo:     event.AccountNo,
		StatisticDate: event.CreateAt,
		MchNo:         event.MchNo,
	}
	handler := false
	switch event.NewStatus {
	case payment.Status_Pending:
		statics.TotalOrder = 1
		statics.TotalAmount = event.Amount
		handler = true
	case payment.Status_Success:
		statics.SuccessOrder = 1
		statics.SuccessAmount = event.Amount
		handler = true
	case payment.Status_Failed:
		if event.OldStatus == payment.Status_Created || event.OldStatus == payment.Status_TempFailed {
			statics.TotalOrder = 1
			statics.TotalAmount = event.Amount
			handler = true
		}
	}
	if handler {
		return statics
	}
	return nil
}

func (t *tradeStatisticsService) HandleRefundEvent(ctx context.Context, event *event2.RefundOrderStatusEvent, tag string) error {
	if !(event.NewStatus == refund.Status_Success) {
		return nil
	}
	exits, _ := t.eventRecordRepo.Get(ctx, event.EventId, event.EventType, tag)
	if exits != nil {
		return nil
	}
	statics, err := t.buildRefundTradeStatistic(ctx, event)
	if err != nil {
		return err
	}
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		er := &model.EventProcessRecord{
			EventID:     event.EventId,
			EventType:   event.EventType,
			ProcessedAt: time.Now().Unix(),
			CreatedAt:   time.Now(),
			Consumer:    tag,
		}
		// 唯一键防重复统计
		if err := t.eventRecordRepo.Save(ctx, er, tx); err != nil {
			return err
		}
		return t.tradeStatisticsRepo.Save(ctx, statics, tx)
	})
	return err
}

func (t *tradeStatisticsService) buildRefundTradeStatistic(ctx context.Context, event *event2.RefundOrderStatusEvent) (*model.TradeStatistic, error) {
	payOrder, err := t.paymentOrderRepo.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		OrderNo: event.OrderNo,
		AppNo:   event.AppNo,
		TradeNo: "",
		MchNo:   event.MchNo,
	})
	if err != nil {
		return nil, err
	}
	statics := &model.TradeStatistic{
		AppNo:         event.AppNo,
		AccountNo:     event.AccountNo,
		StatisticDate: payOrder.CreatedAt,
		MchNo:         event.MchNo,
	}
	switch event.NewStatus {
	case refund.Status_Success:
		statics.RefundOrder = 1
		statics.RefundAmount = event.Amount
	}
	return statics, nil
}
