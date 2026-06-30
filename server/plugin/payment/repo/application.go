package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ApplicationRepo interface {
	Get(ctx context.Context, appId int64) (*model.Application, error)
	Save(ctx context.Context, app *dto.ApplicationCreateRequest) (*model.Application, error)
	GetByAppNoFormCache(ctx context.Context, appNo string) (*model.Application, error)
	Search(ctx context.Context, req *dto.ApplicationQueryRequest) ([]*model.Application, error)
	Count(ctx context.Context, req *dto.ApplicationQueryRequest) (int64, error)
	ChangeStatus(ctx context.Context, appId int64, status enum.MchStatus) error
	GetAppByNoAndMchNoCache(ctx context.Context, appNo string, mchNo string) (*dto.ApplicationCacheInfo, error)
	DelAppCache(ctx context.Context, appNo string, mchNo string) error
	GetAllAppCacheKeys(ctx context.Context, mchNo string) ([]string, error)
	DelAllAppCacheKeys(ctx context.Context, mchNo string) error
	UpdateApplicationAccountVersion(ctx context.Context, appNo string, version string, tx *gorm.DB) error
	GetByAppNo(ctx context.Context, appNo string) (*model.Application, error)
}

type applicationRepoImpl struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewApplicationRepo(db *gorm.DB, rdb *redis.Client) ApplicationRepo {
	return &applicationRepoImpl{db: db, rdb: rdb}
}

func (a *applicationRepoImpl) GetByAppNo(ctx context.Context, appNo string) (*model.Application, error) {
	var app model.Application
	err := a.db.WithContext(ctx).Where("app_no = ?", appNo).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (a *applicationRepoImpl) UpdateApplicationAccountVersion(ctx context.Context, appNo string, version string, tx *gorm.DB) error {
	if appNo == "" {
		return errors.New("appNo is empty")
	}
	if version == "" {
		return errors.New("version is empty")
	}
	var app model.Application
	err := tx.WithContext(ctx).Model(&model.Application{}).Where("app_no = ?", appNo).First(&app).Error
	if err != nil {
		return err
	}
	err = tx.WithContext(ctx).Model(&model.Application{}).Where("app_no = ?", appNo).Update("account_version", version).Error
	if err != nil {
		return err
	}
	err = a.DelAppCache(ctx, app.AppNo, app.MchNo)
	if err != nil {
		return err
	}
	// 延迟删除
	go func() {
		time.Sleep(time.Second)
		_ = a.DelAppCache(ctx, app.AppNo, app.MchNo)
	}()
	return err
}

func (a *applicationRepoImpl) Get(ctx context.Context, appId int64) (*model.Application, error) {
	var app model.Application
	err := a.db.WithContext(ctx).Where("id = ?", appId).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (a *applicationRepoImpl) Save(ctx context.Context, app *dto.ApplicationCreateRequest) (*model.Application, error) {
	if app.Validate() != nil {
		return nil, app.Validate()
	}
	var modelApp *model.Application
	var err error
	if app.ID > 0 {
		// 更新现有应用
		modelApp, err = a.Get(ctx, app.ID)
		if err != nil {
			return nil, err
		}
		// 创建新应用
		modelApp.AppName = app.AppName
		modelApp.Desc = app.Desc
		modelApp.PaymentTitle = app.PaymentTitle
		modelApp.PayIcon = app.PayIcon
		modelApp.Status = int64(app.Status)
		modelApp.IsCustomerDomain = app.IsCustomerDomain
		modelApp.CustomerDomain = app.CustomerDomain
		modelApp.ProxyHost = app.ProxyHost
		modelApp.ProxyPort = app.ProxyPort
		modelApp.ProxyUser = app.ProxyUser
		modelApp.ProxyPwd = app.ProxyPwd
		modelApp.MultiChannel = app.MultiChannel
		modelApp.UpdatedAt = time.Now()
	} else {
		secret, _ := utils.GenAesBase64Str(32)
		// 创建新应用
		modelApp = &model.Application{
			AppName:          app.AppName,
			Secret:           secret,
			AppNo:            "A" + utils.GenDeviceID.Generate0X(),
			MchNo:            app.MchNo,
			Desc:             app.Desc,
			PaymentTitle:     app.PaymentTitle,
			PayIcon:          app.PayIcon,
			Status:           int64(app.Status),
			IsCustomerDomain: app.IsCustomerDomain,
			CustomerDomain:   app.CustomerDomain,
			ProxyHost:        app.ProxyHost,
			ProxyPort:        app.ProxyPort,
			ProxyUser:        app.ProxyUser,
			ProxyPwd:         app.ProxyPwd,
		}
	}
	err = a.db.WithContext(ctx).Save(modelApp).Error
	if err != nil {
		return nil, err
	}
	err = a.DelAppCache(ctx, modelApp.AppNo, modelApp.MchNo)
	if err != nil {
		return nil, err
	}
	return modelApp, nil
}

func (a *applicationRepoImpl) GetByAppNoFormCache(ctx context.Context, appNo string) (*model.Application, error) {
	key := a.getAppNoCacheKeys(appNo)
	var app model.Application
	if a.rdb.Exists(ctx, key).Val() > 0 {
		cacheData, _ := a.rdb.Get(ctx, key).Result()
		if cacheData != "" {
			err := json.Unmarshal([]byte(cacheData), &app)
			if err != nil {
				return nil, err
			}
			if app.ID > 0 {
				return &app, nil
			}
		}
	}
	err := a.db.WithContext(ctx).Where("app_no = ?", appNo).First(&app).Error
	if err != nil {
		return nil, err
	}
	a.rdb.Set(ctx, key, app, enum.ApplicationInfoCacheExpire)
	return &app, nil
}

func (a *applicationRepoImpl) Search(ctx context.Context, req *dto.ApplicationQueryRequest) ([]*model.Application, error) {
	var apps []*model.Application
	query := a.buildQuery(ctx, req)

	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.Page > 0 {
		query = query.Offset((req.Page - 1) * req.PageSize)
	}
	query = query.Limit(req.PageSize)

	err := query.Find(&apps).Error
	return apps, err
}

func (a *applicationRepoImpl) Count(ctx context.Context, req *dto.ApplicationQueryRequest) (int64, error) {
	var count int64
	query := a.buildQuery(ctx, req)
	err := query.Count(&count).Error
	return count, err
}

func (a *applicationRepoImpl) ChangeStatus(ctx context.Context, appId int64, status enum.MchStatus) error {
	result := a.db.WithContext(ctx).Model(&model.Application{}).
		Where("id = ?", appId).
		Update("status", status)

	return result.Error
}

func (a *applicationRepoImpl) GetAppByNoAndMchNoCache(ctx context.Context, appNo string, mchNo string) (*dto.ApplicationCacheInfo, error) {
	cacheKey := a.getAppNoCacheKeys(appNo)
	var app *model.Application
	var result *dto.ApplicationCacheInfo
	// 尝试从缓存获取
	cacheData, err := a.rdb.Get(ctx, cacheKey).Result()
	if err == nil && cacheData != "" {
		err = json.Unmarshal([]byte(cacheData), result)
		if err == nil {
			return result, nil
		}
	}
	// 缓存未命中，从数据库查询
	app, err = a.GetByAppNoAndMchNo(ctx, appNo, mchNo)
	if err != nil {
		return nil, err
	}
	var mch *model.Merchant
	err = a.db.WithContext(ctx).Model(&model.Merchant{}).Where("mch_no = ?", app.MchNo).First(&mch).Error
	if err != nil {
		return nil, err
	}

	result = &dto.ApplicationCacheInfo{
		ID:               app.ID,
		AppName:          app.AppName,
		Secret:           app.Secret,
		AppNo:            app.AppNo,
		MchNo:            app.MchNo,
		PaymentTitle:     app.PaymentTitle,
		PayIcon:          app.PayIcon,
		Status:           enum.MchStatus(app.Status),
		IsCustomerDomain: app.IsCustomerDomain,
		CustomerDomain:   app.CustomerDomain,
		ProxyHost:        app.ProxyHost,
		ProxyPort:        app.ProxyPort,
		ProxyUser:        app.ProxyUser,
		ProxyPwd:         app.ProxyPwd,
		MchInfo: &dto.ApplicationMchCacheInfo{
			ID:      mch.ID,
			MchName: mch.MchName,
			MchNo:   mchNo,
			Status:  enum.MchStatus(mch.Status),
		},
	}
	// 存入缓存
	jsonData, _ := json.Marshal(result)
	a.rdb.Set(ctx, cacheKey, string(jsonData), enum.ApplicationInfoCacheExpire)

	return result, nil
}

// 辅助方法：根据应用号和商户号查询
func (a *applicationRepoImpl) GetByAppNoAndMchNo(ctx context.Context, appNo string, mchNo string) (*model.Application, error) {
	var app model.Application
	err := a.db.WithContext(ctx).Where("app_no = ? AND mch_no = ?", appNo, mchNo).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (a *applicationRepoImpl) DelAllAppCacheKeys(ctx context.Context, mchNo string) error {
	cacheKeys, err := a.GetAllAppCacheKeys(ctx, mchNo)
	if err != nil {
		return err
	}
	if len(cacheKeys) == 0 {
		return nil
	}
	err = a.rdb.Del(ctx, cacheKeys...).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("删除应用缓存失败: %s", err.Error()))
	}

	return nil
}

func (a *applicationRepoImpl) DelAppCache(ctx context.Context, appNo string, mchNo string) error {
	cacheKeys := a.GetAppCacheKeys(appNo, mchNo)
	if len(cacheKeys) == 0 {
		return nil
	}
	err := a.rdb.Del(ctx, cacheKeys...).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("删除应用缓存失败: %s", err.Error()))
	}

	return nil
}

func (a *applicationRepoImpl) GetAppCacheKeys(appNo string, mchNo string) []string {
	return []string{
		a.getAppNoCacheKeys(appNo),
	}
}

func (a *applicationRepoImpl) getAppNoCacheKeys(appNo string) string {
	return fmt.Sprintf(enum.ApplicationInfoCacheKeys, appNo)
}

func (a *applicationRepoImpl) GetAllAppCacheKeys(ctx context.Context, mchNo string) ([]string, error) {
	var apps []*model.Application
	err := a.db.WithContext(ctx).Where("mch_no = ?", mchNo).Find(&apps).Error
	if err != nil {
		return nil, err
	}
	var cacheKeys = make([]string, len(apps))
	for i, app := range apps {
		cacheKeys[i] = a.getAppNoCacheKeys(app.AppNo)
	}
	return cacheKeys, nil
}

// 构建查询条件
func (a *applicationRepoImpl) buildQuery(ctx context.Context, req *dto.ApplicationQueryRequest) *gorm.DB {
	query := a.db.WithContext(ctx).Model(&model.Application{})

	if req.AppNo != "" {
		query = query.Where("app_no = ?", "%"+req.AppNo+"%")
	}
	if req.AppName != "" {
		query = query.Where("app_name LIKE ?", "%"+req.AppName+"%")
	}
	if req.MchNo != "" {
		query = query.Where("mch_no = ?", req.MchNo)
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}
	if len(req.AppNos) > 0 {
		query = query.Where("app_no in (?)", req.AppNos)
	}
	return query
}
