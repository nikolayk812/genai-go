package main

import (
	"context"
	"fmt"
	"log"

	tcollama "github.com/testcontainers/testcontainers-go/modules/ollama"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"

	"github.com/mdelapenya/genai-testcontainers-go/rag/weaviate"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run() error {
	embeddingLLM, err := buildEmbeddingModel()
	if err != nil {
		return fmt.Errorf("build embedding model: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(embeddingLLM)
	if err != nil {
		return fmt.Errorf("new embedder: %w", err)
	}

	store, err := buildEmbeddingStore(embedder)
	if err != nil {
		return fmt.Errorf("build embedding store: %w", err)
	}

	if err := ingestion(store); err != nil {
		return fmt.Errorf("ingestion: %w", err)
	}

	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.80), // use for precision, when you want to get only the most relevant documents
		//vectorstores.WithNameSpace(""),            // use for set a namespace in the storage
		//vectorstores.WithFilters(map[string]interface{}{"language": "en"}), // use for filter the documents
		vectorstores.WithEmbedder(embedder), // use when you want add documents or doing similarity search
		//vectorstores.WithDeduplicater(vectorstores.NewSimpleDeduplicater()), //  This is useful to prevent wasting time on creating an embedding
	}

	relevantDocs, err := store.SimilaritySearch(context.Background(), "What is my favorite sport?", 1, optionsVector...)
	if err != nil {
		return fmt.Errorf("similarity search: %w", err)
	}

	if len(relevantDocs) == 0 {
		fmt.Println("No relevant content found")
		return nil
	}

	chatLLM, err := buildChatModel()
	if err != nil {
		return fmt.Errorf("build chat model: %w", err)
	}

	response := fmt.Sprintf(`
What is your favourite sport?

Answer the question considering the following relevant content:
%s
`, relevantDocs[0].PageContent)

	fmt.Println(response)

	ctx := context.Background()
	originalContent := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, response),
	}

	_, err = chatLLM.GenerateContent(
		ctx, originalContent,
		llms.WithTemperature(0.0001),
		llms.WithTopK(1),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("llm generate original content: %w", err)
	}

	return nil
}

func buildChatModel() (*ollama.LLM, error) {
	c, err := tcollama.Run(context.Background(), "ilopezluna/llama3.2:0.3.13-1b")
	if err != nil {
		return nil, err
	}

	ollamaURL, err := c.ConnectionString(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connection string: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel("llama3.2:1b"), ollama.WithServerURL(ollamaURL))
	if err != nil {
		return nil, fmt.Errorf("ollama new: %w", err)
	}

	return llm, nil
}

func buildEmbeddingModel() (*ollama.LLM, error) {
	c, err := tcollama.Run(context.Background(), "ilopezluna/all-minilm:0.3.13-22m")
	if err != nil {
		return nil, err
	}

	ollamaURL, err := c.ConnectionString(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connection string: %w", err)
	}

	llm, err := ollama.New(ollama.WithModel("all-minilm:22m"), ollama.WithServerURL(ollamaURL))
	if err != nil {
		return nil, fmt.Errorf("ollama new: %w", err)
	}

	return llm, nil
}

func buildEmbeddingStore(embedder embeddings.Embedder) (vectorstores.VectorStore, error) {
	store, err := weaviate.NewStore(context.Background(), embedder)
	if err != nil {
		return nil, fmt.Errorf("weaviate new store: %w", err)
	}

	return store, nil
}

func ingestion(store vectorstores.VectorStore) error {
	docs := []schema.Document{
		{
			PageContent: "I like football",
		},
		{
			PageContent: "The weather is good today.",
		},
	}

	_, err := store.AddDocuments(context.Background(), docs)
	if err != nil {
		return fmt.Errorf("add documents: %w", err)
	}

	return nil
}
