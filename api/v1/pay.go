package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func OrderPaymentHandler(c *gin.Context) {
	var req types.PaymentDownReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("支付失败")
		response.FailWithMessage("支付失败", c)
		return
	}
	fmt.Println(req)
	l := service.GetPaymentSrv()
	err := l.PayDown(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("支付失败")
		response.FailWithMessage("支付失败", c)
		return
	}
	response.OkWithMessage("支付成功", c)
}
