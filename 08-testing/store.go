package main

import (
	"context"
	"os"

	"github.com/mdelapenya/genai-testcontainers-go/testing/pgvector"
	"github.com/mdelapenya/genai-testcontainers-go/testing/weaviate"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores"
)

func selectStore(ctx context.Context, embedder embeddings.Embedder) (vectorstores.VectorStore, error) {
	storeTypeEnv := os.Getenv("VECTOR_STORE")

	switch storeTypeEnv {
	case "pgvector":
		return pgvector.NewStore(ctx, embedder)
	default:
		return weaviate.NewStore(ctx, embedder)
	}
}
