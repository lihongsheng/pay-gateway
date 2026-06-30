package service

import (
	"context"
	"errors"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
)

type ApplicationService interface {
	Get(ctx context.Context, appNo string) (*model.Application, error)
	Save(ctx context.Context, app *dto.ApplicationCreateRequest) (*model.Application, error)
	Search(ctx context.Context, req *dto.ApplicationQueryRequest) ([]*model.Application, error)
	Count(ctx context.Context, req *dto.ApplicationQueryRequest) (int64, error)
	ChangeStatus(ctx context.Context, appId int64, status enum.MchStatus, user *dto.User) error
}

type applicationService struct {
	appRepo repo.ApplicationRepo
	mchRepo repo.MchRepo
}

func NewApplicationService(appRepo repo.ApplicationRepo, mchRepo repo.MchRepo) ApplicationService {
	return &applicationService{
		appRepo: appRepo,
		mchRepo: mchRepo,
	}
}

// DefaultApplication 包级单例
var DefaultApplication ApplicationService

// Get 根据应用ID获取应用信息
func (s *applicationService) Get(ctx context.Context, appNo string) (*model.Application, error) {
	app, err := s.appRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// Save 保存应用信息
func (s *applicationService) Save(ctx context.Context, app *dto.ApplicationCreateRequest) (*model.Application, error) {
	if app == nil {
		return nil, errors.New("应用信息不能为空")
	}
	if err := app.Validate(); err != nil {
		return nil, err
	}
	var err error
	// 如果是商户用户，限制只能查询自己的应用
	_, err = s.mchRepo.GetByMchNo(ctx, app.MchNo)
	if err != nil {
		return nil, errors.New("指定的商户不存在")
	}

	return s.appRepo.Save(ctx, app)
}

// Search 搜索应用
func (s *applicationService) Search(ctx context.Context, req *dto.ApplicationQueryRequest) ([]*model.Application, error) {

	return s.appRepo.Search(ctx, req)
}

// Count 统计符合条件的应用数量
func (s *applicationService) Count(ctx context.Context, req *dto.ApplicationQueryRequest) (int64, error) {
	//if err := user.Validate(); err != nil {
	//	return 0, err
	//}
	//// 如果是商户用户，限制只能统计自己的应用
	//if user.ISMch() {
	//	// 验证商户是否存在
	//	mch, err := s.mchRepo.GetByMchNo(ctx, user.HaveMchNo)
	//	if err != nil {
	//		return 0, errors.New("指定的商户不存在")
	//	}
	//	req.MchNo = mch.MchNo
	//}

	return s.appRepo.Count(ctx, req)
}

// ChangeStatus 更改应用状态
func (s *applicationService) ChangeStatus(ctx context.Context, appId int64, status enum.MchStatus, user *dto.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	app, err := s.appRepo.Get(ctx, appId)
	if err != nil {
		return err
	}
	// 如果是商户用户，限制只能统计自己的应用
	if user.ISMch() {
		// 验证商户是否存在
		mch, err := s.mchRepo.GetByMchNo(ctx, user.HaveMchNo)
		if err != nil {
			return errors.New("指定的商户不存在")
		}
		// 验证用户是否有权限修改此应用的状态
		if app.MchNo != mch.MchNo {
			return errors.New("无权限修改此应用的状态")
		}
	}
	// 更新应用状态
	err = s.appRepo.ChangeStatus(ctx, appId, status)
	if err != nil {
		return err
	}
	// 清除相关缓存
	return s.appRepo.DelAppCache(ctx, app.AppNo, app.MchNo)
}
