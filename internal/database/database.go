package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"notifications/internal/tgBot/util"
	"github.com/jackc/pgx/v5"
	db "github.com/Maden-in-haven/crmlib/pkg/database"
)

func CheckTelegramLink(userID string) (string, error) {
	query := `SELECT telegram_id FROM telegram_user_links WHERE user_id = $1`

	var telegramID sql.NullInt64
	err := db.DB.Pool.QueryRow(context.Background(), query, userID).Scan(&telegramID)

	if err != nil {
		// Проверяем, если ошибка о том, что не найдено ни одной строки
		if err ==  pgx.ErrNoRows {
			// Если запись не найдена, возвращаем пустую строку
			return "", nil
		}
		// Если произошла другая ошибка, возвращаем её
		return "", err
	}

	// Если `telegramID` найден и валиден, получаем username
	if telegramID.Valid {
		usernameTG, err := util.GetUsernameByTelegramID(telegramID.Int64)
		if err != nil {
			return "", err
		}
		return usernameTG, nil
	}

	// Если `telegram_id` отсутствует, возвращаем пустую строку
	return "", nil // заглушка
}

var ErrAlreadyExists = errors.New("связь с данным telegram_id уже существует")

func SaveUserLink(userTgID int64, userID string) (string, error) {
	// Проверка, существует ли уже связь
	err := db.DB.Pool.QueryRow(context.Background(), "SELECT user_id FROM telegram_user_links WHERE telegram_id = $1", userTgID).Scan(new(string))

	if err != nil {
		if err == pgx.ErrNoRows {
			// Если записи нет, продолжаем вставку
			_, insertErr := db.DB.Pool.Exec(context.Background(), "INSERT INTO telegram_user_links (user_id, telegram_id) VALUES ($1, $2)", userID, userTgID)
			if insertErr != nil {
				log.Printf("Ошибка при сохранении связи пользователя %s с telegram_id %d: %s", userID, userTgID, insertErr)
				return "", insertErr
			}
			log.Printf("Пользователь %s привязан к telegram_id %d", userID, userTgID)

			// Получаем имя пользователя
			username, err := db.DB.GetUserByID(context.Background(), userID)
			if err != nil {
				log.Printf("Ошибка при получении имени пользователя для userID %s: %s", userID, err)
				return "", err
			}
			return username.Username, nil
		}
		// Ошибка при выполнении запроса
		log.Printf("Ошибка при проверке существующей связи для telegram_id %d: %s", userTgID, err)
		return "", err
	}

	// Если запись уже существует, возвращаем ошибку
	return "", ErrAlreadyExists
}


func GetTelegramID(userID string) (int64, error) {
    var telegramID int64
    query := "SELECT telegram_id FROM telegram_user_links WHERE user_id = $1"

    err := db.DB.Pool.QueryRow(context.Background(), query, userID).Scan(&telegramID)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, nil // Не найден
        }
        return 0, err // Возвращаем ошибку
    }
    return telegramID, nil
}
