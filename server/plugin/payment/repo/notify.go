package repo

import (
	"context"
	"time"

	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"gorm.io/gorm"
)

type NotifyRepo interface {
	Get(ctx context.Context, mchNO, appNo string, notifyType enum.NotifyType, outNo string) (*model.NotifyRecord, error)
	Save(ctx context.Context, record *model.NotifyRecord) error
	Retry(ctx context.Context, notifyID int64, nextTime time.Time) error
	Delete(ctx context.Context, lastTime time.Time) error
	UpdateStatus(ctx context.Context, id int64, status enum.NotifyStatus, result string) error
	GetRetry(ctx context.Context, start time.Time, end time.Time, lastId int64, limit int, notifyType enum.NotifyType) ([]*model.NotifyRecord, error)
	GetById(ctx context.Context, id int64) (*model.NotifyRecord, error)
	CountRetry(ctx context.Context, start time.Time, end time.Time, notifyType enum.NotifyType) (int64, error)
}
type notifyRepoImpl struct {
	db *gorm.DB
}

func NewNotifyRepo(db *gorm.DB) NotifyRepo {
	return &notifyRepoImpl{db: db}
}

func (n *notifyRepoImpl) CountRetry(ctx context.Context, start time.Time, end time.Time, notifyType enum.NotifyType) (int64, error) {
	var count int64
	err := n.db.WithContext(ctx).Model(&model.NotifyRecord{}).
		Where("last_notify_time BETWEEN ? AND ?", start, end).
		Where("notify_status = ?", int64(enum.NotifyStatus_Init)).
		Where("notify_type = ?", int64(notifyType)).
		Count(&count).Error
	return count, err
}

func (n *notifyRepoImpl) GetById(ctx context.Context, id int64) (*model.NotifyRecord, error) {
	var m model.NotifyRecord
	err := n.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (n *notifyRepoImpl) GetRetry(ctx context.Context, start time.Time, end time.Time, lastId int64, limit int, notifyType enum.NotifyType) ([]*model.NotifyRecord, error) {
	var m []*model.NotifyRecord
	err := n.db.WithContext(ctx).
		Where("last_notify_time BETWEEN ? AND ?", start, end).
		Where("id > ?", lastId).
		Where("notify_status = ?", int64(enum.NotifyStatus_Init)).
		Where("notify_type = ?", int64(notifyType)).
		Order("last_notify_time ASC, id ASC").Limit(limit).Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (n *notifyRepoImpl) Get(ctx context.Context, mchNO, appNo string, notifyType enum.NotifyType, outNo string) (*model.NotifyRecord, error) {
	var m model.NotifyRecord
	err := n.db.WithContext(ctx).
		Where("mch_no = ? AND app_no = ? AND notify_type = ? AND out_no = ?", mchNO, appNo, int64(notifyType), outNo).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (n *notifyRepoImpl) Save(ctx context.Context, record *model.NotifyRecord) error {
	return n.db.WithContext(ctx).Create(record).Error
}

func (n *notifyRepoImpl) UpdateStatus(ctx context.Context, id int64, status enum.NotifyStatus, result string) error {
	return n.db.WithContext(ctx).Model(&model.NotifyRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"notify_status": status,
		"res_result":    result,
	}).Error
}

func (n *notifyRepoImpl) Retry(ctx context.Context, notifyID int64, nextTime time.Time) error {
	return n.db.WithContext(ctx).Model(&model.NotifyRecord{}).Where("id = ?", notifyID).Updates(map[string]interface{}{
		"notify_count":     gorm.Expr("notify_count + 1"),
		"last_notify_time": nextTime,
	}).Error
}

func (n *notifyRepoImpl) Delete(ctx context.Context, lastTime time.Time) error {
	return n.db.WithContext(ctx).Where("created_at <= ?", lastTime).Delete(&model.NotifyRecord{}).Error
}
