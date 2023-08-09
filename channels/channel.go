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
	// GenerateResponse(ctx context.Context, senderID string, message string) (string, error)
	Connect(ctx context.Context)
	Disconnect(ctx context.Context)
}
