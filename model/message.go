package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	TopicId         uint   `json:"topic_id" gorm:"comment:主题id"`
	UserId          uint   `json:"user_id" gorm:"comment:用户id"`
	ProductId       uint   `json:"product_id" gorm:"comment:商品Id"`
	CategoryId      uint   `json:"category_id" gorm:"comment:商品分类Id"`
	Info            string `json:"info" gorm:"size:1000"`
	Title           string `json:"title" gorm:"comment:商品标题"`
	ImgPath         string `json:"img_path" gorm:"comment:商品图片"`
	TopicName       string `json:"topic_name" gorm:"comment:主题名字"`
	ProductName     string `json:"product_name" gorm:"comment:商品名字"`
	ProductCategory string `json:"product_category" gorm:"comment:商品分类"`
}
