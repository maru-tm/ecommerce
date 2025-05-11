FROM golang:1.24 AS builder

WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Компилируем приложение в бинарник
RUN go build -o user-service ./cmd/main.go

# Используем минимальный образ Debian для исполнения
FROM debian:bookworm-slim

# Устанавливаем сертификаты для HTTPS (если это нужно для работы с внешними API)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Указываем рабочую директорию
WORKDIR /app

# Копируем скомпилированный бинарник из предыдущего шага
COPY --from=builder /app/user-service .

# Копируем файл .env с настройками
COPY .env .

# Открываем порт для сервиса
EXPOSE 50051

# Запускаем скомпилированный бинарник
CMD ["./user-service"]
