package service

import (
	"context"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/config"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain"
	enum2 "github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	"github.com/lihongsheng/payment-sdk/enum"
	"github.com/lihongsheng/payment-sdk/enum/payment"
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

	global.DB = db
	zap.L() = zap.NewExample()
	global.Cfg.Payment.Salt = "1234561234561234"
	global.Redis = redis.NewClient(&redis.Options{
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
	global.Cfg.Payment.Salt = "17K9mP2n08rT4vY9"
	global.Cfg.Payment.WebHost = "https://test.web.payment.jianxindianzi.com"
	global.Cfg.Payment.ApiHost = "https://test.api.payment.jianxindianzi.com"
	global.Cfg.Payment.ProxyNotifyPrefix = "/sslab"
}

func TestAggregateService_GetRedirectUrl(t *testing.T) {
	testInit()
	c := NewAggregateService(svc.ServiceContextApp, domain.Service)
	ctx := context.Background()
	// "PDP5VEDL6QMD42LLUA4L7W22AJXPLWWTU576P7FHSRZFPW7FWBV2HD62H7UCQMKTG2BBOXF2G27DFW2F", payment.Payment_Wechat, payment.PaymentProduct_JSAPI
	req := &public.AggregateRedirectUrlRequest{
		Token:          "PDP5VEDL6QMD42LLUA4L7W22AJXPLWWTU576P7FHSRZFPW7FWBV2HD62H7UCQMKTG2BBOXF2G27DFW2F",
		PaymentMethod:  payment.Payment_Wechat,
		PaymentProduct: payment.PaymentProduct_JSAPI,
		Device:         enum.Device_Wechat,
	}
	r, err := c.GetRedirectUrl(ctx, req, enum2.PaymentBaseIndexPath)
	assert.NoError(t, err)
	t.Log(r)
}

func TestAggregateService_GetUserOpenID(t *testing.T) {
	testInit()
	c := NewAggregateService(svc.ServiceContextApp, domain.Service)
	ctx := context.Background()
	r, err := c.GetUserOpenID(ctx, "PDP5VEDL6QMD42LLUA4L7W22AJXPLWWTU576P7FHSRZFPW7FWBV2HD62H7UCQMKTG2BBOXF2G27DFW2F", "PDP5VEDL6QMD42LLUA4L7W22AJXPLWWTU576P7FHSRZFPW7", "P4966c7e946c00")
	assert.NoError(t, err)
	t.Log(r)
}

func TestAggregateService_GetApplication(t *testing.T) {
	testInit()
	c := NewAggregateService(svc.ServiceContextApp, domain.Service)
	ctx := context.Background()
	r, err := c.GetApplication(ctx, "PDP5VEDL6QMD42LLUA4L7W22AJXPLWWTU576P7FHSRZFPW7FWBV2HD62H7UCQMKTG2BBOXF2G27DFW2F")
	assert.NoError(t, err)
	t.Log(r)
}

func TestAggregateService_Payment(t *testing.T) {
	testInit()
	//c := NewAggregateService(svc.ServiceContextApp, domain.Service)
	//ctx := context.Background()
	//token := "PDP5VEDL6QMD42LLUA4L7W22AJXPLWWTU576P7FHSRZFPW7FWBV2HD62H7UCQMKTG2BBOXF2G27DFW2F"
	//r, err := c.Payment(ctx, &public.AggregateOrder{
	//	Token:          token,
	//	OrderNo:        fmt.Sprintf("P%d%s%s", time.Now().UnixMilli(), utils.RandomString(4), utils.HashCrc(token, 4)),
	//	Amount:         1,
	//	Device:         enum.Device_Wechat,
	//	PaymentMethod:  payment.Payment_Wechat,
	//	PaymentProduct: payment.PaymentProduct_JSAPI,
	//	OpenID:         "ogdvH6h9jPp5R3f1fyLsQjdB-fAc",
	//	AccountNo:      "P4966c7e946c00",
	//	RequestID:      fmt.Sprintf("P%d%s%s", time.Now().UnixMilli(), utils.RandomString(4), utils.HashCrc(token, 4)),
	//})
	//assert.NoError(t, err)
	//t.Log(r)
}
