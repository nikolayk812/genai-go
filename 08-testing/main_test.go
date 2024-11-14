package main

import (
	"context"
	"strings"
	"testing"

	"github.com/chewxy/math32"
	"github.com/tmc/langchaingo/embeddings"
)

func Test1(t *testing.T) {
	chatModel, err := buildChatModel()
	if err != nil {
		t.Fatalf("build chat model: %s", err)
	}

	t.Run("straight-answer", func(t *testing.T) {
		answer, err := straightAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight chat: %s", err)
		}

		if !strings.Contains(answer, "cloud.logs.verbose = true") {
			t.Fatalf("straight chat: %s", answer)
		}
	})

	t.Run("ragged-answer", func(t *testing.T) {
		answer, err := raggedAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight chat: %s", err)
		}

		if !strings.Contains(answer, "cloud.logs.verbose = true") {
			t.Fatalf("ragged chat: %s", answer)
		}
	})
}

func Test2(t *testing.T) {
	chatModel, err := buildChatModel()
	if err != nil {
		t.Fatalf("build chat model: %s", err)
	}

	embeddingModel, err := buildEmbeddingModel()
	if err != nil {
		t.Fatalf("build embedding model: %s", err)
	}

	embedder, err := embeddings.NewEmbedder(embeddingModel)
	if err != nil {
		t.Fatalf("new embedder: %s", err)
	}

	reference, err := embedder.EmbedDocuments(context.Background(), []string{
		"To enable verbose logging in Testcontainers Desktop you can set the property cloud.logs.verbose to true in the ~/.testcontainers.properties file or add the --verbose flag when running the cli",
	})
	if err != nil {
		t.Fatal("embed query", err)
	}

	t.Run("straight-answer", func(t *testing.T) {
		answer, err := straightAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight chat: %s", err)
		}

		answerVector, err := embedder.EmbedDocuments(context.Background(), []string{answer})
		if err != nil {
			t.Fatal("embed answer", err)
		}

		sim := cosineSimilarity(t, reference[0], answerVector[0])
		if sim <= 0.80 {
			t.Fatalf("straight chat: %s", answer)
		}
	})

	t.Run("ragged-answer", func(t *testing.T) {
		answer, err := raggedAnswer(chatModel)
		if err != nil {
			t.Fatalf("ragged chat: %s", err)
		}

		answerVector, err := embedder.EmbedDocuments(context.Background(), []string{answer})
		if err != nil {
			t.Fatal("embed answer", err)
		}

		sim := cosineSimilarity(t, reference[0], answerVector[0])
		if sim <= 0.80 {
			t.Fatalf("ragged chat: %s", answer)
		}
	})
}

// cosineSimilarity calculates the cosine similarity between two vectors.
// See https://github.com/tmc/langchaingo/blob/238d1c713de3ca983e8f6066af6b9080c9b0e088/examples/cybertron-embedding-example/cybertron-embedding.go#L19
func cosineSimilarity(t *testing.T, x, y []float32) float32 {
	t.Helper()

	if len(x) != len(y) {
		t.Fatal("x and y have different lengths")
	}

	var dot, nx, ny float32

	for i := range x {
		nx += x[i] * x[i]
		ny += y[i] * y[i]
		dot += x[i] * y[i]
	}

	return dot / (math32.Sqrt(nx) * math32.Sqrt(ny))
}
