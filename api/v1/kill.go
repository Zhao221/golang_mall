package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)


// InitSkillProductHandler 初始化秒杀商品信息
func InitSkillProductHandler(c *gin.Context) {
	var req types.ListSkillProductReq
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("初始化秒杀商品信息失败")
		response.FailWithMessage("初始化秒杀商品信息失败", c)
		return
	}

	l := service.GetSkillProductSrv()
	resp, err := l.InitSkillGoods(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("初始化秒杀商品信息失败")
		response.FailWithMessage("初始化秒杀商品信息失败", c)
		return
	}
	response.OkWithDetailed(resp, "初始化秒杀商品信息成功", c)
}

// ListSkillProductHandler 获取秒杀商品列表
func ListSkillProductHandler(c *gin.Context) {
	var req types.ListSkillProductReq
	if err := c.ShouldBindQuery(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取秒杀商品列表失败")
		response.FailWithMessage("获取秒杀商品列表失败", c)
		return
	}

	l := service.GetSkillProductSrv()
	resp, err := l.ListSkillGoods(c.Request.Context(),req)
	if err != nil {
		global.GVA_LOG.Error("获取秒杀商品列表失败")
		response.FailWithMessage("获取秒杀商品列表失败", c)
		return
	}
	response.OkWithDetailed(resp, "获取秒杀商品列表成功", c)
}

// GetSkillProductHandler 获取秒杀商品的详情
func GetSkillProductHandler(c *gin.Context) {
	var req types.GetSkillProductReq
	if err := c.ShouldBindQuery(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("获取秒杀商品的详情失败")
		response.FailWithMessage("获取秒杀商品的详情失败", c)
		return
	}

	l := service.GetSkillProductSrv()
	resp, err := l.GetSkillGoods(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取秒杀商品的详情失败")
		response.FailWithMessage("获取秒杀商品的详情失败", c)
		return
	}
	response.OkWithDetailed(resp, "获取秒杀商品的详情成功", c)
}

func SkillProductHandler(c *gin.Context) {
	var req types.SkillProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数校验
		global.GVA_LOG.Error("秒杀商品下单失败")
		response.FailWithMessage("秒杀商品下单失败", c)
		return
	}
	// 对数量进行校验
	skillProduct,err :=service.GetSkillProductSrv().CheckNums(c.Request.Context(),req)
	if err!=nil{
		if skillProduct.Num>skillProduct.SkillStruct.Num{
			global.GVA_LOG.Error("购买数量超限")
			response.FailWithMessage("您买商品的数量超限了，请调整数量", c)
		}
		global.GVA_LOG.Error("此商品已卖完，可以逛逛其他商品哦")
		response.FailWithMessage("此商品已卖完，可以逛逛其他商品哦", c)
		return
	}
	l := service.GetSkillProductSrv()
	resp, err := l.SkillProduct(c.Request.Context(), skillProduct)
	if err != nil {
		global.GVA_LOG.Error("秒杀商品下单失败")
		response.FailWithMessage("秒杀商品下单失败", c)
		return
	}
	response.OkWithDetailed(resp, "秒杀商品下单成功", c)
}
