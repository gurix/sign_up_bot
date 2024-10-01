package llm

import (
	"context"
	"fmt"

	"github.com/gurix/sign_up_bot/models"
	"github.com/tmc/langchaingo/llms"
)

type LLM interface {
	GenerateResponse(ctx context.Context, input string) *llms.ContentResponse
}

// Generate a response using the LLaMA model.
func (llm Ai) GenerateResponse(ctx context.Context, input string, dialog models.ChatDialog) (*llms.ContentResponse, error) {

	// Rebuild history
	content := []llms.MessageContent{}
	for _, message := range dialog.Messages {
		// Append user inputs
		content = append(content,
			llms.TextParts(llms.ChatMessageTypeHuman, message.Message))

		// Append Ai responds
		content = append(content,
			llms.TextParts(llms.ChatMessageTypeAI, message.Result))
	}

	// Append finally the current user input
	content = append(content,
		llms.TextParts(llms.ChatMessageTypeHuman, input))

	resp, err := llm.Client.GenerateContent(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	return resp, nil
}

func GetFirstContent(response *llms.ContentResponse) string {
	return response.Choices[0].Content
}
