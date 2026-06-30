package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	payDto "github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/utils/jwt"
)

func GetUserInfo(c *gin.Context) *payDto.User {
	u, err := jwt.GetUser(c.Request.Context())
	if err != nil {
		return &payDto.User{}
	}
	var haveMchNo string
	if u.MchID > 0 {
		haveMchNo = fmt.Sprintf("%d", u.MchID)
	}
	return &payDto.User{
		UserID:    int64(u.ID),
		UserType:  u.SystemType,
		HaveMchNo: haveMchNo,
	}
}
