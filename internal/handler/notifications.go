package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"notifications/internal/database"
	"notifications/internal/tgBot/util"

	db "github.com/Maden-in-haven/crmlib/pkg/database"
	"github.com/go-playground/validator/v10"
)

type Message struct {
	UserID  string   `json:"user_id" validate:"required,uuid"`
	Service []string `json:"service" validate:"required,dive,oneof=email push telegram vk"`
	Subject string   `json:"subject" validate:"required,min=1,max=255"`
	Type    string   `json:"type" validate:"required,oneof=notification error critical deadline service"`
	Message string   `json:"message" validate:"required,min=1"`
}
type response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var validate = validator.New()

func Notifications(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON из тела запроса
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		respondWithError(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(msg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("Поле '%s' не прошло валидацию: %s", err.Field(), err.Tag())
		}
		respondWithError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	user, _ := db.DB.GetUserByID(context.Background(), msg.UserID)
	if user.ID == "" {
		respondWithError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	for _, service := range msg.Service {
		if service == "push" {

		} else if service == "email" {

		} else if service == "telegram" {
			usernameTG, _ := database.CheckTelegramLink(msg.UserID)
			if usernameTG == "" {
				respondWithError(w, "Пользователь не привязан к телеграмму", http.StatusNotFound)
				return
			}
		} else {

		}
	}

	for _, service := range msg.Service {
		if service == "push" {

		} else if service == "email" {

		} else if service == "telegram" {
			telegramID, err := database.GetTelegramID(msg.UserID)
			if err != nil {
				respondWithError(w, "Ошибка сервера", http.StatusInternalServerError)
				return
			} else if telegramID == 0 {
				respondWithError(w, "Пользователь не привязан к телеграмму", http.StatusNotFound)
				return
			}
			text := fmt.Sprintf(
				"==== Уведомление ====\n\n"+
					"📌 **Тема:** %s\n"+
					"🔔 **Тип:** %s\n"+
					"📝 **Сообщение:** %s\n\n"+
					"Спасибо за внимание.",
				msg.Subject, msg.Type, msg.Message)
			if err = util.SendMessage(telegramID, text); err != nil {
				log.Printf("Ошибка при отправке сообщения пользователю %d: %s", telegramID, err)
				respondWithError(w, "Не удалось отправить сообщение", http.StatusInternalServerError)
				return
			}
		} else {

		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{
		Status:  "success",
		Message: "Уведомление успешно отправлено.",
	})
}
func respondWithError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response{
		Status:  "error",
		Message: message,
	})
}
