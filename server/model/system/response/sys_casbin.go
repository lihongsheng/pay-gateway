package response

import (
	"github.com/lihongsheng/pay-gateway/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
