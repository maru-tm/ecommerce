FROM golang:1.24 AS builder

WORKDIR /app

# Копируем зависимости и загружаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код и собираем приложение
COPY . .
RUN go build -o order-service ./cmd/main.go

# Используем минимальный образ Debian для исполнения
FROM debian:bookworm-slim

# Устанавливаем необходимые сертификаты для HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Указываем рабочую директорию
WORKDIR /app

# Копируем бинарник и настройки из предыдущего шага
COPY --from=builder /app/order-service .
COPY .env .env

# Открываем порт
EXPOSE 50052

# Запускаем приложение
CMD ["./order-service"]
