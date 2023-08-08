package llm

import (
	"github.com/minskylab/brain/config"
	openai "github.com/sashabaranov/go-openai"
)

type LLMEngine struct {
	Client *openai.Client
}

func NewLLMEngine(config *config.Config) (*LLMEngine, error) {
	return &LLMEngine{
		Client: openai.NewClient(config.OpenAIKey),
	}, nil
}
