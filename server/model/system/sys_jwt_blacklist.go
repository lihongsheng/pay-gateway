package system

import (
	"github.com/lihongsheng/pay-gateway/global"
)

type JwtBlacklist struct {
	global.GVA_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
