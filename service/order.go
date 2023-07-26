package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang_mall/conf"
	"golang_mall/consts"
	"golang_mall/dao/mysql"
	cache "golang_mall/dao/redis"
	"golang_mall/model"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/types"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const OrderTimeKey = "OrderTime"

var OrderSrvIns *OrderSrv
var OrderSrvOnce sync.Once

type OrderSrv struct {
}

func GetOrderSrv() *OrderSrv {
	OrderSrvOnce.Do(func() {
		OrderSrvIns = &OrderSrv{}
	})
	return OrderSrvIns
}

func (o *OrderSrv) OrderCreate(c context.Context, req types.OrderCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	order := &model.Order{
		UserID:    u.Id,
		ProductID: req.ProductID,
		BossID:    req.BossID,
		AddressID: req.AddressID,
		Num:       int(req.Num),
		Money:     float64(req.Money),
		Type:      req.Type, // 1,未支付；2已支付
	}
	number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	productNum := strconv.Itoa(int(req.ProductID))
	userNum := strconv.Itoa(int(u.Id))
	number = number + productNum + userNum
	orderNum, _ := strconv.ParseUint(number, 10, 64)
	order.OrderNum = orderNum
	orderDao := mysql.NewOrderDao(c)
	err = orderDao.CreateOrder(order)
	if err != nil {
		return nil ,err
	}
	// 订单号存入Redis中，设置过期时间
	data := redis.Z{
		Score:  float64(time.Now().Unix()) + 15*time.Minute.Seconds(),
		Member: orderNum,
	}
	cache.RedisClient.ZAdd(cache.RedisContext, OrderTimeKey, data)
	resp = orderNum
	return resp, err
}

func (o *OrderSrv) OrderList(c context.Context, req types.OrderListReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	orders, total, err := mysql.NewOrderDao(c).ListOrderByCondition(u.Id, &req)
	if err != nil {
		return nil ,err
	}
	for i := range orders {
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			orders[i].ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + orders[i].ImgPath
		}
	}
	resp = types.DataListResp{
		Item:  orders,
		Total: total,
	}
	return
}

func (o *OrderSrv) OrderShow(c context.Context, req *types.OrderShowReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil ,err
	}
	order, err := mysql.NewOrderDao(c).ShowOrderById(req.OrderId, u.Id)
	if err != nil {
		return nil,err
	}
	resp = order
	return resp, err
}

func (o *OrderSrv) OrderDelete(c context.Context, req types.OrderDeleteReq)( err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return err
	}
	err = mysql.NewOrderDao(c).DeleteOrderById(req.OrderId, u.Id)
	if err != nil {
		return err
	}
	return err
}
