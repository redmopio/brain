package brain

import (
	"context"
	"os"

	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
)

type AgentBuilder struct {
	name             string
	constitutionFile *string
	constitution     *string

	beforeResponseHandlers []AgentBeforeResponseFunction
	afterResponseHandlers  []AgentAfterResponseFunction[any]

	brain *Brain
}

type AgentBuilderWithPayload[P any] struct {
	*AgentBuilder
	afterResponseHandlers []AgentAfterResponseFunction[P]
}

func (b *Brain) NewAgentBuilder(name string) *AgentBuilder {
	return &AgentBuilder{
		name:  name,
		brain: b,
	}
}

func NewAgentBuilderWithPayload[P any](b *Brain, name string) *AgentBuilderWithPayload[P] {
	return &AgentBuilderWithPayload[P]{
		AgentBuilder: b.NewAgentBuilder(name),
	}
}

func (ab *AgentBuilder) WithConstitutionFromFile(path string) *AgentBuilder {
	ab.constitutionFile = &path

	return ab
}

func (ab *AgentBuilder) WithConstitution(constitution string) *AgentBuilder {
	ab.constitution = &constitution

	return ab
}

func (ab *AgentBuilderWithPayload[P]) WithAfterResponseFunction(callback AgentAfterResponseFunction[P]) *AgentBuilderWithPayload[P] {
	ab.afterResponseHandlers = append(ab.afterResponseHandlers, callback)

	return ab
}

func (ab *AgentBuilder) WithBeforeResponseFunction(callback AgentBeforeResponseFunction) *AgentBuilder {
	ab.beforeResponseHandlers = append(ab.beforeResponseHandlers, callback)

	return ab
}

func (ab *AgentBuilder) Build(ctx context.Context) (*Agent, error) {
	constitution := ""

	if ab.constitutionFile != nil {
		constitutionBytes, err := os.ReadFile(*ab.constitutionFile)
		if err != nil {
			return nil, errors.Wrap(err, "error reading constitution file")
		}

		constitution = string(constitutionBytes)
	} else if ab.constitution != nil {
		constitution = *ab.constitution
	}

	agentModel, err := ab.brain.System.DatabaseClient.UpsertAgent(ctx, models.UpsertAgentParams{
		Name:         ab.name,
		Constitution: constitution,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error registering agent")
	}

	return &Agent{
		Agent:                  &agentModel,
		beforeResponseHandlers: ab.beforeResponseHandlers,
		afterResponseHandlers:  ab.afterResponseHandlers,
		brain:                  ab.brain,
	}, nil
}
