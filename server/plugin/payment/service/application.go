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

type Application interface {
	Get(ctx context.Context, appNo string) (*model.Application, error)
	Save(ctx context.Context, app *admin.ApplicationCreateRequest) (*model.Application, error)
	Search(ctx context.Context, req *admin.ApplicationQueryRequest) ([]*model.Application, error)
	Count(ctx context.Context, req *admin.ApplicationQueryRequest) (int64, error)
	ChangeStatus(ctx context.Context, appId int64, status enum.MchStatus, user *dto.User) error
}

type ApplicationService struct {
	svc *svc.ServiceContext
}

func NewApplicationService(svc *svc.ServiceContext) Application {
	return &ApplicationService{
		svc: svc,
	}
}

// Get 根据应用ID获取应用信息
func (s *ApplicationService) Get(ctx context.Context, appNo string) (*model.Application, error) {
	app, err := s.svc.AppRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// Save 保存应用信息
func (s *ApplicationService) Save(ctx context.Context, app *admin.ApplicationCreateRequest) (*model.Application, error) {
	if app == nil {
		return nil, errors.New("应用信息不能为空")
	}
	if err := app.Validate(); err != nil {
		return nil, err
	}
	var err error
	// 如果是商户用户，限制只能查询自己的应用
	_, err = s.svc.MchRepo.GetByMchNo(ctx, app.MchNo)
	if err != nil {
		return nil, errors.New("指定的商户不存在")
	}

	return s.svc.AppRepo.Save(ctx, app)
}

// Search 搜索应用
func (s *ApplicationService) Search(ctx context.Context, req *admin.ApplicationQueryRequest) ([]*model.Application, error) {

	return s.svc.AppRepo.Search(ctx, req)
}

// Count 统计符合条件的应用数量
func (s *ApplicationService) Count(ctx context.Context, req *admin.ApplicationQueryRequest) (int64, error) {
	//if err := user.Validate(); err != nil {
	//	return 0, err
	//}
	//// 如果是商户用户，限制只能统计自己的应用
	//if user.UserType.ISMch() {
	//	// 验证商户是否存在
	//	mch, err := s.svc.MchRepo.GetByMchNo(ctx, user.HaveMchNo)
	//	if err != nil {
	//		return 0, errors.New("指定的商户不存在")
	//	}
	//	req.MchNo = mch.MchNo
	//}

	return s.svc.AppRepo.Count(ctx, req)
}

// ChangeStatus 更改应用状态
func (s *ApplicationService) ChangeStatus(ctx context.Context, appId int64, status enum.MchStatus, user *dto.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	app, err := s.svc.AppRepo.Get(ctx, appId)
	if err != nil {
		return err
	}
	// 如果是商户用户，限制只能统计自己的应用
	if user.UserType.ISMch() {
		// 验证商户是否存在
		mch, err := s.svc.MchRepo.GetByMchNo(ctx, user.HaveMchNo)
		if err != nil {
			return errors.New("指定的商户不存在")
		}
		// 验证用户是否有权限修改此应用的状态
		if app.MchNo != mch.MchNo {
			return errors.New("无权限修改此应用的状态")
		}
	}
	// 更新应用状态
	err = s.svc.AppRepo.ChangeStatus(ctx, appId, status)
	if err != nil {
		return err
	}
	// 清除相关缓存
	return s.svc.AppRepo.DelAppCache(ctx, app.AppNo, app.MchNo)
}
