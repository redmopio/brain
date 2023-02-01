package self

import (
	"fmt"

	"github.com/minskylab/brain/models"
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
	return prompt, nil
}
