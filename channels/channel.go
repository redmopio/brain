package channels

import (
	"context"
)

type ChannelName string

const (
	TelegramChannel ChannelName = "telegram"
	WhatsAppChannel ChannelName = "whatsapp"
)

type Channel interface {
	Name() ChannelName
	SendMessage(ctx context.Context, senderID string, message string) (string, error)
	Connect(ctx context.Context)
	Disconnect(ctx context.Context)
}
