package self

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
)

func (brain *BrainEngine) GenerateConversationResponse(ctx context.Context, channel channels.ChannelType, sender string, message string) (string, error) {
	var user models.User
	var err error

	if channel == channels.WhatsAppChannel {
		user, err = brain.DatabaseClient.GetUserByJID(ctx, sql.NullString{
			String: sender,
			Valid:  true,
		})
		if err != nil {
			return "", errors.WithStack(err)
		}
	} else if channel == channels.TelegramChannel {
		user, err = brain.DatabaseClient.GetUserByTelegramID(ctx, sql.NullString{
			String: sender,
			Valid:  true,
		})
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	fmt.Println("User: ", user.UserName.String)

	lastMessages, err := brain.DatabaseClient.GetMessagesByUserID(ctx, uuid.NullUUID{
		UUID:  user.ID,
		Valid: true,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	userMessage := buildUserMessage(user.ID, message, lastMessages)
	userMessage, err = brain.storeMessage(ctx, &userMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	openAiResponse, err := brain.ProcessMessageResponse(ctx, user, lastMessages, userMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	openAiMessage := buildChatbotMessage(user.ID, openAiResponse.Choices[0].Message.Content, userMessage.ID)
	responseMessage, err := brain.storeMessage(ctx, &openAiMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return responseMessage.Content.String, nil
}
