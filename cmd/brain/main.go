package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/config"
	"github.com/minskylab/brain/self"
	"github.com/pkg/errors"
	"go.mau.fi/whatsmeow/types"
)

func main() {
	config, err := config.NewLoadedConfig()
	if err != nil {
		panic(err)
	}

	brain, err := self.NewBrainEngine(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	whatsAppChannel := channels.NewWhatsAppConnector(config, func(ctx context.Context, sender types.JID, message string) (string, error) {
		response, err := brain.GenerateConversationResponse(ctx, sender.String(), message)
		if err != nil {
			return "", errors.WithStack(err)
		}

		return strings.ReplaceAll(response, "<|im_end|>", ""), nil
	})

	fmt.Println(config.WhatsAppDisable)

	if !config.WhatsAppDisable {
		whatsAppChannel.Connect(ctx)
	}

	if config.TelegramAPIKey != "" {
		telegramChannel := channels.NewTelegramConnector(config, func(ctx context.Context, sender string, message string) (string, error) {
			return brain.GenerateConversationResponse(ctx, sender, message)
		})

		telegramChannel.Connect(ctx)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	whatsAppChannel.Disconnect(ctx)
	// telegramChannel.Disconnect(ctx)
}
