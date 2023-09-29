package main

import (
	"github.com/Vivirinter/telegram-bot-url/internal/telegram"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if _, ok := os.LookupEnv("DOCKER_ENV"); !ok {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	tgBotToken, ok := os.LookupEnv("TELEGRAM_TOKEN")
	if !ok {
		log.Fatal("Telegram Bot Token not provided")
	}

	if err := startBot(tgBotToken); err != nil {
		log.Fatal(err)
	}
}

func startBot(tgBotToken string) error {
	bot, err := telegram.NewBot(tgBotToken)
	if err != nil {
		return err
	}

	return bot.Start()
}
