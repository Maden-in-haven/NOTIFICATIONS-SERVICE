// package tgbot

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// var bot *tgbotapi.BotAPI

// type MessageRequest struct {
// 	ChatID  int64  `json:"chat_id"`
// 	Message string `json:"message"`
// }

// func sendMessage(chatID int64, message string) {
// 	msg := tgbotapi.NewMessage(chatID, message)
// 	bot.Send(msg)
// }

// func messageHandler(w http.ResponseWriter, r *http.Request) {
// 	var req MessageRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	sendMessage(req.ChatID, req.Message)
// 	w.WriteHeader(http.StatusOK)
// }

// func main() {
// 	var err error
// 	bot, err = tgbotapi.NewBotAPI("7526416842:AAG2GlhK4R1UHzp6STpRHdzda-Mggbl8-6w")
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	http.HandleFunc("/send", messageHandler)

// 	log.Println("Starting server on :8080")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatal(err)
// 	}
// }

package service

import (
	"errors"
	"log"
	"notifications/internal/database"
	"notifications/internal/tgBot"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartTgBot() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tgBot.Bot.GetUpdatesChan(u)
	log.Println("Телеграмм бот запущен")
	for update := range updates {
		if update.Message != nil {
			if strings.HasPrefix(update.Message.Text, "/start") {
				handleStartCommand(update.Message)
			}
		}
	}
}

func handleStartCommand(message *tgbotapi.Message) {
	// Извлекаем uuid из команды
	parts := strings.Split(message.Text, " ")
	if len(parts) > 1 {
		uuid := parts[1] // uuid будет вторым элементом
		// Здесь вы можете сохранить uuid и ID пользователя в базе данных
		usernsme, err := database.SaveUserLink(message.From.ID, uuid)
		if err != nil {
			// Логируем ошибку
			log.Printf("Ошибка при сохранении связи для telegram_id %d: %s", message.From.ID, err)

			// В зависимости от ошибки, возвращаем разные сообщения
			if errors.Is(err, database.ErrAlreadyExists) {
				// Если связь уже существует
				reply := "Вы уже связаны с этим аккаунтом."
				msg := tgbotapi.NewMessage(message.Chat.ID, reply)
				tgBot.Bot.Send(msg)
			} else {
				// Общее сообщение для других ошибок
				reply := "Не удалось сохранить связь с вашим аккаунтом. Свяжитесь с разработчиками"
				msg := tgbotapi.NewMessage(message.Chat.ID, reply)
				tgBot.Bot.Send(msg)
			}
			return
		}
		reply := "Привет! Ваш аккаунт привязан к " + usernsme
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		tgBot.Bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Используйте ссылку, которую получили в приложении")
		tgBot.Bot.Send(msg)
	}
}
