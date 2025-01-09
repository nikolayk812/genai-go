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
	const question string = "Which pokemon has more moves, Haunter or Gengar?"

	log.Printf("Question: %s", question)

	var c *tcollama.OllamaContainer

	// 3b model version is required to use Tools.
	// See https://ollama.com/library/llama3.2
	c, err = tcollama.Run(context.Background(), "mdelapenya/llama3.2:0.5.4-3b", testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name: "chat-model",
		},
		Reuse: true,
	}))
	if err != nil {
		return err
	}
	defer func() {
		terminateErr := testcontainers.TerminateContainer(c)
		if terminateErr != nil {
			err = fmt.Errorf("terminate container: %w", terminateErr)
		}
	}()

	ollamaURL, connErr := c.ConnectionString(context.Background())
	if connErr != nil {
		err = fmt.Errorf("connection string: %w", connErr)
		return
	}

	llm, ollamaErr := ollama.New(
		ollama.WithModel("llama3.2:3b"),
		ollama.WithServerURL(ollamaURL),
	)
	if ollamaErr != nil {
		err = fmt.Errorf("ollama new: %w", ollamaErr)
		return
	}

	var msgs []llms.MessageContent

	// system message defines the available tools.
	msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeSystem, systemMessage()))
	msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeHuman, question))

	ctx := context.Background()

	response := ""

	for retries := 3; retries > 0; retries = retries - 1 {
		resp, err := llm.GenerateContent(ctx, msgs)
		if err != nil {
			log.Fatal(err)
		}

		choice1 := resp.Choices[0]
		msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeAI, choice1.Content))

		if cs := unmarshalCall(choice1.Content); cs != nil {
			for _, c := range cs {
				msg, cont := dispatchCall(&c)
				if !cont {
					retries = 0
					response = msg.Parts[0].(llms.TextContent).Text
					break
				}
				msgs = append(msgs, msg)
			}
		} else {
			// Ollama doesn't always respond with a function call, let it try again.
			msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeHuman, "Sorry, I don't understand. Please try again."))
		}
	}

	if response == "" {
		log.Fatal("response is empty")
	}

	log.Printf("Final response: %v", response)

	return nil
}
