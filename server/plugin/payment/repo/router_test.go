package repo

import (
	"context"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRouterRepoImpl_Save(t *testing.T) {
	testInit()
	repo := NewRouterRepo(global.DB, global.Redis)
	// 关键：指定时区解析时间（比如东八区）
	loc, _ := time.LoadLocation("Asia/Shanghai") // 东八区
	// 方式1：带时区解析
	tt, _ := time.ParseInLocation(time.DateTime, "2026-02-27 16:36:39", loc)

	err := repo.Save(context.Background(), &model.RouterAccountStatistic{
		ID:              0,
		AppNo:           "A4966c7e946c00",
		AccountNo:       "P4966c7e946c00",
		StatisticDate:   tt,
		TotalRequests:   1,
		SuccessRequests: 1,
		FailureRequests: 0,
		UserLimit:       0,
		SuccessAmount:   1,
		PaymentLimit:    0,
	})
	assert.NoError(t, err)
}
