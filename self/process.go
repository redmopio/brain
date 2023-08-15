package self

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/minskylab/brain/models"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

func (brain *BrainEngine) ProcessMessageResponse(ctx context.Context, user models.User, lastMessages []models.Message, inputMessage models.Message) (string, error) {
	openAiResponse, err := brain.processMessageWithOpenAI(ctx, user, lastMessages, inputMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	responseWithDataKeywords := []string{"tus datos han sido registrados de manera exitosa"}
	openAiMessageContent := openAiResponse.Choices[0].Message.Content

	for _, keyword := range responseWithDataKeywords {
		if strings.Contains(strings.ToLower(openAiMessageContent), keyword) {
			parsedDataFromOpenAI, err := brain.processDataMessageWithOpenAI(ctx, user, openAiMessageContent)
			if err != nil {
				fmt.Printf("Error processing data message: %s\n", err.Error())
				return "", errors.WithStack(err)
			}

			parsedDataFromOpenAIContent := parsedDataFromOpenAI.Choices[0].Message.Content

			if strings.Contains(parsedDataFromOpenAIContent, "Hubo un error al parsear la data") {
				return "Disculpa, creo que no ingresaste todos los datos. Revisa el ejemplo y vuelve a intentarlo.", nil
			}

			// parse content from openAiResponse to DataStruct
			var data DataStruct
			err = json.Unmarshal([]byte(parsedDataFromOpenAI.Choices[0].Message.Content), &data)
			if err != nil {
				fmt.Printf("Error parsing data: %s\n", err.Error())
				return "Disculpa, creo que no ingresaste todos los datos. Revisa el ejemplo y vuelve a intentarlo.", nil
			}

			fmt.Printf("\t -> Calling Hasura endpoint with data: %s\n", parsedDataFromOpenAIContent)
			_, err = brain.callHasuraEndpoint(parsedDataFromOpenAIContent)
			if err != nil {
				return "Lo siento, hubo un error al guardar tus datos. Por favor, intentalo nuevamente en unos minutos.", nil
			}

			return openAiMessageContent, nil
		}
	}

	fmt.Printf("\t -> Processing normal message: %s\n", openAiMessageContent)

	return openAiMessageContent, nil
}

func (brain *BrainEngine) processDataMessageWithOpenAI(ctx context.Context, user models.User, messageContent string) (*openai.ChatCompletionResponse, error) {
	agentParseData, err := brain.getAgent(ctx, models.AgentTypeAgentParseData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	messages := brain.prepareOpenAiMessagesForData(&agentParseData, messageContent)

	response, err := brain.LLMEngine.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (brain *BrainEngine) processMessageWithOpenAI(ctx context.Context, user models.User, lastMessages []models.Message, inputMessage models.Message) (*openai.ChatCompletionResponse, error) {
	agentStoreData, err := brain.getAgent(ctx, models.AgentTypeAgentStoreData)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	messages := brain.prepareMessagesForConversation(&agentStoreData, user, lastMessages, inputMessage)

	fmt.Printf("Total messages: %d\n", len(messages))

	for _, msg := range messages {
		fmt.Printf("\tMessage: [%s] %s\n", msg.Role, firstN(msg.Content, 100))
	}

	response, err := brain.LLMEngine.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (brain *BrainEngine) prepareMessagesForConversation(agent *models.Agent, user models.User, lastMessages []models.Message, inputMessage models.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: agent.Description,
	})

	userName := user.UserName.String
	userName = strings.ReplaceAll(userName, " ", "_")

	for _, msg := range lastMessages {
		messageName := userName
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

	messages = append(messages, openai.ChatCompletionMessage{
		Name:    userName,
		Role:    inputMessage.Role.String,
		Content: inputMessage.Content.String,
	})

	return messages
}

func (brain *BrainEngine) prepareOpenAiMessagesForData(agent *models.Agent, messageContent string) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: agent.Description,
	})

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: messageContent,
	})

	return messages
}

func (brain *BrainEngine) getAgent(ctx context.Context, agentName models.AgentType) (models.Agent, error) {
	agent, err := brain.DatabaseClient.GetAgentByName(ctx, agentName)
	if err != nil {
		return models.Agent{}, errors.WithStack(err)
	}

	fmt.Printf("Agent: %s\n", agent.Name)

	return agent, nil
}

func (brain *BrainEngine) storeMessage(ctx context.Context, message *models.Message) (models.Message, error) {
	storedMessage, err := brain.DatabaseClient.CreateMessage(ctx, models.CreateMessageParams{
		UserID:   message.UserID,
		Role:     message.Role,
		Content:  message.Content,
		ParentID: message.ParentID,
	})
	if err != nil {
		return models.Message{}, errors.WithStack(err)
	}

	fmt.Printf("Stored message [%s][%s]: %s\n", storedMessage.Role.String, storedMessage.UserID.UUID.String(), firstN(storedMessage.Content.String, 100))

	return storedMessage, nil
}

func buildUserMessage(userId uuid.UUID, messageContent string, lastMessages []models.Message) models.Message {
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

func buildChatbotMessage(userId uuid.UUID, messageContent string, parentId uuid.UUID) models.Message {
	responseMessage := models.Message{
		Role:     sql.NullString{String: openai.ChatMessageRoleAssistant, Valid: true},
		UserID:   uuid.NullUUID{UUID: userId, Valid: true},
		Content:  sql.NullString{String: messageContent, Valid: true},
		ParentID: uuid.NullUUID{UUID: parentId, Valid: true},
	}

	return responseMessage
}
