package channels

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/minskylab/brain/config"
	"github.com/pkg/errors"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"google.golang.org/protobuf/proto"
)

type (
	WhatsAppResponseFunc func(ctx context.Context, sender types.JID, message string) (string, error)
	WhatsAppConnector    struct {
		DatabaseName string
		response     WhatsAppResponseFunc
		client       *whatsmeow.Client
	}
)

func NewWhatsAppConnector(config *config.Config, response WhatsAppResponseFunc) *WhatsAppConnector {
	return &WhatsAppConnector{
		DatabaseName: config.WhatsAppDatabaseName,
		response:     response,
	}
}

func (w *WhatsAppConnector) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if time.Since(v.Info.Timestamp).Minutes() > 1 { // filter out old messages
			return
		}

		ctx := context.Background()
		sender := v.Info.Sender
		message := v.Message.GetConversation()

		fmt.Println("Received a message:", message)
		fmt.Println("Sender:", sender)

		senderJID := sender.ToNonAD()
		response, err := w.response(ctx, senderJID, message)
		if err != nil {
			panic(err)
		}

		msg := &waProto.Message{Conversation: proto.String(strings.Join([]string{response}, " "))}

		resp, err := w.client.SendMessage(context.Background(), senderJID, msg)
		if err != nil {
			panic(err)
		}

		fmt.Printf("> Message sent: %s\n", resp.ID)
	}
}

func (w *WhatsAppConnector) Connect(ctx context.Context) {
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

	w.client = client // recursive?

	// return client
}

func (w *WhatsAppConnector) Disconnect(ctx context.Context) {
	w.client.Disconnect()
}

func (w *WhatsAppConnector) SendMessage(ctx context.Context, sender string, message string) (string, error) {
	senderJID, err := types.ParseJID(sender)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return w.response(ctx, senderJID, message)
}
