package repo

import (
	"context"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"gorm.io/gorm"
	"time"
)

type EventRecordRepo interface {
	Save(ctx context.Context, record *model.EventProcessRecord, tx *gorm.DB) error
	Get(ctx context.Context, eventID string, eventType string, consumer string) (*model.EventProcessRecord, error)
	Delete(ctx context.Context, lastTime time.Time) error
}

type eventRecordRepoImpl struct {
}

func NewEventRecordRepo() EventRecordRepo {
	return &eventRecordRepoImpl{}
}

func (e *eventRecordRepoImpl) Save(ctx context.Context, record *model.EventProcessRecord, tx *gorm.DB) error {
	return tx.WithContext(ctx).Create(record).Error
}

func (e *eventRecordRepoImpl) Get(ctx context.Context, eventID string, eventType string, consumer string) (*model.EventProcessRecord, error) {
	var record *model.EventProcessRecord
	err := global.GVA_PAY_DB.WithContext(ctx).Where("event_id = ? and event_type = ? and consumer = ?", eventID, eventType, consumer).First(&record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (e *eventRecordRepoImpl) Delete(ctx context.Context, lastTime time.Time) error {
	//var record model.EventProcessRecord
	//err := global.GVA_PAY_DB.WithContext(ctx).Where("created_at < ?", lastTime).Order("id DESC").First(&record).Error
	//if err != nil {
	//	return err
	//}
	return global.GVA_PAY_DB.WithContext(ctx).Where("created_at < ?", lastTime).Delete(&model.EventProcessRecord{}).Error
}
