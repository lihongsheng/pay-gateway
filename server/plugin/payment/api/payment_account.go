package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
	"strconv"
)

type AccountApi struct {
}

// Get
// @Tags      AccountApiAdmin
// @Summary   获取应用支付渠道配置信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param id query int true "用id查询商户应用信息"
// @Success 200 {object} response.Response{data=admin.PaymentAccountDetail,msg=string} "查询成功"
// @Router /private/v1/payment_account [get]
func (a *AccountApi) Get(c *gin.Context) {
	ID := c.Query("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	user := GetUserInfo(c)
	info, err := enter.ServiceApiApp.PaymentAccount.Get(c.Request.Context(), int64(id))
	if err := user.ISHaveMchID(info.MchNo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(info, c)
}

// Save
// @Tags      AccountApiAdmin
// @Summary   保存应用支付渠道配置信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param data query admin.PaymentAccountCreateRequest true "保存商户信息"
// @Success 200 {object} response.Response{data=,msg=string} "查询成功"
// @Router /private/v1/payment_account [post]
func (a *AccountApi) Save(c *gin.Context) {
	var req admin.PaymentAccountCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	err = enter.ServiceApiApp.PaymentAccount.Save(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("保存失败!", zap.Error(err), zap.Any("data", req))
		response.FailWithMessage("保存失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}

// AccountList
// @Tags      AccountApiAdmin
// @Summary   获取应用下支付渠道列表
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param id query int true "用id查询商户应用信息"
// @Success 200 {object} response.Response{data=[]model.PaymentAccount,msg=string} "查询成功"
// @Router /private/v1/payment_account/list [get]
func (a *AccountApi) AccountList(c *gin.Context) {
	appNO := c.Query("app_no")
	if appNO == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	user := GetUserInfo(c)

	info, err := enter.ServiceApiApp.PaymentAccount.GetAccountByAppNo(c.Request.Context(), appNO)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	for _, v := range info {
		if err := user.ISHaveMchID(v.MchNo); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	response.OkWithData(info, c)
}

// GetChannelConfig
// @Tags      AccountApiAdmin
// @Summary   获取应用可用的支付渠道及其配置信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param id query int true "用id查询商户应用信息"
// @Success 200 {object} response.Response{data=map[string]*params2.Option,msg=string} "查询成功"
// @Router /private/v1/payment_account/channel/config [get]
func (a *AccountApi) GetChannelConfig(c *gin.Context) {
	appNO := c.Query("app_no")
	if appNO == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	info, err := enter.ServiceApiApp.PaymentAccount.GetApplicationChannelConfig(c.Request.Context(), appNO)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(info, c)
}

// GetChannelPaymentMethod
// @Tags      AccountApiAdmin
// @Summary   获取应用可用的支付渠道及其配置信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param id query int true "用id查询商户应用信息"
// @Success 200 {object} response.Response{data=map[string]*params2.Option,msg=string} "查询成功"
// @Router /private/v1/payment_account/channel/payment [get]
func (a *AccountApi) GetChannelPaymentMethod(c *gin.Context) {
	channel := c.Query("channel")
	if channel == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	info, err := enter.ServiceApiApp.PaymentAccount.GetPaymentProduct(channel)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(info, c)
}

// GetTestQrCode
// @Tags      AccountApiAdmin
// @Summary   获取当前支付账号的测试二维码
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param app_no query string true "应用编号"
// @Param account_no query int true "支付编号"
// @Success 200 {object} response.Response{data=map[string]*params2.Option,msg=string} "查询成功"
// @Router /private/v1/payment_account/channel/payment [get]
func (a *AccountApi) GetTestQrCode(c *gin.Context) {
	appNO := c.Query("app_no")
	accNO := c.Query("account_no")
	if appNO == "" || accNO == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	u, order, err := enter.ServiceApiApp.Aggregate.GenTestQrCode(c.Request.Context(), appNO, accNO)
	if err != nil {
		global.GVA_LOG.Error("GetQrCode", zap.Error(err))
		response.FailWithMessage("获取二维码失败："+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"qr_link": u, "order_no": order}, c)
}
