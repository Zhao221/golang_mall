package model

import "gorm.io/gorm"

type SkillProduct struct {
	gorm.Model
	ProductId  uint `gorm:"not null"`
	BossId     uint `gorm:"not null"`
	Num        uint `gorm:"not null"`
	CustomId   uint
	Money      float64
	Title      string
	CustomName string
}

type SkillProduct2MQ struct {
	SkillProductId uint    `json:"skill_good_id"`
	ProductId      uint    `json:"product_id"`
	BossId         uint    `json:"boss_id"`
	UserId         uint    `json:"user_id"`
	AddressId      uint    `json:"address_id"`
	Money          float64 `json:"money"`
	Key            string  `json:"key"`
}
