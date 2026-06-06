// pkg/contextkey/context_key.go
package contextkey

import (
	"context"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
)

// 私有类型，避免外部冲突
type appInfoKey struct{}

// WithAppInfo 向标准 context.Context 注入应用信息
func WithAppInfo(ctx context.Context, appInfo *model.Application) context.Context {
	return context.WithValue(ctx, appInfoKey{}, appInfo)
}

// GetAppInfo 从标准 context.Context 提取应用信息
func GetAppInfo(ctx context.Context) *model.Application {
	val, ok := ctx.Value(appInfoKey{}).(*model.Application)
	if !ok {
		return nil
	}
	return val
}
