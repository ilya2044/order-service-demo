package cache

import (
	"log"
	"sync"

	"github.com/ilya2044/order-service-demo/internal/models"
	"gorm.io/gorm"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]models.Order),
	}
}

func (c *Cache) Set(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *Cache) Get(id string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orders[id]
	return order, ok
}

func (c *Cache) LoadFromDB(db *gorm.DB) error {
	var orders []models.Order
	if err := db.Preload("Delivery").Preload("Payment").Preload("Items").
		Find(&orders).Error; err != nil {
		return err
	}

	for _, o := range orders {
		c.Set(o)
	}
	log.Printf("Кэш восстановлен из бд (%d заказов)", len(orders))
	return nil
}
