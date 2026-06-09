package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"go.uber.org/zap"
)

type TradeApi struct{}

// GetTradeOrderApi
// @Tags      TradeAdmin
// @Summary   获取支付订单详情信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param mch_no query int true "商户应用信息"
// @Param order_no query int true "订单号"
// @Success 200 {object} response.Response{data=admin.TradeOrderDetail,msg=string} "查询成功"
// @Router /private/v1/trade [get]
func (t *TradeApi) GetTradeOrderApi(c *gin.Context) {
	mchNo := c.Query("mch_no")
	orderNo := c.Query("order_no")
	appNo := c.Query("app_no")
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		mchNo = user.HaveMchNo
	}
	if mchNo == "" || orderNo == "" {
		response.FailWithMessage("mch_no or order_no is empty", c)
		return
	}

	detail, err := enter.ServiceApiApp.TradeOrderService.Detail(c.Request.Context(), mchNo, appNo, orderNo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(detail, c)
}

// SearchTradeOrderApi
// @Tags      TradeAdmin
// @Summary   搜索支付订单信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param admin.TradeSearchRequest
// @Success 200 {object} response.Response{data=admin.TradeOrderDetail,msg=string} "查询成功"
// @Router /private/v1/trade/search [get]
func (t *TradeApi) SearchTradeOrderApi(c *gin.Context) {
	var req admin.TradeSearchRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	detail, err := enter.ServiceApiApp.TradeOrderService.Search(c.Request.Context(), &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	count, err := enter.ServiceApiApp.TradeOrderService.Count(c.Request.Context(), &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     detail,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}

// TradeRefundApi
// @Tags      TradeAdmin
// @Summary   获取支付订单信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param mch_no query int true "商户应用信息"
// @Param order_no query int true "订单号"
// @Success 200 {object} response.Response{data=admin.TradeOrderDetail,msg=string} "查询成功"
// @Router /private/v1/refund [post]
func (t *TradeApi) TradeRefundApi(c *gin.Context) {
	l := log.WithCtx(c.Request.Context())
	var req admin.Refund
	err := c.ShouldBindJSON(&req)
	if err != nil {
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
	detail, err := enter.ServiceApiApp.TradeRefundService.Refund(c.Request.Context(), &req)
	if err != nil {
		l.Error("refundErr", zap.Error(err), zap.Any("req", req))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(detail, c)
}

// TradeRefundDetailApi
// @Tags      TradeAdmin
// @Summary   获取退款单详情信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param mch_no query int true "商户应用信息"
// @Param order_no query int true "订单号"
// @Success 200 {object} response.Response{data=admin.RefundOrderDetail,msg=string} "查询成功"
// @Router /private/v1/refund [get]
func (t *TradeApi) TradeRefundDetailApi(c *gin.Context) {
	mchNo := c.Query("mch_no")
	tradeNo := c.Query("trade_no")
	appNO := c.Query("app_no")
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		mchNo = user.HaveMchNo
	}
	if mchNo == "" || tradeNo == "" || appNO == "" {
		response.FailWithMessage("mch_no or trade_no is empty", c)
		return
	}
	detail, err := enter.ServiceApiApp.TradeRefundService.Detail(c.Request.Context(), mchNo, appNO, tradeNo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(detail, c)
}

// TradeRefundSearchApi
// @Tags      TradeAdmin
// @Summary   搜索退款单详情信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param admin.TradeSearchRequest
// @Success 200 {object} response.Response{data=admin.RefundOrderDetail,msg=string} "查询成功"
// @Router /private/v1/refund/search [get]
func (t *TradeApi) TradeRefundSearchApi(c *gin.Context) {
	var req admin.RefundSearchRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	user := GetUserInfo(c)
	if user.UserType.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	detail, err := enter.ServiceApiApp.TradeRefundService.Search(c.Request.Context(), &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	count, err := enter.ServiceApiApp.TradeRefundService.Count(c.Request.Context(), &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     detail,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}

func (t *TradeApi) GetAvailableRefundAmount(c *gin.Context) {
	mchNo := c.Query("mch_no")
	orderNo := c.Query("order_no")
	appNo := c.Query("app_no")
	if mchNo == "" || orderNo == "" || appNo == "" {
		response.FailWithMessage("mch_no or order_no is empty", c)
		return
	}
	amount, err := enter.ServiceApiApp.TradeRefundService.AvailableRefundAmount(c.Request.Context(), mchNo, appNo, orderNo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(amount, c)
}
