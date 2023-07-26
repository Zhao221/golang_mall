package service

import (
	"context"
	"errors"
	"fmt"
	errors2 "github.com/pkg/errors"
	"golang_mall/conf"
	"golang_mall/consts"
	"golang_mall/dao/mysql"
	"golang_mall/global"
	"golang_mall/model"
	"golang_mall/model/common/request"
	"golang_mall/pkg/utils/ctl"
	Email "golang_mall/pkg/utils/email"
	"golang_mall/pkg/utils/jwt"
	"golang_mall/pkg/utils/upload"
	"golang_mall/types"
	"mime/multipart"
)

func CheckUserName(userRegister request.UserRegisterReq) (err error) {
	var sum int64
	global.GVA_DB.Table("user").Where("user_name", userRegister.UserName).Count(&sum)
	if sum != 0 {
		return errors2.New("用户名重复")
	}
	return err
}

// Register 用户注册
func Register(userRegister request.UserRegisterReq) (err error) {
	uR := &model.User{
		NickName: userRegister.NickName,
		UserName: userRegister.UserName,
		Status:   model.Active,
		Money:    consts.UserInitMoney,
		Avatar:   "rxkjiutfo.hn-bkt.clouddn.com/Fmk8lKuP1nCpAkzteoAnBBAQ29-a",
	}
	// 加密密码
	if err = uR.SetPassword(userRegister.Password); err != nil {
		return errors2.New("密码加密失败")
	}
	// 加密money
	money, err := uR.EncryptMoney(userRegister.Key)
	if err != nil {
		return errors2.New("金额加密失败")
	}
	uR.Money = money
	err = global.GVA_DB.Model(&model.User{}).Create(&uR).Error
	return err
}

// UserLogin 用户登录
func UserLogin(c context.Context, uLogin request.UserLoginReq) (resp types.UserTokenData, err error) {
	var user model.User
	userDao := mysql.NewUserDao(c)
	user, exist, err := userDao.ExistOrNotByUserName(uLogin.UserName)
	if !exist { // 如果查询不到，返回相应的错误
		return resp, errors.New("用户不存在")
	}
	// 密码解密
	user.CheckPassword(uLogin.Password)
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, uLogin.UserName)
	if err != nil {
		return resp, err
	}
	userResp := &types.UserInfoResp{
		UserName: user.UserName,
		NickName: user.NickName,
		Email:    user.Email,
		Status:   user.Status,
		Avatar:   user.AvatarURL(),
		CreateAt: user.CreatedAt.Unix(),
	}
	resp = types.UserTokenData{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return resp, err
}

// UserUpdateInfo 用户修改信息
func UserUpdateInfo(c context.Context, update request.UserUpdate) (nickName string, err error) {
	u, _ := ctl.GetUserInfo(c)
	if update.NickName != "" {
		err = global.GVA_DB.Table("user").Where("id", u.Id).Update("nick_name", update.NickName).Error
		if err != nil {
			return nickName, err
		}
	}
	nickName = update.NickName
	return nickName, err
}

func UserInfoShow(c context.Context) (resp interface{}, err error) {
	user, err := ctl.GetUserInfo(c)
	var UInfo model.User
	global.GVA_DB.Table("user").Where("id", user.Id).Find(&UInfo)
	resp = types.UserInfoResp{
		ID:       UInfo.ID,
		UserName: UInfo.UserName,
		NickName: UInfo.NickName,
		Email:    UInfo.Email,
		Status:   UInfo.Status,
		Avatar:   UInfo.AvatarURL(),
		CreateAt: UInfo.CreatedAt.Unix(),
	}
	return resp, err
}

func SendEmail(c context.Context, email types.SendEmailServiceReq) (resp interface{}, err error) {
	user, err := ctl.GetUserInfo(c)
	var address string
	token, err := jwt.GenerateEmailToken(user.Id, email.OperationType, email.Email, email.Password)
	if err != nil {
		return nil, err
	}
	global.GVA_DB.Table("user").Where("id", user.Id).Update("email", email.Email)
	sender := Email.NewEmailSender()
	address = conf.Config.Email.ValidEmail + token
	mailText := fmt.Sprintf(consts.EmailOperationMap[email.OperationType], address)
	if err = sender.Send(mailText, email.Email, "FanOneMall"); err != nil {
		return nil,err
	}
	return resp, err
}

// Valid 验证内容
func Valid(c context.Context, v types.ValidEmailServiceReq) (resp interface{}, err error) {
	var (
		userId        uint
		operationType uint
		email         string
		password      string
	)
	if v.Token == "" {
		return nil, errors2.New("token不存在")
	}
	claims, err := jwt.ParseEmailToken(v.Token)
	if err != nil {
		return nil, err
	}
	email = claims.Email
	userId = claims.UserID
	password = claims.Password
	operationType = claims.OperationType
	userDao := mysql.NewUserDao(c)
	user, err := userDao.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	switch operationType {
	case consts.EmailOperationBinding:
		user.Email = email
	case consts.EmailOperationNoBinding:
		user.Email = ""
	case consts.EmailOperationUpdatePassword:
		err = user.SetPassword(password)
		if err != nil {
			return nil, errors.New("密码加密错误")
		}
	default:
		return nil, errors.New("操作不符合")
	}
	err = userDao.UpdateUserById(userId, user)
	if err != nil {
		return nil, err
	}
	resp = &types.UserInfoResp{
		ID:       user.ID,
		UserName: user.UserName,
		NickName: user.NickName,
		Email:    user.Email,
		Status:   user.Status,
		Avatar:   user.AvatarURL(),
		CreateAt: user.CreatedAt.Unix(),
	}
	return resp, err
}

// UserFollow 关注用户
func UserFollow(c context.Context, req types.UserFollowingReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	err = mysql.NewUserDao(c).FollowUser(u.Id, req.Id)
	return resp, err
}

// UserUnFollow 取消关注
func UserUnFollow(c context.Context, req types.UserUnFollowingReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	err = mysql.NewUserDao(c).UnFollowUser(u.Id, req.Id)
	return resp, err
}

// UserAvatarUpload 上传头像
func UserAvatarUpload(c context.Context, file multipart.File, fileSize int64, name string) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	uId := u.Id
	userDao := mysql.NewUserDao(c)
	user, err := userDao.GetUserById(uId)
	if err != nil {
		return nil, err
	}

	var path string
	path, err = upload.UploadToQiNiu(file, fileSize, name)
	if err != nil {
		return nil, err
	}
	user.Avatar = path
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		return nil, err
	}
	return path, err
}
