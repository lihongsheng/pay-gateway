package repo

import (
	"context"
	"errors"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/dao"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Statistics interface {
	Save(ctx context.Context, statics *model.Statistic, tx *gorm.DB) error
	Get(ctx context.Context, mchNo, appNo, statisticTag, statisticKey string, t time.Time) (*model.Statistic, error)
	CountGroupMch(ctx context.Context, mchNo, appNo, statisticTag, statisticKey string, start, end time.Time) (*model.Statistic, error)
}

type statistics struct {
}

func NewStatistics() Statistics {
	return &statistics{}
}

func (s *statistics) Save(ctx context.Context, statics *model.Statistic, tx *gorm.DB) error {
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
				{Name: "statistic_tag"},
				{Name: "statistic_key"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"statistic":  gorm.Expr("statistic + ?", statics.Statistic),
				"updated_at": time.Now(),
			}),
		},
	).Save(statics)

	return result.Error
}

func (s *statistics) Get(ctx context.Context, mchNo, appNo, statisticTag, statisticKey string, t time.Time) (*model.Statistic, error) {
	return dao.Statistic.WithContext(ctx).Where(dao.Statistic.MchNo.Eq(mchNo),
		dao.Statistic.AppNo.Eq(appNo),
		dao.Statistic.StatisticTag.Eq(statisticTag),
		dao.Statistic.StatisticKey.Eq(statisticKey),
		dao.Statistic.StatisticDate.Eq(t)).First()
}

func (s *statistics) CountGroupMch(ctx context.Context, mchNo, appNo, statisticTag, statisticKey string, start, end time.Time) (*model.Statistic, error) {
	var statics *model.Statistic
	return statics, global.GVA_PAY_DB.WithContext(ctx).
		Model(&model.Statistic{}).
		Select("app_no", "mch_no", "statistic_date", "statistic_tag", "statistic_key",
			"SUM(statistic) as statistic").
		Where("mch_no = ?", mchNo).
		Where("app_no = ?", appNo).
		Where("statistic_tag = ?", statisticTag).
		Where("statistic_key = ?", statisticKey).
		Where("statistic_date >= ?", start.Format("2006-01-02")).
		Where("statistic_date <= ?", end.Format("2006-01-02")).
		Group("mch_no").Group("app_no").Group("statistic_date").Group("statistic_tag").Group("statistic_key").
		First(&statics).Error
}
