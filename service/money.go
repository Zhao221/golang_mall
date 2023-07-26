package service

import (
	"context"
	"github.com/spf13/cast"
	"golang_mall/dao/mysql"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/types"
	"sync"
)

type  MoneySrv struct {}

var MoneySrvIns *MoneySrv
var MoneySrvOnce sync.Once

func GetMoneySrv() *MoneySrv{
	MoneySrvOnce.Do(func() {
		MoneySrvIns=&MoneySrv{}
	})
	return MoneySrvIns
}

func(m *MoneySrv)MoneyShow(c context.Context, req types.MoneyShowReq)(resp interface{},err error){
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil ,err
	}
	user, err := mysql.NewUserDao(c).GetUserById(u.Id)
	if err != nil {
		return nil ,err
	}
	money, err := user.DecryptMoney(req.Key)
	if err != nil {
		return nil ,err
	}
	resp = &types.MoneyShowResp{
		UserID:    user.ID,
		UserName:  user.UserName,
		UserMoney: cast.ToString(money),
	}
	return
}
