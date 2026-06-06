package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lihongsheng/pay-gateway/model/common/response"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/enter"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"go.uber.org/zap"
	"net/http"
)

type PublicRefundApi struct {
}

func (a *PublicRefundApi) Refund(c *gin.Context) {
	req := public.RefundRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	req.RefundFrom = enum.RefundFrom_API
	resp, err := enter.ServiceApiApp.RefundService.Refund(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OkWithData(resp, c)
}

func (a *PublicRefundApi) Query(c *gin.Context) {
	req := public.RefundQueryRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	resp, err := enter.ServiceApiApp.RefundService.Query(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OkWithData(resp, c)
}

func FailMessage(code int, message string, c *gin.Context) {
	response.Result(code, nil, message, c)
}

func FailErrMessage(err error, c *gin.Context) {
	if paymentErr, ok := err.(*errors.Error); ok {
		c.JSON(http.StatusOK, response.Response{
			Code: paymentErr.Code,
			Msg:  paymentErr.Message,
		})
	} else {
		c.JSON(http.StatusOK, response.Response{
			Code: paymentErr.Code,
			Msg:  paymentErr.Message,
		})
	}
}

func (a *PublicRefundApi) Callback(c *gin.Context) {
	channel := c.Param("channel")
	mchNo := c.Param("mchNo")
	appNo := c.Param("appNo")
	tradeNo := c.Param("tradeNo")
	l := log.WithCtx(c.Request.Context())
	body, _ := utils.GetRequestBody(c.Request)
	l.Info("PublicRefundApiCallback", zap.String("mchNo", mchNo), zap.String("appNo", appNo), zap.String("tradeNo", tradeNo), zap.String("body", string(body)))
	if mchNo == "" || appNo == "" || tradeNo == "" || channel == "" {
		FailMessage(errors.ErrCodeInvalidParam, "参数错误", c)
		return
	}

	resp, err := enter.ServiceApiApp.RefundService.PublishCallback(c.Request.Context(), c.Request, channel, mchNo, appNo, tradeNo)
	if err != nil {
		l.Error("PublicRefundApiCallback", zap.Error(err), zap.String("mchNo", mchNo), zap.String("appNo", appNo), zap.String("tradeNo", tradeNo))
		FailErrMessage(err, c)
		return
	}
	c.String(http.StatusOK, resp)
	return
}
