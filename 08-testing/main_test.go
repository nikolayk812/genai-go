package main

import (
	"context"
	"encoding/json"
	"github.com/nikolayk812/genai-go/08-testing/ai"
	internalhttp "github.com/nikolayk812/genai-go/internal/http"
	"net/http"
	"strings"
	"testing"

	"github.com/chewxy/math32"
	"github.com/tmc/langchaingo/embeddings"
)

var httpCli = &http.Client{
	Transport: internalhttp.NewLoggingRoundTripper(),
}

func Test1_oldSchool(t *testing.T) {
	chatModel, err := buildChatModel(httpCli)
	if err != nil {
		t.Fatalf("build chat model: %s", err)
	}

	containsFn := func(innerT *testing.T, answer string) {
		innerT.Helper()

		if !strings.Contains(answer, "cloud.logs.verbose") {
			innerT.Fatalf("%s does not contain 'cloud.logs.verbose'", answer)
		}
	}

	t.Run("pgvector", func(t *testing.T) {
		t.Setenv("VECTOR_STORE", "pgvector")

		t.Run("straight-answer", func(tt *testing.T) {
			answer, err := straightAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("straight chat: %s", err)
			}

			containsFn(tt, answer)
		})

		t.Run("ragged-answer", func(tt *testing.T) {
			answer, err := raggedAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("straight chat: %s", err)
			}

			containsFn(tt, answer)
		})
	})

	t.Run("weaviate", func(t *testing.T) {
		t.Run("straight-answer", func(tt *testing.T) {
			answer, err := straightAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("straight chat: %s", err)
			}

			containsFn(tt, answer)
		})

		t.Run("ragged-answer", func(tt *testing.T) {
			answer, err := raggedAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("straight chat: %s", err)
			}

			containsFn(tt, answer)
		})
	})
}

func Test2_embeddings(t *testing.T) {
	chatModel, err := buildChatModel(httpCli)
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

	similarityFn := func(innerT *testing.T, answer string) {
		innerT.Helper()

		answerVector, err := embedder.EmbedDocuments(context.Background(), []string{answer})
		if err != nil {
			innerT.Fatal("embed answer", err)
		}

		sim := cosineSimilarity(innerT, reference[0], answerVector[0])
		if sim <= 0.80 {
			innerT.Fatalf("similarity is %f: %s", sim, answer)
		}
		t.Logf("similarity is %f: %s", sim, answer)
	}

	t.Run("pgvector", func(t *testing.T) {
		t.Setenv("VECTOR_STORE", "pgvector")

		t.Run("straight-answer", func(tt *testing.T) {
			answer, err := straightAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("straight answer: %s", err)
			}

			similarityFn(tt, answer)
		})

		t.Run("ragged-answer", func(tt *testing.T) {
			answer, err := raggedAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("ragged answer: %s", err)
			}

			similarityFn(tt, answer)
		})
	})

	t.Run("weaviate", func(t *testing.T) {
		t.Run("straight-answer", func(tt *testing.T) {
			answer, err := straightAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("straight answer: %s", err)
			}

			similarityFn(tt, answer)
		})

		t.Run("ragged-answer", func(tt *testing.T) {
			answer, err := raggedAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("ragged answer: %s", err)
			}

			similarityFn(tt, answer)
		})
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

func Test3_evaluatorAgent(t *testing.T) {
	reference := `
There 2 things which answer must contain:
- Answer must indicate that you can enable verbose logging in Testcontainers Desktop by setting the property cloud.logs.verbose to true in the ~/.testcontainers.properties file
- Answer must indicate that you can enable verbose logging in Testcontainers Desktop by adding the --verbose flag when running the cli
`

	chatModel, err := buildChatModel(httpCli)
	if err != nil {
		t.Fatalf("build chat model: %s", err)
	}

	evaluatorAgent := ai.NewEvaluatorAgent(chatModel)

	evaluateFn := func(innerT *testing.T, answer string) {
		innerT.Helper()

		resp, err := evaluatorAgent.Evaluate(question, answer, reference)
		if err != nil {
			innerT.Fatalf("validate: %s", err)
		}

		type r struct {
			Response string `json:"response"`
			Reason   string `json:"reason"`
		}

		var jsonResp r
		err = json.Unmarshal([]byte(resp), &jsonResp)
		if err != nil {
			innerT.Fatalf("json unmarshal: %s", err)
		}

		if jsonResp.Response != "yes" {
			innerT.Fatalf("chat: %+v", jsonResp)
		}
	}

	t.Run("pgvector", func(t *testing.T) {
		t.Setenv("VECTOR_STORE", "pgvector")

		t.Run("straight-answer", func(tt *testing.T) {
			answer, err := straightAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("straight answer: %s", err)
			}

			evaluateFn(tt, answer)
		})

		t.Run("ragged-answer", func(tt *testing.T) {
			answer, err := raggedAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("ragged answer: %s", err)
			}

			evaluateFn(tt, answer)
		})
	})

	t.Run("weaviate", func(t *testing.T) {
		t.Run("straight-answer", func(tt *testing.T) {
			answer, err := straightAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("straight answer: %s", err)
			}

			evaluateFn(tt, answer)
		})

		t.Run("ragged-answer", func(tt *testing.T) {
			answer, err := raggedAnswer(t.Context(), chatModel)
			if err != nil {
				tt.Fatalf("ragged answer: %s", err)
			}

			evaluateFn(tt, answer)
		})
	})
}
