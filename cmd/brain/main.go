package main

import (
	"context"
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

	whatsAppClient := channels.NewWhatsAppConnector(config, func(ctx context.Context, sender types.JID, message string) (string, error) {
		response, err := brain.GenerateConversationResponse(ctx, sender, message)
		if err != nil {
			return "", errors.WithStack(err)
		}

		return strings.ReplaceAll(response, "<|im_end|>", ""), nil
	}).Connect()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	whatsAppClient.Disconnect()
}
