package response

import "github.com/lihongsheng/pay-gateway/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
