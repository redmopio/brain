package self

import (
	_ "github.com/lib/pq"
	"github.com/minskylab/brain/config"
	"github.com/minskylab/brain/llm"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

type BrainEngine struct {
	DatabaseClient *models.Queries
	LLMEngine      *llm.LLMEngine
	Name           string
}

func NewBrainEngine(config *config.Config) (*BrainEngine, error) {
	db, err := dburl.Open(config.DatabaseURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// data, err := os.ReadFile("database/schema.sql")
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// dataString := string(data)
	// fmt.Println(dataString)

	// if _, err = db.Exec(dataString); err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	llmEngine, err := llm.NewLLMEngine(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client := models.New(db)

	return &BrainEngine{
		DatabaseClient: client,
		LLMEngine:      llmEngine,
		Name:           config.Name,
	}, nil
}
