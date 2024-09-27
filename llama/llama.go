package llama

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var (
	client *ollama.LLM
	once   sync.Once
)

func Connect() *ollama.LLM {
	once.Do(func() {
		var err error
		client, err = ollama.New(ollama.WithModel("llama3.1"))
		log.Printf("Initialize %v", client)
		if err != nil {
			log.Fatal(err)
		}
	})
	return client
}

// Generate a response using the LLaMA model.
func GenerateResponse(ctx context.Context, input string) (*llms.ContentResponse, error) {
	client := Connect()

	log.Printf("Use %v", client)
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	resp, err := client.GenerateContent(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	return resp, nil
}

func GetFirstContent(response *llms.ContentResponse) string {
	return response.Choices[0].Content
}
