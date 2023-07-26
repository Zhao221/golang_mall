package service

import (
	"context"
	"golang_mall/dao/mysql"
	"golang_mall/model"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/types"
	"sync"
)

type AddressSrv struct{}

var AddressServ *AddressSrv
var AddressServOnce sync.Once

func GetAddressSrv() *AddressSrv {
	AddressServOnce.Do(func() {
		AddressServ = &AddressSrv{}
	})
	return AddressServ
}

func (a *AddressSrv) AddressCreate(c context.Context, req types.AddressCreateReq) (err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return err
	}
	addressDao := mysql.NewAddressDao(c)
	address := &model.Address{
		UserID:  u.Id,
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
	}
	err = addressDao.CreateAddress(address)
	if err != nil {
		return err
	}
	return err
}

func (a *AddressSrv) AddressShow(c context.Context, req types.AddressGetReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil ,err
	}

	address, err := mysql.NewAddressDao(c).GetAddressByAid(req.Id, u.Id)
	if err != nil {
		return nil ,err
	}
	resp = &types.AddressResp{
		ID:        address.ID,
		UserID:    address.UserID,
		Name:      address.Name,
		Phone:     address.Phone,
		Address:   address.Address,
		CreatedAt: address.CreatedAt.Unix(),
	}
	return resp, err
}

func (a *AddressSrv) AddressList(c context.Context, req types.AddressListReq) (resp interface{}, err error) {
	u, _ := ctl.GetUserInfo(c)
	resp, err = mysql.NewAddressDao(c).ListAddressByUid(u.Id)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (a *AddressSrv) AddressUpdate(c context.Context, req types.AddressServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	addressDao := mysql.NewAddressDao(c)
	address := &model.Address{
		UserID:  u.Id,
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
	}
	err = addressDao.UpdateAddressById(req.Id, address)
	if err != nil {
		return resp, err
	}
	resp = address
	return resp, err
}

func (a *AddressSrv) AddressDelete(c context.Context, req types.AddressDeleteReq) (err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return  err
	}
	err = mysql.NewAddressDao(c).DeleteAddressById(req.Id, u.Id)
	if err != nil {
		return err
	}
	return err
}
