package main

import (
	"context"
	"fmt"
	internalhttp "github.com/nikolayk812/genai-go/internal/http"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
	"net/http"
)

// ollama run llama3.2

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("run: %v", err)
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

	//llm, err := openai.New(
	//	openai.WithModel("gpt-4"),
	//	openai.WithToken(os.Getenv("OPENAI_API_KEY")),
	//	openai.WithHTTPClient(httpCli),
	//)
	//if err != nil {
	//	return fmt.Errorf("openai.New: %w", err)
	//}

	content := []llms.MessageContent{
		// does not work for "system" type with Ollama
		llms.TextParts(llms.ChatMessageTypeHuman, "Give me a detailed and long explanation of why Testcontainers for Go is great"),
	}

	// Streaming is needed because models are usually slow in responding, so showing progress is important.
	_, err = llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		return fmt.Errorf("llm.GenerateContent: %w", err)
	}

	return nil
}
