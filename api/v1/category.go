package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func ListCategoryHandler(c *gin.Context){
	var req types.ListCategoryReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取商品分类失败")
		response.FailWithMessage("获取商品分类失败",c)
		return
	}

	l := service.GetCategorySrv()
	resp, err := l.CategoryList(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取商品分类失败")
		response.FailWithMessage("获取商品分类失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取商品分类成功",c)

}