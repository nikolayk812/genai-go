package main

import (
	"context"
	_ "embed"
	"fmt"
	internalhttp "github.com/nikolayk812/genai-go/internal/http"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"log"
	"net/http"
	"os"

	"github.com/tmc/langchaingo/llms"
)

//go:embed images/cat.jpeg
var catImage []byte

func main() {
	ctx := context.Background()

	if err := run(ctx, true); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run(ctx context.Context, useOllama bool) error {
	httpCli := &http.Client{
		Transport: internalhttp.NewLoggingRoundTripper(),
	}

	var llm llms.Model

	llm, err := ollama.New(
		ollama.WithModel("moondream:1.8b"),
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithHTTPClient(httpCli),
	)
	if err != nil {
		return fmt.Errorf("ollama.New: %w", err)
	}

	if !useOllama {
		llm, err = openai.New(
			openai.WithModel("gpt-4-turbo"),
			openai.WithToken(os.Getenv("OPENAI_API_KEY")),
			openai.WithHTTPClient(httpCli),
		)
		if err != nil {
			return fmt.Errorf("openai.New: %w", err)
		}
	}

	var content []llms.MessageContent

	imagePart := imagePart(useOllama)

	content = append(content, llms.MessageContent{
		Role: llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{
			llms.TextPart("Please tell me what you see in this image"),
			imagePart,
		},
	})

	_, err = llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		return fmt.Errorf("llm.GenerateContent: %w", err)
	}

	return nil
}

func imagePart(useOllama bool) llms.ContentPart {
	if !useOllama {
		return llms.ImageURLPart("https://media.istockphoto.com/id/1511923057/photo/cute-ginger-cat-lying-on-carpet-at-home-closeup.jpg?s=612x612&w=0&k=20&c=2mliyfSTOrFUW4O2vU2W8fMF220Fu0TOoNL3Kkiabes=")
	}

	return llms.BinaryPart("image/jpeg", catImage)
}
