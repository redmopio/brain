package brain

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/config"
	"github.com/minskylab/brain/models"
	"github.com/minskylab/brain/system"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow/types"
)

type (
	AgentBeforeResponseFunction func(ctx context.Context, agent *Agent, messages []Message) (*Message, error)
	AgentAfterResponseFunction  func(ctx context.Context, agent *Agent, messages []Message, toResponse Message) (*Message, error)
)

type Agent struct {
	// ID           string
	// Name         string
	// Constitution string
	*models.Agent

	beforeResponseHandlers []AgentBeforeResponseFunction
	afterResponseHandlers  []AgentAfterResponseFunction

	brain *Brain
}

// func (a *Agent) Model() *models.Agent {
// 	return &models.Agent{
// 		// TODO: Fix it
// 		ID:           uuid.MustParse(a.Agent.ID.String()),
// 		Name:         a.Name,
// 		Constitution: a.Constitution,
// 	}
// }

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

func UserMessages(user *models.User, messages ...string) []Message {
	brainMessages := []Message{}

	for _, msg := range messages {
		brainMessages = append(brainMessages, Message{
			Message: &models.Message{
				Role: sql.NullString{
					String: openai.ChatMessageRoleUser,
					Valid:  true,
				},
				Content: sql.NullString{
					String: msg,
					Valid:  true,
				},
			},
			User: user,
		})
	}

	return brainMessages
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
	lastMessage, err := lastMessage(ctx, messages)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user := lastMessage.User

	openAIMessages := prepareMessagesForConversation(a.Agent, lastMessage.Content.String, user, messages)

	fmt.Printf("Total messages: %d\n", len(openAIMessages))

	for _, msg := range openAIMessages {
		fmt.Printf("\tMessage: [%s] %s\n", msg.Role, firstN(msg.Content, 100))
	}

	openAIResponse, err := a.brain.System.LLMEngine.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: openAIMessages,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	openAIMessageContent := openAIResponse.Choices[0].Message.Content

	messageResponse, err := a.brain.NewAssistantMessage(ctx, openAIMessageContent)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Message{
		Message: &messageResponse,
		// User:    user,
	}, nil
}

func lastMessage(ctx context.Context, messages []Message) (*Message, error) {
	if len(messages) == 0 {
		return nil, errors.New("no messages")
	}

	return &messages[len(messages)-1], nil
}

// func messagesToModelMessages(ctx context.Context, messages []Message) ([]*models.Message, error) {
// 	modelMessages := []*models.Message{}

// 	for _, msg := range messages {
// 		modelMessages = append(modelMessages, msg.Message)
// 	}

// 	return modelMessages, nil
// }

func firstN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func (b *Brain) ObtainUserByTelegramID(ctx context.Context, telegramID string) (models.User, error) {
	return b.System.DatabaseClient.GetUserByTelegramID(ctx, sql.NullString{
		String: telegramID,
		Valid:  true,
	})
}

func (b *Brain) ObtainUserByWhatsAppID(ctx context.Context, whatsAppID string) (models.User, error) {
	return b.System.DatabaseClient.GetUserByJID(ctx, sql.NullString{
		String: whatsAppID,
		Valid:  true,
	})
}
