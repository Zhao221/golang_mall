package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"strconv"
)

type Product struct {
	Id    int
	Name  string
	Title string
}

func NewProduct() error {
	brokers := []string{"localhost:9092"}
	producer, err := sarama.NewAsyncProducer(brokers, Kcfg)
	if err != nil {
		return err
	}
	p := &Product{
		Id: 1,
		Name: "钻戒",
		Title: "那戒指的质地似乎是钻石制成的吧，闪闪发光又不失内敛，清雅又不失高贵，阳光洒下来，发出淡淡的光，和淡淡的清香，有着像是通了灵般的仙气",
	}

	key := sarama.StringEncoder(strconv.Itoa(p.Id))
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: "new-products",
		Key:   key,
		Value: sarama.ByteEncoder(value),
	}
	producer.Input() <- msg
	return nil
}
