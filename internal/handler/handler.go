package handler

import (
	"encoding/json"
	"log"

	"github.com/ilya2044/order-service-demo/internal/models"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	cache map[string]models.Order
}

func NewHandler(db *gorm.DB, cache map[string]models.Order) *Handler {
	return &Handler{db: db, cache: cache}
}

func (h *Handler) ProcessMessage(msg []byte) {
	var order models.Order
	if err := json.Unmarshal(msg, &order); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}

	if err := h.db.Create(&order).Error; err != nil {
		log.Println("db error:", err)
		return
	}

	h.cache[order.OrderUID] = order

	log.Printf("Заказ %s сохранен в бд и кэширован\n", order.OrderUID)
}
