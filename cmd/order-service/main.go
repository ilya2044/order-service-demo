package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ilya2044/order-service-demo/internal/api"
	"github.com/ilya2044/order-service-demo/internal/cache"
	"github.com/ilya2044/order-service-demo/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Не удалось загрузить .env файл")
	}

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
		log.Fatal("Ошибка загрузки кэша:", err)
	}

	apiHandler := api.NewAPI(dbConn, c)

	r := gin.Default()

	r.GET("/order/:id", apiHandler.GetOrder)

	r.StaticFile("/", "./web/index.html")

	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
