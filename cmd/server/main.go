package main

import (
	"log"
	"os"

	"github.com/Lumi4s/PromptRelay/internal/bot"
	"github.com/Lumi4s/PromptRelay/internal/comfy"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
	}

	workflows, err := comfy.LoadWorkflows()
	if err != nil {
		log.Fatalf("failed to load workflows: %v", err)
	}

	client, err := comfy.New(os.Getenv("COMFY_URL"))
	if err != nil {
		log.Fatalf("comfy is unreachable: %v", err)
	}

	tgBot, err := bot.New(os.Getenv("BOT_TOKEN"), workflows, client)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	log.Println("Program started")

	tgBot.Run()
}
