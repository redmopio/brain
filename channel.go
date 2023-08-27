package brain

import (
	"context"
)

type ChannelName string

type Channel interface {
	Name() string
	// GetUserByID(ctx context.Context, brain *Brain, userID string) (*models.User, error)
	SendMessage(ctx context.Context, senderID string, message string) (string, error)
	Connect(ctx context.Context)
	Disconnect(ctx context.Context)
}

// system *self.SystemEngine,
