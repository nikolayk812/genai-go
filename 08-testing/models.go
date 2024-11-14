package main

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	tcollama "github.com/testcontainers/testcontainers-go/modules/ollama"
	"github.com/tmc/langchaingo/llms/ollama"
)

func buildChatModel() (*ollama.LLM, error) {
	c, err := tcollama.Run(context.Background(), "mdelapenya/"+model+":0.3.13-"+tag, testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name: "chat-model",
		},
		Reuse: true,
	}))
	if err != nil {
		return nil, err
	}

	ollamaURL, err := c.ConnectionString(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connection string: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel(modelName), ollama.WithServerURL(ollamaURL))
	if err != nil {
		return nil, fmt.Errorf("ollama new: %w", err)
	}

	return llm, nil
}

func buildEmbeddingModel() (*ollama.LLM, error) {
	c, err := tcollama.Run(context.Background(), "mdelapenya/all-minilm:0.3.13-22m", testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name: "embeddings-model",
		},
		Reuse: true,
	}))
	if err != nil {
		return nil, err
	}

	ollamaURL, err := c.ConnectionString(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connection string: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel("all-minilm:22m"), ollama.WithServerURL(ollamaURL))
	if err != nil {
		return nil, fmt.Errorf("ollama new: %w", err)
	}

	return llm, nil
}
