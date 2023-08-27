package brain

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/config"
	"github.com/minskylab/brain/models"
	"github.com/minskylab/brain/system"
	"github.com/pkg/errors"
	"go.mau.fi/whatsmeow/types"
)

type (
	AgentBeforeResponseFunction func(ctx context.Context, agent *Agent, messages []Message) (*Message, error)
	AgentAfterResponseFunction  func(ctx context.Context, agent *Agent, messages []Message, toResponse Message) (*Message, error)
)

type Agent struct {
	ID           string
	Name         string
	Constitution string

	beforeResponseHandlers []AgentBeforeResponseFunction
	afterResponseHandlers  []AgentAfterResponseFunction

	brain *Brain
}

type Brain struct {
	System *system.SystemEngine
	Agents map[string]*Agent
	// Channels map[string]Channel
}

type Message struct {
	*models.Message
	User *models.User
	// Group *models.Group
}

type BrainBuilder struct {
	config *config.Config
	// channels []ChannelName
}

func NewMessages(messages ...models.Message) []Message {
	return []Message{}
}

// func (b *Brain) lastMessage(ctx context.Context, messages []Message) (*Message, error) {
// 	if len(messages) == 0 {
// 		return nil, nil
// 	}

// 	return &messages[len(messages)-1], nil
// }

func (b *Brain) Interact(ctx context.Context, agent *Agent, messages []Message) (*Message, error) {
	// lastMessage, _ := b.lastMessage(ctx, messages)
	// b.System.GenerateConversationResponse(ctx)

	return nil, nil
}

func NewBrainBuilder(config *config.Config) *BrainBuilder {
	return &BrainBuilder{
		config: config,
	}
}

// func (bb *BrainBuilder) WithChannel(channel ChannelName) *BrainBuilder {
// 	bb.channels = append(bb.channels, channel)

// 	return bb
// }

func (bb *BrainBuilder) Build(ctx context.Context) (*Brain, error) {
	system, err := system.NewSystemEngine(bb.config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Brain{
		System: system,
		Agents: map[string]*Agent{},
		// Channels: map[string]Channel{},
	}, nil
}

func (b *Brain) Run(ctx context.Context) error {
	whatsAppChannel := channels.NewWhatsAppConnector(b.System.Config, func(ctx context.Context, sender types.JID, message string) (string, error) {
		return b.System.GenerateConversationResponse(ctx, string(channels.WhatsAppChannelName), sender.String(), message)
	})

	if !b.System.Config.WhatsAppDisabled {
		go whatsAppChannel.Connect(ctx)
	}

	var telegramChannel *channels.TelegramConnector

	if b.System.Config.TelegramAPIKey != "" {
		telegramChannel = channels.NewTelegramConnector(b.System.Config, func(ctx context.Context, sender string, message string) (string, error) {
			return b.System.GenerateConversationResponse(ctx, string(channels.TelegramChannelName), sender, message)
		})

		telegramChannel.Connect(ctx)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	whatsAppChannel.Disconnect(ctx)

	if telegramChannel != nil {
		telegramChannel.Disconnect(ctx)
	}

	return nil
}

func (a *Agent) Interact(ctx context.Context, messages []Message) (*Message, error) {
	return nil, nil
}
