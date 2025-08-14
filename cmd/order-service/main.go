package main

import (
	"fmt"
	"os"

	"github.com/ilya2044/order-service-demo/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Не удалось загрузить .env файл")
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
	fmt.Println("БД подключена:", dbConn != nil)
}
