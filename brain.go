package brain

import (
	"context"
	"os"

	"github.com/minskylab/brain/config"
	"github.com/minskylab/brain/models"
	"github.com/minskylab/brain/system"
	"github.com/pkg/errors"
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
}

type Brain struct {
	System   *system.SystemEngine
	Agents   map[string]*Agent
	Channels map[string]Channel
}

type Message struct {
	*models.Message
	User *models.User
	// Group *models.Group
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

func (b *Brain) Interact(ctx context.Context, agentName string, messages []Message) (*Message, error) {
	// lastMessage, _ := b.lastMessage(ctx, messages)
	// b.System.GenerateConversationResponse(ctx)

	return nil, nil
}

func NewBrain(ctx context.Context, config *config.Config) (*Brain, error) {
	system, err := system.NewSystemEngine(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	registeredAgents, err := system.DatabaseClient.GetAllAgents(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	agents := map[string]Agent{}
	for _, agent := range registeredAgents {
		agents[agent.Name] = Agent{
			ID:           agent.ID.String(),
			Name:         agent.Name,
			Constitution: agent.Constitution,
		}
	}

	return &Brain{
		System:   system,
		Agents:   map[string]*Agent{},
		Channels: map[string]Channel{},
	}, nil
}

func (b *Brain) RegisterChannel(channel Channel) {
	b.Channels[string(channel.Name())] = channel
}

func (b *Brain) RegisterBeforeResponseFunction(agentName string, f AgentBeforeResponseFunction) {
	b.Agents[agentName].beforeResponseHandlers = append(b.Agents[agentName].beforeResponseHandlers, f)
}

func (b *Brain) RegisterAfterResponseFunction(agentName string, f AgentAfterResponseFunction) {
	b.Agents[agentName].afterResponseHandlers = append(b.Agents[agentName].afterResponseHandlers, f)
}

func (b *Brain) RegisterAgent(ctx context.Context, name string, constitution string) (models.Agent, error) {
	return b.System.DatabaseClient.UpsertAgent(ctx, models.UpsertAgentParams{
		Name:         name,
		Constitution: constitution,
	})
}

func (b *Brain) RegisterAgentFromFile(ctx context.Context, name string, constitutionPath string) (models.Agent, error) {
	constitution, err := os.ReadFile(constitutionPath)
	if err != nil {
		return models.Agent{}, errors.WithStack(err)
	}

	return b.RegisterAgent(ctx, name, string(constitution))
}
