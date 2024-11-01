package tgBot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Создайте экземпляр бота
var Bot *tgbotapi.BotAPI

func init() {
	var err error
	Bot, err = tgbotapi.NewBotAPI("7526416842:AAG2GlhK4R1UHzp6STpRHdzda-Mggbl8-6w")
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота: %s", err)
	}
}
