package main

import (
	"bufio"
	"context"
	"fmt"
	internalhttp "github.com/nikolayk812/genai-go/internal/http"
	internalllms "github.com/nikolayk812/genai-go/internal/llms"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
		Transport: internalhttp.NewLoggingRoundTripper(),
	}

	llm, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithHTTPClient(httpCli),
	)
	if err != nil {
		return fmt.Errorf("ollama.New: %w", err)
	}

	// listen for interrupt signals to end the chat session gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nInterrupt signal received, ending chat session")
		os.Exit(0)
	}()

	var conversation []llms.MessageContent

	reader := bufio.NewReader(os.Stdin)

	// Enter a conversation loop
	for {
		fmt.Print("\nYou: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reader.ReadString: %w", err)
		}

		input = strings.TrimSpace(input)
		switch input {
		case "quit", "exit", "bye":
			fmt.Println("Ending chat session")
			os.Exit(0)
		}

		conversation = append(conversation, llms.TextParts(llms.ChatMessageTypeHuman, input))

		// TODO: skip earlier messages if context length is approaching the limit, or to save costs
		// see: internalllms.TotalTokens(responses)

		responses, err := llm.GenerateContent(ctx, conversation, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))
		if err != nil {
			return fmt.Errorf("llm.GenerateContent: %w", err)
		}

		choiceContents := internalllms.ContentResponseToStrings(responses)
		conversation = append(conversation, llms.TextParts(llms.ChatMessageTypeAI, choiceContents...))
	}
}
