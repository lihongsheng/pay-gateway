package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"go.uber.org/zap"
)

type AppApi struct {
}

func (a *AppApi) GenAppSecret(c *gin.Context) {
	secret, err := utils.GenAesHexStr(32)
	if err != nil {
		response.FailWithMessage("生成密钥失败", c)
		return
	}
	response.OkWithData(secret, c)
}

// GetAppInfo
// @Tags      MchAppAdmin
// @Summary   获取商户应用信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param app_no query string true "用id查询商户应用信息"
// @Success 200 {object} response.Response{data=model.Application,msg=string} "查询成功"
// @Router /private/v1/application [get]
func (a *AppApi) GetAppInfo(c *gin.Context) {
	appNO := c.Query("app_no")
	if appNO == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	user := GetUserInfo(c)
	info, err := enter.ServiceApiApp.AppService.Get(c.Request.Context(), appNO)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	// 验证用户是否有权限访问此应用
	if err := user.ISHaveMchID(info.MchNo); err != nil {
		response.FailWithMessage("无权限", c)
		return
	}

	response.OkWithData(info, c)
}

// Search
// @Tags MchAppAdmin
// @Summary 分页商户应用列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query admin.MchQueryRequest true "分页商户应用列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /private/v1/application/search [get]
func (a *AppApi) Search(c *gin.Context) {
	var req admin.ApplicationQueryRequest
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
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	info, err := enter.ServiceApiApp.AppService.Search(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	count, err := enter.ServiceApiApp.AppService.Count(c.Request.Context(), &req)
	response.OkWithData(response.PageResult{
		List:     info,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}

// Save
// @Tags MchAppAdmin
// @Summary 保存商户应用信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query admin.MchCreateRequest true "保存商户应用信息"
// @Success 200 {object} response.Response{msg=string} "保存商户应用信息成功"
// @Router /private/v1/application [post]
func (a *AppApi) Save(c *gin.Context) {
	var req admin.ApplicationCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	_, err = enter.ServiceApiApp.AppService.Save(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("保存失败!", zap.Error(err))
		response.FailWithMessage("保存失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ChangeStatus
// @Tags MchAppAdmin
// @Summary 保存商户应用状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query admin.MchStatusRequest true "保存商户应用状态"
// @Success 200 {object} response.Response{msg=string} "保存商户应用状态成功"
// @Router /private/v1/application/status [post]
func (a *AppApi) ChangeStatus(c *gin.Context) {
	var req admin.ApplicationStatusRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	err = enter.ServiceApiApp.AppService.ChangeStatus(c.Request.Context(), req.ID, req.Status, user)
	if err != nil {
		global.GVA_LOG.Error("改变失败!", zap.Error(err))
		response.FailWithMessage("改变失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *AppApi) GetQrCode(c *gin.Context) {
	appNO := c.Query("app_no")
	if appNO == "" {
		response.FailWithMessage("参数错误", c)
		return
	}

	u, err := enter.ServiceApiApp.Aggregate.GenQrCode(c.Request.Context(), appNO)
	if err != nil {
		global.GVA_LOG.Error("GetQrCode", zap.Error(err))
		response.FailWithMessage("获取二维码失败："+err.Error(), c)
		return
	}
	response.OkWithData(u, c)
}
