package service

import (
	"github.com/lihongsheng/pay-gateway/service/example"
	"github.com/lihongsheng/pay-gateway/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
}
