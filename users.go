package brain

import (
	"context"

	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
)

func (b *Brain) ObtainUserByChannelAndID(ctx context.Context, channelName string, userID string) (*models.User, error) {
	for name, channel := range b.Channels {
		if name == channelName {
			return channel.GetUserByID(ctx, b, userID)
		}
	}

	return nil, errors.New("channel not found")
}