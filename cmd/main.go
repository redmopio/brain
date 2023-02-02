package main

import (
	"context"
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

	// fmt.Println(config)

	brain, err := self.NewBrainEngine(config)
	if err != nil {
		panic(err)
	}

	client := channels.NewWhatsAppConnector("examplestore.db", func(ctx context.Context, sender types.JID, message string) (string, error) {
		return brain.GenerateConversationResponse(ctx, sender, message)
	}).Connect()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
