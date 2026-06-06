package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type RouterRepo interface {
	GetAppRouterStatisticFromCache(ctx context.Context, appNo string, t time.Time) ([]*model.RouterAccountStatistic, error)
	Save(ctx context.Context, statics *model.RouterAccountStatistic) error
	DeleteByDate(ctx context.Context, t time.Time) error
	UpdateCache(ctx context.Context, appNo string, t time.Time) error
}

type routerRepoImpl struct {
}

func NewRouterRepo() RouterRepo {
	return &routerRepoImpl{}
}

func (r *routerRepoImpl) UpdateCache(ctx context.Context, appNo string, t time.Time) error {
	var app []*model.RouterAccountStatistic
	err := global.GVA_PAY_DB.WithContext(ctx).Where("app_no = ? and statistic_date = ?", appNo, t.Format("2006-01-02")).Find(&app).Error
	if err != nil {
		return err
	}
	key := r.getCacheKey(appNo, t)
	if len(app) > 0 {
		cache, err := json.Marshal(app)
		if err != nil {
			log.WithCtx(ctx).Error("MarshalCacheUpdateRouterCache", zap.Error(err), zap.String("cacheKey", key), zap.String("appNO", appNo), zap.Time("statisticDate", t))
			return err
		}
		return global.GVA_REDIS.Set(ctx, key, string(cache), enum.RouterStatisticCacheExpire).Err()
	}
	return nil
}

func (r *routerRepoImpl) DeleteByDate(ctx context.Context, t time.Time) error {
	return global.GVA_PAY_DB.WithContext(ctx).Where("statistic_date < ?", t.Format("2006-01-02")).Delete(&model.RouterAccountStatistic{}).Error
}

func (r *routerRepoImpl) GetAppRouterStatisticFromCache(ctx context.Context, appNo string, t time.Time) ([]*model.RouterAccountStatistic, error) {
	var app []*model.RouterAccountStatistic
	key := r.getCacheKey(appNo, t)
	lockKey := fmt.Sprintf(enum.RouterStatisticLockCacheKeys, appNo, t.Format("2006-01-02"))

	cacheStr, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		log.WithCtx(ctx).Error("GetAppRouterStatisticFromCacheError", zap.Error(err), zap.String("appNO", appNo))
	}
	if cacheStr != "" {
		err := json.Unmarshal([]byte(cacheStr), &app)
		if err != nil {
			log.WithCtx(ctx).Error("UnmarshalCacheRouterStatistic", zap.Error(err), zap.String("cacheKey", cacheStr), zap.String("appNO", appNo))
		}
		if len(app) > 0 {
			return app, nil
		}
	}
	lock, err := global.GVA_REDIS.SetNX(ctx, lockKey, time.Now(), time.Second*2).Result()
	if err != nil && err != redis.Nil {
		log.WithCtx(ctx).Error("GetAppRouterStatisticFromCacheLockError", zap.Error(err), zap.String("appNO", appNo))
	}
	defer func() {
		_ = global.GVA_REDIS.Del(ctx, lockKey)
	}()
	retry := 3
	if !lock && err == nil {
		for i := 0; i < retry; i++ {
			time.Sleep(200 * time.Millisecond)
			cacheStr, err = global.GVA_REDIS.Get(ctx, key).Result()
			if err == nil && len(cacheStr) > 0 {
				_ = json.Unmarshal([]byte(cacheStr), &app)
				if len(app) > 0 {
					return app, nil
				}
			}
		}
	}
	err = global.GVA_PAY_DB.WithContext(ctx).Where("app_no = ? and statistic_date = ?", appNo, t.Format("2006-01-02")).Find(&app).Error
	if err != nil {
		return nil, err
	}
	if len(app) > 0 {
		cache, err := json.Marshal(app)
		if err != nil {
			log.WithCtx(ctx).Error("MarshalCacheRouterStatistic", zap.Error(err), zap.String("cacheKey", key), zap.String("appNO", appNo))
		} else {
			_ = global.GVA_REDIS.Set(ctx, key, string(cache), enum.RouterStatisticCacheExpire)
		}
	}
	return app, nil
}

func (r *routerRepoImpl) getCacheKey(appNo string, t time.Time) string {
	return fmt.Sprintf(enum.RouterStatisticCacheKeys, appNo, t.Format("2006-01-02"))
}

func (r *routerRepoImpl) Save(ctx context.Context, statics *model.RouterAccountStatistic) error {
	var updates = map[string]interface{}{
		"total_requests":   gorm.Expr("total_requests + ?", statics.TotalRequests),
		"success_requests": gorm.Expr("success_requests + ?", statics.SuccessRequests),
		"failure_requests": gorm.Expr("failure_requests + ?", statics.FailureRequests),
		"success_amount":   gorm.Expr("success_amount + ?", statics.SuccessAmount),
		"updated_at":       time.Now(),
	}
	if statics.PaymentLimit != 0 {
		updates["payment_limit"] = statics.PaymentLimit
	}
	if statics.UserLimit != 0 {
		updates["user_limit"] = statics.UserLimit
	}
	// 使用 GORM 的 OnConflict 子句实现插入或更新，并累加 total_requests
	result := global.GVA_PAY_DB.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "app_no"},
				{Name: "account_no"},
				{Name: "statistic_date"},
			},
			DoUpdates: clause.Assignments(updates),
		},
	).Save(statics)
	global.GVA_REDIS.Del(ctx, r.getCacheKey(statics.AppNo, statics.StatisticDate))
	return result.Error
}
