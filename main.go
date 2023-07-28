package main

import (
	"golang_mall/conf"
	"golang_mall/dao/mysql"
	"golang_mall/dao/redis"
	"golang_mall/pkg/utils/log"
	"golang_mall/pkg/utils/timer"
	"golang_mall/repository/es"
	"golang_mall/routes"
)

func main() {
	loading() // 加载配置
	r := routes.NewRouter()
	_ = r.Run(conf.Config.System.HttpPort)
}

func loading() {
	conf.InitConfig()
	mysql.InitMySQL()
	redis.InitRedis()
	timer.InitDayCheckinDayTimer() // 初始化每日定时任务
	timer.InitMonthCheckinTimer() // 初始化每月定时任务
	log.InitLog()
	log.InitLogger() // 如果接入ELK请进入这个func打开注释
	es.InitEs()      // 如果需要接入ELK可以打开这个注释

	// rabbitmq.InitRabbitMQ() // 如果需要接入RabbitMQ可以打开这个注释
	// kafka.InitKafka()
	// track.InitJaeger( )
	go scriptStarting()
}

func scriptStarting() {
	// 启动一些脚本
}
