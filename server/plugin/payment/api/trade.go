package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	servicePay "github.com/lihongsheng/pay-gateway/plugin/payment/service"
	"github.com/lihongsheng/pay-gateway/utils/response"
)

// GetTradeOrder
// @Tags      TradeAdmin
// @Summary   获取支付订单详情信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param mch_no query int true "商户应用信息"
// @Param order_no query int true "订单号"
// @Success 200 {object} response.Body{data=dto.TradeOrderDetail,msg=string} "查询成功"
// @Router /private/v1/trade [get]
func GetTradeOrder(c *gin.Context) {
	mchNo := c.Query("mch_no")
	orderNo := c.Query("order_no")
	appNo := c.Query("app_no")
	user := GetUserInfo(c)
	if user.ISMch() {
		mchNo = user.HaveMchNo
	}
	if mchNo == "" || orderNo == "" {
		response.Fail(c, "mch_no or order_no is empty")
		return
	}

	detail, err := servicePay.DefaultTradeOrder.Detail(c.Request.Context(), mchNo, appNo, orderNo)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, detail)
}

// SearchTradeOrder
// @Tags      TradeAdmin
// @Summary   搜索支付订单信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param dto.TradeSearchRequest
// @Success 200 {object} response.Body{data=dto.TradeOrderDetail,msg=string} "查询成功"
// @Router /private/v1/trade/search [get]
func SearchTradeOrder(c *gin.Context) {
	var req dto.TradeSearchRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Fail(c, err.Error())
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	detail, err := servicePay.DefaultTradeOrder.Search(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	count, err := servicePay.DefaultTradeOrder.Count(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, PageResult{
		List:     detail,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// TradeRefund
// @Tags      TradeAdmin
// @Summary   获取支付订单信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param mch_no query int true "商户应用信息"
// @Param order_no query int true "订单号"
// @Success 200 {object} response.Body{data=dto.TradeOrderDetail,msg=string} "查询成功"
// @Router /private/v1/refund [post]
func TradeRefund(c *gin.Context) {
	l := log.WithCtx(c.Request.Context())
	var req dto.Refund
	err := c.ShouldBindJSON(&req)
	if err != nil {
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
	detail, err := servicePay.DefaultTradeRefund.Refund(c.Request.Context(), &req)
	if err != nil {
		l.Error("refundErr", zap.Error(err), zap.Any("req", req))
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, detail)
}

// TradeRefundDetail
// @Tags      TradeAdmin
// @Summary   获取退款单详情信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param mch_no query int true "商户应用信息"
// @Param order_no query int true "订单号"
// @Success 200 {object} response.Body{data=dto.RefundOrderDetail,msg=string} "查询成功"
// @Router /private/v1/refund [get]
func TradeRefundDetail(c *gin.Context) {
	mchNo := c.Query("mch_no")
	tradeNo := c.Query("trade_no")
	appNO := c.Query("app_no")
	user := GetUserInfo(c)
	if user.ISMch() {
		mchNo = user.HaveMchNo
	}
	if mchNo == "" || tradeNo == "" || appNO == "" {
		response.Fail(c, "mch_no or trade_no is empty")
		return
	}
	detail, err := servicePay.DefaultTradeRefund.Detail(c.Request.Context(), mchNo, appNO, tradeNo)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, detail)
}

// TradeRefundSearch
// @Tags      TradeAdmin
// @Summary   搜索退款单详情信息
// @Security  ApiKeyAuth
// @accept  application/json
// @Produce application/json
// @Param dto.RefundSearchRequest
// @Success 200 {object} response.Body{data=dto.RefundOrderDetail,msg=string} "查询成功"
// @Router /private/v1/refund/search [get]
func TradeRefundSearch(c *gin.Context) {
	var req dto.RefundSearchRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Fail(c, err.Error())
	}
	user := GetUserInfo(c)
	if user.ISMch() {
		req.MchNo = user.HaveMchNo
	}
	if err := user.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	detail, err := servicePay.DefaultTradeRefund.Search(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	count, err := servicePay.DefaultTradeRefund.Count(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, PageResult{
		List:     detail,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

func GetAvailableRefundAmount(c *gin.Context) {
	mchNo := c.Query("mch_no")
	orderNo := c.Query("order_no")
	appNo := c.Query("app_no")
	if mchNo == "" || orderNo == "" || appNo == "" {
		response.Fail(c, "mch_no or order_no is empty")
		return
	}
	amount, err := servicePay.DefaultTradeRefund.AvailableRefundAmount(c.Request.Context(), mchNo, appNo, orderNo)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, amount)
}
