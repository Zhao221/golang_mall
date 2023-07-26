package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/consts"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func CreateFavoriteHandler(c *gin.Context){
	var req types.FavoriteCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("创建收藏夹失败")
		response.FailWithMessage("创建收藏夹失败",c)
		return
	}
	l := service.GetFavoriteSrv()
	resp, err := l.FavoriteCreate(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("创建收藏夹失败")
		response.FailWithMessage("创建收藏夹失败",c)
		return
	}
	response.OkWithDetailed(resp,"创建收藏夹成功",c)
}

func ListFavoritesHandler(c *gin.Context){
	var req types.FavoritesServiceReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取收藏夹列表失败")
		response.FailWithMessage("获取收藏夹列表失败",c)
		return
	}
	if req.PageSize==0{
		req.PageSize=consts.BasePageSize
	}
	l := service.GetFavoriteSrv()
	resp, err := l.FavoriteList(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取收藏夹列表失败")
		response.FailWithMessage("获取收藏夹列表失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取收藏夹列表成功",c)
}


func DeleteFavoriteHandler(c *gin.Context){
	var req types.FavoriteDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("删除收藏夹列表失败")
		response.FailWithMessage("删除收藏夹列表失败",c)
		return
	}
	l := service.GetFavoriteSrv()
	resp, err := l.FavoriteDelete(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("删除收藏夹列表失败")
		response.FailWithMessage("删除收藏夹列表失败",c)
		return
	}
	response.OkWithDetailed(resp,"删除收藏夹列表成功",c)
}