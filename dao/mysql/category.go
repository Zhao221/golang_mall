package mysql

import (
	"context"
	"golang_mall/global"

	"gorm.io/gorm"

	"golang_mall/model"
)

type CategoryDao struct {
	*gorm.DB
}

func NewCategoryDao(ctx context.Context) *CategoryDao {
	return &CategoryDao{NewDBClient(ctx)}
}

func NewCategoryDaoByDB() *CategoryDao {
	return &CategoryDao{DB: global.GVA_DB}
}

// ListCategory 分类列表
func (dao *CategoryDao) ListCategory() (r []*model.Category, err error) {
	err = dao.DB.Model(&model.Category{}).Find(&r).Error
	return
}
