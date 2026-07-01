package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/utils/response"
)

func CashierRedirectUrl(c *gin.Context) {
	var req public.AggregateRedirectUrlRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	r, err := enter.ServiceApiApp.Aggregate.GetRedirectUrl(c.Request.Context(), &req, enum.CashierPayment)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, r)
}

func CashierCreate(c *gin.Context) {
	var req public.CashierCreateOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	r, err := enter.ServiceApiApp.CashierService.Create(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, r)
}

func CashierPayment(c *gin.Context) {
	var req public.CashierPaymentOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	r, err := enter.ServiceApiApp.CashierService.Payment(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, r)
}

func CashierQuery(c *gin.Context) {
	var req public.CashierQueryOrder
	if err := c.ShouldBindQuery(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	r, err := enter.ServiceApiApp.CashierService.Query(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, r)
}
