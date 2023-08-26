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

	b, err := brain.NewBrainBuilder(config).
		WithChannel(channels.TelegramChannelName).
		WithChannel(channels.WhatsAppChannelName).
		Build(ctx)
	if err != nil {
		panic(err)
	}

	agent, err := b.NewAgentBuilder("default").
		WithConstitutionFromFile("agents/default.md").
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

	fmt.Println(agent)

	user, err := b.ObtainUserByChannelAndID(ctx, string(channels.WhatsAppChannelName), "1234567890")
	if err != nil {
		panic(err)
	}

	message, err := b.NewUserMessage(ctx, user, "Hola, soy un mensaje de prueba")
	if err != nil {
		panic(err)
	}

	b.Interact(ctx, agent, brain.NewMessages(message))
	fmt.Println(b)
}
