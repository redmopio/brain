package channels

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minskylab/brain/config"
)

type (
	TelegramResponseFunc func(ctx context.Context, senderID string, message string) (string, error)
	TelegramConnector    struct {
		apiKey   string
		response TelegramResponseFunc
	}
)

func NewTelegramConnector(config *config.Config, response TelegramResponseFunc) *TelegramConnector {
	return &TelegramConnector{
		apiKey:   config.TelegramAPIKey,
		response: response,
	}
}

func (t *TelegramConnector) Connect(ctx context.Context) {
	bot, err := tgbotapi.NewBotAPI(t.apiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Println(update)

		if update.ChannelPost != nil {
			log.Printf("[%s] %s", update.ChannelPost.From.UserName, update.ChannelPost.Text)
		}
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			responseMessage, err := t.response(ctx, update.Message.From.UserName, update.Message.Text)
			if err != nil {
				log.Println(err)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseMessage)
			// msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func (t *TelegramConnector) Disconnect(ctx context.Context) {
	// w.client.Disconnect()
}

func (t *TelegramConnector) SendMessage(ctx context.Context, sender string, message string) (string, error) {
	return t.response(ctx, sender, message)
}
