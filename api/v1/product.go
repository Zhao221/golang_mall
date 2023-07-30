package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/consts"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func CreateProductHandler(c *gin.Context){
	var req types.ProductCreateReq
	if err:=c.ShouldBind(&req);err!=nil{
		global.GVA_LOG.Error("创建商品绑定参数时，失败")
		response.FailWithMessage("创建商品绑定参数时，失败",c)
		return
	}
	form, _ := c.MultipartForm()
	files := form.File["file"]
	l := service.GetProductSrv()
	err := l.ProductCreate(c.Request.Context(),files, req)
	if err!=nil{
		global.GVA_LOG.Error("创建商品失败")
		response.FailWithMessage("创建商品失败",c)
		return
	}
	response.OkWithMessage("创建商品成功",c)
}

func UpdateProductHandler(c *gin.Context){
	var req types.ProductUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("更新商品绑定参数时，失败")
		response.FailWithMessage("更新商品绑定参数时，失败",c)
		return
	}
	l := service.GetProductSrv()
	err := l.ProductUpdate(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("更新商品信息失败")
		response.FailWithMessage("更新商品信息失败",c)
		return
	}
	response.OkWithMessage("更新商品信息成功",c)
}

func DeleteProductHandler(c *gin.Context){
	var req types.ProductDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("删除商品绑定参数时，失败")
		response.FailWithMessage("删除商品绑定参数时，失败",c)
		return
	}

	l := service.GetProductSrv()
	err := l.ProductDelete(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("删除商品失败")
		response.FailWithMessage("删除商品失败",c)
		return
	}
	response.OkWithMessage("删除商品信息成功",c)
}

// ListProductsHandler 商品列表
func ListProductsHandler(c *gin.Context){
	var req types.ProductListReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取商品绑定参数时，失败")
		response.FailWithMessage("获取商品绑定参数时，失败",c)
		return
	}
	if req.PageSize == 0 {
		req.PageSize = consts.BaseProductPageSize
	}
	l:=service.GetProductSrv()
	resp, err := l.ProductList(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取商品列表失败")
		response.FailWithMessage("获取商品列表失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取商品列表成功",c)
}

func ShowProductHandler(c *gin.Context){
	var req types.ProductShowReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取商品详情失败")
		response.FailWithMessage("获取商品详情失败",c)
		return
	}

	l := service.GetProductSrv()
	resp, err := l.ProductShow(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取商品详情失败")
		response.FailWithMessage("获取商品详情失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取商品列表成功",c)
}

func SearchProductsHandler(c *gin.Context){
	var req types.ProductSearchReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("搜索商品失败")
		response.FailWithMessage("搜索商品失败",c)
		return
	}
	if req.PageSize == 0 {
		req.PageSize = consts.BasePageSize
	}

	l := service.GetProductSrv()
	resp, err := l.ProductSearch(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("搜索商品失败")
		response.FailWithMessage("搜索商品失败",c)
		return
	}
	response.OkWithDetailed(resp,"搜索商品成功",c)
}

func ListProductImgHandler(c *gin.Context){
	var req types.ListProductImgReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取商品图片失败")
		response.FailWithMessage("获取商品图片失败",c)
		return
	}

	l := service.GetProductSrv()
	resp, err := l.ProductImgList(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取商品图片失败")
		response.FailWithMessage("获取商品图片失败",c)
		return
	}
	response.OkWithDetailed(resp,"获取商品图片成功",c)
}

