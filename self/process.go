package self

import (
	"context"
	"fmt"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
)

func (brain *BrainEngine) prepareConversationPrompt(conversation *models.Conversation, message string) string {
	return fmt.Sprintf("%s\n\nCurrent conversation:\n\n%s\n\n%s\n%s: %s\n%s:",
		conversation.Context.String,
		conversation.ConversationSummary.String,
		conversation.ConversationBuffer.String,
		conversation.UserName.String,
		message,
		brain.Name,
	)
}

func (brain *BrainEngine) Predict(conversation *models.Conversation, message string) (string, error) {
	prompt := brain.prepareConversationPrompt(conversation, message)

	ctx := context.Background()

	fmt.Println("PROMPT:", prompt)

	engineName := "text-davinci-003"

	stopSequences := []string{
		// "\n",
	}

	// res, err := brain.LLMEngine.Client.Completion()
	res, err := brain.LLMEngine.Client.CompletionWithEngine(ctx, engineName, gpt3.CompletionRequest{
		Prompt:          []string{prompt},
		Stop:            stopSequences,
		N:               gpt3.IntPtr(1),
		Temperature:     gpt3.Float32Ptr(0.9),
		TopP:            gpt3.Float32Ptr(1.0),
		PresencePenalty: 0.6,
		MaxTokens:       gpt3.IntPtr(1000),
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return strings.TrimSpace(res.Choices[0].Text), nil
}

type MessageResponse struct {
	PredictedResponse string
	NewBuffer         string
}

func (brain *BrainEngine) ProcessMessageResponse(conversation *models.Conversation, message string) (*MessageResponse, error) {
	predicted, err := brain.Predict(conversation, message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newBuffer := fmt.Sprintf("%s\n%s:%s\n%s:%s",
		conversation.ConversationBuffer.String,
		conversation.UserName.String,
		message,
		brain.Name,
		predicted,
	)

	maxLines := 5

	lines := strings.Split(newBuffer, "\n")

	if len(lines) > maxLines {
		newBuffer = strings.Join(lines[len(lines)-maxLines:], "\n")
	} else {
		newBuffer = strings.Join(lines, "\n")
	}

	return &MessageResponse{
		PredictedResponse: predicted,
		NewBuffer:         newBuffer,
	}, nil
}
