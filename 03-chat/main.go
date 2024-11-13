package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	c, err := tcollama.Run(context.Background(), "mdelapenya/llama3.2:0.3.13-1b")
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
			return fmt.Errorf("read string: %w", err)
		}

		input = strings.TrimSpace(input)
		switch input {
		case "quit", "exit":
			fmt.Println("Ending chat session")
			os.Exit(0)
		}

		conversation = append(conversation, llms.TextParts(llms.ChatMessageTypeHuman, input))

		ctx := context.Background()
		_, err = llm.GenerateContent(ctx, conversation, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))
		if err != nil {
			return fmt.Errorf("llm generate content: %w", err)
		}
	}
}
