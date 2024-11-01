package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notifications/internal/database"
)

type Response struct {
	URL string `json:"url"`
}

func SubscribeTG(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	log.Printf("Начало подписки пользователя %s", userID)

	usernameTG, err := database.CheckTelegramLink(userID)
	if err != nil {
		log.Printf("Ошибка при проверке Telegram-связи для пользователя %s: %s", userID, err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if usernameTG != "" {
		log.Printf("Пользователь %s уже связан с Telegram: %s", userID, usernameTG)
		http.Error(w, fmt.Sprintf("Пользователь уже связан с Telegram: %s", usernameTG), http.StatusConflict)
		return
	}
	url := fmt.Sprintf("https://t.me/tirazhi_crm_bot?start=%s", userID)
	response := Response{URL: url}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при кодировании ответа для пользователя %s: %s", userID, err)
	}
	log.Printf("Ссылка на Telegram для пользователя %s: %s", userID, url)
	log.Printf("Обработка подписки для пользователя %s завершена успешно", userID)
}


// func GetProfile(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Header.Get("X-User-ID")
// 	user, err := database.DB.GetUserByID(context.Background(), userID)
// 	if err != nil {
// 		http.Error(w, "Пользователь не найден", http.StatusNotFound)
// 		return
// 	}
// 	universalRequest := UniversalRequest{
// 		ID:       userID,
// 		Username: user.Username,
// 		Role:     user.Role,
// 	}
// 	if user.Role == "Admin" {
// 		admin, err := database.DB.GetAdminByID(context.Background(), userID)
// 		if err != nil {
// 			http.Error(w, "Ошибка получения данных администратора", http.StatusInternalServerError)
// 			return
// 		}
// 		universalRequest.Permissions = &admin.Permissions

// 	} else if user.Role == "Manager" {
// 		manager, err := database.DB.GetManagerByID(context.Background(), userID)
// 		if err != nil {
// 			http.Error(w, "Ошибка получения данных менеджера", http.StatusInternalServerError)
// 			return
// 		}
// 		universalRequest.FullName = &manager.FullName
// 		universalRequest.HireDate = &manager.HireDate
// 	} else if user.Role == "Client" {
// 		client, err := database.DB.GetClientByID(context.Background(), userID)
// 		if err != nil {
// 			http.Error(w, "Ошибка получения данных клиента", http.StatusInternalServerError)
// 			return
// 		}
// 		universalRequest.FullName = &client.FullName
// 		universalRequest.PhoneNumber = &client.PhoneNumber
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	if err := json.NewEncoder(w).Encode(universalRequest); err != nil {
// 		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
// 		return
// 	}
// }
