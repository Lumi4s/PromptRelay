package bot

import (
	"log"

	"github.com/Lumi4s/PromptRelay/internal/comfy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api       *tgbotapi.BotAPI
	workflows map[uint8][]byte
	client    *comfy.Client
}

func New(token string, workflows map[uint8][]byte, comfy *comfy.Client) (*Bot, error) {
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
	var selectedWorkflow uint8 = 1

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			workflowBytes, err := comfy.BuildWorkflow("prompt", b.workflows[selectedWorkflow])
			if err != nil {
				log.Printf("failed to build workflow: %v", err)
				continue
			}

			if err = comfy.ParseWorkflow(workflowBytes); err != nil {
				log.Printf("failed to parse workflow: %v", err)
				continue
			}

			resp, err := b.client.SendWorkflow(workflowBytes)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
			msg.ReplyToMessageID = update.Message.MessageID

			if err != nil {
				log.Printf("failed to sendWorkflow: %v", err)
			}

			_, err = b.api.Send(msg)
			if err != nil {
				log.Printf("failed to send message: %v", err)
			}
		}
	}
}
