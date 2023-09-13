package channels

import (
	"context"
)

type ChannelName string

const (
	TelegramChannel ChannelName = "telegram"
	WhatsAppChannel ChannelName = "whatsapp"
)

func (c ChannelName) String() string {
	return string(c)
}

type Channel interface {
	Name() ChannelName
	SendMessage(ctx context.Context, groupId string, senderID string, message string) (string, error)
	Connect(ctx context.Context)
	Disconnect(ctx context.Context)
}
