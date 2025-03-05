package main

import (
	"context"
	"fmt"
	"github.com/nikolayk812/genai-go/08-testing/ai"
	"log"
	"net/http"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/vectorstores"

	internalhttp "github.com/nikolayk812/genai-go/internal/http"
)

const (
	question  string = "How I can enable verbose logging in Testcontainers Desktop?"
	model     string = "llama3.2"
	tag       string = "3b"
	modelName string = model + ":" + tag
)

func main() {
	ctx := context.Background()

	log.Println(question)

	if err := run(ctx); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run(ctx context.Context) error {
	httpCli := &http.Client{
		Transport: internalhttp.NewLoggingRoundTripper(),
	}

	chatModel, err := buildChatModel(httpCli)
	if err != nil {
		return fmt.Errorf("buildChatModel: %w", err)
	}

	resp, err := straightAnswer(ctx, chatModel)
	if err != nil {
		log.Fatalf("straight chat: %s", err)
	}
	fmt.Println(">> Straight answer:\n", resp)

	resp, err = raggedAnswer(ctx, chatModel)
	if err != nil {
		return fmt.Errorf("ragged chat: %s", err)
	}
	fmt.Println(">> Ragged answer:\n", resp)

	return nil
}

func straightAnswer(ctx context.Context, chatModel llms.Model) (string, error) {
	chatter := ai.NewChat(chatModel)

	return chatter.Chat(ctx, question)
}

func raggedAnswer(ctx context.Context, chatModel llms.Model) (string, error) {
	chatter, err := buildRaggedChat(ctx, chatModel)
	if err != nil {
		return "", fmt.Errorf("build ragged chat: %s", err)
	}

	return chatter.Chat(ctx, question)
}

func buildRaggedChat(ctx context.Context, chatModel llms.Model) (ai.Chatter, error) {
	embeddingModel, err := buildEmbeddingModel()
	if err != nil {
		return nil, fmt.Errorf("buildEmbeddingModel: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(embeddingModel)
	if err != nil {
		return nil, fmt.Errorf("embeddings.NewEmbedder: %w", err)
	}

	store, err := selectStore(ctx, embedder)
	if err != nil {
		return nil, fmt.Errorf("selectStore: %w", err)
	}

	//if err := ingestion(ctx, store); err != nil {
	//	return nil, fmt.Errorf("ingestion: %w", err)
	//}

	// Enrich the response with the relevant documents after the ingestion
	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.60), // use for precision, when you want to get only the most relevant documents
		//vectorstores.WithNameSpace("default"),            // use for set a namespace in the storage
		//vectorstores.WithFilters(map[string]interface{}{"language": "en"}), // use for filter the documents
		vectorstores.WithEmbedder(embedder), // use when you want add documents or doing similarity search
		//vectorstores.WithDeduplicater(vectorstores.NewSimpleDeduplicater()), //  This is useful to prevent wasting time on creating an embedding
	}

	maxResults := 3 // Number of relevant documents to return

	relevantDocs, err := store.SimilaritySearch(ctx, "cloud.logs.verbose", maxResults, optionsVector...)
	if err != nil {
		return nil, fmt.Errorf("similarity search: %w", err)
	}
	log.Printf("Relevant documents for RAG: %d\n", len(relevantDocs))

	return ai.NewChat(chatModel, ai.WithRAGContext(relevantDocs)), nil
}
