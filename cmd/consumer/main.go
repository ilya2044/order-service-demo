package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ilya2044/order-service-demo/internal/db"
	"github.com/ilya2044/order-service-demo/internal/handler"
	"github.com/ilya2044/order-service-demo/internal/kafka"
	"github.com/ilya2044/order-service-demo/internal/models"
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

	cache := make(map[string]models.Order)

	h := handler.NewHandler(dbConn, cache)

	var orders []models.Order
	dbConn.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders)
	for _, o := range orders {
		cache[o.OrderUID] = o
	}
	log.Printf("Кэш восстановлен, %d заказов загружено", len(cache))

	consumer, err := kafka.NewConsumer(
		os.Getenv("KAFKA_GROUP_ID"),
		os.Getenv("KAFKA_TOPIC"),
	)
	if err != nil {
		log.Fatal("Kafka consumer error:", err)
	}

	go consumer.Start(h.ProcessMessage)

	select {}
}
