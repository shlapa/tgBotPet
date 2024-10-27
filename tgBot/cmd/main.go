package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func mustToken() string {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("Bot token not found")
	}
	return botToken
}

func main() {
	bot, err := tgbotapi.NewBotAPI(mustToken())
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
}
