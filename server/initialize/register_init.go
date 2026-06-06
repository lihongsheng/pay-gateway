package initialize

import (
	_ "github.com/lihongsheng/pay-gateway/source/example"
	_ "github.com/lihongsheng/pay-gateway/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
