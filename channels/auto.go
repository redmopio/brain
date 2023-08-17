package channels

import (
	"context"

	"github.com/minskylab/brain/config"
)

func InferChannelsFromConfig(ctx context.Context, config *config.Config) []Channel {
	channels := []Channel{}

	// whatsAppChannel := NewWhatsAppConnector(config, func(ctx context.Context, sender types.JID, message string) (string, error) {
	// 	return brain.GenerateConversationResponse(ctx, WhatsAppChannel, sender.String(), message)
	// })

	return channels
}
