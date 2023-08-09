package self

import (
	"github.com/minskylab/brain/models"
	"github.com/sashabaranov/go-openai"
)

func (brain *BrainEngine) prepareConversationPrompt(user *models.User, lastMessages []*models.Message, message *models.Message) openai.ChatCompletionRequest {
	messages := make([]openai.ChatCompletionMessage, 0, len(lastMessages)+1)

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: user.Context.String,
	})

	for _, msg := range lastMessages {
		role := openai.ChatMessageRoleUser

		if msg.Role.String == openai.ChatMessageRoleAssistant {
			role = openai.ChatMessageRoleAssistant
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content.String,
		})
	}

	return openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	}
}

func (brain *BrainEngine) Predict(user *models.User, message *models.Message) (string, error) {
	// prompt := brain.prepareConversationPrompt(conversation, message)

	// ctx := context.Background()

	// fmt.Println("PROMPT:", prompt)

	// engineName := "text-davinci-003"

	// stopSequences := []string{
	// 	// "\n",
	// }

	// res, err := brain.LLMEngine.Client.Completion()
	// res, err := brain.LLMEngine.Client.CompletionWithEngine(ctx, engineName, gpt3.CompletionRequest{
	// 	Prompt:          []string{prompt},
	// 	Stop:            stopSequences,
	// 	N:               gpt3.IntPtr(1),
	// 	Temperature:     gpt3.Float32Ptr(0.9),
	// 	TopP:            gpt3.Float32Ptr(1.0),
	// 	PresencePenalty: 0.6,
	// 	MaxTokens:       gpt3.IntPtr(1000),
	// })
	// if err != nil {
	// 	return "", errors.WithStack(err)
	// }

	// return strings.TrimSpace(res.Choices[0].Text), nil

	return "", nil
}

func (brain *BrainEngine) ProcessMessageResponse(user *models.User, message *models.Message) (*models.Message, error) {
	// predicted, err := brain.Predict(conversation, message)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// newBuffer := fmt.Sprintf("%s\n%s:%s\n%s:%s",
	// 	conversation.ConversationBuffer.String,
	// 	conversation.UserName.String,
	// 	message,
	// 	brain.Name,
	// 	predicted,
	// )

	// maxLines := 5

	// lines := strings.Split(newBuffer, "\n")

	// if len(lines) > maxLines {
	// 	newBuffer = strings.Join(lines[len(lines)-maxLines:], "\n")
	// } else {
	// 	newBuffer = strings.Join(lines, "\n")
	// }

	// return &MessageResponse{
	// 	PredictedResponse: predicted,
	// 	NewBuffer:         newBuffer,
	// }, nil

	return nil, nil
}
