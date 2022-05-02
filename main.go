package main

import (
	"log"

	tgClient "TelegramBot/MyFirstTgBot/clients/telegram"
	"TelegramBot/MyFirstTgBot/config"
	event_consumer "TelegramBot/MyFirstTgBot/consumer/event-consumer"
	"TelegramBot/MyFirstTgBot/events/telegram"
	"TelegramBot/MyFirstTgBot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	cfg := config.MustLoad()
	storage := files.New(storagePath)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TgBotHost),
		storage,
	)

	log.Print("service started") //сообщение, что сервис запущен

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stoped", err)
	}
}
