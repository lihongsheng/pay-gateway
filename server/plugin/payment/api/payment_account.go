package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	servicePay "github.com/lihongsheng/pay-gateway/plugin/payment/service"
	"github.com/lihongsheng/pay-gateway/utils/response"
)

// GetPaymentAccount
// @Tags      AccountApiAdmin
// @Summary   获取应用支付渠道配置信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param id query int true "用id查询商户应用信息"
// @Success 200 {object} response.Body{data=dto.PaymentAccountDetail,msg=string} "查询成功"
// @Router /private/v1/payment_account [get]
func GetPaymentAccount(c *gin.Context) {
	ID := c.Query("id")
	id, err := strconv.Atoi(ID)
	if err != nil {
		zap.L().Error("参数错误!", zap.Error(err))
		response.Fail(c, "参数错误")
		return
	}
	user := GetUserInfo(c)
	info, err := servicePay.DefaultPaymentAccount.Get(c.Request.Context(), int64(id))
	if err := user.ISHaveMchID(info.MchNo); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err != nil {
		zap.L().Error("查询失败!", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}
	response.OK(c, info)
}

// Save
// @Tags      AccountApiAdmin
// @Summary   保存应用支付渠道配置信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param data query dto.PaymentAccountCreateRequest true "保存商户信息"
// @Success 200 {object} response.Body{data=,msg=string} "查询成功"
// @Router /private/v1/payment_account [post]
func SavePaymentAccount(c *gin.Context) {
	var req dto.PaymentAccountCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	err = servicePay.DefaultPaymentAccount.Save(c.Request.Context(), &req)
	if err != nil {
		zap.L().Error("保存失败!", zap.Error(err), zap.Any("data", req))
		response.Fail(c, "保存失败"+err.Error())
		return
	}
	response.OKMsg(c, "保存成功")
}

// AccountList
// @Tags      AccountApiAdmin
// @Summary   获取应用下支付渠道列表
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param id query int true "用id查询商户应用信息"
// @Success 200 {object} response.Body{data=[]model.PaymentAccount,msg=string} "查询成功"
// @Router /private/v1/payment_account/list [get]
func AccountList(c *gin.Context) {
	appNO := c.Query("app_no")
	if appNO == "" {
		response.Fail(c, "参数错误")
		return
	}
	user := GetUserInfo(c)

	info, err := servicePay.DefaultPaymentAccount.GetAccountByAppNo(c.Request.Context(), appNO)
	if err != nil {
		zap.L().Error("查询失败!", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}
	for _, v := range info {
		if err := user.ISHaveMchID(v.MchNo); err != nil {
			response.Fail(c, err.Error())
			return
		}
	}
	response.OK(c, info)
}

// GetChannelConfig
// @Tags      AccountApiAdmin
// @Summary   获取应用可用的支付渠道及其配置信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param id query int true "用id查询商户应用信息"
// @Success 200 {object} response.Body{data=map[string]*params2.Option,msg=string} "查询成功"
// @Router /private/v1/payment_account/channel/config [get]
func GetChannelConfig(c *gin.Context) {
	appNO := c.Query("app_no")
	if appNO == "" {
		response.Fail(c, "参数错误")
		return
	}
	info, err := servicePay.DefaultPaymentAccount.GetApplicationChannelConfig(c.Request.Context(), appNO)
	if err != nil {
		zap.L().Error("查询失败!", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}
	response.OK(c, info)
}

// GetChannelPaymentMethod
// @Tags      AccountApiAdmin
// @Summary   获取应用可用的支付渠道及其配置信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param id query int true "用id查询商户应用信息"
// @Success 200 {object} response.Body{data=map[string]*params2.Option,msg=string} "查询成功"
// @Router /private/v1/payment_account/channel/payment [get]
func GetChannelPaymentMethod(c *gin.Context) {
	channel := c.Query("channel")
	if channel == "" {
		response.Fail(c, "参数错误")
		return
	}
	info, err := servicePay.DefaultPaymentAccount.GetPaymentProduct(channel)
	if err != nil {
		zap.L().Error("查询失败!", zap.Error(err))
		response.Fail(c, "查询失败")
		return
	}
	response.OK(c, info)
}

// GetTestQrCode
// @Tags      AccountApiAdmin
// @Summary   获取当前支付账号的测试二维码
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param app_no query string true "应用编号"
// @Param account_no query int true "支付编号"
// @Success 200 {object} response.Body{data=map[string]*params2.Option,msg=string} "查询成功"
// @Router /private/v1/payment_account/channel/payment [get]
func GetTestQrCode(c *gin.Context) {
	appNO := c.Query("app_no")
	accNO := c.Query("account_no")
	if appNO == "" || accNO == "" {
		response.Fail(c, "参数错误")
		return
	}
	u, order, err := servicePay.DefaultAggregate.GenTestQrCode(c.Request.Context(), appNO, accNO)
	if err != nil {
		zap.L().Error("GetQrCode", zap.Error(err))
		response.Fail(c, "获取二维码失败："+err.Error())
		return
	}
	response.OK(c, gin.H{"qr_link": u, "order_no": order})
}
