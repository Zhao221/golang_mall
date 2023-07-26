package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/consts"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func CreateCartHandler(c *gin.Context){
	var req types.CartCreateReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("添加购物车商品失败")
		response.FailWithMessage("添加购物车商品失败",c)
		return
	}

	l := service.GetCartSrv()
	err := l.CartCreate(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("添加购物车商品失败")
		response.FailWithMessage("添加购物车商品失败",c)
		return
	}
	response.OkWithMessage("添加购物车商品成功",c)
}

func ListCartHandler(c *gin.Context){
	var req types.CartListReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取购物车列表失败")
		response.FailWithMessage("获取购物车列表失败",c)
		return
	}
	if req.PageSize == 0 {
		req.PageSize = consts.BasePageSize
	}

	l := service.GetCartSrv()
	resp, err := l.CartList(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取购物车列表失败")
		response.FailWithMessage("获取购物车列表失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取购物车列表成功",c)
}

func UpdateCartHandler(c *gin.Context){
	var req types.UpdateCartServiceReq
	if err:=c.ShouldBindJSON(&req);err!=nil{
		global.GVA_LOG.Error("更新购物车列表失败")
		response.FailWithMessage("更新购物车列表失败",c)
	}
	l:=service.GetCartSrv()
	err:=l.CartUpdate(c.Request.Context(),req)
	if err!=nil{
		global.GVA_LOG.Error("更新购物车列表失败")
		response.FailWithMessage("更新购物车列表失败",c)
	}
	response.OkWithMessage("更新购物车列表成功",c)
}

func DeleteCartHandler(c *gin.Context){
	var req types.CartDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("删除购物车中商品失败")
		response.FailWithMessage("删除购物车中商品失败",c)
		return
	}

	l := service.GetCartSrv()
	err := l.CartDelete(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("删除购物车中商品失败")
		response.FailWithMessage("删除购物车中商品失败",c)
		return
	}
	response.OkWithMessage("删除购物车中商品成功",c)
}
