package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run(ctx context.Context) error {
	llm, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithServerURL("http://localhost:11434"))
	if err != nil {
		return fmt.Errorf("ollama.New: %w", err)
	}

	//llm, err := openai.New(
	//	openai.WithModel("gpt-4"),
	//	openai.WithToken(os.Getenv("OPENAI_API_KEY")),
	//)
	//if err != nil {
	//	return fmt.Errorf("openai.New: %w", err)
	//}

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a fellow Go developer."),
		llms.TextParts(llms.ChatMessageTypeHuman, "Provide 3 short bullet points explaining why Go is awesome"),
	}

	// The response from the model happens when the model finishes processing the input, which it's usually slow.
	completion, err := llm.GenerateContent(ctx, content)
	if err != nil {
		return fmt.Errorf("llm.GenerateContent: %w", err)
	}

	for _, choice := range completion.Choices {
		fmt.Println(choice.Content)
	}

	return nil
}
