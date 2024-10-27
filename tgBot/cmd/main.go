package main

import (
	"log"
	"os"
	telegram2 "tgBot/clients/telegram"
	event_consumer "tgBot/consumer/event-consumer"
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
	eventsProcessor := telegram.New(telegram2.NewClient("api.telegram.org", mustToken()), files.NewStorage("storage"))

	log.Print("Starting bot...")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		log.Fatal("Starting consumer failed: ", err)
	}
}
