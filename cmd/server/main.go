package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/Lumi4s/PromptRelay/internal/bot"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
	}

	log.Println("Program started")

	tgBot, err := bot.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	tgBot.Run()
}
