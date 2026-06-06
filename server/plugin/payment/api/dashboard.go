package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
)

type DashboardApi struct {
}

// TotalRequest 获取商户总订单数
// @Tags      DashboardApi
// @Summary   首页统计数据
// @Security  ApiKeyAuth
// @Param admin.MchTradeStatisticRequest
// @Success 200 {object} response.Response{data=admin.TradeStatistic,msg=string} "查询成功"
// @Router /private/v1/dashboard/total_request [get]
func (d *DashboardApi) TotalRequest(c *gin.Context) {
	var req admin.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.GetAllRequestOrder(c.Request.Context(), req.MchNo, req.StartTime, req.EndTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// Index 获取商户交易统计
// @Tags      DashboardApi
// @Summary   首页统计数据
// @Security  ApiKeyAuth
// @Param admin.MchTradeStatisticRequest
// @Success 200 {object} response.Response{data=admin.TradeStatistic,msg=string} "查询成功"
// @Router /private/v1/dashboard/index [get]
func (d *DashboardApi) Index(c *gin.Context) {
	var req admin.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		response.FailWithMessage("无权限", c)
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.CountGroup(c.Request.Context(), req.StartTime, req.EndTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

func (d *DashboardApi) MchIndex(c *gin.Context) {
	var req admin.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.CountGroupMch(c.Request.Context(), req.MchNo, req.StartTime, req.EndTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}
func (d *DashboardApi) MchAppIndex(c *gin.Context) {
	var req admin.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.CountGroupMchApp(c.Request.Context(), req.MchNo, req.AppNo, req.StartTime, req.EndTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

func (d *DashboardApi) MchAppAccountAllIndex(c *gin.Context) {
	var req admin.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.CountGroupMchAppAccount(c.Request.Context(), req.MchNo, req.AppNo, req.StartTime, req.EndTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

func (d *DashboardApi) SearchIndex(c *gin.Context) {
	var req admin.MchTradeStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := enter.ServiceApiApp.TradeStatisticsService.SearchDashboard(c.Request.Context(), &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}
