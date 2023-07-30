package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"golang_mall/pkg/utils/log"
	"strings"
)

func Consumer(ctx context.Context, key, topic string, fn func(message *sarama.ConsumerMessage) error) (err error) {
	kafka, err := GetClient(key)
	if err != nil {
		return
	}

	partitions, err := kafka.Consumer.Partitions(topic)
	if err != nil {
		return
	}

	for _, partition := range partitions {
		// 针对每个分区创建一个对应的分区消费者
		offset, errx := kafka.Client.GetOffset(topic, partition, sarama.OffsetNewest)
		if errx != nil {
			log.LogrusObj.Infoln("获取Offset失败:", errx)
			continue
		}

		pc, errx := kafka.Consumer.ConsumePartition(topic, partition, offset)
		if errx != nil {
			log.LogrusObj.Infoln("获取Offset失败:", errx)
			return
		}

		// 从每个分区都消费消息
		go func(consumer sarama.PartitionConsumer) {
			defer func() {
				// if err := recover(); err != nil {
				// 	log.LogrusObj.Error("消费kafka信息发生panic,err:%s", err)
				// }
			}()

			defer func() {
				err := pc.Close()
				if err != nil {
					log.LogrusObj.Infoln("消费kafka信息发生panic,err:%s", err)
				}
			}()

			for {
				select {
				case msg := <-pc.Messages():
					err := middlewareConsumerHandler(fn)(msg)
					if err != nil {
						return
					}
				case <-ctx.Done():
					return
				}
			}

		}(pc)
	}
	return nil
}

func ConsumerGroup(ctx context.Context, key, groupId, topics string, msgHandler ConsumerGroupHandler) (err error) {
	kafka, err := GetClient(key)
	if err != nil {
		return
	}

	if isConsumerDisabled(kafka) {
		return
	}

	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupId, kafka.Client)
	if err != nil {
		return
	}

	go func() {
		defer func() {
			// if err := recover(); err != nil {
			// 	log.LogrusObj.Error("消费kafka发生panic", zap.Any("panic", err))
			// }
		}()

		defer func() {
			err := consumerGroup.Close()
			if err != nil {
				log.LogrusObj.Error("close err", zap.Any("panic", err))
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := consumerGroup.Consume(ctx, strings.Split(topics, ","), ConsumerGroupHandler(func(msg *sarama.ConsumerMessage) error {
					return middlewareConsumerHandler(msgHandler)(msg)
				}))
				if err != nil {
					log.LogrusObj.Error("消费kafka失败 err", zap.Any("panic", err))

				}
			}
		}

	}()
	return
}

func isConsumerDisabled(in *Kafka) bool {
	if in.DisableConsumer {
		log.LogrusObj.Infoln("kafka consumer disabled,key:%s", in.Key)
	}

	return in.DisableConsumer
}

func middlewareConsumerHandler(fn func(message *sarama.ConsumerMessage) error) func(message *sarama.ConsumerMessage) error {
	return func(msg *sarama.ConsumerMessage) error {
		return fn(msg)
	}
}

type ConsumerGroupHandler func(message *sarama.ConsumerMessage) error

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h(msg); err != nil {
			log.LogrusObj.Infoln("消息处理失败",
				zap.String("topic", msg.Topic),
				zap.String("value", string(msg.Value)))
			continue
		}
		sess.MarkMessage(msg, "")
	}

	return nil
}
