package service

import (
	"context"
	"errors"
	errors2 "github.com/pkg/errors"
	"golang_mall/dao/mysql"
	"golang_mall/pkg/e"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/types"
	"sync"
)

var CartSrvIns *CartSrv
var  CartSrvOnce sync.Once

type CartSrv struct {}

func GetCartSrv() *CartSrv{
	CartSrvOnce.Do(func() {
		CartSrvIns=&CartSrv{}
	})
	return CartSrvIns
}

func(C *CartSrv)CartCreate(c context.Context, req types.CartCreateReq) (err error){
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return  err
	}
	product, _ := mysql.NewProductDao(c).GetProductById(req.ProductId)
	if product == nil {
		return  errors2.New("没有此商品")
	}
	// 创建购物车
	cartDao := mysql.NewCartDao(c)
	_, status, _ := cartDao.CreateCart(req.ProductId, u.Id, req.BossID)
	if status == e.ErrorProductMoreCart {
		return errors.New(e.GetMsg(status))
	}
	return err
}

func(C *CartSrv)CartList(c context.Context, req types.CartListReq) (resp interface{}, err error){
	u , err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	carts, err := mysql.NewCartDao(c).ListCartByUserId(u.Id)
	if err != nil {
		return nil,err
	}
	resp = &types.DataListResp{
		Item:  carts, // TODO 无分页，之后考虑要不要加
		Total: int64(len(carts)),
	}
	return resp,err
}

func(C *CartSrv)CartUpdate(c context.Context, req types.UpdateCartServiceReq) (err error){
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return  err
	}
	err=mysql.NewCartDao(c).UpdateCartNumById(req.Id,u.Id,req.Num)
	if err != nil {
		return  err
	}
	return err
}

func(C *CartSrv)CartDelete(c context.Context, req types.CartDeleteReq) (err error){
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return  err
	}
	err = mysql.NewCartDao(c).DeleteCartById(req.Id, u.Id)
	if err != nil {
		return err
	}
	return err
}