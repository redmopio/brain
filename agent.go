package brain

import "context"

type AgentBuilder struct {
	name             string
	constitutionFile *string
	constitution     *string

	beforeResponseHandlers []AgentBeforeResponseFunction
	afterResponseHandlers  []AgentAfterResponseFunction

	brain *Brain
}

func (b *Brain) NewAgentBuilder(name string) *AgentBuilder {
	return &AgentBuilder{
		name:  name,
		brain: b,
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

func (ab *AgentBuilder) WithAfterResponseFunction(callback AgentAfterResponseFunction) *AgentBuilder {
	ab.afterResponseHandlers = append(ab.afterResponseHandlers, callback)

	return ab
}

func (ab *AgentBuilder) WithBeforeResponseFunction(callback AgentBeforeResponseFunction) *AgentBuilder {
	ab.beforeResponseHandlers = append(ab.beforeResponseHandlers, callback)

	return ab
}

func (ab *AgentBuilder) Build(ctx context.Context) (*Agent, error) {
	return nil, nil
}
