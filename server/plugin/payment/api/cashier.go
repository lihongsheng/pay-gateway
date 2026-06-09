package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
)

// CashierApi 收银台支付
type CashierApi struct {
}

func (a *CashierApi) RedirectUrl(c *gin.Context) {
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
	response.OkWithData(r, c)
}

func (a *CashierApi) Create(c *gin.Context) {
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
	response.OkWithData(r, c)
}

func (a *CashierApi) Payment(c *gin.Context) {
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
	response.OkWithData(r, c)
}

func (a *CashierApi) Query(c *gin.Context) {
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
	response.OkWithData(r, c)
}
