package main

import (
	"log"
	"net/http"
	"notifications/internal/handler"
	"notifications/internal/middlewares"
	"notifications/internal/tgBot/service"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализируем новый маршрутизатор
	go service.StartTgBot()
	r := mux.NewRouter()

	r.Handle("/api/notifications/subscribe/telegram", middlewares.JWTAuthentication(http.HandlerFunc(handler.SubscribeTG))).Methods("GET")
	r.Handle("/api/notifications/send", middlewares.JWTAuthentication(http.HandlerFunc(handler.Notifications))).Methods("POST")

	// Запускаем сервер на порту 8080
	port := ":8080"
	log.Printf("Сервер запущен на порту %s", port)

	// Логируем ошибки, если таковые возникнут при запуске
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
