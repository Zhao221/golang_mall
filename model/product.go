package model

import (
	"strconv"

	"github.com/jinzhu/gorm"

	cache "golang_mall/dao/redis"
)

// Product 商品模型
type Product struct {
	gorm.Model
	Name          string `gorm:"size:255;index"`
	Num           int
	CategoryID    uint `gorm:"not null"`
	BossID        uint
	MessageID     uint
	Title         string
	Info          string `gorm:"size:1000"`
	ImgPath       string
	Price         string
	DiscountPrice string
	BossName      string
	BossAvatar    string
	OnSale        bool    `gorm:"default:false"`
	Mg            Message `gorm:"foreignKey:ProductId"`
}

// View 获取点击数
func (product *Product) View() uint64 {
	countStr, _ := cache.RedisClient.Get(cache.RedisContext, cache.ProductViewKey(product.ID)).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView 商品游览
func (product *Product) AddView() {
	// 增加视频点击数
	cache.RedisClient.Incr(cache.RedisContext, cache.ProductViewKey(product.ID))
	// 增加排行点击数
	cache.RedisClient.ZIncrBy(cache.RedisContext, cache.RankKey, 1, strconv.Itoa(int(product.ID)))
}
