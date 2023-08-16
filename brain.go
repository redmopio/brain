package brain

import (
	"context"

	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/models"
	"github.com/minskylab/brain/self"
)

type (
	AgentBeforeResponseFunction func(ctx context.Context, agentName string, messages []Message) (*Message, error)
	AgentAfterResponseFunction  func(ctx context.Context, agentName string, messages []Message, toResponse Message) (*Message, error)
)

type Agent struct {
	ID           string
	Name         string
	Constitution string

	BeforeResponse AgentBeforeResponseFunction
	AfterResponse  AgentAfterResponseFunction
}

type Channel struct {
	ID   string
	Name string
}

type Brain struct {
	System   *self.SystemEngine
	Agents   map[string]Agent
	Channels channels.Channel
}

type Message struct {
	*models.Message
	User *models.User
	// Group *models.Group
}

func (b *Brain) lastMessage(ctx context.Context, messages []Message) (*Message, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	return &messages[len(messages)-1], nil
}

func (b *Brain) Interact(ctx context.Context, agentName string, messages []Message) (*Message, error) {
	// lastMessage, _ := b.lastMessage(ctx, messages)
	// b.System.GenerateConversationResponse(ctx)

	return nil, nil
}
