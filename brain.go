package brain

import (
	"context"

	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/models"
)

type Agent struct {
	ID           string
	Name         string
	Constitution string
}

type Channel struct {
	ID   string
	Name string
}

type Brain struct {
	Agents   []Agent
	Channels channels.Channel
}

type Message struct {
	*models.Message
	User *models.User
	// Group *models.Group
}

func (b *Brain) Interact(ctx context.Context, message []Message) (*Message, error) {
	return nil, nil
}
