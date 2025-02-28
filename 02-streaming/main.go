package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"io"
	"log"
	"net/http"
)

// ollama run llama3.2

func main() {
	ctx := context.Background()

	httpClient := &http.Client{
		Transport: &CustomRoundTripper{
			Transport: http.DefaultTransport,
		},
	}

	if err := run(ctx, httpClient); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func run(ctx context.Context, httpCli *http.Client) error {
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

	return nil
}

// CustomRoundTripper logs each request
type CustomRoundTripper struct {
	Transport http.RoundTripper
}

func (c *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log method and URL
	log.Printf("Request: %s %s", req.Method, req.URL)

	if req.Body == nil {
		return http.DefaultTransport.RoundTrip(req)
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	// Pretty-print the JSON body
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err != nil {
		return nil, fmt.Errorf("json.Indent: %w", err)
	}

	log.Printf("Request Body: %s", prettyJSON.String())

	// Reset the request body so it can be read again
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return c.Transport.RoundTrip(req)
}
