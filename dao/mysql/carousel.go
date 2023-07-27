package mysql

import (
	"context"
	"golang_mall/global"
	"golang_mall/model"
	"golang_mall/types"
	"gorm.io/gorm"
)

type CarouselDao struct {
	*gorm.DB
}

func NewCarouselDao(ctx context.Context) *CarouselDao {
	return &CarouselDao{NewDBClient(ctx)}
}

func NewNewCarouselDao() *CarouselDao {
	return &CarouselDao{DB: global.GVA_DB}
}

func (dao *CarouselDao) ListCarousel() (r []*types.ListCarouselResp, err error) {
	err = dao.DB.Model(&model.Carousel{}).
		Select("id, img_path, product_id, UNIX_TIMESTAMP(created_at)").
		Find(&r).Error
	return r,err
}
