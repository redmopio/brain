package self

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

func (brain *BrainEngine) prepareMessagesForConversation(user *models.User, lastMessages []models.Message, message *models.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: user.Context.String,
	})

	for _, msg := range lastMessages {
		role := openai.ChatMessageRoleUser
		name := user.UserName.String

		if msg.Role.String == openai.ChatMessageRoleAssistant {
			role = openai.ChatMessageRoleAssistant
			name = ""
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Name:    name,
			Content: msg.Content.String,
		})
	}

	return messages
}

func (brain *BrainEngine) ProcessMessageResponse(ctx context.Context, user *models.User, lastMessages []models.Message, message *models.Message) (*models.Message, error) {
	// predicted, err := brain.Predict(conversation, message)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// newBuffer := fmt.Sprintf("%s\n%s:%s\n%s:%s",
	// 	conversation.ConversationBuffer.String,
	// 	conversation.UserName.String,
	// 	message,
	// 	brain.Name,
	// 	predicted,
	// )

	// maxLines := 5

	// lines := strings.Split(newBuffer, "\n")

	// if len(lines) > maxLines {
	// 	newBuffer = strings.Join(lines[len(lines)-maxLines:], "\n")
	// } else {
	// 	newBuffer = strings.Join(lines, "\n")
	// }

	// return &MessageResponse{
	// 	PredictedResponse: predicted,
	// 	NewBuffer:         newBuffer,
	// }, nil

	messages := brain.prepareMessagesForConversation(user, lastMessages, message)

	response, err := brain.LLMEngine.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	responseMessage := models.Message{
		Role: sql.NullString{
			String: openai.ChatMessageRoleAssistant,
			Valid:  true,
		},
		Content: sql.NullString{
			String: response.Choices[0].Message.Content,
			Valid:  true,
		},
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		ParentID: uuid.NullUUID{
			UUID:  message.ID,
			Valid: true,
		},
	}

	return &responseMessage, nil
}
