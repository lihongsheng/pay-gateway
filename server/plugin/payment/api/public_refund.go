package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lihongsheng/pay-gateway/plugin/payment/dto/public"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/log"
	servicePay "github.com/lihongsheng/pay-gateway/plugin/payment/service"
	"github.com/lihongsheng/pay-gateway/plugin/payment/utils"
	"github.com/lihongsheng/pay-gateway/utils/response"
)

func Refund(c *gin.Context) {
	req := public.RefundRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	req.RefundFrom = enum.RefundFrom_API
	resp, err := servicePay.DefaultRefund.Refund(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, resp)
}

func QueryRefund(c *gin.Context) {
	req := public.RefundQueryRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		FailMessage(errors.ErrCodeInvalidParam, err.Error(), c)
		return
	}
	resp, err := servicePay.DefaultRefund.Query(c.Request.Context(), &req)
	if err != nil {
		FailErrMessage(err, c)
		return
	}
	response.OK(c, resp)
}

// FailMessage 返回错误响应（带错误码）
func FailMessage(code int, message string, c *gin.Context) {
	response.FailCode(c, code, message)
}

// FailErrMessage 从 payment error 返回错误响应
func FailErrMessage(err error, c *gin.Context) {
	if paymentErr, ok := err.(*errors.Error); ok {
		response.FailCode(c, paymentErr.Code, paymentErr.Message)
	} else {
		response.Fail(c, err.Error())
	}
}

func RefundCallback(c *gin.Context) {
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

	resp, err := servicePay.DefaultRefund.PublishCallback(c.Request.Context(), c.Request, channel, mchNo, appNo, tradeNo)
	if err != nil {
		l.Error("PublicRefundApiCallback", zap.Error(err), zap.String("mchNo", mchNo), zap.String("appNo", appNo), zap.String("tradeNo", tradeNo))
		FailErrMessage(err, c)
		return
	}
	c.String(http.StatusOK, resp)
	return
}
