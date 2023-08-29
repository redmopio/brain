package brain

import (
	"strings"

	"github.com/minskylab/brain/models"
	"github.com/sashabaranov/go-openai"
)

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/minskylab/brain/config"
// 	"github.com/minskylab/brain/llm"
// 	"github.com/minskylab/brain/models"
// 	"github.com/pkg/errors"
// 	"github.com/xo/dburl"
// )

// type SystemEngine struct {
// 	DatabaseClient *models.Queries
// 	LLMEngine      *llm.LLMEngine
// 	Config         *config.Config
// }

// func NewSystemEngine(config *config.Config) (*SystemEngine, error) {
// 	url, err := dburl.Parse(config.DatabaseURL)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	db, err := dburl.Open(config.DatabaseURL)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	prelude := ""

// 	if url.Driver == "sqlite3" || url.Driver == "sqlite" {
// 		data, err := os.ReadFile("database/schema-sqlite.sql")
// 		if err != nil {
// 			return nil, errors.WithStack(err)
// 		}

// 		prelude = string(data)
// 	} else {
// 		data, err := os.ReadFile("database/schema.sql")
// 		if err != nil {
// 			return nil, errors.WithStack(err)
// 		}

// 		prelude = string(data)
// 	}

// 	if _, err = db.Exec(prelude); err != nil {
// 		if err != nil && !strings.Contains(err.Error(), "already exists") {
// 			return nil, errors.WithStack(err)
// 		}

// 		fmt.Println("Database already exists, skipping prelude")
// 	}

// 	llmEngine, err := llm.NewLLMEngine(config)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	client := models.New(db)

// 	return &SystemEngine{
// 		DatabaseClient: client,
// 		LLMEngine:      llmEngine,
// 		Config:         config,
// 		// Name:           config.Name,
// 		// HasuraToken:    config.HasuraToken,
// 	}, nil
// }

func prepareMessagesForConversation(agent *models.Agent, messageContent string, user *models.User, lastMessages []Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: agent.Constitution,
	})

	for i := len(lastMessages) - 1; i >= 0; i-- {
		msg := lastMessages[i]
		messageName := msg.User.UserName.String
		messageName = strings.ReplaceAll(messageName, " ", "_")
		role := openai.ChatMessageRoleUser

		if msg.Role.String == openai.ChatMessageRoleAssistant {
			role = openai.ChatMessageRoleAssistant
			messageName = ""
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Name:    messageName,
			Role:    role,
			Content: msg.Content.String,
		})
	}

	// lastMessageToAppend := openai.ChatCompletionMessage{
	// 	Role:    openai.ChatMessageRoleUser,
	// 	Content: messageContent,
	// }

	// if user != nil {
	// userName := user.UserName.String
	// userName = strings.ReplaceAll(userName, " ", "_")

	// lastMessageToAppend.Name = userName
	// }

	// messages = append(messages, lastMessageToAppend)

	return messages
}
