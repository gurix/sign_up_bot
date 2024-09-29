package llm

import (
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

type Ai struct {
	Client *ollama.LLM
}

// Connect contects to the client.
func OllamaConnect() *ollama.LLM {
	client, err := ollama.New(ollama.WithModel("llama3.1"))
	if err != nil {
		log.Fatal(err)
	}

	return client
}
