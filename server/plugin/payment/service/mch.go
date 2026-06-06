package service

import (
	"context"
	"errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
)

type MchService interface {
	Get(ctx context.Context, mchNo string) (*model.Merchant, error)
	Save(ctx context.Context, mch admin.MchCreateRequest, user *dto.User) error
	Search(ctx context.Context, mch admin.MchQueryRequest) ([]*model.Merchant, error)
	ChangeStatus(ctx context.Context, mchNo string, status enum.MchStatus) error
	Count(ctx context.Context, mch admin.MchQueryRequest) (int64, error)
}

type mchService struct {
	svc *svc.ServiceContext
}

func NewMchService(svc *svc.ServiceContext) MchService {
	return &mchService{
		svc: svc,
	}
}

// Get 根据商户ID获取商户信息
func (s *mchService) Get(ctx context.Context, mchNo string) (*model.Merchant, error) {
	m, err := s.svc.MchRepo.GetByMchNo(ctx, mchNo)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Save 保存商户信息
func (s *mchService) Save(ctx context.Context, mch admin.MchCreateRequest, user *dto.User) error {
	if mch.Validate() != nil {
		return mch.Validate()
	}
	if user.UserType.ISMch() {
		return errors.New("用户类型错误")
	}

	return s.svc.MchRepo.Save(ctx, mch)
}

// Search 搜索商户
func (s *mchService) Search(ctx context.Context, mch admin.MchQueryRequest) ([]*model.Merchant, error) {
	if err := mch.Validate(); err != nil {
		return nil, err
	}
	//if err := user.Validate(); err != nil {
	//	return nil, err
	//}
	//if user.UserType.ISMch() {
	//	mch.MchNo = user.HaveMchNo
	//}
	return s.svc.MchRepo.Search(ctx, mch)
}

// ChangeStatus 更改商户状态
func (s *mchService) ChangeStatus(ctx context.Context, mchNo string, status enum.MchStatus) error {
	mch, err := s.svc.MchRepo.GetByMchNo(ctx, mchNo)
	if err != nil {
		return err
	}
	err = s.svc.MchRepo.ChangeStatus(ctx, mch.ID, status)
	if err != nil {
		return err
	}
	return s.svc.AppRepo.DelAllAppCacheKeys(ctx, mch.MchNo)
}

// Count 统计符合条件的商户数量
func (s *mchService) Count(ctx context.Context, mch admin.MchQueryRequest) (int64, error) {

	return s.svc.MchRepo.Count(ctx, mch)
}
