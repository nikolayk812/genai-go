package main

import (
	"context"
	"fmt"
	"github.com/chewxy/math32"
	internalhttp "github.com/nikolayk812/genai-go/internal/http"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	docs := []string{
		"A cat is a small domesticated carnivorous mammal",
		"A tiger is a large carnivorous feline mammal",
		"Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container",
		"Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code.",
	}

	if err := run(ctx, docs); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run(ctx context.Context, docs []string) error {
	httpCli := &http.Client{
		Transport: internalhttp.NewLoggingRoundTripper(),
	}

	llm, err := ollama.New(
		ollama.WithModel("nomic-embed-text:v1.5"),
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithHTTPClient(httpCli),
	)
	if err != nil {
		return fmt.Errorf("ollama.New: %w", err)
	}

	//llm, err := openai.New(
	//	openai.WithModel("text-embedding-3-small"),
	//	openai.WithToken(os.Getenv("OPENAI_API_KEY")),
	//	openai.WithHTTPClient(httpCli),
	//)
	//if err != nil {
	//	return fmt.Errorf("openai.New: %w", err)
	//}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return fmt.Errorf("embeddings.NewEmbedder: %w", err)
	}

	vectors, err := embedder.EmbedDocuments(ctx, docs)
	if err != nil {
		return fmt.Errorf("embeddings.EmbedDocuments: %w", err)
	}

	fmt.Println("Similarities:")
	docLen := len(docs)
	for i := 0; i < docLen; i++ {
		for j := i; j < docLen; j++ {
			fmt.Printf("%d ~ %d = %0.2f\n", i, j, cosineSimilarity(vectors[i], vectors[j]))
		}
	}

	return nil
}

// cosineSimilarity calculates the cosine similarity between two vectors.
// See https://github.com/tmc/langchaingo/blob/238d1c713de3ca983e8f6066af6b9080c9b0e088/examples/cybertron-embedding-example/cybertron-embedding.go#L19
func cosineSimilarity(x, y []float32) float32 {
	if len(x) != len(y) {
		log.Fatal("x and y have different lengths")
	}

	var dot, nx, ny float32

	for i := range x {
		nx += x[i] * x[i]  // Sum of squares of x's elements
		ny += y[i] * y[i]  // Sum of squares of y's elements
		dot += x[i] * y[i] // Dot product of x and y
	}

	// return dot / float32(math.Sqrt(float64(nx)) * math.Sqrt(float64(ny)))
	return dot / (math32.Sqrt(nx) * math32.Sqrt(ny))
}
