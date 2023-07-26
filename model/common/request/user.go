package request

type UserRegisterReq struct {
	NickName string `json:"nick_name" gorm:"comment:昵称"`
	UserName string `json:"user_name" gorm:"comment:用户名称"`
	Password string `json:"password" gorm:"comment:密码"`
	Key      string `json:"key" gorm:"comment:key"` // 前端进行判断
}

type UserLoginReq struct {
	UserName string `json:"user_name" gorm:"comment:用户名称"`
	Password string `json:"password" gorm:"comment:密码"`
}

type UserUpdate struct {
	NickName string `json:"nick_name" gorm:"comment:昵称"`
}

type UserInfo struct {

}
