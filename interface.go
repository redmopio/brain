package brain

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"

	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
)

func (b *Brain) generateConversationResponse(ctx context.Context, agent *Agent, channelName string, sender string, message string) (string, error) {
	var user models.User
	var err error

	if channelName == string("channels.WhatsAppChannelName") {
		user, err = b.System.DatabaseClient.GetUserByJID(ctx, sql.NullString{
			String: sender,
			Valid:  true,
		})
		if err != nil {
			return "", errors.WithStack(err)
		}
	} else if channelName == string("channels.TelegramChannel") {
		user, err = b.System.DatabaseClient.GetUserByTelegramID(ctx, sql.NullString{
			String: sender,
			Valid:  true,
		})
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	fmt.Println("User: ", user.UserName.String)

	lastMessages, err := b.System.DatabaseClient.GetMessagesByUserID(ctx, models.GetMessagesByUserIDParams{
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Limit:  20,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	userMessage := buildUserMessage(user.ID, message, lastMessages)
	userMessage, err = b.storeMessage(ctx, &userMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// brainMessage, agent, err := b.processMessageResponse(ctx, &user, lastMessages, userMessage)
	// if err != nil {
	// 	return "", errors.WithStack(err)
	// }

	brainMessage, err := agent.Interact(ctx, UserMessages(
		&user,
		message,
	))
	if err != nil {
		return "", errors.WithStack(err)
	}

	chatbotMessage := buildChatbotMessage(user.ID, brainMessage.Content.String, userMessage.ID, agent.Agent)
	responseMessage, err := b.storeMessage(ctx, &chatbotMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return responseMessage.Content.String, nil
}

func (b *Brain) storeMessage(ctx context.Context, message *models.Message) (models.Message, error) {
	storedMessage, err := b.System.DatabaseClient.CreateMessage(ctx, models.CreateMessageParams{
		UserID:   message.UserID,
		Role:     message.Role,
		Content:  message.Content,
		ParentID: message.ParentID,
		AgentID:  message.AgentID,
	})
	if err != nil {
		return models.Message{}, errors.WithStack(err)
	}

	fmt.Printf("Stored message [%s][%s]: %s\n", storedMessage.Role.String, storedMessage.UserID.UUID.String(), firstN(storedMessage.Content.String, 100))

	return storedMessage, nil
}

func buildUserMessage(userId uuid.UUID, messageContent string, lastMessages []models.GetMessagesByUserIDRow) models.Message {
	userMessage := models.Message{
		Role:    sql.NullString{String: openai.ChatMessageRoleUser, Valid: true},
		UserID:  uuid.NullUUID{UUID: userId, Valid: true},
		Content: sql.NullString{String: messageContent, Valid: true},
	}

	if len(lastMessages) > 0 {
		userMessage.ParentID = uuid.NullUUID{UUID: lastMessages[len(lastMessages)-1].ID, Valid: true}
	}

	return userMessage
}

func buildChatbotMessage(userId uuid.UUID, messageContent string, parentId uuid.UUID, agent *models.Agent) models.Message {
	responseMessage := models.Message{
		Role:     sql.NullString{String: openai.ChatMessageRoleAssistant, Valid: true},
		UserID:   uuid.NullUUID{UUID: userId, Valid: true},
		Content:  sql.NullString{String: messageContent, Valid: true},
		ParentID: uuid.NullUUID{UUID: parentId, Valid: true},
	}

	if agent != nil {
		responseMessage.AgentID = uuid.NullUUID{UUID: agent.ID, Valid: true}
	}

	return responseMessage
}
