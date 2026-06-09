package router

import (
	"github.com/lihongsheng/pay-gateway/router/example"
	"github.com/lihongsheng/pay-gateway/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
}
