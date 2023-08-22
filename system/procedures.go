package system

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/minskylab/brain/models"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

const (
	AgentNameDefault             AgentName = "default"
	AgentNameAgentWriteParseData AgentName = "agent_write_parse_data"
	AgentNameAgentWriteStoreData AgentName = "agent_write_store_data"
	AgentNameAgentReadParseData  AgentName = "agent_read_parse_data"
	AgentNameAgentReadResponse   AgentName = "agent_read_response"
)

func firstN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

type AgentName string

func (system *SystemEngine) processMessageResponse(ctx context.Context, user *models.User, lastMessages []models.GetMessagesByUserIDRow, inputMessage models.Message) (string, *models.Agent, error) {
	openAiResponse, agent, err := system.processMessageWithOpenAI(ctx, user, lastMessages, inputMessage)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	responseWithDataKeywords := []string{"tus datos han sido registrados de manera exitosa"}
	openAiMessageContent := openAiResponse.Choices[0].Message.Content

	for _, keyword := range responseWithDataKeywords {
		if strings.Contains(strings.ToLower(openAiMessageContent), keyword) {
			parsedDataFromOpenAI, agent, err := system.processDataMessageWithOpenAI(ctx, user, openAiMessageContent)
			if err != nil {
				fmt.Printf("Error processing data message: %s\n", err.Error())
				return "", nil, errors.WithStack(err)
			}

			parsedDataFromOpenAIContent := parsedDataFromOpenAI.Choices[0].Message.Content

			if strings.Contains(parsedDataFromOpenAIContent, "Hubo un error al parsear la data") {
				return "Disculpa, creo que no ingresaste todos los datos. Revisa el ejemplo y vuelve a intentarlo.", nil, nil
			}

			// parse content from openAiResponse to DataStruct
			var data DataStruct
			err = json.Unmarshal([]byte(parsedDataFromOpenAI.Choices[0].Message.Content), &data)
			if err != nil {
				fmt.Printf("Error parsing data: %s\n", err.Error())
				return "Disculpa, creo que no ingresaste todos los datos. Revisa el ejemplo y vuelve a intentarlo.", nil, nil
			}

			fmt.Printf("\t -> Calling Hasura endpoint with data: %s\n", parsedDataFromOpenAIContent)
			_, err = system.callHasuraEndpoint(parsedDataFromOpenAIContent)
			if err != nil {
				return "Lo siento, hubo un error al guardar tus datos. Por favor, intentalo nuevamente en unos minutos.", nil, nil
			}

			return openAiMessageContent, agent, nil
		}
	}

	fmt.Printf("\t -> Processing normal message: %s\n", openAiMessageContent)

	return openAiMessageContent, agent, nil
}

func (system *SystemEngine) processDataMessageWithOpenAI(ctx context.Context, user *models.User, messageContent string) (*openai.ChatCompletionResponse, *models.Agent, error) {
	agentWriteParseData, err := system.getAgent(ctx, string(AgentNameAgentWriteParseData))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	messages := system.prepareMessagesForConversation(&agentWriteParseData, messageContent, nil, []models.GetMessagesByUserIDRow{})

	response, err := system.LLMEngine.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &response, &agentWriteParseData, nil
}

func (system *SystemEngine) processMessageWithOpenAI(ctx context.Context, user *models.User, lastMessages []models.GetMessagesByUserIDRow, inputMessage models.Message) (*openai.ChatCompletionResponse, *models.Agent, error) {
	agentWriteStoreData, err := system.getAgent(ctx, string(AgentNameAgentWriteStoreData))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	messages := system.prepareMessagesForConversation(&agentWriteStoreData, inputMessage.Content.String, user, lastMessages)

	fmt.Printf("Total messages: %d\n", len(messages))

	for _, msg := range messages {
		fmt.Printf("\tMessage: [%s] %s\n", msg.Role, firstN(msg.Content, 100))
	}

	response, err := system.LLMEngine.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &response, &agentWriteStoreData, nil
}

func (system *SystemEngine) prepareMessagesForConversation(agent *models.Agent, messageContent string, user *models.User, lastMessages []models.GetMessagesByUserIDRow) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: agent.Constitution,
	})

	for i := len(lastMessages) - 1; i >= 0; i-- {
		msg := lastMessages[i]
		messageName := msg.Username.String
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

	lastMessageToAppend := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: messageContent,
	}

	if user != nil {
		userName := user.UserName.String
		userName = strings.ReplaceAll(userName, " ", "_")

		lastMessageToAppend.Name = userName
	}

	messages = append(messages, lastMessageToAppend)

	return messages
}

func (system *SystemEngine) getAgent(ctx context.Context, agentName string) (models.Agent, error) {
	agent, err := system.DatabaseClient.GetAgentByName(ctx, agentName)
	if err != nil {
		return models.Agent{}, errors.WithStack(err)
	}

	fmt.Printf("Agent: %s\n", agent.Name)

	return agent, nil
}

func (system *SystemEngine) storeMessage(ctx context.Context, message *models.Message) (models.Message, error) {
	storedMessage, err := system.DatabaseClient.CreateMessage(ctx, models.CreateMessageParams{
		UserID:   message.UserID,
		Role:     message.Role,
		Content:  message.Content,
		ParentID: message.ParentID,
		AgentID:  message.AgentID,
	})
	if err != nil {
		return models.Message{}, errors.WithStack(err)
	}

	fmt.Printf("Stored message [%s][%s]: %s\n", storedMessage.Role.String, storedMessage.UserID.UUID.String(), firstN(storedMessage.Content.String, 100))

	return storedMessage, nil
}

func buildUserMessage(userId uuid.UUID, messageContent string, lastMessages []models.GetMessagesByUserIDRow) models.Message {
	userMessage := models.Message{
		Role:    sql.NullString{String: openai.ChatMessageRoleUser, Valid: true},
		UserID:  uuid.NullUUID{UUID: userId, Valid: true},
		Content: sql.NullString{String: messageContent, Valid: true},
	}

	if len(lastMessages) > 0 {
		userMessage.ParentID = uuid.NullUUID{UUID: lastMessages[len(lastMessages)-1].ID, Valid: true}
	}

	return userMessage
}

func buildChatbotMessage(userId uuid.UUID, messageContent string, parentId uuid.UUID, agent *models.Agent) models.Message {
	responseMessage := models.Message{
		Role:     sql.NullString{String: openai.ChatMessageRoleAssistant, Valid: true},
		UserID:   uuid.NullUUID{UUID: userId, Valid: true},
		Content:  sql.NullString{String: messageContent, Valid: true},
		ParentID: uuid.NullUUID{UUID: parentId, Valid: true},
	}

	if agent != nil {
		responseMessage.AgentID = uuid.NullUUID{UUID: agent.ID, Valid: true}
	}

	return responseMessage
}

// map json to a struct
type DataStruct struct {
	StationCode    string  `json:"station_code"`
	StreamName     string  `json:"stream_name"`
	RainIntensity  int     `json:"rain_intensity"`
	RainLevel      float64 `json:"rain_level"`
	CurrentWeather string  `json:"current_weather"`
	Datetime       string  `json:"date_time"`
}

func (system *SystemEngine) callHasuraEndpoint(bodyContent string) (bool, error) {
	// HTTP endpoint
	posturl := "https://redmop.practical-action.minsky.cc/api/rest/v1/instrument-records"

	// JSON body
	body := []byte(bodyContent)

	// Create a HTTP post request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-hasura-admin-secret", system.Config.HasuraToken)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}
