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

func (brain *SystemEngine) HandleGroup(ctx context.Context, channel channels.ChannelName, groupId string, groupName string) (string, error) {
	dbGroup, err := brain.DatabaseClient.GetGroupByID(ctx, groupId)
	if err != nil {
		// scope if err is due to group not found
		if !errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("Error getting group %s: %s\n", groupId, err.Error())
			return "", errors.WithStack(err)
		}

		fmt.Printf("Group %s not found, creating it...\n", groupId)

		connector, err := brain.DatabaseClient.GetConnectorByName(ctx, channel.String())
		if err != nil {
			fmt.Printf("Error getting connector from channel '%s': %s\n", channel.String(), err.Error())
			return "", errors.WithStack(err)
		}

		createdGroup, err := brain.DatabaseClient.CreateGroup(ctx, models.CreateGroupParams{
			ID:          groupId,
			Name:        sql.NullString{String: groupName, Valid: true},
			Description: sql.NullString{String: groupName, Valid: true},
			ConnectorID: sql.NullString{String: connector.ID.String(), Valid: true},
		})
		if err != nil {
			fmt.Printf("Error creating group %s: %s\n", groupId, err.Error())
			return "", errors.WithStack(err)
		}

		fmt.Printf("Group %s created!\n", groupId)

		dbGroup = createdGroup
	}

	return dbGroup.Name.String, nil
}

func (brain *SystemEngine) HandleGroupSender(ctx context.Context, channel channels.ChannelName, groupId string, sender string) (string, error) {
	user, err := brain.getUserFromSender(ctx, channel, sender)
	if err != nil {
		return "", errors.WithStack(err)
	}

	userId := user.ID.String()

	groupUsers, err := brain.DatabaseClient.GetUsersFromGroup(ctx, groupId)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// check if user is already in group
	for _, userFromGroup := range groupUsers {
		if userFromGroup.ID.String() == userId {
			return userId, nil
		}
	}

	_, err = brain.DatabaseClient.AddUserToGroup(ctx, models.AddUserToGroupParams{
		UserID:  user.ID,
		GroupID: groupId,
	})

	if err != nil {
		return "", errors.WithStack(err)
	}

	return userId, nil
}

func (brain *SystemEngine) GenerateConversationResponse(ctx context.Context, channel channels.ChannelName, sender string, message string) (string, error) {
	user, err := brain.getUserFromSender(ctx, channel, sender)
	if err != nil {
		return "", errors.WithStack(err)
	}

	fmt.Println("User: ", user.UserName.String)

	lastMessages, err := brain.DatabaseClient.GetMessagesByUserID(ctx, models.GetMessagesByUserIDParams{
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Limit:  20,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	userMessage := buildUserMessage(user.ID, message, lastMessages)
	userMessage, err = brain.storeMessage(ctx, &userMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	brainMessage, agent, err := brain.processMessageResponse(ctx, user, lastMessages, userMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	chatbotMessage := buildChatbotMessage(user.ID, brainMessage, userMessage.ID, agent)
	responseMessage, err := brain.storeMessage(ctx, &chatbotMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return responseMessage.Content.String, nil
}

func (brain *SystemEngine) getUserFromSender(ctx context.Context, channel channels.ChannelName, sender string) (*models.User, error) {
	var user models.User
	var err error

	if channel == channels.WhatsAppChannel {
		user, err = brain.DatabaseClient.GetUserByJID(ctx, sql.NullString{
			String: sender,
			Valid:  true,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else if channel == channels.TelegramChannel {
		user, err = brain.DatabaseClient.GetUserByTelegramID(ctx, sql.NullString{
			String: sender,
			Valid:  true,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &user, nil
}
