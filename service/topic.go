package service

import (
	"context"
	errors2 "github.com/pkg/errors"
	"golang_mall/dao/mysql"
	"golang_mall/model"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/types"
	"sync"
)

type TopicSrv struct{}

var TopicSrvIns *TopicSrv
var TopicSrvOnce sync.Once

func GetTopicSrv() *TopicSrv {
	TopicSrvOnce.Do(func() {
		TopicSrvIns = &TopicSrv{}
	})
	return TopicSrvIns
}

func (t *TopicSrv) GetTopicListHandle(c context.Context, req types.BasePage) (resp interface{}, total int64, err error) {
	limit := req.PageSize
	offset := limit * (req.PageNum - 1)
	var topic []model.Topic
	err = mysql.NewTopicDao(c).Table("topic").Count(&total).
		Limit(limit).Offset(offset).Find(&topic).Error
	if err != nil {
		return nil, total, errors2.New("err = mysql.NewTopicDao(c).Table(\"topic\").Count(&total).\n\t\tLimit(limit).Offset(offset).Find(&topic).Error")
	}
	return resp, total, err
}

func (t *TopicSrv) GetUserTopicListHandle(c context.Context, req types.BasePage) (resp interface{}, total int64, err error) {
	limit := req.PageSize
	offset := limit * (req.PageNum - 1)
	u, err := ctl.GetUserInfo(c)
	var topic []model.Topic
	err = mysql.NewTopicDao(c).Table("topic").Where("user_id=?", u.Id).
		Count(&total).Limit(limit).Offset(offset).Find(&topic).Error
	if err != nil {
		return nil, total, errors2.New("err = mysql.NewTopicDao(c).Table(\"topic\").Where(\"user_id=?\",u.Id).\n\t\tCount(&total).Limit(limit).Offset(offset).Find(&topic).Error")
	}
	return resp, total, err
}

func (t *TopicSrv) SubscribeTopic(c context.Context, req types.SubscribeTopicReq) error {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return err
	}
	createTopic := model.Topic{
		Name:   req.Name,
		UserId: u.Id,
	}
	err = mysql.NewTopicDao(c).Table("topic").Create(&createTopic).Error
	if err != nil {
		return errors2.New("err = mysql.NewTopicDao(c).Table(\"topic\").Create(&createTopic).Error")

	}
	return err
}

func (t *TopicSrv) GetMessageListHandle(c context.Context, req types.GetMessageList) (resp interface{}, total int64, err error) {
	limit := req.PageSize
	offset := limit * (req.PageNum - 1)
	u, err := ctl.GetUserInfo(c)
	var message []types.MessageResp
	err = mysql.NewMessageDao(c).Table("message").Where("user_id=?", u.Id).Where("topic_name=?", req.Name).
		Count(&total).Limit(limit).Offset(offset).Find(&message).Error
	if err != nil {
		return nil, total, errors2.New("err = mysql.NewTopicDao(c).Table(\"topic\").Where(\"user_id=?\",u.Id).\n\t\tCount(&total).Limit(limit).Offset(offset).Find(&topic).Error")
	}
	resp = message
	return resp, total, err
}

func (t *TopicSrv) DeleteMessage(c context.Context, req types.DeleteMessage) (err error) {
	err = mysql.NewMessageDao(c).Table("message").Where("id=?", req.Id).Delete(&model.Message{}).Error
	if err != nil {
		return errors2.New("err = mysql.NewMessageDao(c).Table(\"message\").Where(\"id=?\", req.Id).Delete(&model.Message{}).Error")
	}
	return err
}
