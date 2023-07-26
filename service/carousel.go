package service

import (
	"context"
	"golang_mall/dao/mysql"
	"golang_mall/types"
	"sync"
)

var CarouselSrvIns *CarouselSrv
var CarouselSrvOnce sync.Once

type CarouselSrv struct {
}

func GetCarouselSrv() *CarouselSrv {
	CarouselSrvOnce.Do(func() {
		CarouselSrvIns = &CarouselSrv{}
	})
	return CarouselSrvIns
}

func (s *CarouselSrv) ListCarousel(ctx context.Context) (resp interface{}, err error) {
	carousels, err := mysql.NewCarouselDao(ctx).ListCarousel()
	if err != nil {
		return nil, err
	}
	resp = &types.DataListResp{
		Item:  carousels,
		Total: int64(len(carousels)),
	}
	return resp, err
}
