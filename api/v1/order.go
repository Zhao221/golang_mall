package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/consts"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func CreateOrderHandler(c *gin.Context){
	var req types.OrderCreateReq
	if err := c.ShouldBindJSON(&req); err!= nil {
		global.GVA_LOG.Error("创建订单失败")
		response.FailWithMessage("创建订单失败",c)
		return
	}
	l := service.GetOrderSrv()
	resp, err := l.OrderCreate(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("创建订单失败")
		response.FailWithMessage("创建订单失败",c)
		return
	}
	response.OkWithDetailed(resp,"创建订单成功",c)
}

func ListOrdersHandler(c *gin.Context){
	var req types.OrderListReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取订单列表失败")
		response.FailWithMessage("获取订单列表失败",c)
		return
	}
	if req.PageSize == 0 {
		req.PageSize = consts.BasePageSize
	}

	l := service.GetOrderSrv()
	resp, err := l.OrderList(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取订单列表失败")
		response.FailWithMessage("获取订单失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取订单列表成功",c)
}

func ShowOrderHandler(c *gin.Context){
	var req types.OrderShowReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取订单详情失败")
		response.FailWithMessage("获取订单详情失败",c)
		return
	}

	l := service.GetOrderSrv()
	resp, err := l.OrderShow(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("获取订单详情失败")
		response.FailWithMessage("获取订单详情失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取订单详情成功",c)
}

func DeleteOrderHandler(c *gin.Context){
	var req types.OrderDeleteReq
	if err:=c.ShouldBind(&req);err!=nil{
		global.GVA_LOG.Error("删除订单失败")
		response.FailWithMessage("删除订单失败",c)
		return
	}
	l := service.GetOrderSrv()
	err := l.OrderDelete(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("删除订单失败")
		response.FailWithMessage("删除订单失败",c)
		return
	}
	response.OkWithData("删除订单成功",c)
}