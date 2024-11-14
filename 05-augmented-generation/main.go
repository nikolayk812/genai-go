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
	c, err := tcollama.Run(context.Background(), "mdelapenya/llama3.2:0.3.13-1b", testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name: "chat-model",
		},
		Reuse: true,
	}))
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

	llm, err := ollama.New(ollama.WithModel("llama3.2:1b"), ollama.WithServerURL(ollamaURL))
	if err != nil {
		return fmt.Errorf("ollama new: %w", err)
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

	ctx := context.Background()
	originalContent := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, originalMessage),
	}

	originalCompletion, err := llm.GenerateContent(
		ctx, originalContent,
		llms.WithTemperature(0.0001),
		llms.WithTopK(1),
	)
	if err != nil {
		return fmt.Errorf("llm generate original content: %w", err)
	}

	fmt.Println("\nOriginal completion:")
	for _, choice := range originalCompletion.Choices {
		fmt.Println(choice.Content)
	}

	augmentedContent := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, augmentedMessage),
	}

	augmentedCompletion, err := llm.GenerateContent(
		ctx, augmentedContent,
		llms.WithTemperature(0.0001),
		llms.WithTopK(1),
	)
	if err != nil {
		return fmt.Errorf("llm generate original content: %w", err)
	}

	fmt.Println("\nAugmented completion:")
	for _, choice := range augmentedCompletion.Choices {
		fmt.Println(choice.Content)
	}

	return nil
}
