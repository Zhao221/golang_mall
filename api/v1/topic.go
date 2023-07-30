package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/global"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

func GetTopicHandler(c *gin.Context) {
	var req types.BasePage
	if err := c.ShouldBindQuery(&req); err != nil {
		global.GVA_LOG.Error("获取主题列表传参失败")
		response.FailWithMessage("获取主题列表传参失败", c)
		return
	}
	l := service.GetTopicSrv()
	resp, total, err := l.GetTopicListHandle(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取主题列表失败")
		response.FailWithMessage("获取主题列表失败", c)
		return
	}
	topicList := make(map[string]interface{})
	topicList["total"] = total
	topicList["pageSize"] = req.PageSize
	topicList["pageNum"] = req.PageNum
	topicList["resp"] = resp
	response.OkWithDetailed(topicList, "获取主题列表成功", c)
}

func GetUserTopicHandler(c *gin.Context) {
	var req types.BasePage
	if err := c.ShouldBindQuery(&req); err != nil {
		global.GVA_LOG.Error("获取主题列表传参失败")
		response.FailWithMessage("获取主题列表传参失败", c)
		return
	}
	l := service.GetTopicSrv()
	resp, total, err := l.GetUserTopicListHandle(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取主题列表失败")
		response.FailWithMessage("获取主题列表失败", c)
		return
	}
	topicList := make(map[string]interface{})
	topicList["total"] = total
	topicList["pageSize"] = req.PageSize
	topicList["pageNum"] = req.PageNum
	topicList["resp"] = resp
	response.OkWithDetailed(topicList, "获取主题列表成功", c)
}

func SubscribeTopicHandler(c *gin.Context) {
	var req types.SubscribeTopicReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("订阅主题传参失败")
		response.FailWithMessage("订阅主题传参失败", c)
		return
	}
	l := service.GetTopicSrv()
	if err := l.SubscribeTopic(c.Request.Context(),req); err != nil {
		global.GVA_LOG.Error("订阅主题失败")
		response.FailWithMessage("订阅主题失败", c)
	}
	response.OkWithMessage("订阅主题成功", c)
}



func GetMessageHandler(c *gin.Context) {
	var req types.GetMessageList
	if err := c.ShouldBindQuery(&req); err != nil {
		global.GVA_LOG.Error("获取消息列表传参失败")
		response.FailWithMessage("获取消息列表传参失败", c)
		return
	}
	l := service.GetTopicSrv()
	resp, total, err := l.GetMessageListHandle(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取消息列表失败")
		response.FailWithMessage("获取消息列表失败", c)
		return
	}
	MessageList := make(map[string]interface{})
	MessageList["total"] = total
	MessageList["pageSize"] = req.PageSize
	MessageList["pageNum"] = req.PageNum
	MessageList["resp"] = resp
	response.OkWithDetailed(MessageList, "获取消息列表成功", c)
}

func DeleteMessageHandler(c *gin.Context) {
	var req types.DeleteMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("删除消息传参失败")
		response.FailWithMessage("删除消息传参失败", c)
		return
	}
	l := service.GetTopicSrv()
	 err := l.DeleteMessage(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("删除消息失败")
		response.FailWithMessage("删除消息失败", c)
		return
	}
	response.OkWithMessage("删除消息成功", c)
}
