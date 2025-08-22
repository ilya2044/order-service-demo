package handler

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ilya2044/order-service-demo/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	db    *gorm.DB
	cache Cache
}

type Cache interface {
	Set(models.Order)
}

func NewHandler(db *gorm.DB, cache Cache) *Handler {
	return &Handler{db: db, cache: cache}
}

func validateOrder(o *models.Order) error {
	if o.OrderUID == "" {
		return errors.New("order_uid is empty")
	}
	if o.Payment.Transaction == "" {
		return errors.New("payment.transaction is empty")
	}
	if o.DateCreated.IsZero() {
		o.DateCreated = time.Now().UTC()
	}
	if len(o.Items) == 0 {
		return errors.New("items is empty")
	}
	return nil
}

func (h *Handler) ProcessMessage(msg []byte) error {
	var order models.Order
	if err := json.Unmarshal(msg, &order); err != nil {
		log.Println("Invalid JSON:", err)
		return err
	}
	if err := validateOrder(&order); err != nil {
		log.Println("Validation error:", err)
		return err
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "order_uid"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"track_number", "entry", "locale", "internal_signature",
					"customer_id", "delivery_service", "shard_key", "sm_id",
					"date_created", "oof_shard",
				}),
			},
		).Create(&order).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "order_uid"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "phone", "zip", "city", "address", "region", "email"}),
		}).Create(&order.Delivery).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "order_uid"}},
			DoUpdates: clause.AssignmentColumns([]string{"transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"}),
		}).Create(&order.Payment).Error; err != nil {
			return err
		}
		if err := tx.Where("order_uid = ?", order.OrderUID).Delete(&models.Item{}).Error; err != nil {
			return err
		}
		for i := range order.Items {
			order.Items[i].OrderUID = order.OrderUID
		}
		if len(order.Items) > 0 {
			if err := tx.Create(&order.Items).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Println("db tx error:", err)
		return err
	}

	h.cache.Set(order)
	log.Printf("Заказ %s сохранён в бд и кэширован\n", order.OrderUID)
	return nil
}
