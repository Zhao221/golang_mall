package types

import "golang_mall/model"

type SkillProductImportReq struct {
}

type SkillProductReq struct {
	SkillProductId uint   `json:"skill_product_id"`
	ProductId      uint   `json:"product_id"`
	BossId         uint   `json:"boss_id"`
	AddressId      uint   `json:"address_id"`
	Num            uint   `json:"num"`
	Key            string `json:"key"`
}

type ListSkillProductReq struct {
	PageSize int `form:"page_size" json:"page_size"`
	PageNum  int `form:"page_num" json:"page_num"`
}

type GetSkillProductReq struct {
	SkillProductId uint `json:"skill_product_id" form:"skill_product_id"`
}

type SkillProductResp struct {
	SkillList []model.SkillProduct
	Total     int64 `json:"total"`
}

type SkillProduct struct {
	SkillStruct model.SkillProduct
	Num         uint   `json:"num"`
	Address     uint   `json:"address"`
	Key         string `json:"key"`
	LockKey     string `json:"lock_key"`
	LockValue   string `json:"lock_value"`
}
