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

	fmt.Println("User:", user)

	// conversation, err := brain.DatabaseClient.GetConversationByJid(ctx, sql.NullString{
	// 	String: sender,
	// 	Valid:  true,
	// })
	// if err != nil {
	// 	return "", errors.WithStack(err)
	// }

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

	response, err := brain.ProcessMessageResponse(ctx, &user, lastMessages, &inputMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// brain.DatabaseClient.UpdateConversationBuffer(ctx, models.UpdateConversationBufferParams{
	// 	ID: conversation.ID,
	// 	ConversationBuffer: sql.NullString{
	// 		String: response.NewBuffer,
	// 		Valid:  true,
	// 	},
	// })

	return response.Content.String, nil
}
