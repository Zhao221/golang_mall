package v1

import (
	"github.com/gin-gonic/gin"
	"golang_mall/consts"
	"golang_mall/global"
	"golang_mall/model/common/request"
	"golang_mall/model/common/response"
	"golang_mall/service"
	"golang_mall/types"
)

// UserRegisterHandler 用户注册
func UserRegisterHandler(c *gin.Context) {
	var userRegister request.UserRegisterReq
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		global.GVA_LOG.Error("用户注册，绑定参数出错了")
		response.FailWithMessage("用户注册，绑定参数出错了", c)
		return
	}
	if userRegister.Key == "" || len(userRegister.Key) != consts.EncryptMoneyKeyLength {
		global.GVA_LOG.Error("key长度错误,必须是6位数")
		response.FailWithMessage("key长度错误,必须是6位数", c)
		return
	}
	if err := service.CheckUserName(userRegister); err != nil {
		global.GVA_LOG.Error("用户名重复")
		response.FailWithMessage("用户名重复", c)
		return
	}
	if err := service.Register(userRegister); err != nil {
		global.GVA_LOG.Error("用户注册存入数据库出错了")
		response.FailWithMessage("用户注册存入数据库出错了", c)
		return
	}
	response.OkWithDetailed(userRegister, "用户注册成功", c)
}

// UserLoginHandler 用户登录
func UserLoginHandler(c *gin.Context) {
	var login request.UserLoginReq
	if err := c.ShouldBindJSON(&login); err != nil {
		global.GVA_LOG.Error("用户登录，绑定参数出错了")
		response.FailWithMessage("用户登录，绑定参数出错了", c)
		return
	}
	resp, err := service.UserLogin(c.Request.Context(), login)
	if err != nil {
		global.GVA_LOG.Error("用户登录出错了")
		response.FailWithMessage("用户登录出错了", c)
		return
	}
	response.OkWithDetailed(resp, "用户登录成功", c)
}

// UserUpdateHandler 更新用户信息
func UserUpdateHandler(c *gin.Context) {
	var update request.UserUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		global.GVA_LOG.Error("用户更新，绑定参数出错了")
		response.FailWithMessage("用户更新，绑定参数出错了", c)
		return
	}
	nickName, err := service.UserUpdateInfo(c.Request.Context(), update)
	if err != nil {
		global.GVA_LOG.Error("用户更新出错了")
		response.FailWithMessage("用户更新出错了", c)
		return
	}
	response.OkWithDetailed(nickName, "用户更新信息成功", c)
}

// ShowUserInfoHandler 展示用户信息
func ShowUserInfoHandler(c *gin.Context) {
	resp, err := service.UserInfoShow(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("获取用户信息，出错了")
		response.FailWithMessage("获取用户信息，出错了", c)
		return
	}
	response.OkWithDetailed(resp, "用户更新信息成功", c)

}

// SendEmailHandler 发送邮箱
func SendEmailHandler(c *gin.Context) {
	var email types.SendEmailServiceReq
	if err := c.ShouldBindJSON(&email); err != nil {
		global.GVA_LOG.Error("绑定邮箱，绑定参数出错了")
		response.FailWithMessage("绑定邮箱，绑定参数出错了", c)
		return
	}
	resp, err := service.SendEmail(c.Request.Context(), email)
	if err != nil {
		global.GVA_LOG.Error("绑定邮箱，出错了")
		response.FailWithMessage("绑定邮箱，出错了", c)
		return
	}
	response.OkWithDetailed(resp, "绑定邮箱成功", c)
}

// ValidEmailHandler 验证邮箱
func ValidEmailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ValidEmailServiceReq
		if err := c.ShouldBindJSON(&req); err != nil {
			global.GVA_LOG.Error("验证邮箱绑定参数时，出错了")
			response.FailWithMessage("验证邮箱绑定参数时，出错了", c)
			return
		}
		resp, err := service.Valid(c.Request.Context(), req)
		if err != nil {
			global.GVA_LOG.Error("验证邮箱，出错了")
			response.FailWithMessage("验证邮箱，出错了", c)
			return
		}
		response.OkWithDetailed(resp, "验证邮箱成功", c)
	}
}

// UserFollowingHandler 关注用户
func UserFollowingHandler(c *gin.Context) {
	var req types.UserFollowingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("关注用户绑定参数时，出错了")
		response.FailWithMessage("关注用户绑定参数时，出错了", c)
		return
	}
	resp, err := service.UserFollow(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("关注用户失败")
		response.FailWithMessage("关注用户失败", c)
		return
	}
	response.OkWithDetailed(resp, "关注用户成功", c)
}

// UserUnFollowingHandler 取关
func UserUnFollowingHandler(c *gin.Context) {
	var req types.UserUnFollowingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("取关用户绑定参数时，出错了")
		response.FailWithMessage("取关用户绑定参数时，出错了", c)
		return
	}
	resp, err := service.UserUnFollow(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("取关用户，出错了")
		response.FailWithMessage("取关用户，出错了", c)
		return
	}
	response.OkWithDetailed(resp, "取关用户成功", c)
}

// UploadAvatarHandler 上传头像
func UploadAvatarHandler(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")
	fileSize := fileHeader.Size
	resp, err := service.UserAvatarUpload(c.Request.Context(), file, fileSize, fileHeader.Filename)
	if err != nil {
		global.GVA_LOG.Error("上传头像时，出错了")
		response.FailWithMessage("上传头像时，出错了", c)
		return
	}
	response.OkWithDetailed(resp, "上传头像成功", c)
}

func UserCheckinHandler(c *gin.Context) {
	var req types.UserCheckin
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("用户签到传参有误")
		response.FailWithMessage("用户签到传参有误", c)
		return
	}
	if err = service.UserCheckinService(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("用户签到失败")
		response.FailWithMessage("用户签到失败", c)
		return
	}
	response.OkWithMessage("用户签到成功",c)
}
