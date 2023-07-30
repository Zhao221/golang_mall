package mysql

import (
	"golang_mall/global"
	model2 "golang_mall/model"
)

func migrate() (err error) {
	err = global.GVA_DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model2.User{}, &model2.Favorite{},
			&model2.Order{}, &model2.Admin{}, &model2.Address{},
			&model2.Cart{}, &model2.Category{}, &model2.Carousel{},
			&model2.Notice{}, &model2.Notice{}, &model2.Product{},
			&model2.ProductImg{}, &model2.SkillProduct{},
			&model2.SkillProduct2MQ{}, &model2.FileUpload{},
			&model2.Topic{}, &model2.Message{},
		)

	return
}
