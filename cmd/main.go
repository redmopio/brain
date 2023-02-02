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

	fmt.Println(config)

	brain, err := self.NewBrainEngine(config)
	if err != nil {
		panic(err)
	}

	// jid := "51986253867@s.whatsapp.net"
	// ctx := context.Background()

	// conversation, err := brain.DatabaseClient.GetConversationByJid(ctx, sql.NullString{
	// 	String: jid,
	// 	Valid:  true,
	// })
	// if err != nil {
	// 	panic(errors.WithStack(err))
	// }

	// fmt.Println(conversation)

	client := channels.NewWhatsAppConnector("examplestore.db", func(ctx context.Context, sender types.JID, message string) (string, error) {
		return brain.ResponseWhatsAppMessage(ctx, sender, message)
	}).Connect()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
