package kafka

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(group, topic string) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("KAFKA_BROKER"),
		"group.id":           group,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, err
	}

	return &Consumer{consumer: c}, nil
}

func (c *Consumer) Start(handler func([]byte)) {
	log.Println("Consumer started")
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			log.Println("Kafka error:", err)
			continue
		}
		handler(msg.Value)
	}
}
