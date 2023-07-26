package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	errors2 "github.com/pkg/errors"
	"golang_mall/dao/mysql"
	"golang_mall/dao/redis"
	"golang_mall/model"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/types"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type KillSrv struct{}

var killSrvIns *KillSrv
var killSrvOnce sync.Once

func GetSkillProductSrv() *KillSrv {
	killSrvOnce.Do(func() {
		killSrvIns = &KillSrv{}
	})
	return killSrvIns
}

// 每个商品过期时间只设置一次
type skillProductHash struct{}

var skillProduct *skillProductHash
var skillExpire sync.Once

func GetSkillExpire(c context.Context, key string) *skillProductHash {
	skillExpire.Do(func() {
		redis.RedisClient.Expire(c, key, time.Hour*6)
	})
	return skillProduct
}

func (k *KillSrv) InitSkillGoods(c context.Context) (resp interface{}, err error) {
	spList := make([]*model.SkillProduct, 0)
	for i := 1; i <= 3; i++ {
		spList = append(spList, &model.SkillProduct{
			ProductId: uint(i),
			BossId:    1,
			Title:     "秒杀商品测试使用",
			Money:     100,
			Num:       2,
		})
	}
	err = mysql.NewSkillGoodsDao(c).BatchCreate(spList)
	if err != nil {
		return nil, err
	}
	// 导入数据库的同时，初始化缓存
	for i := range spList {
		jsonBytes, errx := json.Marshal(spList[i])
		if errx != nil {
			return nil, errx
		}
		jsonString := string(jsonBytes)
		key := strconv.FormatUint(uint64(spList[i].ID), 10)
		_, errx = redis.RedisClient.HSet(c, key, redis.SkillProductHash, jsonString).Result()
		if errx != nil {
			return nil, errx
		}
	}
	resp = spList
	return resp, err
}

// ListSkillGoods 获取秒杀商品列表
func (k *KillSrv) ListSkillGoods(c context.Context, req types.ListSkillProductReq) (resp interface{}, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.PageNum - 1)
	var skillProductList []model.SkillProduct
	sum := int64(0)
	mysql.NewCartDao(c).Model(&model.SkillProduct{}).Where("num > ?", 0).
		Count(&sum).Limit(limit).Offset(offset).Find(&skillProductList)
	skill := types.SkillProductResp{
		SkillList: skillProductList,
		Total:     sum,
	}
	resp = skill
	return resp, err
}

// GetSkillGoods 获取秒杀商品商品详细信息
func (k *KillSrv) GetSkillGoods(c context.Context, req types.GetSkillProductReq) (resp interface{}, err error) {
	key := strconv.FormatUint(uint64(req.SkillProductId), 10)
	skillInfo, err := redis.RedisClient.HGet(c, key, redis.SkillProductHash).Result()
	var skillStruct model.SkillProduct
	skill := []byte(skillInfo)
	err = json.Unmarshal(skill, &skillStruct)
	resp = skillStruct
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (k *KillSrv) CheckNums(c context.Context, req types.SkillProductReq) (resp types.SkillProduct, err error) {
	key := strconv.FormatUint(uint64(req.SkillProductId), 10)
	skillInfo, err := redis.RedisClient.HGet(c, key, redis.SkillProductHash).Result()
	skill := []byte(skillInfo)
	var skillStruct model.SkillProduct
	err = json.Unmarshal(skill, &skillStruct)
	if err != nil {
		return resp, err
	}
	if skillStruct.Num < req.Num {
		return resp, errors2.New("您买商品的数量超限了，请调整数量")
	}
	// SetNX 加锁
	lockKey := "product_lock"
	// 为了确保在释放锁时不会错误地释放其他客户端的锁
	// 可以在获取锁时为每个客户端分配一个唯一的值（如 UUID）。
	// 当释放锁时，只有锁的值与客户端的唯一值匹配时，才会释放锁。
	// 这可以避免因为锁过期而被其他客户端获取，
	// 导致原客户端错误地释放了新客户端的锁的问题。
	lockValue := uuid.New().String()
	lockSuccess, err := redis.RedisClient.SetNX(c, lockKey, lockValue, 5*time.Second).Result()
	if err != nil || !lockSuccess {
		return resp, errors2.New("Failed to acquire the lock. Please try again later.")
	}
	if skillStruct.Num <= 0 {
		redis.RedisClient.Del(c, lockKey, lockValue)
		go func() {
			GetSkillExpire(c, key)
		}()
		return resp, errors2.New("此商品已经卖完了，可以逛逛其他好货哦")
	}
	resp = types.SkillProduct{
		SkillStruct: skillStruct,
		Key:         req.Key,
		Num:         req.Num,
		Address:     req.AddressId,
		LockKey:     lockKey,
		LockValue:   lockValue,
	}
	return resp, err
}

func (k *KillSrv) SkillProduct(c context.Context, req types.SkillProduct) (resp interface{}, err error) {
	// 获取此用户id
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	err = mysql.NewSkillGoodsDao(c).Transaction(func(tx *gorm.DB) error {
		// 更新秒杀商品数量
		req.SkillStruct.Num -= req.Num
		skillProductDao := mysql.NewSkillGoodsDao(c)
		updateSkill, err := skillProductDao.UpdateSkillProduct(req.SkillStruct.ID, req.SkillStruct.Num)
		if err != nil { // 更新商品数量减少失败，回滚
			return err
		}
		jsonBytes, errX := json.Marshal(updateSkill)
		if errX != nil {
			return errX
		}
		jsonString := string(jsonBytes)
		key := strconv.FormatUint(uint64(updateSkill.ID), 10)
		_, errX = redis.RedisClient.HSet(c, key, redis.SkillProductHash, jsonString).Result()
		if errX != nil {
			return errX // 缓存更新失败，回滚
		}

		// 创建自己的购物订单信息
		number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
		productNum := strconv.Itoa(int(req.SkillStruct.ProductId))
		userNum := strconv.Itoa(int(u.Id))
		number = number + productNum + userNum
		orderNum, _ := strconv.ParseUint(number, 10, 64)
		order := model.Order{
			UserID:    u.Id,
			ProductID: req.SkillStruct.ProductId,
			BossID:    req.SkillStruct.BossId,
			AddressID: req.Address,
			Num:       int(req.Num),
			OrderNum:  orderNum,
			Type:      2,
			Money:     req.SkillStruct.Money,
		}
		err = mysql.NewOrderDao(c).CreateOrder(&order)
		redis.RedisClient.Del(c, req.LockKey,req.LockValue)
		return err
	})
	if err != nil {
		return nil, err
	}
	return resp, err
}
