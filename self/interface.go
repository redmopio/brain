package self

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
)

func (brain *BrainEngine) GenerateConversationResponse(ctx context.Context, sender string, message string) (string, error) {
	fmt.Println("Message from", sender, ":\n", message)

	conversation, err := brain.DatabaseClient.GetConversationByJid(ctx, sql.NullString{
		String: sender,
		Valid:  true,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	response, err := brain.ProcessMessageResponse(&conversation, message)
	if err != nil {
		return "", errors.WithStack(err)
	}

	brain.DatabaseClient.UpdateConversationBuffer(ctx, models.UpdateConversationBufferParams{
		ID: conversation.ID,
		ConversationBuffer: sql.NullString{
			String: response.NewBuffer,
			Valid:  true,
		},
	})

	return response.PredictedResponse, nil
}
