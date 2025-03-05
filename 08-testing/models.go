package main

import (
	"fmt"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"net/http"

	"github.com/tmc/langchaingo/llms/ollama"
)

func buildChatModel(httpCli *http.Client) (llms.Model, error) {
	llm, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithHTTPClient(httpCli),
	)
	if err != nil {
		return nil, fmt.Errorf("ollama.New: %w", err)
	}

	return llm, nil
}

func buildEmbeddingModel() (embeddings.EmbedderClient, error) {
	llm, err := ollama.New(
		ollama.WithModel("nomic-embed-text:v1.5"),
		ollama.WithServerURL("http://localhost:11434"),
		//ollama.WithHTTPClient(httpCli),
	)
	if err != nil {
		return nil, fmt.Errorf("ollama.New: %w", err)
	}

	return llm, nil
}
