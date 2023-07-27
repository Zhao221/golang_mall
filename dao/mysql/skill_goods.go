package mysql

import (
	"context"
	"golang_mall/consts"
	"golang_mall/global"
	"golang_mall/model"
	"gorm.io/gorm"
)

type SkillGoodsDao struct {
	*gorm.DB
}

func NewSkillGoodsDao(ctx context.Context) *SkillGoodsDao {
	return &SkillGoodsDao{NewDBClient(ctx)}
}

func NewSkillGoodsDaoByDB() *SkillGoodsDao {
	return &SkillGoodsDao{DB: global.GVA_DB}
}

func (dao *SkillGoodsDao) Create(in *model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).Create(&in).Error
}

func (dao *SkillGoodsDao) BatchCreate(in []*model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).CreateInBatches(&in, consts.ProductBatchCreate).Error
}

func (dao *SkillGoodsDao) CreateByList(in []*model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).Create(&in).Error
}

func (dao *SkillGoodsDao) ListSkillGoods() (resp []*model.SkillProduct, err error) {
	err = dao.Model(&model.SkillProduct{}).
		Where("num > 0").Find(&resp).Error

	return
}

func (dao *SkillGoodsDao) GetSkillProductById(SkillProductId uint) (resp model.SkillProduct, err error) {
	err = dao.Model(&model.SkillProduct{}).Where("id=?", SkillProductId).First(&resp).Error
	return resp, err
}

func (dao *SkillGoodsDao) UpdateSkillProduct(SkillProductId, num uint) (resp model.SkillProduct, err error) {
	err = dao.Model(&model.SkillProduct{}).Where("id = ?", SkillProductId).Update("num", num).Error
	err = dao.Model(&model.SkillProduct{}).Where("id = ?", SkillProductId).First(&resp).Error
	if err!=nil{
		return resp,err
	}
	return resp,err
}
