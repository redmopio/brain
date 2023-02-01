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

	res, err := brain.LLMEngine.Client.Completion(ctx, gpt3.CompletionRequest{
		Prompt: []string{prompt},
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return strings.TrimSpace(res.Choices[0].Text), nil
}
