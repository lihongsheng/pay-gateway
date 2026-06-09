package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
)

type AggregateApi struct {
}

func (a *AggregateApi) RedirectUrl(c *gin.Context) {
	var req public.AggregateRedirectUrlRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	r, err := enter.ServiceApiApp.Aggregate.GetRedirectUrl(c.Request.Context(), &req, enum.PaymentBaseIndexPath)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OkWithData(r, c)
}

func (a *AggregateApi) TestRedirectUrl(c *gin.Context) {
	var req public.AggregateRedirectUrlRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	if req.OrderNo == "" {
		FailMessage(errors.ErrCodeInvalidParam, "order_no is empty", c)
		return
	}
	if req.AccountNo == "" {
		FailMessage(errors.ErrCodeInvalidParam, "account_no is empty", c)
		return
	}
	if req.Token == "" {
		FailMessage(errors.ErrCodeInvalidParam, "token is empty", c)
		return
	}
	r, err := enter.ServiceApiApp.Aggregate.GetRedirectUrl(c.Request.Context(), &req, enum.PaymentTestPath)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OkWithData(r, c)
}

func (a *AggregateApi) TestPayment(c *gin.Context) {
	var req public.AggregateOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	req.Amount = 1
	r, err := enter.ServiceApiApp.Aggregate.Payment(c.Request.Context(), &req, enum.PaymentTestNotify)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OkWithData(r, c)
}

func (a *AggregateApi) GetUserOpenID(c *gin.Context) {
	var req public.AggregateUserOpenIDRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	r, err := enter.ServiceApiApp.Aggregate.GetUserOpenID(c.Request.Context(), req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OkWithData(r, c)
}

func (a *AggregateApi) GetApplication(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		FailMessage(errors.ErrCodeInvalidParam, "token is empty", c)
		return
	}
	r, err := enter.ServiceApiApp.Aggregate.GetApplication(c.Request.Context(), token)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	multiChannel := false
	if r.MultiChannel == enum.MultiChannel_Enable {
		multiChannel = true
	}
	response.OkWithData(public.AggregateApplication{
		AppNo:        r.AppNo,
		PayTitle:     r.PaymentTitle,
		PayIcon:      r.PayIcon,
		MultiChannel: multiChannel,
	}, c)
}

func (a *AggregateApi) Payment(c *gin.Context) {
	var req public.AggregateOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	r, err := enter.ServiceApiApp.Aggregate.Payment(c.Request.Context(), &req, enum.PaymentNotify)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OkWithData(r, c)
}

func (a *AggregateApi) Query(c *gin.Context) {
	token := c.Query("token")
	orderNo := c.Query("order_no")
	if token == "" {
		FailMessage(errors.ErrCodeInvalidParam, "token is empty", c)
		return
	}
	if orderNo == "" {
		FailMessage(errors.ErrCodeInvalidParam, "order_no is empty", c)
		return
	}
	r, err := enter.ServiceApiApp.Aggregate.Query(c.Request.Context(), token, orderNo)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OkWithData(r, c)
}
