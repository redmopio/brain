package self

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	// "google.golang.org/appengine/log"
	// "google.golang.org/appengine/log"
	"google.golang.org/protobuf/proto"
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

func (brain *BrainEngine) ResponseWhatsAppMessage() (string, error) {
	message := "Hello World"
	recipient, _ := parseJID("")
	msg := &waProto.Message{Conversation: proto.String(strings.Join([]string{message}, " "))}
	resp, err := brain.WhatsAppClient.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return resp.ID, nil
}
