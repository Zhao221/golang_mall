package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	logging "github.com/sirupsen/logrus"
	"golang_mall/conf"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client
var RedisContext = context.Background()

// InitRedis 在中间件中初始化redis链接
func InitRedis() {
	rConfig := conf.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rConfig.RedisHost, rConfig.RedisPort),
		Username: rConfig.RedisUsername,
		Password: rConfig.RedisPassword,
		DB:       rConfig.RedisDbName,
	})
	_, err := client.Ping(RedisContext).Result()
	if err != nil {
		logging.Info(err)
		fmt.Println("redis初始化配置错误：",err)
		return
	}
	RedisClient = client
}
