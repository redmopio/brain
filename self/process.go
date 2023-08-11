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
	greetingKeywords := []string{"hola", "hello", "hi", "hey"}
	helpKeywords := []string{"ayuda", "help"}
	byeKeywords := []string{"adios", "bye"}
	thanksKeywords := []string{"gracias", "thanks"}
	inputDataKeywords := []string{"datos", "data", "registrar", "ingresar"}
	actualDataKeywords := []string{"tiempo actual", "nivel de lluvia", "intensidad de la lluvia", "quebrada y distrito", "fecha y hora"}

	for _, keyword := range greetingKeywords {
		if strings.Contains(strings.ToLower(inputMessage.Content.String), keyword) {
			return "Hola, soy el asistente virtual del proyecto Redmop. ¿En qué puedo ayudarte?", nil
		}
	}

	for _, keyword := range helpKeywords {
		if strings.Contains(strings.ToLower(inputMessage.Content.String), keyword) {
			return "Soy parte del proyecto Redmop, un sistema de monitoreo climático distribuido que utiliza una red de voluntarios para recopilar datos sobre las condiciones meteorológicas en una zona geográfica específica. ¿Deseas registrar datos? Por favor, escribe: *Registar datos*", nil
		}
	}

	for _, keyword := range byeKeywords {
		if strings.Contains(strings.ToLower(inputMessage.Content.String), keyword) {
			return "¡Hasta pronto!", nil
		}
	}

	for _, keyword := range thanksKeywords {
		if strings.Contains(strings.ToLower(inputMessage.Content.String), keyword) {
			return "¡De nada!", nil
		}
	}

	for _, keyword := range inputDataKeywords {
		if strings.Contains(strings.ToLower(inputMessage.Content.String), keyword) {
			return `Fantástico! Estos son los datos que necesito (con ejemplos):
- Código de la estación: RC-Cy-CP-34
- Quebrada y distrito: CULTURA Y PROGRESO CHACLACAYO
- Intensidad de la lluvia: 3 (en escala de 1 a 5)
- Nivel de lluvia acumulada: 0.5 (en mm)
- Tiempo actual: Cielo nublado sigue llovizna
- Fecha y hora de la revisión del pluviómetro: 09-06-2023 7:31am

Deben estar en un solo mensaje, puedes copiar el mensaje y modificarlo para más facilidad :D`, nil
		}
	}

	for _, keyword := range actualDataKeywords {
		if strings.Contains(strings.ToLower(inputMessage.Content.String), keyword) {
			openAiResponse, err := brain.processDataMessageWithOpenAI(ctx, user, inputMessage)
			if err != nil {
				return "", errors.WithStack(err)
			}

			// parse content from openAiResponse to DataStruct
			var data DataStruct
			err = json.Unmarshal([]byte(openAiResponse.Choices[0].Message.Content), &data)
			if err != nil {
				fmt.Printf("Error parsing data: %s\n", err.Error())

				return "Disculpa, creo que no ingresaste todos los datos. Revisa el ejemplo y vuelve a intentarlo.", nil
			}

			fmt.Printf("\n\n\tData: %+v\n\n", data)

			return openAiResponse.Choices[0].Message.Content, nil
		}
	}

	return "Disculpa, no tengo una respuesta para eso. ¿Deseas registrar datos? Por favor, escribe: *Registar datos*", nil

	// openAiResponse, err := brain.processMessageWithOpenAI(ctx, user, lastMessages, inputMessage)
	// if err != nil {
	// 	return "", errors.WithStack(err)
	// }

	// return openAiResponse.Choices[0].Message.Content, nil
}

func (brain *BrainEngine) processDataMessageWithOpenAI(ctx context.Context, user models.User, inputMessage models.Message) (*openai.ChatCompletionResponse, error) {
	openAiContext := `El siguiente mensaje es un listado de datos climáticos del proyecto Redmop. Redmop es un sistema de monitoreo climático distribuido que utiliza una red de voluntarios para recopilar datos sobre las condiciones meteorológicas en una zona geográfica específica. Los datos que el mensaje debe contener son:
- Código de la estación (ej.: RC-Cy-CP-34)
- Quebrada y distrito (ej.: CULTURA Y PROGRESO CHACLACAYO)
- Intensidad de la lluvia (en una escala de 1 a 5, ej.: 3)
- Nivel de lluvia acumulada (en mm, ej.: 0.5mm)
- Tiempo actual (ej.: Cielo nublado sigue llovizna)
- Fecha y hora de la revisión del pluviómetro (ej.: 09-06-2023 7:31am)

No añadas explicación, solo los datos.
La respuesta debe ser en formato JSON, con los siguientes campos:
- "station_code": string,
- "stream_name": string,
- "rain_intensity": int,
- "rain_level": float,
- "current_weather": string,
- "timedate": string (ISO format)`

	messages := brain.prepareOpenAiMessagesForData(openAiContext, user, inputMessage)

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
	messages := brain.prepareMessagesForConversation(user, lastMessages, inputMessage)

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

func (brain *BrainEngine) prepareMessagesForConversation(user models.User, lastMessages []models.Message, inputMessage models.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: user.Context.String,
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

func (brain *BrainEngine) prepareOpenAiMessagesForData(systemContext string, user models.User, inputMessage models.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemContext,
	})

	userName := user.UserName.String
	userName = strings.ReplaceAll(userName, " ", "_")

	messages = append(messages, openai.ChatCompletionMessage{
		Name:    userName,
		Role:    inputMessage.Role.String,
		Content: inputMessage.Content.String,
	})

	return messages
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
