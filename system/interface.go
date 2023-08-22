package system

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
)

func (system *SystemEngine) GenerateConversationResponse(ctx context.Context, channelName string, sender string, message string) (string, error) {
	var user models.User
	var err error

	if channelName == string("channels.WhatsAppChannelName") {
		user, err = system.DatabaseClient.GetUserByJID(ctx, sql.NullString{
			String: sender,
			Valid:  true,
		})
		if err != nil {
			return "", errors.WithStack(err)
		}
	} else if channelName == string("channels.TelegramChannel") {
		user, err = system.DatabaseClient.GetUserByTelegramID(ctx, sql.NullString{
			String: sender,
			Valid:  true,
		})
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	fmt.Println("User: ", user.UserName.String)

	lastMessages, err := system.DatabaseClient.GetMessagesByUserID(ctx, models.GetMessagesByUserIDParams{
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Limit:  20,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	userMessage := buildUserMessage(user.ID, message, lastMessages)
	userMessage, err = system.storeMessage(ctx, &userMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	brainMessage, agent, err := system.processMessageResponse(ctx, &user, lastMessages, userMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	chatbotMessage := buildChatbotMessage(user.ID, brainMessage, userMessage.ID, agent)
	responseMessage, err := system.storeMessage(ctx, &chatbotMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return responseMessage.Content.String, nil
}
