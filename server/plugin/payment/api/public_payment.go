package api

import (
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/utils/response"
)

func Pay(c *gin.Context) {
	var req public.PaymentOrder
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	resp, err := enter.ServiceApiApp.PaymentService.Payment(c.Request.Context(), &req, uuid.NewString())
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, resp)
}

func QueryPayment(c *gin.Context) {
	var req public.QueryPaymentRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	resp, err := enter.ServiceApiApp.PaymentService.Query(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, resp)
}

func ClosePayment(c *gin.Context) {
	var req public.QueryPaymentRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	err = enter.ServiceApiApp.PaymentService.Close(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, nil)
}

func PaymentCallback(c *gin.Context) {
	channel := c.Param("channel")
	mchNo := c.Param("mchNo")
	appNo := c.Param("appNo")
	orderNo := c.Param("orderNo")
	if mchNo == "" || appNo == "" || orderNo == "" || channel == "" {
		FailMessage(errors.ErrCodeInvalidParam, "参数错误", c)
		return
	}
	resp, err := enter.ServiceApiApp.PaymentService.PublishCallback(c.Request.Context(), c.Request, channel, mchNo, appNo, orderNo, false)
	if err != nil {
		FailErrMessage(err, c)
		return
	}

	c.String(http.StatusOK, resp)
	return
}

func PaymentTestCallback(c *gin.Context) {
	channel := c.Param("channel")
	mchNo := c.Param("mchNo")
	appNo := c.Param("appNo")
	orderNo := c.Param("orderNo")
	if mchNo == "" || appNo == "" || orderNo == "" || channel == "" {
		FailMessage(errors.ErrCodeInvalidParam, "参数错误", c)
		return
	}
	resp, err := enter.ServiceApiApp.PaymentService.PublishCallback(c.Request.Context(), c.Request, channel, mchNo, appNo, orderNo, true)
	if err != nil {
		FailErrMessage(err, c)
		return
	}

	c.String(http.StatusOK, resp)
	return
}
