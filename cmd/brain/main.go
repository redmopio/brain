package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/config"
	"github.com/minskylab/brain/self"
	"go.mau.fi/whatsmeow/types"
)

func main() {
	config, err := config.NewLoadedConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("Brain Engine is starting...")

	brain, err := self.NewBrainEngine(config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Brain Engine is ready to work!")

	ctx := context.Background()

	whatsAppChannel := channels.NewWhatsAppConnector(config,
		func(ctx context.Context, chatJID types.JID, groupName string) (string, error) {
			return brain.HandleGroup(ctx, channels.WhatsAppChannel, chatJID.String(), groupName)
		},
		func(ctx context.Context, chatJID types.JID, senderJID types.JID) (string, error) {
			return brain.HandleGroupSender(ctx, channels.WhatsAppChannel, chatJID.String(), senderJID.String())
		},
		func(ctx context.Context, senderJID types.JID, message string) (string, error) {
			return brain.GenerateConversationResponse(ctx, channels.WhatsAppChannel, "", senderJID.String(), message)
		})

	if !config.WhatsAppDisabled {
		go whatsAppChannel.Connect(ctx)
	}

	var telegramChannel *channels.TelegramConnector

	if config.TelegramAPIKey != "" {
		telegramChannel = channels.NewTelegramConnector(config,
			func(ctx context.Context, groupId int64, groupName string) (string, error) {
				return brain.HandleGroup(ctx, channels.TelegramChannel, fmt.Sprintf("%d", groupId), groupName)
			},
			func(ctx context.Context, groupId int64, sender string) (string, error) {
				return brain.HandleGroupSender(ctx, channels.TelegramChannel, fmt.Sprintf("%d", groupId), sender)
			},
			func(ctx context.Context, groupId int64, sender string, message string) (string, error) {
				return brain.GenerateConversationResponse(ctx, channels.TelegramChannel, fmt.Sprintf("%d", groupId), sender, message)
			})

		telegramChannel.Connect(ctx)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	whatsAppChannel.Disconnect(ctx)

	if telegramChannel != nil {
		telegramChannel.Disconnect(ctx)
	}
}
