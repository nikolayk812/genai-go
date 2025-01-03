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
	var c *tcollama.OllamaContainer

	const (
		modelName string = "DavidAU/DistiLabelOrca-TinyLLama-1.1B-Q8_0-GGUF"
		modelFile string = "distilabelorca-tinyllama-1.1b.Q8_0.gguf"

		// the name of the image that we will commit to the registry.
		// It's using a significant and valid image name in order to be
		// identified by this project.
		imageName string = "distilabelorca-tinyllama-guff"
	)

	opts := []testcontainers.ContainerCustomizer{
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name:          "huggingface-model",
				ImagePlatform: "linux/amd64",
			},
			Reuse: true,
		}),
	}

	// first try to run the image with the huggingface tooling and model already installed
	c, err = tcollama.Run(context.Background(), imageName, opts...)
	if err != nil {
		// the image does not exist: build it including the huggingface tooling and model
		c, err = tcollama.Run(context.Background(),
			"ollama/ollama:0.5.4",
			append(opts, WithHuggingfaceModel(modelName, modelFile))...,
		)
		if err != nil {
			return err
		}

		// commit the image to the registry so that we can be reuse it in subsequent runs
		// without having to build it again.
		err = c.Commit(context.Background(), imageName)
		if err != nil {
			return err
		}
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
		// IMPORTANT: the Ollama model is the model file.
		ollama.WithModel(modelFile),
		ollama.WithServerURL(ollamaURL),
	)
	if ollamaErr != nil {
		err = fmt.Errorf("ollama new: %w", ollamaErr)
		return
	}

	ctx := context.Background()
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a fellow Go developer."),
		llms.TextParts(llms.ChatMessageTypeHuman, "Provide 3 short bullet points explaining why Go is awesome"),
	}

	// The response from the model happens when the model finishes processing the input, which it's usually slow.
	completion, generateContentErr := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if generateContentErr != nil {
		err = fmt.Errorf("llm generate content: %w", generateContentErr)
		return
	}

	for _, choice := range completion.Choices {
		fmt.Println(choice.Content)
	}

	return nil
}
