package repo

import (
	"context"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"gorm.io/gorm"
	"time"
)

type MchRepo interface {
	Get(ctx context.Context, mchId int64) (*model.Merchant, error)
	Save(ctx context.Context, mch admin.MchCreateRequest) error
	Search(ctx context.Context, mch admin.MchQueryRequest) ([]*model.Merchant, error)
	Count(ctx context.Context, mch admin.MchQueryRequest) (int64, error)
	ChangeStatus(ctx context.Context, mchId int64, status enum.MchStatus) error
	GetByMchNo(ctx context.Context, mchNo string) (*model.Merchant, error)
}

type mchRepoImpl struct {
}

func NewMchRepo() MchRepo {
	return &mchRepoImpl{}
}

func (m *mchRepoImpl) Get(ctx context.Context, mchId int64) (*model.Merchant, error) {
	var mdl model.Merchant
	err := global.GVA_PAY_DB.WithContext(ctx).Where("id = ?", mchId).First(&mdl).Error
	if err != nil {
		return nil, err
	}
	return &mdl, nil
}

func (m *mchRepoImpl) Save(ctx context.Context, mch admin.MchCreateRequest) (err error) {
	var mdl *model.Merchant
	if mch.Validate() != nil {
		return mch.Validate()
	}
	if mch.ID > 0 {
		mdl, err = m.Get(ctx, mch.ID)
		if err != nil {
			return err
		}
		mdl.MchName = mch.MchName
		mdl.Linker = mch.Linker
		mdl.Phone = mch.Phone
		mdl.Email = mch.Email
		mdl.Address = mch.Address
		mdl.Reason = mch.Reason
		mdl.Status = int64(mch.Status)
		mdl.UpdatedAt = time.Now()
	} else {
		mdl = &model.Merchant{
			MchNo:   "M" + utils.GenDeviceID.Generate0X(),
			MchName: mch.MchName,
			Linker:  mch.Linker,
			Phone:   mch.Phone,
			Email:   mch.Email,
			Address: mch.Address,
			Reason:  mch.Reason,
			Status:  int64(mch.Status),
		}
	}
	return global.GVA_PAY_DB.WithContext(ctx).Save(mdl).Error
}

func (m *mchRepoImpl) Search(ctx context.Context, mch admin.MchQueryRequest) ([]*model.Merchant, error) {
	// 实现搜索功能，根据请求参数进行查询
	var merchants []*model.Merchant
	query := m.buildQuery(ctx, mch)
	if mch.PageSize < 1 {
		mch.PageSize = 10
	}
	if mch.Page > 0 {
		query = query.Offset(mch.PageSize * (mch.Page - 1))
	}
	if mch.PageSize > 0 {
		query = query.Limit(mch.PageSize)
	}
	err := query.Find(&merchants).Error
	return merchants, err
}
func (m *mchRepoImpl) buildQuery(ctx context.Context, mch admin.MchQueryRequest) *gorm.DB {
	query := global.GVA_PAY_DB.WithContext(ctx).Model(&model.Merchant{})
	if mch.MchNo != "" {
		query = query.Where("mch_no = ?", "%"+mch.MchNo+"%")
	}
	if mch.MchName != "" {
		query = query.Where("mch_name LIKE ?", "%"+mch.MchName+"%")
	}
	if mch.ID != 0 {
		query = query.Where("id = ?", mch.ID)
	}
	if mch.Status > 0 {
		query = query.Where("status = ?", mch.Status)
	}
	if mch.IDList != nil {
		query = query.Where("id IN ?", mch.IDList)
	}
	if len(mch.MchNos) > 0 {
		query = query.Where("mch_no IN ?", mch.MchNos)
	}
	return query
}
func (m *mchRepoImpl) Count(ctx context.Context, mch admin.MchQueryRequest) (int64, error) {
	var count int64
	query := m.buildQuery(ctx, mch)
	err := query.Count(&count).Error
	return count, err
}

func (m *mchRepoImpl) ChangeStatus(ctx context.Context, mchId int64, status enum.MchStatus) error {
	mdl, err := m.Get(ctx, mchId)
	if err != nil {
		return err
	}
	if mdl.Status == int64(status) {
		return nil
	}
	result := global.GVA_PAY_DB.WithContext(ctx).Model(&model.Merchant{}).
		Where("id = ?", mchId).
		Update("status", status)

	return result.Error
}

func (m *mchRepoImpl) GetByMchNo(ctx context.Context, mchNo string) (*model.Merchant, error) {
	var mdl model.Merchant
	err := global.GVA_PAY_DB.WithContext(ctx).Where("mch_no = ?", mchNo).First(&mdl).Error
	if err != nil {
		return nil, err
	}
	return &mdl, nil
}
