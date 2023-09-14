package channels

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minskylab/brain/config"
)

type (
	TelegramGroupHandler     func(ctx context.Context, groupId int64, groupName string) (string, error)
	TelegramUserGroupHandler func(ctx context.Context, groupId int64, sender string) (string, error)
	TelegramResponseFunc     func(ctx context.Context, groupId int64, senderId string, message string) (string, error)
	TelegramConnector        struct {
		apiKey           string
		groupHandler     TelegramGroupHandler
		userGroupHandler TelegramUserGroupHandler
		response         TelegramResponseFunc
	}
)

func NewTelegramConnector(config *config.Config, groupHandler TelegramGroupHandler, userGroupHandler TelegramUserGroupHandler, response TelegramResponseFunc) *TelegramConnector {
	return &TelegramConnector{
		apiKey:           config.TelegramAPIKey,
		groupHandler:     groupHandler,
		userGroupHandler: userGroupHandler,
		response:         response,
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
		if update.ChannelPost != nil {
			log.Printf("[%s] %s", update.ChannelPost.From.UserName, update.ChannelPost.Text)
		}
		if update.Message != nil { // If we got a message
			log.Printf("\n--> [Group:%d-%s][User:%s] %s", update.Message.Chat.ID, update.Message.Chat.Title, update.Message.From.UserName, update.Message.Text)

			_, err := t.groupHandler(ctx, update.Message.Chat.ID, update.Message.Chat.Title)
			if err != nil {
				log.Printf("\nError handling group: %s\n", err)
			}

			_, err = t.userGroupHandler(ctx, update.Message.Chat.ID, update.Message.From.UserName)
			if err != nil {
				log.Printf("\nError handling user group: %s\n", err)
			}

			responseMessage, err := t.response(ctx, update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)
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
	return t.response(ctx, 0, sender, message)
}
