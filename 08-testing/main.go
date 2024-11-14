package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/mdelapenya/genai-testcontainers-go/testing/ai"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

const (
	question  string = "How I can enable verbose logging in Testcontainers Desktop?"
	model     string = "llama3.2"
	tag       string = "3b"
	modelName string = model + ":" + tag
)

//go:embed knowledge
var knowledge embed.FS

func main() {
	log.Println(question)
	if err := run(); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run() error {
	chatModel, err := buildChatModel()
	if err != nil {
		return fmt.Errorf("build chat model: %s", err)
	}

	resp, err := straightAnswer(chatModel)
	if err != nil {
		log.Fatalf("straight chat: %s", err)
	}
	fmt.Println(">> Straight answer:\n", resp)

	resp, err = raggedAnswer(chatModel)
	if err != nil {
		return fmt.Errorf("ragged chat: %s", err)
	}
	fmt.Println(">> Ragged answer:\n", resp)

	return nil
}

func straightAnswer(chatModel *ollama.LLM) (string, error) {
	chatter := buildStraightChat(chatModel)

	return chatter.Chat(question)
}

func raggedAnswer(chatModel *ollama.LLM) (string, error) {
	chatter, err := buildRaggedChat(chatModel)
	if err != nil {
		return "", fmt.Errorf("build ragged chat: %s", err)
	}

	return chatter.Chat(question)
}

func buildStraightChat(chatModel llms.Model) ai.Chatter {
	return ai.New(chatModel)
}

func buildRaggedChat(chatModel llms.Model) (ai.Chatter, error) {
	embeddingModel, err := buildEmbeddingModel()
	if err != nil {
		return nil, fmt.Errorf("build embedding model: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(embeddingModel)
	if err != nil {
		return nil, fmt.Errorf("new embedder: %w", err)
	}

	store, err := selectStore(context.Background(), embedder)
	if err != nil {
		return nil, fmt.Errorf("new store: %w", err)
	}

	if err := ingestion(store); err != nil {
		return nil, fmt.Errorf("ingestion: %w", err)
	}

	// Enrich the response with the relevant documents after the ingestion
	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.60), // use for precision, when you want to get only the most relevant documents
		//vectorstores.WithNameSpace("default"),            // use for set a namespace in the storage
		//vectorstores.WithFilters(map[string]interface{}{"language": "en"}), // use for filter the documents
		vectorstores.WithEmbedder(embedder), // use when you want add documents or doing similarity search
		//vectorstores.WithDeduplicater(vectorstores.NewSimpleDeduplicater()), //  This is useful to prevent wasting time on creating an embedding
	}

	maxResults := 3 // Number of relevant documents to return

	relevantDocs, err := store.SimilaritySearch(context.Background(), "cloud.logs.verbose", maxResults, optionsVector...)
	if err != nil {
		return nil, fmt.Errorf("similarity search: %w", err)
	}
	log.Printf("Relevant documents for RAG: %d\n", len(relevantDocs))

	return ai.New(chatModel, ai.WithRAGContext(relevantDocs)), nil
}

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

		var loader documentloaders.Loader

		if strings.HasSuffix(d.Name(), ".txt") {
			loader = documentloaders.NewText(file)
		} else {
			return fmt.Errorf("unsupported file type: %s", d.Name())
		}

		doc, err := loader.Load(context.Background())
		if err != nil {
			return fmt.Errorf("load document (%s): %w", path, err)
		}

		docs = append(docs, doc...)

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
