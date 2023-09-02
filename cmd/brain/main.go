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
	"go.mau.fi/whatsmeow/types"
)

func main() {
	config, err := config.NewLoadedConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("Brain Engine is starting...")

	// engine, err := system.NewSystemEngine(config)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Brain Engine is ready to work!")

	ctx := context.Background()

	whatsAppChannel := channels.NewWhatsAppConnector(config, func(ctx context.Context, sender types.JID, message string) (string, error) {
		// return engine.GenerateConversationResponse(ctx, nil, string(channels.WhatsAppChannelName), sender.String(), message)
		return "", nil
	})

	if !config.WhatsAppDisabled {
		go whatsAppChannel.Connect(ctx)
	}

	var telegramChannel *channels.TelegramConnector

	if config.TelegramAPIKey != "" {
		telegramChannel = channels.NewTelegramConnector(config, func(ctx context.Context, sender string, message string) (string, error) {
			// return engine.GenerateConversationResponse(ctx, nil, string(channels.TelegramChannelName), sender, message)
			return "", nil
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
