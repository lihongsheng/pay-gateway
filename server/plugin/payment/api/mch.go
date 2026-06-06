package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	_ "github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

type MchApi struct {
}

// Get
// @Tags      MchAdmin
// @Summary   获取商户信息
// @Security  ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id query int true "用id查询商户信息"
// @Success 200 {object} response.Response{data=model.Merchant,msg=string} "查询成功"
// @Router /private/v1/merchant [get]
func (m *MchApi) Get(c *gin.Context) {
	mchNO := c.Query("mch_no")
	user := GetUserInfo(c)
	if user.UserType.ISPlatform() && mchNO == "" {
		response.FailWithMessage("商户号必传", c)
		return
	}
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if user.UserType.ISMch() {
		mchNO = user.HaveMchNo
	}
	info, err := enter.ServiceApiApp.MchService.Get(c.Request.Context(), mchNO)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(info, c)
}

// Search
// @Tags MchAdmin
// @Summary 分页商户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query admin.MchQueryRequest true "分页商户列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /private/v1/merchant/search [get]
func (m *MchApi) Search(c *gin.Context) {
	var req admin.MchQueryRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}
	user := GetUserInfo(c)
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	info, err := enter.ServiceApiApp.MchService.Search(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	count, err := enter.ServiceApiApp.MchService.Count(c.Request.Context(), req)
	response.OkWithData(response.PageResult{
		List:     info,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}

// Save
// @Tags MchAdmin
// @Summary 保存商户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query admin.MchCreateRequest true "保存商户信息"
// @Success 200 {object} response.Response{msg=string} "保存商户信息成功"
// @Router /private/v1/merchant [post]
func (m *MchApi) Save(c *gin.Context) {
	var req admin.MchCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	user := GetUserInfo(c)
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = enter.ServiceApiApp.MchService.Save(c.Request.Context(), req, user)
	if err != nil {
		global.GVA_LOG.Error("保存失败!", zap.Error(err))
		response.FailWithMessage("保存失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ChangeStatus
// @Tags MchAdmin
// @Summary 保存商户状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query admin.MchStatusRequest true "保存商户状态"
// @Success 200 {object} response.Response{msg=string} "保存商户状态成功"
// @Router /private/v1/merchant/status [post]
func (m *MchApi) ChangeStatus(c *gin.Context) {
	var req admin.MchStatusRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := user.ISHaveMchID(req.MchNo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = enter.ServiceApiApp.MchService.ChangeStatus(c.Request.Context(), req.MchNo, req.Status)
	if err != nil {
		global.GVA_LOG.Error("改变失败!", zap.Error(err))
		response.FailWithMessage("改变失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
