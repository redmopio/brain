package main

import (
	"context"
	"fmt"

	"github.com/minskylab/brain"
	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/config"
)

func main() {
	config, err := config.NewLoadedConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	b, err := brain.NewBrain(ctx, config)
	if err != nil {
		panic(err)
	}

	agent, err := b.RegisterAgentFromFile(ctx, "default", "agents/default.md")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", agent)

	b.RegisterBeforeResponseFunction("default", func(ctx context.Context, agent *brain.Agent, messages []brain.Message) (*brain.Message, error) {
		fmt.Println("default agent: before response")

		return nil, nil
	})

	b.RegisterAfterResponseFunction("default", func(ctx context.Context, agent *brain.Agent, messages []brain.Message, toResponse brain.Message) (*brain.Message, error) {
		fmt.Println("default agent: after response")

		return nil, nil
	})

	user, err := b.ObtainUserByChannelAndID(ctx, string(channels.WhatsAppChannelName), "1234567890")
	if err != nil {
		panic(err)
	}

	message, err := b.NewUserMessage(ctx, user, "Hola, soy un mensaje de prueba")
	if err != nil {
		panic(err)
	}

	b.Interact(ctx, "default", brain.NewMessages(message))
	fmt.Println(b)
}
