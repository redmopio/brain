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
			agent.Interact(ctx, messages)
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

	fmt.Printf("%+v\n", minskyBrain)

	fmt.Printf("%+v\n", agentWriteParse)
	fmt.Printf("%+v\n", agentWriteStoreData)

	user, err := minskyBrain.ObtainUserByWhatsAppID(ctx, "51957821858@s.whatsapp.net")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", user)

	response, err := agentWriteParse.Interact(ctx, brain.UserMessages(
		&user,
		"Hola",
	))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", response)
	fmt.Printf("%+v\n", response.Content)
}
