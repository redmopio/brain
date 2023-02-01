package self

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"go.mau.fi/whatsmeow/types"
)

func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			// log.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			// log.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}

func (brain *BrainEngine) ResponseWhatsAppMessage(ctx context.Context, sender types.JID, message string) (string, error) {
	conversation, err := brain.DatabaseClient.GetConversationByJid(ctx, sql.NullString{
		String: sender.String(),
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	response, err := brain.ProcessMessageResponse(&conversation, message)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// message := "Hello World"
	// recipient, _ := parseJID("")
	// msg := &waProto.Message{Conversation: proto.String(strings.Join([]string{response}, " "))}
	// resp, err := brain.WhatsAppClient.SendMessage(context.Background(), sender, msg)
	// if err != nil {
	// 	return "", errors.WithStack(err)
	// }

	return response, nil
}
