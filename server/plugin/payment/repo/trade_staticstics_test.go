package repo

import (
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/config"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
	"time"
)

func testInit() {
	db, _ := gorm.Open(mysql.Open("root:J4l71RPO14@tcp(47.110.88.165:2624)/payment?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"), &gorm.Config{
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})

	global.GVA_PAY_DB = db
	global.GVA_LOG = zap.NewExample()
	config.Config.Salt = "1234561234561234"
	global.GVA_REDIS = redis.NewClient(&redis.Options{
		Addr:            "47.110.88.165:6813",
		Username:        "",
		Password:        "f392Vt3e7O", // no password set
		DB:              1,            // use default DB
		MaxIdleConns:    30,
		MinIdleConns:    5,
		ReadTimeout:     2 * time.Second,
		WriteTimeout:    3 * time.Second,
		MaxRetries:      1,
		MaxRetryBackoff: 128 * time.Millisecond,
		DialTimeout:     3 * time.Second,
	})
	config.Config.Salt = "17K9mP2n08rT4vY9"
	config.Config.WebHost = "https://test.web.payment.jianxindianzi.com"
	config.Config.ApiHost = "https://test.api.payment.jianxindianzi.com"
	config.Config.ProxyNotifyPrefix = "/sslab"
}
func TestTradeStaticsRepoImpl_Save(t *testing.T) {
	testInit()
	// 关键：指定时区解析时间（比如东八区）
	loc, _ := time.LoadLocation("Asia/Shanghai") // 东八区
	// 方式1：带时区解析
	tt, _ := time.ParseInLocation(time.DateTime, "2026-02-27 16:36:39", loc)
	tradeStaticsRepo := NewTradeStaticsRepo()
	fmt.Println(tt)
	d := &model.TradeStatistic{
		ID:              0,
		AppNo:           "A4966c7e946c00",
		AccountNo:       "P4966c7e946c00",
		MchNo:           "M4954d4475f800",
		StatisticDate:   tt,
		TotalRequests:   0,
		SuccessRequests: 1,
		TotalAmount:     0,
		SuccessAmount:   1,
		RefundOrder:     0,
		RefundAmount:    0,
	}
	err := tradeStaticsRepo.Save(nil, d, global.GVA_PAY_DB)
	assert.NoError(t, err)
}

func TestTradeStaticsRepoImpl_Get(t *testing.T) {
	testInit()
	tradeStaticsRepo := NewTradeStaticsRepo()
	d, err := tradeStaticsRepo.Get(nil, "M4954d4475f800", "A4966c7e946c00", time.Now().Add(-time.Hour*24*7), time.Now())
	assert.NoError(t, err)
	fmt.Println(d)
}

func TestTradeStaticsRepoImpl_CountGroup(t *testing.T) {
	testInit()
	tradeStaticsRepo := NewTradeStaticsRepo()
	d, err := tradeStaticsRepo.CountGroup(nil, time.Now().Add(-time.Hour*24*7), time.Now())
	assert.NoError(t, err)
	fmt.Println(fmt.Sprintf("d: %+v", d))
	if len(d) > 0 {
		fmt.Println(fmt.Sprintf("d[0]: %+v", d[0]))
	}
}
