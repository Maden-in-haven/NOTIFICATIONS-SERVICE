package util

import (
	"notifications/internal/tgBot"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetUsernameByTelegramID получает username пользователя по Telegram ID
func GetUsernameByTelegramID(telegramID int64) (string, error) {
	chat, err := tgBot.Bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: telegramID,
		},
	})
	if err != nil {
		return "", err
	}

	if chat.UserName != "" {
		return chat.UserName, nil
	}

	return strconv.FormatInt(telegramID, 10), nil
}

func SendMessage(chatID int64, text string) error {
    msg := tgbotapi.NewMessage(chatID, text)
    _, err := tgBot.Bot.Send(msg)
    return err
}