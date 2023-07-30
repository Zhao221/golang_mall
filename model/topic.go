package model

import "gorm.io/gorm"

type Topic struct {
	gorm.Model
	Name   string `json:"name" gorm:"comment:主题名字"`
	UserId uint   `json:"user_id" gorm:"comment:用户Id"`
}
