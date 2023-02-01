package channels

import (
	"context"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppResponseFunc func(ctx context.Context, sender types.JID, message string) (string, error)

type WhatsAppConnector struct {
	// Brain        *self.BrainEngine
	DatabaseName string
	Response     WhatsAppResponseFunc
	// Client       *whatsmeow.Client
}

func NewWhatsAppConnector(databaseName string, response WhatsAppResponseFunc) *WhatsAppConnector {
	return &WhatsAppConnector{
		DatabaseName: databaseName,

		Response: response,
		// Brain:        brain,
	}
}

func (w *WhatsAppConnector) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		ctx := context.Background()
		sender := v.Info.Sender
		message := v.Message.GetConversation()

		fmt.Println("Received a message:", message)
		fmt.Println("Sender:", sender)

		_, err := w.Response(ctx, sender, message)
		if err != nil {
			panic(err)
		}
	}
}

func (w *WhatsAppConnector) Connect() *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	storeName := fmt.Sprintf("file:%s?_foreign_keys=on", w.DatabaseName)

	container, err := sqlstore.New("sqlite3", storeName, dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)

	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(w.eventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	return client
}
