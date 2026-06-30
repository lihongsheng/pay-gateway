package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	enum2 "github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/payment-sdk/enum/refund"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RefundRepo interface {
	GetRefunds(ctx context.Context, mchNO, appNo, orderNo string, status []refund.Status) ([]*model.RefundOrder, error)
	GetFromCache(ctx context.Context, request public.RefundQueryRequest) (*model.RefundOrder, error)
	Save(ctx context.Context, refundOrder *model.RefundOrder) error
	CountAmount(ctx context.Context, mchNO, appNo, orderNo string, status []refund.Status) (int64, error)
	Search(ctx context.Context, req *dto.RefundSearchRequest) ([]*model.RefundOrder, error)
	Count(ctx context.Context, req *dto.RefundSearchRequest) (int64, error)
	ConfirmNotifyStatus(ctx context.Context, id int64, status enum2.NotifyStatus) error
	GetRefund(ctx context.Context, mchNO string, appNO string, refundTradeNo string) (*model.RefundOrder, error)
	//SaveWithTransaction(ctx context.Context, refundOrder *model.RefundOrder, tx *gorm.DB) error
}

type refundRepoImpl struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRefundRepo(db *gorm.DB, rdb *redis.Client) RefundRepo {
	return &refundRepoImpl{db: db, rdb: rdb}
}

func (r *refundRepoImpl) GetRefund(ctx context.Context, mchNO string, appNO string, refundTradeNo string) (*model.RefundOrder, error) {
	var refundOrder model.RefundOrder
	err := r.db.WithContext(ctx).Where("mch_no = ? AND app_no = ? AND refund_trade_no = ?", mchNO, appNO, refundTradeNo).First(&refundOrder).Error
	if err != nil {
		return nil, err
	}
	return &refundOrder, nil
}

func (r *refundRepoImpl) ConfirmNotifyStatus(ctx context.Context, id int64, status enum2.NotifyStatus) error {
	return r.db.WithContext(ctx).Model(&model.RefundOrder{}).Where("id = ?", id).Update("notify_status", status).Error
}

func (r *refundRepoImpl) Search(ctx context.Context, req *dto.RefundSearchRequest) ([]*model.RefundOrder, error) {
	query := r.buildAdminQuery(ctx, req)
	var refunds []*model.RefundOrder
	if req.PageSize > 0 && req.Page > 0 {
		query = query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
	}
	err := query.Order("id DESC").Find(&refunds).Error
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

func (r *refundRepoImpl) Count(ctx context.Context, req *dto.RefundSearchRequest) (int64, error) {
	query := r.buildAdminQuery(ctx, req)
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *refundRepoImpl) buildAdminQuery(ctx context.Context, req *dto.RefundSearchRequest) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&model.RefundOrder{})
	if req.MchNo != "" {
		query = query.Where("mch_no = ?", req.MchNo)
	}
	if req.AppNo != "" {
		query = query.Where("app_no = ?", req.AppNo)
	}
	if req.OrderNo != "" {
		query = query.Where("order_no = ?", req.OrderNo)
	}
	if req.TradeNo != "" {
		query = query.Where("refund_trade_no = ?", req.TradeNo)
	}
	if req.RefundNo != "" {
		query = query.Where("refund_no = ?", req.RefundNo)
	}
	if !req.StartTime.IsZero() {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		query = query.Where("created_at <= ?", req.EndTime)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", int64(req.Status))
	}
	return query
}

func (r *refundRepoImpl) CountAmount(ctx context.Context, mchNO, appNo, orderNo string, status []refund.Status) (int64, error) {
	var amount int64
	query := r.db.WithContext(ctx).Model(&model.RefundOrder{}).Select("IFNULL(SUM(refund_amount), 0) AS total_refund_amount").Where("mch_no = ? AND app_no = ? AND order_no = ?", mchNO, appNo, orderNo)
	if len(status) > 0 {
		s := []int64{}
		for _, v := range status {
			s = append(s, int64(v))
		}
		query = query.Where("status IN ?", s)
	}
	err := query.First(&amount).Error
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func (r *refundRepoImpl) GetRefunds(ctx context.Context, mchNO, appNo, orderNo string, status []refund.Status) ([]*model.RefundOrder, error) {
	var refunds []*model.RefundOrder
	query := r.db.WithContext(ctx).Where("mch_no = ? AND app_no = ? AND order_no = ?", mchNO, appNo, orderNo)
	if len(status) > 0 {
		s := []int64{}
		for _, v := range status {
			s = append(s, int64(v))
		}
		query = query.Where("status IN ?", s)
	}
	err := query.Find(&refunds).Error
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

func (r *refundRepoImpl) GetFromCache(ctx context.Context, req public.RefundQueryRequest) (*model.RefundOrder, error) {
	key := r.genRefundKey(req)
	var m model.RefundOrder
	cache, err := r.rdb.Get(ctx, key).Result()
	if err == nil && len(cache) > 0 {
		_ = json.Unmarshal([]byte(cache), &m)
		if m.ID > 0 {
			return &m, nil
		}
	}

	query := r.db.WithContext(ctx).Where("mch_no = ? AND app_no = ?", req.MchNo, req.AppNo)
	if req.RefundNo != "" {
		query = query.Where("refund_no = ?", req.RefundNo)
	}
	if req.RefundTradeNo != "" {
		query = query.Where("refund_trade_no = ?", req.RefundTradeNo)
	}
	err = query.First(&m).Error
	if err != nil {
		return nil, err
	}
	bytes, _ := json.Marshal(m)
	if len(bytes) > 0 {
		_ = r.rdb.Set(ctx, key, string(bytes), time.Minute*10)
	}
	return &m, nil
}

func (r *refundRepoImpl) genRefundKey(req public.RefundQueryRequest) string {
	return fmt.Sprintf("refund:%s:%s:%s:%s", req.MchNo, req.AppNo, req.RefundNo, req.RefundTradeNo)
}

func (r *refundRepoImpl) getRefundKey(req public.RefundQueryRequest) []string {
	return []string{
		fmt.Sprintf("refund:%s:%s:%s:%s", req.MchNo, req.AppNo, req.RefundNo, req.RefundTradeNo),
		fmt.Sprintf("refund:%s:%s:%s:%s", req.MchNo, req.AppNo, "", req.RefundTradeNo),
		fmt.Sprintf("refund:%s:%s:%s:%s", req.MchNo, req.AppNo, req.RefundNo, ""),
	}
}

func (r *refundRepoImpl) Save(ctx context.Context, refundOrder *model.RefundOrder) error {
	cacheKeys := r.getRefundKey(public.RefundQueryRequest{
		MchNo:         refundOrder.MchNo,
		AppNo:         refundOrder.AppNo,
		RefundNo:      refundOrder.RefundNo,
		RefundTradeNo: refundOrder.RefundTradeNo,
	})
	r.rdb.Del(ctx, cacheKeys...)
	if refundOrder.NotifyURL != "" && refundOrder.NotifyStatus == 0 {
		refundOrder.NotifyStatus = int64(enum2.NotifyStatus_Init)
	}
	refundOrder.UpdatedAt = time.Now()
	err := r.db.WithContext(ctx).Save(refundOrder).Error
	r.rdb.Del(ctx, cacheKeys...)
	return err
}
