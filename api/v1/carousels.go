package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
)

func ListCarouselsHandler(c *gin.Context){
		l := service.GetCarouselSrv()
		resp, err := l.ListCarousel(c.Request.Context())
		if err != nil {
			global.GVA_LOG.Error("获取收货地址列表失败")
			response.FailWithMessage("获取轮播图失败",c)
			return
		}

	response.OkWithDetailed(resp,"获取轮播图成功",c)
}
