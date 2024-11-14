package main

import (
	"context"
	"strings"
	"testing"

	"github.com/chewxy/math32"
	"github.com/mdelapenya/genai-testcontainers-go/testing/ai"
	"github.com/tmc/langchaingo/embeddings"
)

func Test1_oldSchool(t *testing.T) {
	chatModel, err := buildChatModel()
	if err != nil {
		t.Fatalf("build chat model: %s", err)
	}

	containsFn := func(t *testing.T, answer string) {
		t.Helper()

		if !strings.Contains(answer, "cloud.logs.verbose") {
			t.Fatalf("contains: %s", answer)
		}
	}

	t.Run("straight-answer", func(t *testing.T) {
		answer, err := straightAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight chat: %s", err)
		}

		containsFn(t, answer)
	})

	t.Run("ragged-answer", func(t *testing.T) {
		answer, err := raggedAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight chat: %s", err)
		}

		containsFn(t, answer)
	})
}

func Test2_embeddings(t *testing.T) {
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

	similarityFn := func(t *testing.T, answer string) {
		t.Helper()

		answerVector, err := embedder.EmbedDocuments(context.Background(), []string{answer})
		if err != nil {
			t.Fatal("embed answer", err)
		}

		sim := cosineSimilarity(t, reference[0], answerVector[0])
		if sim <= 0.80 {
			t.Fatalf("similarity: %s", answer)
		}
	}

	t.Run("straight-answer/pgvector", func(t *testing.T) {
		t.Setenv("VECTOR_STORE", "pgvector")

		answer, err := straightAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight answer: %s", err)
		}

		similarityFn(t, answer)
	})

	t.Run("ragged-answer/pgvector", func(t *testing.T) {
		t.Setenv("VECTOR_STORE", "pgvector")

		answer, err := raggedAnswer(chatModel)
		if err != nil {
			t.Fatalf("ragged answer: %s", err)
		}

		similarityFn(t, answer)
	})

	t.Run("straight-answer/weaviate", func(t *testing.T) {
		answer, err := straightAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight answer: %s", err)
		}

		similarityFn(t, answer)
	})

	t.Run("ragged-answer/weaviate", func(t *testing.T) {
		answer, err := raggedAnswer(chatModel)
		if err != nil {
			t.Fatalf("ragged answer: %s", err)
		}

		similarityFn(t, answer)
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

func Test3_validatorAgent(t *testing.T) {
	reference := `
- Answer must indicate that you can enable verbose logging in Testcontainers Desktop by setting the property cloud.logs.verbose to true in the ~/.testcontainers.properties file
- Answer must indicate that you can enable verbose logging in Testcontainers Desktop by adding the --verbose flag when running the cli
`

	chatModel, err := buildChatModel()
	if err != nil {
		t.Fatalf("build chat model: %s", err)
	}

	validatorAgent := ai.NewValidatorAgent(chatModel)

	validateFn := func(t *testing.T, answer string) {
		t.Helper()

		resp, err := validatorAgent.Validate(question, answer, reference)
		if err != nil {
			t.Fatalf("validate: %s", err)
		}

		if resp != "yes" {
			t.Fatalf("chat: %s", validatorAgent.Response())
		}
	}

	t.Run("straight-answer/pgvector", func(t *testing.T) {
		t.Setenv("VECTOR_STORE", "pgvector")

		answer, err := straightAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight answer: %s", err)
		}

		validateFn(t, answer)
	})

	t.Run("ragged-answer/pgvector", func(t *testing.T) {
		t.Setenv("VECTOR_STORE", "pgvector")

		answer, err := raggedAnswer(chatModel)
		if err != nil {
			t.Fatalf("ragged answer: %s", err)
		}

		validateFn(t, answer)
	})

	t.Run("straight-answer/weaviate", func(t *testing.T) {
		answer, err := straightAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight answer: %s", err)
		}

		validateFn(t, answer)
	})

	t.Run("ragged-answer/weaviate", func(t *testing.T) {
		answer, err := raggedAnswer(chatModel)
		if err != nil {
			t.Fatalf("ragged answer: %s", err)
		}

		validateFn(t, answer)
	})
}
