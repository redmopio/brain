package llm

import (
	"github.com/PullRequestInc/go-gpt3"
	"github.com/minskylab/brain/config"
)

type LLMEngine struct {
	Client gpt3.Client
}

func NewLLMEngine(config *config.Config) (*LLMEngine, error) {
	return &LLMEngine{
		Client: gpt3.NewClient(config.OpenAIKey),
	}, nil
}
