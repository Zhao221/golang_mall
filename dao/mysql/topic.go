package mysql

import (
	"context"
	"golang_mall/global"
	"gorm.io/gorm"
)

type TopicDao struct {
	*gorm.DB
}

func NewTopicDao(ctx context.Context) *TopicDao {
	return &TopicDao{NewDBClient(ctx)}
}

func NewTopicDaoByDB() *TopicDao {
	return &TopicDao{DB: global.GVA_DB}
}