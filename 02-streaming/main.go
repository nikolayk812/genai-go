package main

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	tcollama "github.com/testcontainers/testcontainers-go/modules/ollama"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run() (err error) {
	c, err := tcollama.Run(context.Background(), "mdelapenya/qwen2:0.3.13-0.5b")
	if err != nil {
		return err
	}
	defer func() {
		err = testcontainers.TerminateContainer(c)
		if err != nil {
			err = fmt.Errorf("terminate container: %w", err)
		}
	}()

	ollamaURL, err := c.ConnectionString(context.Background())
	if err != nil {
		return fmt.Errorf("connection string: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel("qwen2:0.5b"), ollama.WithServerURL(ollamaURL))
	if err != nil {
		return fmt.Errorf("ollama new: %w", err)
	}

	ctx := context.Background()
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Give me a detailed and long explanation of why Testcontainers for Go is great"),
	}

	_, err = llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		return fmt.Errorf("llm generate content: %w", err)
	}

	return nil
}
