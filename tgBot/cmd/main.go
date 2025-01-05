package main

import (
	"log"
	"os"
	telegramclient "tgBot/clients/telegram"
	eventconsumer "tgBot/consumer/event-consumer"
	"tgBot/db"
	"tgBot/events/telegram"
	"tgBot/storage/files"
)

func mustToken() string {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("Bot token not found")
	}
	return botToken
}

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	eventsProcessor := telegram.New(telegramclient.NewClient("api.telegram.org", mustToken()), files.NewStorage("storage", db))

	log.Print("Starting bot...")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		log.Fatal("Starting consumer failed: ", err)
	}
}
