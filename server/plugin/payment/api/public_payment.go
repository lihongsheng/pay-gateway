package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"net/http"
)

type PublicPaymentApi struct {
}

func (a *PublicPaymentApi) Pay(c *gin.Context) {
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
	response.OkWithData(resp, c)
}

func (a *PublicPaymentApi) Query(c *gin.Context) {
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
	response.OkWithData(resp, c)
}

func (a *PublicPaymentApi) Close(c *gin.Context) {
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
	response.Ok(c)
}

func (a *PublicPaymentApi) Callback(c *gin.Context) {
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

func (a *PublicPaymentApi) TestCallback(c *gin.Context) {
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
