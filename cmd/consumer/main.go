package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ilya2044/order-service-demo/internal/cache"
	"github.com/ilya2044/order-service-demo/internal/db"
	"github.com/ilya2044/order-service-demo/internal/handler"
	"github.com/ilya2044/order-service-demo/internal/kafka"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	dbConn := db.InitDB(dsn)

	c := cache.NewCache()

	if err := c.LoadFromDB(dbConn); err != nil {
		log.Println("Cache warmup error:", err)
	}

	h := handler.NewHandler(dbConn, c)

	consumer, err := kafka.NewConsumer(
		os.Getenv("KAFKA_GROUP_ID"),
		os.Getenv("KAFKA_TOPIC"),
	)
	if err != nil {
		log.Fatal("Kafka consumer error:", err)
	}

	consumer.Start(func(b []byte) error {
		return h.ProcessMessage(b)
	})
}
