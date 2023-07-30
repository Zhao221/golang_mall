package mysql

import (
	"context"
	"golang_mall/global"
	"gorm.io/gorm"
)

type MessageDao struct {
	*gorm.DB
}

func NewMessageDao(ctx context.Context) *MessageDao {
	return &MessageDao{NewDBClient(ctx)}
}

func NewMessageDaoByDB() *MessageDao {
	return &MessageDao{DB: global.GVA_DB}
}
