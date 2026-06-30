package repo

import (
	"context"
	"time"

	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"gorm.io/gorm"
)

type EventRecordRepo interface {
	Save(ctx context.Context, record *model.EventProcessRecord, tx *gorm.DB) error
	Get(ctx context.Context, eventID string, eventType string, consumer string) (*model.EventProcessRecord, error)
	Delete(ctx context.Context, lastTime time.Time) error
}

type eventRecordRepoImpl struct {
	db *gorm.DB
}

func NewEventRecordRepo(db *gorm.DB) EventRecordRepo {
	return &eventRecordRepoImpl{db: db}
}

func (e *eventRecordRepoImpl) Save(ctx context.Context, record *model.EventProcessRecord, tx *gorm.DB) error {
	return tx.WithContext(ctx).Create(record).Error
}

func (e *eventRecordRepoImpl) Get(ctx context.Context, eventID string, eventType string, consumer string) (*model.EventProcessRecord, error) {
	var record *model.EventProcessRecord
	err := e.db.WithContext(ctx).Where("event_id = ? and event_type = ? and consumer = ?", eventID, eventType, consumer).First(&record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (e *eventRecordRepoImpl) Delete(ctx context.Context, lastTime time.Time) error {
	return e.db.WithContext(ctx).Where("created_at < ?", lastTime).Delete(&model.EventProcessRecord{}).Error
}
