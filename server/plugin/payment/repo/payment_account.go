package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"sync"
)

type PaymentAccountRepo interface {
	GetFormCache(ctx context.Context, id int64) (*entity.PaymentAccount, error)
	GetByAccountNoFormCache(ctx context.Context, accNo string) (*entity.PaymentAccount, error)
	Save(ctx context.Context, app *dto.PaymentAccountCreateRequest, tx *gorm.DB) error
	GetByAppNo(ctx context.Context, appNo string, status []enum.MchStatus) ([]*entity.PaymentAccount, error)
	GetOnlyAppNO(ctx context.Context, appNo []string) ([]*model.PaymentAccount, error)
	GetOnlyAccountNo(ctx context.Context, accNo []string) ([]*model.PaymentAccount, error)
	SetPayValidateSuccess(ctx context.Context, accNo string) error
	GetByAppNoFormCache(ctx context.Context, appNo string, appAccountVersion string) ([]*entity.PaymentAccount, error)
	RefreshAppAccountsCache(ctx context.Context, appNo string, appAccountVersion string) error
}

type paymentAccountRepoImpl struct {
	db            *gorm.DB
	rdb           *redis.Client
	localMapCache *localMapCache
	singleflight  *singleflight.Group
}

func NewPaymentAccountRepo(db *gorm.DB, rdb *redis.Client) PaymentAccountRepo {
	return &paymentAccountRepoImpl{
		db:            db,
		rdb:           rdb,
		localMapCache: NewLocalMapCache(),
		singleflight:  &singleflight.Group{},
	}
}

func (p *paymentAccountRepoImpl) GetByAppNoFormCache(ctx context.Context, appNo string, appAccountVersion string) ([]*entity.PaymentAccount, error) {
	oldAppAccountVersion := p.localMapCache.GetAppAccountsVersion(ctx, appNo)
	if appAccountVersion != "" && oldAppAccountVersion == appAccountVersion {
		return p.localMapCache.GetAppAccounts(ctx, appNo), nil
	}
	if appAccountVersion != "" {
		err := p.RefreshAppAccountsCache(ctx, appNo, appAccountVersion)
		if err != nil {
			log.WithCtx(ctx).Error("refreshAppAccountsCache", zap.Error(err), zap.String("appNO", appNo), zap.String("appAccountVersion", appAccountVersion))
			return p.GetByAppNo(ctx, appNo, nil)
		}
		return p.localMapCache.GetAppAccounts(ctx, appNo), nil
	}
	return p.GetByAppNo(ctx, appNo, nil)
}

func (p *paymentAccountRepoImpl) RefreshAppAccountsCache(ctx context.Context, appNo string, appAccountVersion string) error {
	_, err, _ := p.singleflight.Do(appNo, func() (interface{}, error) {
		oldAppAccountVersion := p.localMapCache.GetAppAccountsVersion(ctx, appNo)
		if oldAppAccountVersion == appAccountVersion {
			return nil, nil
		}
		accounts, err := p.GetByAppNo(ctx, appNo, nil)
		if err != nil {
			return nil, err
		}
		p.localMapCache.SetAppAccounts(ctx, appNo, accounts)
		p.localMapCache.SetAppAccountsVersion(ctx, appNo, appAccountVersion)
		return nil, nil
	})
	return err
}

func (p *paymentAccountRepoImpl) SetPayValidateSuccess(ctx context.Context, accNo string) error {
	return p.db.WithContext(ctx).Model(&model.PaymentAccount{}).Where("account_no = ?", accNo).Update("validate_status", 1).Error
}

func (p *paymentAccountRepoImpl) GetOnlyAccountNo(ctx context.Context, accNo []string) ([]*model.PaymentAccount, error) {
	var app []*model.PaymentAccount
	err := p.db.WithContext(ctx).Where("account_no in (?)", accNo).Find(&app).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (p *paymentAccountRepoImpl) GetOnlyAppNO(ctx context.Context, appNo []string) ([]*model.PaymentAccount, error) {
	var app []*model.PaymentAccount
	err := p.db.WithContext(ctx).Where("app_no in (?)", appNo).Find(&app).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (p *paymentAccountRepoImpl) GetByAppNo(ctx context.Context, appNo string, status []enum.MchStatus) ([]*entity.PaymentAccount, error) {
	var app []*model.PaymentAccount
	query := p.db.WithContext(ctx).Where("app_no = ?", appNo)
	if len(status) > 0 {
		query = query.Where("status in ?", status)
	}
	err := query.Preload("PaymentMethod").Order("id DESC").Find(&app).Error
	if err != nil {
		return nil, err
	}
	entities := make([]*entity.PaymentAccount, 0, len(app))
	for _, item := range app {
		entities = append(entities, p.modelToEntity(item, item.PaymentMethod))
	}
	return entities, nil
}
func (p *paymentAccountRepoImpl) Save(ctx context.Context, appReq *dto.PaymentAccountCreateRequest, tx *gorm.DB) error {
	var account *model.PaymentAccount
	var err error
	var deleteMethods []int64
	if appReq.ID > 0 {
		account, err = p.get(ctx, appReq.ID)
		if err != nil {
			return err
		}
		account, deleteMethods = p.buildUpdate(account, appReq)
	} else {
		account = p.buildCreate(appReq)
	}
	if len(deleteMethods) > 0 {
		err = tx.Delete(&model.PaymentMethod{}, "id in ?", deleteMethods).Error
		if err != nil {
			return err
		}
		return nil
	}
	payments := account.PaymentMethod
	saveTx := tx.WithContext(ctx).Model(&model.PaymentAccount{})
	create := true
	if appReq.ID > 0 {
		create = false
		account.PaymentMethod = []*model.PaymentMethod{}
		err = saveTx.Where("id = ?", appReq.ID).Save(account).Error
	} else {
		err = saveTx.Create(account).Error
	}
	if err != nil {
		return err
	}
	if len(payments) > 0 && account.ID > 0 && !create {
		for i := range payments {
			payments[i].PaymentAccountID = account.ID
		}
		err = tx.Create(&account.PaymentMethod).Error
		if err != nil {
			return err
		}
		account.PaymentMethod = nil
	}
	p.rdb.Del(ctx, p.getCacheKeys(ctx, account)...)
	return err
}

func (p *paymentAccountRepoImpl) getCacheKeys(ctx context.Context, account *model.PaymentAccount) []string {
	return []string{
		p.getCacheIdKey(account.ID),
		p.getCacheNoKey(account.AccountNo),
	}
}

func (p *paymentAccountRepoImpl) buildUpdate(account *model.PaymentAccount, appReq *dto.PaymentAccountCreateRequest) (*model.PaymentAccount, []int64) {
	// 如果记录已存在，更新现有记录
	account.Name = appReq.Name
	account.Remark = appReq.Remark
	account.AppNo = appReq.AppNo
	account.Channel = appReq.Channel
	account.Status = appReq.Status
	account.Extend = appReq.Extend
	account.UpdatedAt = time.Now()
	account.ChannelConfig = appReq.ChannelConfig
	account.MaxLimit = appReq.MaxLimit
	if account.Extend == "" {
		account.Extend = "{}"
	}
	var deleteMethods = make([]int64, 0, len(appReq.PaymentMethod))
	paymentMethods := make([]*model.PaymentMethod, 0, len(appReq.PaymentMethod))
	for _, item1 := range appReq.PaymentMethod {
		exists := false
		for _, item2 := range account.PaymentMethod {
			if item2.PaymentMethod == item1.Method && item2.PaymentProduct == item1.Product {
				exists = true
				break
			}
		}
		if !exists {
			paymentMethods = append(paymentMethods, &model.PaymentMethod{
				PaymentMethod:  item1.Method,
				PaymentProduct: item1.Product,
			})
		}
	}
	for _, item1 := range account.PaymentMethod {
		exists := false
		for _, item2 := range appReq.PaymentMethod {
			if item1.PaymentMethod == item2.Method && item1.PaymentProduct == item2.Product {
				exists = true
				break
			}
		}
		if !exists {
			deleteMethods = append(deleteMethods, item1.ID)
		}
	}
	account.PaymentMethod = paymentMethods
	return account, deleteMethods
}

func (p *paymentAccountRepoImpl) buildCreate(appReq *dto.PaymentAccountCreateRequest) *model.PaymentAccount {
	account := &model.PaymentAccount{
		Name:           appReq.Name,
		Remark:         appReq.Remark,
		AppNo:          appReq.AppNo,
		MchNo:          appReq.MchNo,
		Channel:        appReq.Channel,
		Status:         appReq.Status,
		Extend:         appReq.Extend,
		CreatedAt:      time.Now(),
		ChannelConfig:  appReq.ChannelConfig,
		ValidateStatus: 0,
		AccountNo:      "P" + utils.GenDeviceID.Generate0X(),
		PaymentMethod:  make([]*model.PaymentMethod, 0, len(appReq.PaymentMethod)),
		MaxLimit:       appReq.MaxLimit,
	}
	if account.Extend == "" {
		account.Extend = "{}"
	}
	for _, method := range appReq.PaymentMethod {
		account.PaymentMethod = append(account.PaymentMethod, &model.PaymentMethod{
			PaymentMethod:  method.Method,
			PaymentProduct: method.Product,
		})
	}
	return account
}

func (p *paymentAccountRepoImpl) get(ctx context.Context, Id int64) (*model.PaymentAccount, error) {
	var app *model.PaymentAccount
	err := p.db.WithContext(ctx).Where("id = ?", Id).Preload("PaymentMethod").First(&app).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (p *paymentAccountRepoImpl) GetFormCache(ctx context.Context, id int64) (*entity.PaymentAccount, error) {
	cacheKey := p.getCacheIdKey(id)
	cacheData, _ := p.rdb.Get(ctx, cacheKey).Result()
	var app entity.PaymentAccount
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &app)
		if app.AppNo != "" {
			return &app, nil
		}
	}
	m, err := p.get(ctx, id)
	if err != nil {
		return nil, err
	}
	result := p.modelToEntity(m, m.PaymentMethod)
	p.rdb.Set(ctx, cacheKey, result, enum.ApplicationInfoCacheExpire)
	return result, nil
}

func (p *paymentAccountRepoImpl) getCacheIdKey(appId int64) string {
	return fmt.Sprintf(enum.AccountIDCacheKeys, appId)
}
func (p *paymentAccountRepoImpl) getCacheNoKey(accNo string) string {
	return fmt.Sprintf(enum.AccountNoCacheKeys, accNo)
}

func (p *paymentAccountRepoImpl) GetByAccountNoFormCache(ctx context.Context, accNo string) (*entity.PaymentAccount, error) {
	cacheKey := p.getCacheNoKey(accNo)
	cacheData, _ := p.rdb.Get(ctx, cacheKey).Result()
	var appEntity entity.PaymentAccount
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &appEntity)
		if appEntity.AppNo != "" {
			return &appEntity, nil
		}
	}
	var app *model.PaymentAccount
	err := p.db.WithContext(ctx).Where("account_no = ?", accNo).First(&app).Error
	if err != nil {
		return nil, err
	}
	var paymentMethods []*model.PaymentMethod
	err = p.db.WithContext(ctx).Where("payment_account_id = ?", app.ID).Find(&paymentMethods).Error
	if err != nil {
		return nil, err
	}
	result := p.modelToEntity(app, paymentMethods)
	p.rdb.Set(ctx, cacheKey, result, enum.ApplicationInfoCacheExpire)
	return result, nil
}

func (p *paymentAccountRepoImpl) modelToEntity(model *model.PaymentAccount, methods []*model.PaymentMethod) *entity.PaymentAccount {
	m := &entity.PaymentAccount{
		ID:             model.ID,
		Name:           model.Name,
		Remark:         model.Remark,
		AppNo:          model.AppNo,
		MchNo:          model.MchNo,
		Channel:        model.Channel,
		Status:         enum.MchStatus(model.Status),
		Extend:         model.Extend,
		ChannelConfig:  model.ChannelConfig,
		MaxLimit:       model.MaxLimit,
		ValidateStatus: model.ValidateStatus,
		AccountNo:      model.AccountNo,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		PaymentMethod:  make([]entity.PaymentMethod, 0, len(methods)),
	}
	for _, method := range methods {
		m.PaymentMethod = append(m.PaymentMethod, entity.PaymentMethod{
			Method:  method.PaymentMethod,
			Product: method.PaymentProduct,
		})
	}
	return m
}

type localMapCache struct {
	accountsMap *sync.Map
}

func NewLocalMapCache() *localMapCache {
	return &localMapCache{
		accountsMap: &sync.Map{},
	}
}

func (l *localMapCache) GetAppAccountsVersion(ctx context.Context, appNo string) string {
	key := fmt.Sprintf(enum.AppAccountsVersion, appNo)
	version, ok := l.accountsMap.Load(key)
	if !ok {
		return ""
	}
	return version.(string)
}

func (l *localMapCache) SetAppAccountsVersion(ctx context.Context, appNo string, version string) {
	key := fmt.Sprintf(enum.AppAccountsVersion, appNo)
	l.accountsMap.Store(key, version)
}

func (l *localMapCache) GetAppAccounts(ctx context.Context, appNo string) []*entity.PaymentAccount {
	key := fmt.Sprintf(enum.AppAccounts, appNo)
	accounts, ok := l.accountsMap.Load(key)
	if !ok {
		return nil
	}
	return accounts.([]*entity.PaymentAccount)
}

func (l *localMapCache) SetAppAccounts(ctx context.Context, appNo string, accounts []*entity.PaymentAccount) {
	key := fmt.Sprintf(enum.AppAccounts, appNo)
	l.accountsMap.Store(key, accounts)
}
