package repo

import (
	"context"
	"encoding/json"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPaymentOrderRepoImpl_GetOrderWithCache(t *testing.T) {
	testInit()
	// 关键：指定时区解析时间（比如东八区）
	time.LoadLocation("Asia/Shanghai") // 东八区
	client := NewPaymentOrderRepo(global.DB, global.Redis)
	ctx := context.Background()
	r, err := client.GetOrderWithCache(ctx, public.QueryPaymentRequest{
		AppNo:   "A4966c7e946c00",
		MchNo:   "M4954d4475f800",
		OrderNo: "PT177261780846304242494",
	})
	assert.NoError(t, err)
	by, _ := json.Marshal(r)
	t.Log(string(by))
}
