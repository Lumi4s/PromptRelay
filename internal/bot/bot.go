package bot

import (
	"log"

	"github.com/Lumi4s/PromptRelay/internal/comfy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api       *tgbotapi.BotAPI
	workflows map[string][]byte
	client    *comfy.Client
}

func New(token string, workflows map[string][]byte, comfy *comfy.Client) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:       api,
		workflows: workflows,
		client:    comfy,
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	selectedWorklow := "Krea2.json"

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			workflowBytes, err := comfy.BuildWorkflow(update.Message.Text, b.workflows[selectedWorklow], selectedWorklow)
			if err != nil {
				log.Printf("failed to build workflow: %v", err)
				continue
			}

			if err = comfy.ParseWorkflow(workflowBytes); err != nil {
				log.Printf("failed to parse workflow: %v", err)
				continue
			}

			resp, err := b.client.SendWorkflow(workflowBytes)
			if err != nil {
				log.Printf("%v", err)
			}

			filename, err := b.client.WaitForResultAndGetName(resp)
			if err != nil {
				log.Printf("%v", err)
			}

			imgBytes, err := b.client.GetImageBytes(filename)
			if err != nil {
				log.Printf("%v", err)
			}

			b.sendImage(update.Message.Chat.ID, imgBytes, filename)

		}
	}
}

func (b *Bot) sendImage(chatID int64, imgBytes []byte, filename string) error {
	file := tgbotapi.FileBytes{
		Name:  filename,
		Bytes: imgBytes,
	}

	msg := tgbotapi.NewPhoto(chatID, file)
	if _, err := b.api.Send(msg); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Success! %v was sent to chatID %v", filename, chatID)

	return nil
}
