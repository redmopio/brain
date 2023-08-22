package channels

// func InferChannelsFromConfig(ctx context.Context, b *brain.Brain) []brain.Channel {
// 	totalChannels := []brain.Channel{}

// 	whatsAppChannel := NewWhatsAppConnector(b.System.Config, func(ctx context.Context, sender types.JID, message string) (string, error) {
// 		// return b.System.GenerateConversationResponse(ctx, WhatsAppChannelName, sender.String(), message)
// 		return "", nil
// 	})

// 	totalChannels = append(totalChannels, whatsAppChannel)

// 	// return channels

// 	return totalChannels
// }
