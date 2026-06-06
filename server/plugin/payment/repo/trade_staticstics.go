package repo

import (
	"context"
	"errors"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TradeStaticsRepo interface {
	Save(ctx context.Context, statics *model.TradeStatistic, tx *gorm.DB) error
	Get(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*model.TradeStatistic, error)
	CountGroup(ctx context.Context, start, end time.Time) ([]*model.TradeStatistic, error)
	CountGroupMch(ctx context.Context, mchNo string, start, end time.Time) ([]*model.TradeStatistic, error)
	CountGroupMchApp(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*model.TradeStatistic, error)
	CountGroupMchAppAccount(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*model.TradeStatistic, error)
	SearchDashboard(ctx context.Context, req *admin.MchTradeStatisticRequest) ([]*model.TradeStatistic, error)
	Count(ctx context.Context, req *admin.MchTradeStatisticRequest) (int64, error)
}

type tradeStaticsRepoImpl struct {
}

func NewTradeStaticsRepo() TradeStaticsRepo {
	return &tradeStaticsRepoImpl{}
}

func (t *tradeStaticsRepoImpl) Count(ctx context.Context, req *admin.MchTradeStatisticRequest) (int64, error) {
	var count int64
	query := t.BuildQuery(ctx, req)
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (t *tradeStaticsRepoImpl) SearchDashboard(ctx context.Context, req *admin.MchTradeStatisticRequest) ([]*model.TradeStatistic, error) {
	var models []*model.TradeStatistic
	query := t.BuildQuery(ctx, req)
	err := query.Model(&model.TradeStatistic{}).Where(query).Order("statistic_date DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (t *tradeStaticsRepoImpl) BuildQuery(ctx context.Context, req *admin.MchTradeStatisticRequest) *gorm.DB {
	query := global.GVA_PAY_DB.WithContext(ctx).Model(&model.TradeStatistic{})
	if req.AppNo != "" {
		query = query.Where("app_no = ?", req.AppNo)
	}
	if req.MchNo != "" {
		query = query.Where("mch_no = ?", req.MchNo)
	}
	if !req.StartTime.IsZero() {
		query = query.Where("statistic_date >= ?", req.StartTime.Format("2006-01-02"))
	}
	if !req.EndTime.IsZero() {
		query = query.Where("statistic_date <= ?", req.EndTime.Format("2006-01-02"))
	}
	if req.AccountNo != "" {
		query = query.Where("account_no = ?", req.AccountNo)
	}
	return query
}

func (t *tradeStaticsRepoImpl) Save(ctx context.Context, statics *model.TradeStatistic, tx *gorm.DB) error {
	// 检查参数合法性
	if statics == nil {
		return errors.New("参数 statics 不能为空")
	}
	// 使用 GORM 的 OnConflict 子句实现插入或更新，并累加 total_order
	result := tx.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "app_no"},
				{Name: "mch_no"},
				{Name: "statistic_date"},
				{Name: "account_no"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"total_order":    gorm.Expr("total_order + ?", statics.TotalOrder),
				"success_order":  gorm.Expr("success_order + ?", statics.SuccessOrder),
				"total_amount":   gorm.Expr("total_amount + ?", statics.TotalAmount),
				"success_amount": gorm.Expr("success_amount + ?", statics.SuccessAmount),
				"refund_order":   gorm.Expr("refund_order + ?", statics.RefundOrder),
				"refund_amount":  gorm.Expr("refund_amount + ?", statics.RefundAmount),
				"updated_at":     time.Now(),
			}),
		},
	).Save(statics)

	return result.Error
}

func (t *tradeStaticsRepoImpl) Get(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*model.TradeStatistic, error) {
	var statics []*model.TradeStatistic
	return statics, global.GVA_PAY_DB.WithContext(ctx).Where("mch_no = ?", mchNo).Where("app_no = ?", appNo).Where("statistic_date >= ?", start.Format("2006-01-02")).Where("statistic_date <= ?", end.Format("2006-01-02")).Find(&statics).Error
}

func (t *tradeStaticsRepoImpl) CountGroup(ctx context.Context, start, end time.Time) ([]*model.TradeStatistic, error) {
	var statics []*model.TradeStatistic
	return statics, global.GVA_PAY_DB.WithContext(ctx).
		Model(&model.TradeStatistic{}).
		Select("statistic_date",
			"SUM(total_order) as total_order",
			"SUM(success_order) as success_order",
			"SUM(total_amount) as total_amount",
			"SUM(success_amount) as success_amount",
			"SUM(refund_order) as refund_order",
			"SUM(refund_amount) as refund_amount").
		Where("statistic_date >= ?", start.Format("2006-01-02")).
		Where("statistic_date <= ?", end.Format("2006-01-02")).
		Group("statistic_date").
		Find(&statics).Error
}

func (t *tradeStaticsRepoImpl) CountGroupMch(ctx context.Context, mchNo string, start, end time.Time) ([]*model.TradeStatistic, error) {
	var statics []*model.TradeStatistic
	return statics, global.GVA_PAY_DB.WithContext(ctx).
		Model(&model.TradeStatistic{}).
		Select("mch_no", "statistic_date",
			"SUM(total_order) as total_order",
			"SUM(success_order) as success_order",
			"SUM(total_amount) as total_amount",
			"SUM(success_amount) as success_amount",
			"SUM(refund_order) as refund_order",
			"SUM(refund_amount) as refund_amount").
		Where("mch_no = ?", mchNo).
		Where("statistic_date >= ?", start.Format("2006-01-02")).
		Where("statistic_date <= ?", end.Format("2006-01-02")).
		Group("mch_no").Group("statistic_date").
		Find(&statics).Error
}

func (t *tradeStaticsRepoImpl) CountGroupMchApp(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*model.TradeStatistic, error) {
	var statics []*model.TradeStatistic
	return statics, global.GVA_PAY_DB.WithContext(ctx).
		Model(&model.TradeStatistic{}).
		Select("app_no", "mch_no", "statistic_date",
			"SUM(total_order) as total_order",
			"SUM(success_order) as success_order",
			"SUM(total_amount) as total_amount",
			"SUM(success_amount) as success_amount",
			"SUM(refund_order) as refund_order",
			"SUM(refund_amount) as refund_amount").
		Where("mch_no = ?", mchNo).
		Where("app_no = ?", appNo).
		Where("statistic_date >= ?", start.Format("2006-01-02")).
		Where("statistic_date <= ?", end.Format("2006-01-02")).
		Group("mch_no").Group("app_no").Group("statistic_date").
		Find(&statics).Error
}

func (t *tradeStaticsRepoImpl) CountGroupMchAppAccount(ctx context.Context, mchNo string, appNo string, start, end time.Time) ([]*model.TradeStatistic, error) {
	var statics []*model.TradeStatistic
	return statics, global.GVA_PAY_DB.WithContext(ctx).
		Model(&model.TradeStatistic{}).
		Select("app_no", "mch_no", "account_no", "statistic_date",
			"SUM(total_order) as total_order",
			"SUM(success_order) as success_order",
			"SUM(total_amount) as total_amount",
			"SUM(success_amount) as success_amount",
			"SUM(refund_order) as refund_order",
			"SUM(refund_amount) as refund_amount").
		Where("mch_no = ?", mchNo).
		Where("app_no = ?", appNo).
		Where("statistic_date >= ?", start.Format("2006-01-02")).
		Where("statistic_date <= ?", end.Format("2006-01-02")).
		Group("mch_no").Group("app_no").Group("statistic_date").Group("account_no").
		Find(&statics).Error
}
