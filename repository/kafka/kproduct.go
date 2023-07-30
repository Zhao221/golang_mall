package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"golang_mall/pkg/utils/log"
	"time"
)


// SendMessage 发送消息
func SendMessage(ctx context.Context, key, topic, value string) error {
	return SendMessagePartitionPar(ctx, key, topic, value, "")
}

// SendMessagePartitionPar 发送消息指定分区
func SendMessagePartitionPar(ctx context.Context, key, topic, value, partitionKey string) error {
	kafka, err := GetClient(key)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(value),
		Timestamp: time.Now(),
	}

	if partitionKey != "" {
		msg.Key = sarama.StringEncoder(partitionKey)
	}

	partition, offset, err := kafka.Producer.SendMessage(msg)

	if err != nil {
		return err
	}

	if kafka.Debug {
		log.LogrusObj.Infoln("发送kafka消息成功",
			zap.Int32("partition", partition),
			zap.Int64("offset", offset))
	}

	return err
}
