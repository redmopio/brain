package main

import (
	"context"
	"fmt"

	"github.com/minskylab/brain"
	"github.com/minskylab/brain/config"
)

func main() {
	config, err := config.NewLoadedConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	minskyBrain, err := brain.NewBrainBuilder(config).
		Build(ctx)
	if err != nil {
		panic(err)
	}

	agentWriteParse, err := minskyBrain.NewAgentBuilder("agent_write_parse_data").
		WithConstitutionFromFile("agents/agent_write_parse_data.md").
		WithBeforeResponseFunction(func(ctx context.Context, agent *brain.Agent, messages []brain.Message) (*brain.Message, error) {
			fmt.Println("default agent: before response")
			return nil, nil
		}).
		WithAfterResponseFunction(func(ctx context.Context, agent *brain.Agent, messages []brain.Message, toResponse brain.Message) (*brain.Message, error) {
			fmt.Println("default agent: after response")
			return nil, nil
		}).
		Build(ctx)
	if err != nil {
		panic(err)
	}

	agentWriteStoreData, err := minskyBrain.NewAgentBuilder("agent_write_store_data").
		WithConstitutionFromFile("agents/agent_write_store_data.md").
		WithBeforeResponseFunction(func(ctx context.Context, agent *brain.Agent, messages []brain.Message) (*brain.Message, error) {
			fmt.Println("default agent: before response")
			return nil, nil
		}).
		WithAfterResponseFunction(func(ctx context.Context, agent *brain.Agent, messages []brain.Message, toResponse brain.Message) (*brain.Message, error) {
			fmt.Println("default agent: after response")
			return nil, nil
		}).
		Build(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(agentWriteParse)
	fmt.Println(agentWriteStoreData)

	// user, err := b.ObtainUserByChannelAndID(ctx, string(channels.WhatsAppChannelName), "1234567890")
	// if err != nil {
	// 	panic(err)
	// }

	// message, err := b.NewUserMessage(ctx, user, "Hola, soy un mensaje de prueba")
	// if err != nil {
	// 	panic(err)
	// }

	// b.Interact(ctx, agent, brain.NewMessages(message))
	fmt.Println(minskyBrain)
}
