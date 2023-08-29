package brain

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/minskylab/brain/models"
	"github.com/sashabaranov/go-openai"
)

func (b *Brain) NewUserMessageWithUser(ctx context.Context, user *models.User, message string) (models.Message, error) {
	return b.System.DatabaseClient.CreateMessage(ctx, models.CreateMessageParams{
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		Role: sql.NullString{
			String: openai.ChatMessageRoleUser,
			Valid:  true,
		},
		Content: sql.NullString{
			String: message,
			Valid:  true,
		},
	})
}

func (b *Brain) NewAssistantMessage(ctx context.Context, message string) (models.Message, error) {
	return b.System.DatabaseClient.CreateMessage(ctx, models.CreateMessageParams{
		Role: sql.NullString{
			String: openai.ChatMessageRoleAssistant,
			Valid:  true,
		},
		Content: sql.NullString{
			String: message,
			Valid:  true,
		},
	})
}

func (b *Brain) NewUserMessageWithParent(ctx context.Context, user *models.User, parentID uuid.UUID, message string) (models.Message, error) {
	return b.System.DatabaseClient.CreateMessage(ctx, models.CreateMessageParams{
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		Role: sql.NullString{
			String: openai.ChatMessageRoleUser,
			Valid:  true,
		},
		Content: sql.NullString{
			String: message,
			Valid:  true,
		},
		ParentID: uuid.NullUUID{
			UUID:  parentID,
			Valid: true,
		},
	})
}
