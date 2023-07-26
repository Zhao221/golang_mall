package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func ShowMoneyHandler(c *gin.Context){
	var req types.MoneyShowReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("展示金额失败")
		response.FailWithMessage("展示金额失败",c)
		return
	}
	l := service.GetMoneySrv()
	resp, err := l.MoneyShow(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("展示金额失败")
		response.FailWithMessage("展示金额失败",c)
		return
	}
	response.OkWithDetailed(resp,"展示金额成功",c)
}
