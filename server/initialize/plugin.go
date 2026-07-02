package initialize

import (
	// 通过空导入触发各插件 init() 自注册
	_ "github.com/lihongsheng/pay-gateway/plugin/example"
	_ "github.com/lihongsheng/pay-gateway/plugin/payment"

	"github.com/lihongsheng/pay-gateway/plugin"
)

// LoadPlugins 仅依赖 init() 完成注册；这里只是给个显式触点
func LoadPlugins() {
	_ = plugin.All() // 占位，确保包被引用
}
