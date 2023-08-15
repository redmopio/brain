package channels

import (
	"context"
)

type ChannelType int

const (
	TelegramChannel ChannelType = iota
	WhatsAppChannel
)

type Channel interface {
	Name() string
	SendMessage(ctx context.Context, senderID string, message string) (string, error)
	Connect(ctx context.Context)
	Disconnect(ctx context.Context)
}
