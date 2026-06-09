package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/dto"
	"github.com/lihongsheng/pay-gateway/utils"
)

type ApiGroup struct {
	MchApi
	AppApi
	AccountApi
	AggregateApi
	PublicPaymentApi
	PublicRefundApi
	TradeApi
	CashierApi
	DashboardApi
}

var ApiGroupApp = new(ApiGroup)

func GetUserInfo(c *gin.Context) *dto.User {
	claims, _ := utils.GetClaims(c)
	return &dto.User{
		UserID:    int64(claims.BaseClaims.ID),
		UserType:  claims.UserType,
		HaveMchNo: claims.HaveMchNO,
	}
}
