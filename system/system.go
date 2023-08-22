package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/minskylab/brain/config"
	"github.com/minskylab/brain/llm"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
	"github.com/xo/dburl"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type SystemEngine struct {
	DatabaseClient *models.Queries
	LLMEngine      *llm.LLMEngine
	Config         *config.Config
}

func NewSystemEngine(config *config.Config) (*SystemEngine, error) {
	url, err := dburl.Parse(config.DatabaseURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db, err := dburl.Open(config.DatabaseURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	prelude := ""

	if url.Driver == "sqlite3" || url.Driver == "sqlite" {
		data, err := os.ReadFile("database/schema-sqlite.sql")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		prelude = string(data)
	} else {
		data, err := os.ReadFile("database/schema.sql")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		prelude = string(data)
	}

	if _, err = db.Exec(prelude); err != nil {
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return nil, errors.WithStack(err)
		}

		fmt.Println("Database already exists, skipping prelude")
	}

	llmEngine, err := llm.NewLLMEngine(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client := models.New(db)

	return &SystemEngine{
		DatabaseClient: client,
		LLMEngine:      llmEngine,
		Config:         config,
	}, nil
}
