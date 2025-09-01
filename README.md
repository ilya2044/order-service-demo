# Order Service with Kafka and PostgreSQL

Демонстрационной сервис заказов, реализованный на Go
Демонстрация - https://drive.google.com/drive/folders/19Mc1TCTGjpQ0y2w5xBrCDaBG-1e7TZhX?usp=sharing

## Использует:
- PostgreSQL для хранения данных,
- Kafka для асинхронной обработки заказов,
- Docker Compose для запуска инфраструктуры.

## Сервис включает:
- Producer — создает заказы и публикует их в Kafka,
- Consumer — обрабатывает сообщения из Kafka и обновляет базу данных,
- API для взаимодействия с заказами.

## Установка и запуск

1. Клонировать репозиторий
```
git clone https://github.com/ilya2044/order-service-demo.git
cd order-service
```
2. Скачать зависимости
```
go mod download
```
3. Запустить инфраструктуру
```
docker-compose up -d
```
7. Запустить сервисы
```
## Producer
go run cmd/producer/main.go
## Consumer
go run cmd/consumer/main.go
## API
go run cmd/order-service/main.go
```
9. Перейти на http://localhost:8081 и выполнить GET запрос

## Структура проекта

```
.
├── cmd
│   ├── order-service   # API сервис
│   │   └── main.go
│   ├── producer        # Kafka producer
│   │   └── main.go
│   └── consumer        # Kafka consumer
│       └── main.go
├── internal
│   ├── api             # Работа с апи
│   ├── cache           # Кэш
│   ├── db              # Работа с БД
│   ├── handler          
│   ├── kafka           # Kafka producer/consumer
│   └── models          # Модели данных
├── web
├── .env
├── docker-compose.yml
├── go.mod
└── README.md
```
## Восстановление кеша при рестарте
Consumer при старте:
- Загружает актуальные данные из PostgreSQL в кеш,
- Продолжает обслуживание запросов без задержек.

Kafka UI:
http://localhost:8080




