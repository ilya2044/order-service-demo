FROM golang:1.24-alpine


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service ./cmd/order-service

CMD ["./order-service"]
