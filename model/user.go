package model

import (
	"github.com/CocaineCong/secret"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"golang_mall/conf"
	"golang_mall/consts"
)

// User 用户模型
type User struct {
	gorm.Model
	UserName       string `json:"user_name" gorm:"unique"`
	Email          string `json:"email"`
	PasswordDigest string `json:"password_digest"`
	NickName       string `json:"nick_name"`
	Status         string `json:"status"`
	Avatar         string `gorm:"size:1000"`
	Money          string `json:"money"`
	Relations      []User `gorm:"many2many:relation;"`
}

const (
	PassWordCost        = 12       // 密码加密难度
	Active       string = "active" // 激活用户
)

// SetPassword 设置密码
//  Go 语言中使用 bcrypt 算法对密码进行哈希加密的函数。bcrypt 算法是一种密码哈希算法
//  它是由 Niels Provos 和 David Mazières 基于 Blowfish 加密算法设计的。
//  它被设计成计算时间较长，以抵抗穷举攻击和彩虹表攻击
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	return err == nil
}

// EncryptMoney 加密金额
func (u *User) EncryptMoney(key string) (money string, err error) {
	aesObj, err := secret.NewAesEncrypt(conf.Config.EncryptSecret.MoneySecret, key, "", secret.AesEncrypt128, secret.AesModeTypeCBC)
	if err != nil {
		return
	}
	money = aesObj.SecretEncrypt(u.Money)
	return money , err
}

// DecryptMoney 解密金额
func (u *User) DecryptMoney(key string) (money float64, err error) {
	aesObj, err := secret.NewAesEncrypt(conf.Config.EncryptSecret.MoneySecret, key, "", secret.AesEncrypt128, secret.AesModeTypeCBC)
	if err != nil {
		return
	}
	money = cast.ToFloat64(aesObj.SecretDecrypt(u.Money))
	return money ,err
}

// AvatarURL 头像地址
func (u *User) AvatarURL() string {
	if conf.Config.System.UploadModel == consts.UploadModelOss {
		return u.Avatar
	}
	return u.Avatar
}

