package self

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/minskylab/brain/channels"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
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

	inputMessage := models.Message{
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		Content: sql.NullString{
			String: message,
			Valid:  true,
		},
		Role: sql.NullString{
			String: openai.ChatMessageRoleUser,
			Valid:  true,
		},
	}

	if len(lastMessages) > 0 {
		inputMessage.ParentID = uuid.NullUUID{
			UUID:  lastMessages[len(lastMessages)-1].ID,
			Valid: true,
		}
	}

	inputMessage, err = brain.DatabaseClient.CreateMessage(ctx, models.CreateMessageParams{
		UserID:   inputMessage.UserID,
		Role:     inputMessage.Role,
		Content:  inputMessage.Content,
		ParentID: inputMessage.ParentID,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	fmt.Printf("Input message [%s]: %s\n", inputMessage.Role.String, inputMessage.Content.String)

	response, err := brain.ProcessMessageResponse(ctx, user, lastMessages, inputMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	responseMessage, err := brain.DatabaseClient.CreateMessage(ctx, models.CreateMessageParams{
		UserID:   response.UserID,
		Role:     response.Role,
		Content:  response.Content,
		ParentID: response.ParentID,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return responseMessage.Content.String, nil
}
