package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/consts"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func CreateAddressHandler(c *gin.Context) {
	var req types.AddressCreateReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("新建收货地址失败")
		response.FailWithMessage("新建收货地址失败",c)
		return
	}

	l := service.GetAddressSrv()
	 err := l.AddressCreate(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("新建收货地址失败")
		response.FailWithMessage("新建收货地址失败",c)
		return
	}
	response.OkWithMessage("新建收货地址成功",c)
}

func ShowAddressHandler(c *gin.Context){
	var req types.AddressGetReq
	if err := c.ShouldBindQuery(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("展示收货地址失败")
		response.FailWithMessage("展示收货地址失败",c)
		return
	}

	l := service.GetAddressSrv()
	resp, err := l.AddressShow(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("展示收货地址失败")
		response.FailWithMessage("展示收货地址失败",c)
		return
	}
	response.OkWithDetailed(resp,"展示收货地址成功",c)
}

func ListAddressHandler(c *gin.Context){
	var req types.AddressListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取收货地址列表失败")
		response.FailWithMessage("获取收货地址列表失败",c)
		return
	}
	if req.PageSize == 0 {
		req.PageSize = consts.BasePageSize
	}
	l := service.GetAddressSrv()
	resp, err := l.AddressList(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取收货地址列表失败")
		response.FailWithMessage("获取收货地址列表失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取收货地址列表成功",c)
}

func UpdateAddressHandler(c *gin.Context){
	var req types.AddressServiceReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("更新收货地址失败")
		response.FailWithMessage("更新收货地址失败",c)
		return
	}

	l := service.GetAddressSrv()
	resp, err := l.AddressUpdate(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("更新收货地址失败")
		response.FailWithMessage("更新收货地址失败",c)
		return
	}
	response.OkWithDetailed(resp,"更新收货地址成功",c)
}

func DeleteAddressHandler(c *gin.Context){
	var req types.AddressDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("删除收货地址失败")
		response.FailWithMessage("删除收货地址失败",c)
		return
	}

	l := service.GetAddressSrv()
	err := l.AddressDelete(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("删除收货地址失败")
		response.FailWithMessage("删除收货地址失败",c)
		return
	}
	response.OkWithMessage("删除收货地址成功",c)
}