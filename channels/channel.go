package channels

import (
	"context"
)

type Channel interface {
	GenerateResponse(ctx context.Context, senderID string, message string) (string, error)
}
