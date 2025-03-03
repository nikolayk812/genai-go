package main

import (
	"context"
	"fmt"
	internalhttp "github.com/nikolayk812/genai-go/internal/http"
	internalllms "github.com/nikolayk812/genai-go/internal/llms"
	"log"
	"net/http"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run(ctx context.Context) error {
	httpCli := &http.Client{
		Transport: internalhttp.NewLoggingRoundTripper(http.DefaultTransport),
	}

	llm, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithHTTPClient(httpCli),
	)
	if err != nil {
		return fmt.Errorf("ollama.New: %w", err)
	}

	originalMessage := `
		What is the current topic of the conference?
	`

	augmentedMessage := fmt.Sprintf(`
		%s

		Use the following bullet points to answer the question:
		- The Conference is about how to leverage Testcontainers for building Generative AI applications.
		- The meeting will explore how Testcontainers can be used to create a seamless development environment for AI projects.

		Do not indicate that you have been given any additional information.
		`, originalMessage)

	originalContent := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, originalMessage),
	}

	originalCompletion, err := llm.GenerateContent(
		ctx, originalContent,
		llms.WithTemperature(0.0001),
		llms.WithTopK(1),
	)
	if err != nil {
		return fmt.Errorf("llm.GenerateContent[original]: %w", err)
	}

	fmt.Println("\nOriginal completion:")
	for _, choiceContent := range internalllms.ContentResponseToStrings(originalCompletion) {
		fmt.Println(choiceContent)
	}

	augmentedContent := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, augmentedMessage),
	}

	augmentedCompletion, err := llm.GenerateContent(
		ctx, augmentedContent,
		llms.WithTemperature(0.0001),
		llms.WithTopK(1),
	)
	if err != nil {
		return fmt.Errorf("llm.GenerateContent[augmented] %w", err)
	}

	fmt.Println("\nAugmented completion:")
	for _, choiceContent := range internalllms.ContentResponseToStrings(augmentedCompletion) {
		fmt.Println(choiceContent)
	}

	return nil
}
