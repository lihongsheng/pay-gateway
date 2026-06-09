package repo

import (
	"context"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/dao"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"time"
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
type notifyRepoImpl struct{}

func NewNotifyRepo() NotifyRepo {
	return &notifyRepoImpl{}
}

func (n *notifyRepoImpl) CountRetry(ctx context.Context, start time.Time, end time.Time, notifyType enum.NotifyType) (int64, error) {
	return dao.NotifyRecord.WithContext(ctx).Where(dao.NotifyRecord.LastNotifyTime.Between(start, end), dao.NotifyRecord.NotifyStatus.Eq(int64(enum.NotifyStatus_Init)), dao.NotifyRecord.NotifyType.Eq(int64(notifyType))).Count()
}

func (n *notifyRepoImpl) GetById(ctx context.Context, id int64) (*model.NotifyRecord, error) {
	return dao.NotifyRecord.WithContext(ctx).Where(dao.NotifyRecord.ID.Eq(id)).First()
}

func (n *notifyRepoImpl) GetRetry(ctx context.Context, start time.Time, end time.Time, lastId int64, limit int, notifyType enum.NotifyType) ([]*model.NotifyRecord, error) {
	return dao.NotifyRecord.WithContext(ctx).
		Where(dao.NotifyRecord.LastNotifyTime.Between(start, end), dao.NotifyRecord.ID.Gt(lastId),
			dao.NotifyRecord.NotifyStatus.Eq(int64(enum.NotifyStatus_Init)), dao.NotifyRecord.NotifyType.Eq(int64(notifyType))).
		Order(dao.NotifyRecord.LastNotifyTime.Asc(), dao.NotifyRecord.ID.Asc()).Limit(limit).Find()
}

func (n *notifyRepoImpl) Get(ctx context.Context, mchNO, appNo string, notifyType enum.NotifyType, outNo string) (*model.NotifyRecord, error) {
	return dao.NotifyRecord.WithContext(ctx).Where(dao.NotifyRecord.MchNo.Eq(mchNO), dao.NotifyRecord.AppNo.Eq(appNo), dao.NotifyRecord.NotifyType.Eq(int64(notifyType)), dao.NotifyRecord.OutNo.Eq(outNo)).First()
}

func (n *notifyRepoImpl) Save(ctx context.Context, record *model.NotifyRecord) error {
	return dao.NotifyRecord.WithContext(ctx).Create(record)
}

func (n *notifyRepoImpl) UpdateStatus(ctx context.Context, id int64, status enum.NotifyStatus, result string) error {
	_, err := dao.NotifyRecord.WithContext(ctx).Where(dao.NotifyRecord.ID.Eq(id)).Updates(map[string]interface{}{
		"notify_status": status,
		"res_result":    result,
	})
	return err
}

//func (n *notifyRepoImpl) ConfirmSuccess(ctx context.Context, id int64) error {
//	_, err := dao.NotifyRecord.WithContext(ctx).Where(dao.NotifyRecord.ID.Eq(id)).Updates(map[string]interface{}{
//		"status": enum.NotifyStatus_Success,
//	})
//	return err
//}

func (n *notifyRepoImpl) Retry(ctx context.Context, notifyID int64, nextTime time.Time) error {
	_, err := dao.NotifyRecord.WithContext(ctx).Where(dao.NotifyRecord.ID.Eq(notifyID)).Updates(map[string]interface{}{
		"notify_count":     dao.NotifyRecord.NotifyCount.Add(1),
		"last_notify_time": nextTime,
	})
	return err
}

func (n *notifyRepoImpl) Delete(ctx context.Context, lastTime time.Time) error {
	//record, err := dao.NotifyRecord.WithContext(ctx).Where(dao.NotifyRecord.CreatedAt.Lte(lastTime)).Order(dao.NotifyRecord.ID.Desc()).First()
	//if err != nil {
	//	return err
	//}
	_, err := dao.NotifyRecord.WithContext(ctx).Where(dao.NotifyRecord.CreatedAt.Lte(lastTime)).Delete()
	return err
}
