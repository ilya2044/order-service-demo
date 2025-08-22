package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/ilya2044/order-service-demo/internal/kafka"
	"github.com/joho/godotenv"
)

const (
	topic = "orders"
)

func main() {
	_ = godotenv.Load()
	p, err := kafka.NewProducer()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		orderUID := uuid.New().String()

		msg := fmt.Sprintf(`
{
  "order_uid": "%s",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "%s",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
`, orderUID, orderUID)

		if i == 0 {
			msg = "test invalid message"
		}

		if err := p.Produce(msg, topic); err != nil {
			log.Print(err)
		}
	}
}
