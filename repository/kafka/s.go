package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"golang_mall/model"
	"log"
)

func Consume() error {
	// 初始化 Kafka 消费者
	brokers := []string{"localhost:9092"}
	consumer, err := sarama.NewConsumer(brokers, Kcfg)
	partitionConsumer, err := consumer.ConsumePartition("newProduct", 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Error consuming partition: %v", err)
		return err
	}
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var product model.Product
			err = json.Unmarshal(msg.Value, &product)
			if err != nil {
				log.Printf("Error unmarshaling product: %v", err)
				return err
			} else {
				fmt.Printf("New product: %+v\n", product)
			}
		case err = <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %v", err)
			return err
		}
	}
}
