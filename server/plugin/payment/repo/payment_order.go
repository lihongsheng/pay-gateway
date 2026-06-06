package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/payment-sdk/enum/payment"
	"gorm.io/gorm"
	"time"
)

type PaymentOrderRepo interface {
	GetOrderWithCache(ctx context.Context, req public.QueryPaymentRequest) (*entity.PaymentOrder, error)
	Save(ctx context.Context, order *entity.PaymentOrder) error
	UpdateStatus(ctx context.Context, req *entity.PaymentOrder, old payment.Status, updates map[string]interface{}) error
	Search(ctx context.Context, req *admin.SearchRequest) ([]*entity.PaymentOrder, error)
	Count(ctx context.Context, req *admin.SearchRequest) (int64, error)
	GetModel(ctx context.Context, req public.QueryPaymentRequest, status []payment.Status) (*model.PaymentOrder, error)
	SearchModel(ctx context.Context, req *admin.SearchRequest) ([]*model.PaymentOrder, error)
	DeletePaymentRequestLog(ctx context.Context, lastTime time.Time) error
	ConfirmNotifyStatus(ctx context.Context, id int64, status enum.NotifyStatus) error
	SavePaymentExpireRecord(ctx context.Context, req *model.PaymentExpireRecord) error
	DeletePaymentExpireRecord(ctx context.Context, lastTime time.Time) error
	GetPaymentExpireRecord(ctx context.Context, start, end time.Time, limit int, lastId int64) ([]*model.PaymentExpireRecord, error)
	DeletePaymentExpireRecordByID(ctx context.Context, id int64) error
	DeletePaymentExpireRecordByIDs(ctx context.Context, ids []int64) error
}
type paymentOrderRepoImpl struct {
}

func NewPaymentOrderRepo() PaymentOrderRepo {
	return &paymentOrderRepoImpl{}
}

func (p *paymentOrderRepoImpl) DeletePaymentExpireRecordByIDs(ctx context.Context, ids []int64) error {
	return global.GVA_PAY_DB.WithContext(ctx).Where("id in (?)", ids).Delete(&model.PaymentExpireRecord{}).Error
}

func (p *paymentOrderRepoImpl) DeletePaymentExpireRecordByID(ctx context.Context, id int64) error {
	return global.GVA_PAY_DB.WithContext(ctx).Where("id = ?", id).Delete(&model.PaymentExpireRecord{}).Error
}

func (p *paymentOrderRepoImpl) SavePaymentExpireRecord(ctx context.Context, req *model.PaymentExpireRecord) error {
	return global.GVA_PAY_DB.WithContext(ctx).Create(req).Error
}

func (p *paymentOrderRepoImpl) GetPaymentExpireRecord(ctx context.Context, start, end time.Time, limit int, lastId int64) ([]*model.PaymentExpireRecord, error) {
	var m []*model.PaymentExpireRecord
	return m, global.GVA_PAY_DB.WithContext(ctx).Where("expire_time >= ? and expire_time <= ? and id > ?", start, end, lastId).Order("expire_time ASC, id ASC").Limit(limit).Find(&m).Error
}

func (p *paymentOrderRepoImpl) DeletePaymentExpireRecord(ctx context.Context, lastTime time.Time) error {
	return global.GVA_PAY_DB.WithContext(ctx).Where("expire_time < ?", lastTime).Delete(&model.PaymentExpireRecord{}).Error
}

func (p *paymentOrderRepoImpl) ConfirmNotifyStatus(ctx context.Context, id int64, status enum.NotifyStatus) error {
	return global.GVA_PAY_DB.WithContext(ctx).Model(&model.PaymentOrder{}).Where("id = ?", id).Update("notify_status", status).Error
}

func (p *paymentOrderRepoImpl) SearchModel(ctx context.Context, req *admin.SearchRequest) ([]*model.PaymentOrder, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var m []*model.PaymentOrder
	query := buildQuery(ctx, req)
	page := 1
	pageSize := 10
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *paymentOrderRepoImpl) GetModel(ctx context.Context, req public.QueryPaymentRequest, status []payment.Status) (*model.PaymentOrder, error) {
	var m *model.PaymentOrder
	if err := req.Validate(); err != nil {
		return nil, err
	}
	query := global.GVA_PAY_DB.WithContext(ctx).Where("mch_no = ? and app_no = ?", req.MchNo, req.AppNo).Preload("PaymentOrderProducts")
	if len(status) > 0 {
		query = query.Where("status in (?)", status)
	}
	if req.OrderNo != "" {
		query = query.Where("order_no = ?", req.OrderNo)
	}
	if req.TradeNo != "" {
		query = query.Where("trade_no = ?", req.TradeNo)
	}
	return m, query.First(&m).Error
}

func (p *paymentOrderRepoImpl) Search(ctx context.Context, req *admin.SearchRequest) ([]*entity.PaymentOrder, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var m []*model.PaymentOrder
	query := buildQuery(ctx, req)
	page := 1
	pageSize := 10
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&m).Error
	if err != nil {
		return nil, err
	}
	var results = make([]*entity.PaymentOrder, 0, len(m))
	for _, item := range m {
		results = append(results, buildEntity(item))
	}
	return results, nil
}

func buildQuery(ctx context.Context, req *admin.SearchRequest) *gorm.DB {
	query := global.GVA_PAY_DB.WithContext(ctx).Model(&model.PaymentOrder{}).Where("mch_no = ?", req.MchNo).
		Where("created_at >= ?", req.StartTime).Where("created_at <= ?", req.EndTime)
	if req.OrderNo != "" {
		query = query.Where("order_no = ?", req.OrderNo)
	}
	if req.AppNo != "" {
		query = query.Where("app_no = ?", req.AppNo)
	}
	if req.TradeNo != "" {
		query = query.Where("trade_no = ?", req.TradeNo)
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}
	return query
}
func (p *paymentOrderRepoImpl) Count(ctx context.Context, req *admin.SearchRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}
	var m = int64(0)
	query := buildQuery(ctx, req)
	err := query.Count(&m).Error
	if err != nil {
		return 0, err
	}
	return m, nil
}

func (p *paymentOrderRepoImpl) UpdateStatus(ctx context.Context, req *entity.PaymentOrder, old payment.Status, updates map[string]interface{}) error {
	keys := p.getOrderKey(public.QueryPaymentRequest{
		MchNo:   req.MchNo,
		AppNo:   req.AppNo,
		OrderNo: req.OrderNo,
		TradeNo: req.TradeNo,
	})
	_ = global.GVA_REDIS.Del(ctx, keys...)
	err := global.GVA_PAY_DB.WithContext(ctx).Model(&model.PaymentOrder{}).Where("id = ? and status = ?", req.ID, old).Updates(updates).Error
	_ = global.GVA_REDIS.Del(ctx, keys...)
	return err
}

func (p *paymentOrderRepoImpl) genOrderKey(req public.QueryPaymentRequest) string {
	return fmt.Sprintf("order:%s:%s:%s:%s", req.MchNo, req.AppNo, req.OrderNo, req.TradeNo)
}

func (p *paymentOrderRepoImpl) getOrderKey(req public.QueryPaymentRequest) []string {
	return []string{
		fmt.Sprintf("order:%s:%s:%s:%s", req.MchNo, req.AppNo, req.OrderNo, req.TradeNo),
		fmt.Sprintf("order:%s:%s:%s:%s", req.MchNo, req.AppNo, "", req.TradeNo),
		fmt.Sprintf("order:%s:%s:%s:%s", req.MchNo, req.AppNo, req.OrderNo, ""),
	}
}

func (p *paymentOrderRepoImpl) GetOrderWithCache(ctx context.Context, req public.QueryPaymentRequest) (*entity.PaymentOrder, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var m model.PaymentOrder
	var order entity.PaymentOrder
	key := p.genOrderKey(req)
	cache, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == nil && len(cache) > 0 {
		_ = json.Unmarshal([]byte(cache), &order)
		if order.ID > 0 {
			return &order, nil
		}
	}
	query := global.GVA_PAY_DB.WithContext(ctx).Where("mch_no = ? and app_no = ?", req.MchNo, req.AppNo)
	if req.OrderNo != "" {
		query = query.Where("order_no = ?", req.OrderNo)
	}
	if req.TradeNo != "" {
		query = query.Where("trade_no = ?", req.TradeNo)
	}
	err = query.Order("id DESC").First(&m).Error
	if err != nil {
		return nil, err
	}
	r := buildEntity(&m)
	bytes, _ := json.Marshal(r)
	if len(bytes) > 0 {
		_ = global.GVA_REDIS.Set(ctx, key, string(bytes), time.Minute*10)
	}
	return r, nil
}

func (p *paymentOrderRepoImpl) Save(ctx context.Context, order *entity.PaymentOrder) error {
	keys := p.getOrderKey(public.QueryPaymentRequest{
		MchNo:   order.MchNo,
		AppNo:   order.AppNo,
		OrderNo: order.OrderNo,
		TradeNo: order.TradeNo,
	})
	_ = global.GVA_REDIS.Del(ctx, keys...)
	var err error
	if order.ID > 0 {
		m := buildModel(order)
		// 不用更新关联关系
		m.PaymentOrderProducts = []*model.PaymentOrderProduct{}
		err = global.GVA_PAY_DB.WithContext(ctx).Model(&model.PaymentOrder{}).Where("id = ?", order.ID).Save(m).Error
	} else {
		err = global.GVA_PAY_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return tx.Model(&model.PaymentOrder{}).Create(buildModel(order)).Error
		})
	}
	global.GVA_REDIS.Del(ctx, keys...)
	return err
}

func buildEntity(req *model.PaymentOrder) *entity.PaymentOrder {
	result := &entity.PaymentOrder{
		ID:                   req.ID,
		PaymentAccountNo:     req.PaymentAccountNo,
		AppNo:                req.AppNo,
		MchNo:                req.MchNo,
		OrderSubject:         req.OrderSubject,
		PaymentMethod:        payment.Payment(payment.Payment_value[req.PaymentMethod]),
		PaymentProduct:       payment.PaymentProduct(payment.PaymentProduct_value[req.PaymentProduct]),
		OrderNo:              req.OrderNo,
		OutMchTradeNo:        req.OutMchTradeNo,
		TradeNo:              req.TradeNo,
		Device:               req.Device,
		System:               req.System,
		OrderAmount:          req.OrderAmount,
		PaymentAmount:        req.PaymentAmount,
		UserPaymentAmount:    0,
		Currency:             req.Currency,
		Status:               payment.Status(req.Status),
		StatusFrom:           req.StatusFrom,
		TimeExpire:           req.TimeExpire,
		OrderCreate:          req.OrderCreate,
		ExtParam:             req.ExtParam,
		RedirectURL:          req.RedirectURL,
		NotifyURL:            req.NotifyURL,
		OrderDesc:            req.OrderDesc,
		ThirdCode:            req.ThirdCode,
		ThirdMsg:             req.ThirdMsg,
		CreatedAt:            req.CreatedAt,
		PaymentOrderProducts: make([]entity.Product, 0, len(req.PaymentOrderProducts)),
		UserOpenid:           req.UserOpenid,
		NotifyStatus:         req.NotifyStatus,
		Retry:                req.Retry,
	}
	if len(req.PaymentOrderProducts) > 0 {
		for _, item := range req.PaymentOrderProducts {
			result.PaymentOrderProducts = append(result.PaymentOrderProducts, entity.Product{
				TradeNo:     req.TradeNo,
				ProductName: item.ProductName,
				AppNo:       req.AppNo,
				MchNo:       req.MchNo,
				Quantity:    int64(item.Quantity),
				Price:       int64(item.Price),
				Sku:         item.Sku,
				ProductDesc: item.ProductDesc,
				URL:         "",
				OrderNo:     req.OrderNo,
				ID:          item.ID,
			})
		}
	}
	if result.NotifyURL != "" && result.NotifyStatus == 0 {
		result.NotifyStatus = int64(enum.NotifyStatus_Init)
	}
	return result
}

func buildModel(req *entity.PaymentOrder) *model.PaymentOrder {
	result := &model.PaymentOrder{
		ID:                   req.ID,
		PaymentAccountNo:     req.PaymentAccountNo,
		AppNo:                req.AppNo,
		MchNo:                req.MchNo,
		OrderSubject:         req.OrderSubject,
		PaymentMethod:        req.PaymentMethod.String(),
		PaymentProduct:       req.PaymentProduct.String(),
		OrderNo:              req.OrderNo,
		OutMchTradeNo:        req.OutMchTradeNo,
		TradeNo:              req.TradeNo,
		Device:               req.Device,
		System:               req.System,
		OrderAmount:          req.OrderAmount,
		PaymentAmount:        req.PaymentAmount,
		UserPaymentAmount:    0,
		Currency:             req.Currency,
		Status:               int64(req.Status),
		StatusFrom:           req.StatusFrom,
		TimeExpire:           req.TimeExpire,
		OrderCreate:          req.OrderCreate,
		ExtParam:             req.ExtParam,
		RedirectURL:          req.RedirectURL,
		NotifyURL:            req.NotifyURL,
		OrderDesc:            req.OrderDesc,
		ThirdCode:            req.ThirdCode,
		ThirdMsg:             req.ThirdMsg,
		CreatedAt:            req.CreatedAt,
		PaymentOrderProducts: make([]*model.PaymentOrderProduct, 0, len(req.PaymentOrderProducts)),
		Retry:                req.Retry,
		UserOpenid:           req.UserOpenid,
	}
	if len(req.PaymentOrderProducts) > 0 {
		for _, item := range req.PaymentOrderProducts {
			result.PaymentOrderProducts = append(result.PaymentOrderProducts, &model.PaymentOrderProduct{
				TradeNo:     req.TradeNo,
				ProductName: item.ProductName,
				AppNo:       req.AppNo,
				MchNo:       req.MchNo,
				Quantity:    int64(item.Quantity),
				Price:       int64(item.Price),
				Sku:         item.Sku,
				ProductDesc: item.ProductDesc,
				URL:         item.URL,
				OrderNo:     req.OrderNo,
				ID:          item.ID,
			})
		}
	}
	return result
}

func (p *paymentOrderRepoImpl) DeletePaymentRequestLog(ctx context.Context, lastTime time.Time) error {
	return global.GVA_PAY_DB.WithContext(ctx).Where("created_at < ?", lastTime).Delete(&model.PaymentRequestLog{}).Error
}
