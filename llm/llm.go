package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type LLM interface {
	GenerateResponse(ctx context.Context, input string) *llms.ContentResponse
}

// Generate a response using the LLaMA model.
func (llm Ai) GenerateResponse(ctx context.Context, input string) (*llms.ContentResponse, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	resp, err := llm.Client.GenerateContent(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	return resp, nil
}

func GetFirstContent(response *llms.ContentResponse) string {
	return response.Choices[0].Content
}
