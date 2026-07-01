package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/utils/response"
)

// TotalRequest 获取商户总订单数
// @Tags      DashboardApi
// @Summary   首页统计数据
// @Security  ApiKeyAuth
// @Param dto.MchTradeStatisticRequest
// @Success 200 {object} response.Body{data=dto.TradeStatistic,msg=string} "查询成功"
// @Router /private/v1/dashboard/total_request [get]
func TotalRequest(c *gin.Context) {
	var req dto.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.GetAllRequestOrder(c.Request.Context(), req.MchNo, req.StartTime, req.EndTime)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, data)
}

// Index 获取商户交易统计
// @Tags      DashboardApi
// @Summary   首页统计数据
// @Security  ApiKeyAuth
// @Param dto.MchTradeStatisticRequest
// @Success 200 {object} response.Body{data=dto.TradeStatistic,msg=string} "查询成功"
// @Router /private/v1/dashboard/index [get]
func DashboardIndex(c *gin.Context) {
	var req dto.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		response.Fail(c, "无权限")
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.CountGroup(c.Request.Context(), req.StartTime, req.EndTime)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, data)
}

func MchIndex(c *gin.Context) {
	var req dto.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.CountGroupMch(c.Request.Context(), req.MchNo, req.StartTime, req.EndTime)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, data)
}

func MchAppIndex(c *gin.Context) {
	var req dto.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.CountGroupMchApp(c.Request.Context(), req.MchNo, req.AppNo, req.StartTime, req.EndTime)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, data)
}

func MchAppAccountAllIndex(c *gin.Context) {
	var req dto.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.CountGroupMchAppAccount(c.Request.Context(), req.MchNo, req.AppNo, req.StartTime, req.EndTime)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, data)
}

func SearchIndex(c *gin.Context) {
	var req dto.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.SearchDashboard(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, data)
}
