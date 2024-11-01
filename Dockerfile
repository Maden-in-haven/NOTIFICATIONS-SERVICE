# Используем официальный образ Go для сборки приложения
FROM golang:1.23.2-alpine AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum отдельно для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Переходим в директорию с main.go
WORKDIR /app/cmd/app

# Собираем бинарный файл Go-приложения
RUN go build -o /app/myapp

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из фазы сборки
COPY --from=build /app/myapp .

# Открываем порт 8080 для работы приложения
EXPOSE 8080

# Запускаем приложение
CMD ["./myapp"]

ENV POSTGRESQL_HOST=db.crm.evil-chan.ru
ENV POSTGRESQL_PORT=5432
ENV POSTGRESQL_USER=gen_user
ENV POSTGRESQL_PASSWORD=m%3A0oC.h%3F3L_WKl
ENV POSTGRESQL_DBNAME=default_db
ENV JWT_SECRET_KEY=пушкабомба!