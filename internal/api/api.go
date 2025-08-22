package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilya2044/order-service-demo/internal/cache"
	"github.com/ilya2044/order-service-demo/internal/models"
	"gorm.io/gorm"
)

type API struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewAPI(db *gorm.DB, cache *cache.Cache) *API {
	return &API{db: db, cache: cache}
}

func (a *API) GetOrder(c *gin.Context) {
	id := c.Param("id")

	if order, ok := a.cache.Get(id); ok {
		c.JSON(http.StatusOK, order)
		return
	}

	var order models.Order
	if err := a.db.Preload("Delivery").Preload("Payment").Preload("Items").
		First(&order, "order_uid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	a.cache.Set(order)
	c.JSON(http.StatusOK, order)
}
