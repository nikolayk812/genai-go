package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/nikolayk812/genai-go/08-testing/weaviate"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
)

//go:embed knowledge
var knowledge embed.FS

func ingestion(ctx context.Context, store vectorstores.VectorStore) error {
	var docs []schema.Document

	err := fs.WalkDir(knowledge, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		log.Printf("Ingesting document: %s\n", path)

		file, err := knowledge.Open(path)
		if err != nil {
			return fmt.Errorf("os.Open[%s]: %w", path, err)
		}
		defer file.Close()

		if !strings.HasSuffix(d.Name(), ".txt") {
			return fmt.Errorf("unsupported file type path[%s]: %s", path, d.Name())
		}

		splitter := textsplitter.NewMarkdownTextSplitter(
			textsplitter.WithChunkSize(1024),
			textsplitter.WithChunkOverlap(100))

		fileDocs, err := documentloaders.NewText(file).LoadAndSplit(ctx, splitter)
		if err != nil {
			return fmt.Errorf("load document [%s]: %w", path, err)
		}

		docs = append(docs, fileDocs...)

		return nil
	})
	if err != nil {
		return fmt.Errorf("fs.WalkDir: %w", err)
	}

	if _, err := store.AddDocuments(ctx, docs); err != nil {
		return fmt.Errorf("store.AddDocuments: %w", err)
	}

	log.Printf("Ingested %d documents\n", len(docs))

	return nil
}

func selectStore(ctx context.Context, embedder embeddings.Embedder) (vectorstores.VectorStore, error) {
	storeTypeEnv := os.Getenv("VECTOR_STORE")

	switch storeTypeEnv {
	//case "pgvector":
	//	return pgvector.NewStore(ctx, embedder)
	default:
		return weaviate.NewStore(embedder)
	}
}
