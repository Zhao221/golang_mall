package service

import (
	"context"
	"errors"
	"fmt"
	"golang_mall/consts"
	"golang_mall/dao/mysql"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/types"
	"gorm.io/gorm"
	"sync"
)

type PaySrv struct{}

var PaySrvIns *PaySrv
var PaySrvOnce sync.Once

func GetPaymentSrv() *PaySrv {
	PaySrvOnce.Do(func() {
		PaySrvIns = &PaySrv{}
	})
	return PaySrvIns
}

func (p *PaySrv) PayDown(c context.Context, req types.PaymentDownReq) (err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return err
	}
	err = mysql.NewOrderDao(c).Transaction(func(tx *gorm.DB) error {
		uId := u.Id
		payment, err := mysql.NewOrderDaoByDB().GetOrderById(req.OrderId, uId)
		if err != nil {
			return err
		}
		money := payment.Money
		num := payment.Num
		money = money * float64(num)
		userDao := mysql.NewUserDaoByDB()
		user, err := userDao.GetUserById(uId)
		if err != nil {
			return err
		}
		// 对钱进行解密。减去订单。再进行加密。
		moneyFloat, err := user.DecryptMoney(req.Key)
		if err != nil {
			return err
		}
		if moneyFloat-money < 0.0 { // 金额不足进行回滚
			return errors.New("金币不足")
		}

		finMoney := fmt.Sprintf("%f", moneyFloat-money)
		user.Money = finMoney
		user.Money, err = user.EncryptMoney(req.Key)
		if err != nil {
			return err
		}

		err = userDao.UpdateUserById(uId, user)
		if err != nil { // 更新用户金额失败，回滚
			return err
		}

		boss, err := userDao.GetUserById(uint(req.BossID))
		if err != nil {
			return err
		}

		moneyFloat, _ = boss.DecryptMoney(req.Key)
		finMoney = fmt.Sprintf("%f", moneyFloat+money)
		boss.Money = finMoney
		boss.Money, err = boss.EncryptMoney(req.Key)
		if err != nil {
			return err
		}

		err = userDao.UpdateUserById(uint(req.BossID), boss)
		if err != nil { // 更新boss金额失败，回滚
			return err
		}

		productDao := mysql.NewProductDaoByDB()
		product, err := productDao.GetProductById(uint(req.ProductID))
		if err != nil {
			return err
		}
		product.Num -= num
		err = productDao.UpdateProduct(uint(req.ProductID), product)
		if err != nil { // 更新商品数量减少失败，回滚
			return err
		}

		// 更新订单状态
		payment.Type = consts.OrderTypePendingShipping
		err = mysql.NewOrderDaoByDB().UpdateOrderById(req.OrderId, uId, payment)
		if err != nil { // 更新订单失败，回滚
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	return
}
