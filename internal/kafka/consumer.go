package kafka

import (
	"log"
	"os"
	"time"

	ck "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *ck.Consumer
	dlqProd  *ck.Producer
	dlqTopic string
}

func NewConsumer(group, topic string) (*Consumer, error) {
	cfg := &ck.ConfigMap{
		"bootstrap.servers":        os.Getenv("KAFKA_BROKER"),
		"group.id":                 group,
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.commit":       false,
		"enable.auto.offset.store": false,
	}
	c, err := ck.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}
	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, err
	}

	p, err := ck.NewProducer(&ck.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER"),
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		dlqProd:  p,
		dlqTopic: envOr("KAFKA_DLQ_TOPIC", "orders_dlq"),
	}, nil
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func (c *Consumer) Start(process func([]byte) error) {
	log.Println("Consumer started")
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		const maxRetries = 3
		ok := false
		var lastErr error

		for attempt := 1; attempt <= maxRetries; attempt++ {
			if err := process(msg.Value); err != nil {
				lastErr = err
				log.Printf("Ошибка обработки (attempt %d/%d): %v\n", attempt, maxRetries, err)
				time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
				continue
			}
			ok = true
			break
		}

		if ok {
			_, _ = c.consumer.StoreMessage(msg)
			_, err := c.consumer.CommitMessage(msg)
			if err != nil {
				log.Println("Commit error:", err)
			}
		} else {
			dlqMsg := &ck.Message{
				TopicPartition: ck.TopicPartition{Topic: &c.dlqTopic, Partition: ck.PartitionAny},
				Value:          msg.Value,
				Headers:        msg.Headers,
				Key:            msg.Key,
			}
			if err := c.dlqProd.Produce(dlqMsg, nil); err != nil {
				log.Println("DLQ produce error:", err)
			}
			_, _ = c.consumer.StoreMessage(msg)
			_, err := c.consumer.CommitMessage(msg)
			if err != nil {
				log.Println("Commit after DLQ error:", err)
			}
			log.Printf("Message sent to DLQ (%s). Last error: %v\n", c.dlqTopic, lastErr)
		}
	}
}
