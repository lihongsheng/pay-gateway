package api

import "github.com/lihongsheng/pay-gateway/plugin/announcement/service"

var (
	Api         = new(api)
	serviceInfo = service.Service.Info
)

type api struct{ Info info }
