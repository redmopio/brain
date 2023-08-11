package self

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

func (brain *BrainEngine) prepareMessagesForConversation(user models.User, lastMessages []models.Message, inputMessage models.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: user.Context.String,
	})

	userName := user.UserName.String
	userName = strings.ReplaceAll(userName, " ", "_")

	for _, msg := range lastMessages {
		messageName := userName
		role := openai.ChatMessageRoleUser

		if msg.Role.String == openai.ChatMessageRoleAssistant {
			role = openai.ChatMessageRoleAssistant
			messageName = ""
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Name:    messageName,
			Role:    role,
			Content: msg.Content.String,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Name:    userName,
		Role:    inputMessage.Role.String,
		Content: inputMessage.Content.String,
	})

	return messages
}

func (brain *BrainEngine) ProcessMessageResponse(ctx context.Context, user models.User, lastMessages []models.Message, inputMessage models.Message) (*models.Message, error) {
	messages := brain.prepareMessagesForConversation(user, lastMessages, inputMessage)

	fmt.Printf("Total messages: %d\n", len(messages))

	for _, msg := range messages {
		fmt.Printf("\tMessage: [%s] %s\n", msg.Role, firstN(msg.Content, 100))
	}

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
			UUID:  inputMessage.ID,
			Valid: true,
		},
	}

	return &responseMessage, nil
}
