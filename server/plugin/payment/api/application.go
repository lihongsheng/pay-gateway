package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	servicePay "github.com/lihongsheng/pay-gateway/plugin/payment/service"
	payUtils "github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"github.com/lihongsheng/pay-gateway/utils/response"
)

func GenAppSecret(c *gin.Context) {
	secret, err := payUtils.GenAesHexStr(32)
	if err != nil {
		response.Fail(c, "生成密钥失败")
		return
	}
	response.OK(c, secret)
}

// GetAppInfo
// @Tags      MchAppAdmin
// @Summary   获取商户应用信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param app_no query string true "用id查询商户应用信息"
// @Success 200 {object} response.Body{data=model.Application,msg=string} "查询成功"
// @Router /private/v1/application [get]
func GetAppInfo(c *gin.Context) {
	appNO := c.Query("app_no")
	if appNO == "" {
		response.Fail(c, "参数错误")
		return
	}
	user := GetUserInfo(c)
	info, err := servicePay.DefaultApplication.Get(c.Request.Context(), appNO)
	if err != nil {
		zap.L().Error("查询失败!", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}
	// 验证用户是否有权限访问此应用
	if err := user.ISHaveMchID(info.MchNo); err != nil {
		response.Fail(c, "无权限")
		return
	}

	response.OK(c, info)
}

// Search
// @Tags MchAppAdmin
// @Summary 分页商户应用列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dto.ApplicationQueryRequest true "分页商户应用列表"
// @Success 200 {object} response.Body{data=PageResult,msg=string} "获取成功"
// @Router /private/v1/application/search [get]
func SearchApplication(c *gin.Context) {
	var req dto.ApplicationQueryRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	info, err := servicePay.DefaultApplication.Search(c.Request.Context(), &req)
	if err != nil {
		zap.L().Error("查询失败!", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}
	count, err := servicePay.DefaultApplication.Count(c.Request.Context(), &req)
	response.OK(c, PageResult{
		List:     info,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// Save
// @Tags MchAppAdmin
// @Summary 保存商户应用信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dto.ApplicationCreateRequest true "保存商户应用信息"
// @Success 200 {object} response.Body{msg=string} "保存商户应用信息成功"
// @Router /private/v1/application [post]
func SaveApplication(c *gin.Context) {
	var req dto.ApplicationCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	_, err = servicePay.DefaultApplication.Save(c.Request.Context(), &req)
	if err != nil {
		zap.L().Error("保存失败!", zap.Error(err))
		response.Fail(c, "保存失败")
		return
	}
	response.OKMsg(c, "更新成功")
}

// ChangeStatus
// @Tags MchAppAdmin
// @Summary 保存商户应用状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dto.ApplicationStatusRequest true "保存商户应用状态"
// @Success 200 {object} response.Body{msg=string} "保存商户应用状态成功"
// @Router /private/v1/application/status [post]
func ChangeAppStatus(c *gin.Context) {
	var req dto.ApplicationStatusRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	err = servicePay.DefaultApplication.ChangeStatus(c.Request.Context(), req.ID, req.Status, user)
	if err != nil {
		zap.L().Error("改变失败!", zap.Error(err))
		response.Fail(c, "改变失败")
		return
	}
	response.OKMsg(c, "更新成功")
}

func GetQrCode(c *gin.Context) {
	appNO := c.Query("app_no")
	if appNO == "" {
		response.Fail(c, "参数错误")
		return
	}

	u, err := servicePay.DefaultAggregate.GenQrCode(c.Request.Context(), appNO)
	if err != nil {
		zap.L().Error("GetQrCode", zap.Error(err))
		response.Fail(c, "获取二维码失败："+err.Error())
		return
	}
	response.OK(c, u)
}
