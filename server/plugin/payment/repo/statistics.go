package repo

import (
	"context"
	"errors"
	"time"

	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Statistics interface {
	Save(ctx context.Context, statics *model.Statistic, tx *gorm.DB) error
	Get(ctx context.Context, mchNo, appNo, statisticTag, statisticKey string, t time.Time) (*model.Statistic, error)
	CountGroupMch(ctx context.Context, mchNo, appNo, statisticTag, statisticKey string, start, end time.Time) (*model.Statistic, error)
}

type statistics struct {
	db *gorm.DB
}

func NewStatistics(db *gorm.DB) Statistics {
	return &statistics{db: db}
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
	var m model.Statistic
	err := s.db.WithContext(ctx).
		Where("mch_no = ? AND app_no = ? AND statistic_tag = ? AND statistic_key = ? AND statistic_date = ?", mchNo, appNo, statisticTag, statisticKey, t).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *statistics) CountGroupMch(ctx context.Context, mchNo, appNo, statisticTag, statisticKey string, start, end time.Time) (*model.Statistic, error) {
	var statics *model.Statistic
	return statics, s.db.WithContext(ctx).
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
		Take(&statics).Error
}
