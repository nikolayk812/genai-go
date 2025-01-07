package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/mdelapenya/genai-testcontainers-go/testing/pgvector"
	"github.com/mdelapenya/genai-testcontainers-go/testing/weaviate"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
)

func ingestion(store vectorstores.VectorStore) error {
	var docs []schema.Document

	err := fs.WalkDir(knowledge, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		log.Printf("Ingesting document: %s\n", path)

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}

		if !strings.HasSuffix(d.Name(), ".txt") {
			return fmt.Errorf("unsupported file type: %s", d.Name())
		}

		fileDocs, err := documentloaders.NewText(file).LoadAndSplit(
			context.Background(),
			textsplitter.NewMarkdownTextSplitter(textsplitter.WithChunkSize(1024), textsplitter.WithChunkOverlap(100)),
		)
		if err != nil {
			return fmt.Errorf("load document (%s): %w", path, err)
		}

		docs = append(docs, fileDocs...)

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk dir: %w", err)
	}

	_, err = store.AddDocuments(context.Background(), docs)
	if err != nil {
		return fmt.Errorf("add documents: %w", err)
	}

	log.Printf("Ingested %d documents\n", len(docs))

	return nil
}

func selectStore(ctx context.Context, embedder embeddings.Embedder) (vectorstores.VectorStore, error) {
	storeTypeEnv := os.Getenv("VECTOR_STORE")

	switch storeTypeEnv {
	case "pgvector":
		return pgvector.NewStore(ctx, embedder)
	default:
		return weaviate.NewStore(ctx, embedder)
	}
}
