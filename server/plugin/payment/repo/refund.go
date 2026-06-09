package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	enum2 "github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/dao"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/payment-sdk/enum/refund"
	"gorm.io/gorm"
	"time"
)

type RefundRepo interface {
	GetRefunds(ctx context.Context, mchNO, appNo, orderNo string, status []refund.Status) ([]*model.RefundOrder, error)
	GetFromCache(ctx context.Context, request public.RefundQueryRequest) (*model.RefundOrder, error)
	Save(ctx context.Context, refundOrder *model.RefundOrder) error
	CountAmount(ctx context.Context, mchNO, appNo, orderNo string, status []refund.Status) (int64, error)
	Search(ctx context.Context, req *admin.RefundSearchRequest) ([]*model.RefundOrder, error)
	Count(ctx context.Context, req *admin.RefundSearchRequest) (int64, error)
	ConfirmNotifyStatus(ctx context.Context, id int64, status enum2.NotifyStatus) error
	GetRefund(ctx context.Context, mchNO string, appNO string, refundTradeNo string) (*model.RefundOrder, error)
	//SaveWithTransaction(ctx context.Context, refundOrder *model.RefundOrder, tx *gorm.DB) error
}

type refundRepoImpl struct {
}

func NewRefundRepo() RefundRepo {
	return &refundRepoImpl{}
}

func (r *refundRepoImpl) GetRefund(ctx context.Context, mchNO string, appNO string, refundTradeNo string) (*model.RefundOrder, error) {
	var refundOrder model.RefundOrder
	err := global.GVA_PAY_DB.WithContext(ctx).Where(dao.RefundOrder.MchNo.Eq(mchNO), dao.RefundOrder.AppNo.Eq(appNO), dao.RefundOrder.RefundTradeNo.Eq(refundTradeNo)).First(&refundOrder).Error
	if err != nil {
		return nil, err
	}
	return &refundOrder, nil
}

func (r *refundRepoImpl) ConfirmNotifyStatus(ctx context.Context, id int64, status enum2.NotifyStatus) error {
	return global.GVA_PAY_DB.WithContext(ctx).Model(&model.RefundOrder{}).Where("id = ?", id).Update("notify_status", status).Error
}

func (r *refundRepoImpl) Search(ctx context.Context, req *admin.RefundSearchRequest) ([]*model.RefundOrder, error) {
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

func (r *refundRepoImpl) Count(ctx context.Context, req *admin.RefundSearchRequest) (int64, error) {
	query := r.buildAdminQuery(ctx, req)
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *refundRepoImpl) buildAdminQuery(ctx context.Context, req *admin.RefundSearchRequest) *gorm.DB {
	query := global.GVA_PAY_DB.WithContext(ctx).Model(&model.RefundOrder{})
	if req.MchNo != "" {
		query = query.Where(dao.RefundOrder.MchNo.Eq(req.MchNo))
	}
	if req.AppNo != "" {
		query = query.Where(dao.RefundOrder.AppNo.Eq(req.AppNo))
	}
	if req.OrderNo != "" {
		query = query.Where(dao.RefundOrder.OrderNo.Eq(req.OrderNo))
	}
	if req.TradeNo != "" {
		query = query.Where(dao.RefundOrder.RefundTradeNo.Eq(req.TradeNo))
	}
	if req.RefundNo != "" {
		query = query.Where(dao.RefundOrder.RefundNo.Eq(req.RefundNo))
	}
	if !req.StartTime.IsZero() {
		query = query.Where(dao.RefundOrder.CreatedAt.Gte(req.StartTime))
	}
	if !req.EndTime.IsZero() {
		query = query.Where(dao.RefundOrder.CreatedAt.Lte(req.EndTime))
	}
	if req.Status != 0 {
		query = query.Where(dao.RefundOrder.Status.Eq(int64(req.Status)))
	}
	return query
}

func (r *refundRepoImpl) CountAmount(ctx context.Context, mchNO, appNo, orderNo string, status []refund.Status) (int64, error) {
	var amount int64
	query := global.GVA_PAY_DB.WithContext(ctx).Model(&model.RefundOrder{}).Select("IFNULL(SUM(refund_amount), 0) AS total_refund_amount").Where(dao.RefundOrder.MchNo.Eq(mchNO), dao.RefundOrder.AppNo.Eq(appNo), dao.RefundOrder.OrderNo.Eq(orderNo))
	if len(status) > 0 {
		s := []int64{}
		for _, v := range status {
			s = append(s, int64(v))
		}
		query = query.Where(dao.RefundOrder.Status.In(s...))
	}
	err := query.First(&amount).Error
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func (r *refundRepoImpl) GetRefunds(ctx context.Context, mchNO, appNo, orderNo string, status []refund.Status) ([]*model.RefundOrder, error) {
	var refunds []*model.RefundOrder
	query := global.GVA_PAY_DB.WithContext(ctx).Where(dao.RefundOrder.MchNo.Eq(mchNO), dao.RefundOrder.AppNo.Eq(appNo), dao.RefundOrder.OrderNo.Eq(orderNo))
	if len(status) > 0 {
		s := []int64{}
		for _, v := range status {
			s = append(s, int64(v))
		}
		query = query.Where(dao.RefundOrder.Status.In(s...))
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
	cache, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == nil && len(cache) > 0 {
		_ = json.Unmarshal([]byte(cache), &m)
		if m.ID > 0 {
			return &m, nil
		}
	}

	query := global.GVA_PAY_DB.WithContext(ctx).Where(dao.RefundOrder.MchNo.Eq(req.MchNo), dao.RefundOrder.AppNo.Eq(req.AppNo))
	if req.RefundNo != "" {
		query = query.Where(dao.RefundOrder.RefundNo.Eq(req.RefundNo))
	}
	if req.RefundTradeNo != "" {
		query = query.Where(dao.RefundOrder.RefundTradeNo.Eq(req.RefundTradeNo))
	}
	err = query.First(&m).Error
	if err != nil {
		return nil, err
	}
	bytes, _ := json.Marshal(m)
	if len(bytes) > 0 {
		_ = global.GVA_REDIS.Set(ctx, key, string(bytes), time.Minute*10)
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
	//return fmt.Sprintf("refund:%s:%s:%s:%s", req.MchNo, req.AppNo, req.RefundNo,req.RefundTradeNo)
}

func (r *refundRepoImpl) Save(ctx context.Context, refundOrder *model.RefundOrder) error {
	cacheKeys := r.getRefundKey(public.RefundQueryRequest{
		MchNo:         refundOrder.MchNo,
		AppNo:         refundOrder.AppNo,
		RefundNo:      refundOrder.RefundNo,
		RefundTradeNo: refundOrder.RefundTradeNo,
	})
	global.GVA_REDIS.Del(ctx, cacheKeys...)
	if refundOrder.NotifyURL != "" && refundOrder.NotifyStatus == 0 {
		refundOrder.NotifyStatus = int64(enum2.NotifyStatus_Init)
	}
	refundOrder.UpdatedAt = time.Now()
	err := global.GVA_PAY_DB.WithContext(ctx).Save(refundOrder).Error
	global.GVA_REDIS.Del(ctx, cacheKeys...)
	return err
}

//func (r *refundRepoImpl) SaveWithTransaction(ctx context.Context, refundOrder *model.RefundOrder, tx *gorm.DB) error {
//	if refundOrder.NotifyURL != "" && refundOrder.NotifyStatus == 0 {
//		refundOrder.NotifyStatus = int64(enum2.NotifyStatus_Init)
//	}
//	refundOrder.UpdatedAt = time.Now()
//	return tx.WithContext(ctx).Model(&model.RefundOrder{}).Save(refundOrder).Error
//}
