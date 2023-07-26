package global

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	GVA_DB    *gorm.DB
	GVA_REDIS *redis.Client
	GVA_LOG   *zap.Logger
)
