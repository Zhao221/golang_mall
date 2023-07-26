package rabbitmq

import (
	"fmt"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"

	conf "github.com/CocaineCong/gin-mall/config"
)

// GlobalRabbitMQ rabbitMQ链接单例
var GlobalRabbitMQ *amqp.Connection

// InitRabbitMQ 在中间件中初始化rabbitMQ链接
func InitRabbitMQ() {
	rConfig := conf.Config.RabbitMq
	pathRabbitMQ := strings.Join([]string{rConfig.RabbitMQ, "://", rConfig.RabbitMQUser, ":", rConfig.RabbitMQPassWord, "@", rConfig.RabbitMQHost, ":", rConfig.RabbitMQPort, "/"}, "")
	conn, err := amqp.Dial(pathRabbitMQ)
	if err != nil {
		fmt.Println("rabbitmq连接失败",err)
		return
	}
	GlobalRabbitMQ = conn
}
