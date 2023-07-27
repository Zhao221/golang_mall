package service

import (
	"context"
	"errors"
	"golang_mall/dao/mysql"
	"golang_mall/model"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/types"
	"sync"
)

var FavoriteSrvIns *FavoriteSrv
var FavoriteSrvOnce sync.Once

type FavoriteSrv struct {
}

func GetFavoriteSrv() *FavoriteSrv {
	FavoriteSrvOnce.Do(func() {
		FavoriteSrvIns = &FavoriteSrv{}
	})
	return FavoriteSrvIns
}

// FavoriteList 涉及用户信息，商品信息，商品老板信息
func (f *FavoriteSrv) FavoriteList(c context.Context, req types.FavoritesServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	favorites, total, err := mysql.NewFavoritesDao(c).ListFavoriteByUserId(u.Id, req.PageSize, req.PageNum)
	if err != nil {
		return nil,err
	}
	resp = &types.DataListResp{
		Item:  favorites,
		Total: total,
	}
	return resp, err
}

// FavoriteCreate 涉及用户信息，商品信息，商品老板信息
func (f *FavoriteSrv) FavoriteCreate(c context.Context, req types.FavoriteCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return nil, err
	}
	// 判断此收藏夹是否存在
	fDao := mysql.NewFavoritesDao(c)
	exist, _ := fDao.FavoriteExistOrNot(req.ProductId, u.Id)
	if exist {
		return nil , errors.New("已经存在了")
	}
	userDao := mysql.NewUserDao(c)
	user, err := userDao.GetUserById(u.Id)
	if err != nil {
		return nil,err
	}

	bossDao := mysql.NewUserDaoByDB()
	boss, err := bossDao.GetUserById(req.BossId)
	if err != nil {
		return nil,err
	}

	product, err := mysql.NewProductDao(c).GetProductById(req.ProductId)
	if err != nil {
		return nil,err
	}

	favorite := &model.Favorite{
		UserID:    u.Id,
		User:      *user,
		ProductID: req.ProductId,
		Product:   *product,
		BossID:    req.BossId,
		Boss:      *boss,
	}
	err = fDao.CreateFavorite(favorite)
	if err != nil {
		return nil,err
	}

	return favorite, err
}

func (f *FavoriteSrv) FavoriteDelete(c context.Context, req types.FavoriteDeleteReq) (resp interface{}, err error) {
	favoriteDao := mysql.NewFavoritesDao(c)
	err = favoriteDao.DeleteFavoriteById(req.Id)
	if err != nil {
		return nil,err
	}
	return
}
