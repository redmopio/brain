package self

import (
	"github.com/minskylab/brain/config"
	"github.com/minskylab/brain/llm"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
	"go.mau.fi/whatsmeow"
)

type BrainEngine struct {
	DatabaseClient *models.Queries
	LLMEngine      *llm.LLMEngine
	WhatsAppClient *whatsmeow.Client
	Name           string
}

func NewBrainEngine(config *config.Config) (*BrainEngine, error) {
	db, err := dburl.Open(config.DatabaseURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client := models.New(db)

	return &BrainEngine{
		DatabaseClient: client,
	}, nil
}
