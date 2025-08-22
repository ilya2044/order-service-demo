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
	return &Cache{orders: make(map[string]models.Order)}
}

func (c *Cache) Set(order models.Order) {
	c.mu.Lock()
	c.orders[order.OrderUID] = order
	c.mu.Unlock()
}

func (c *Cache) Get(id string) (models.Order, bool) {
	c.mu.RLock()
	o, ok := c.orders[id]
	c.mu.RUnlock()
	return o, ok
}

func (c *Cache) LoadFromDB(db *gorm.DB) error {
	var orders []models.Order
	if err := db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error; err != nil {
		return err
	}
	c.mu.Lock()
	for _, o := range orders {
		c.orders[o.OrderUID] = o
	}
	n := len(c.orders)
	c.mu.Unlock()
	log.Printf("Кэш восстановлен из БД (%d заказов)", n)
	return nil
}
