package main

import (
	"context"
	"fmt"
	internalhttp "github.com/nikolayk812/genai-go/internal/http"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
	"log"
	"net/http"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("run: %s", err)
	}
}

func run(ctx context.Context) error {
	httpCli := &http.Client{
		Transport: internalhttp.NewLoggingRoundTripper(),
	}

	embeddingLLM, err := buildEmbeddingModel(httpCli)
	if err != nil {
		return fmt.Errorf("buildEmbeddingModel: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(embeddingLLM)
	if err != nil {
		return fmt.Errorf("embeddings.NewEmbedder: %w", err)
	}

	store, err := buildEmbeddingStore(embedder)
	if err != nil {
		return fmt.Errorf("buildEmbeddingStore: %w", err)
	}

	if err := ingestion(ctx, store); err != nil {
		return fmt.Errorf("ingestion: %w", err)
	}

	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.70), // use for precision, when you want to get only the most relevant documents
		//vectorstores.WithNameSpace(""),            // use for set a namespace in the storage
		//vectorstores.WithFilters(map[string]interface{}{"language": "en"}), // use for filter the documents
		// vectorstores.WithEmbedder(embedder), // use when you want add documents or doing similarity search
		//vectorstores.WithDeduplicater(vectorstores.NewSimpleDeduplicater()), //  This is useful to prevent wasting time on creating an embedding
	}

	originalQuestion := "What is my favorite sport?"

	relevantDocs, err := store.SimilaritySearch(ctx, originalQuestion, 1, optionsVector...)
	if err != nil {
		return fmt.Errorf("store.SimilaritySearch: %w", err)
	}

	if len(relevantDocs) == 0 {
		fmt.Println("No relevant content found")
		return nil
	}

	chatLLM, err := buildChatModel(httpCli)
	if err != nil {
		return fmt.Errorf("build chat model: %w", err)
	}

	raggedQuestion := fmt.Sprintf(`
%s

Answer the question considering the following relevant content, be very confident:

%s
`, originalQuestion, relevantDocs[0].PageContent)

	fmt.Println(raggedQuestion)

	raggedContent := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, raggedQuestion),
	}

	if _, err := chatLLM.GenerateContent(ctx, raggedContent,
		llms.WithTemperature(0.0001),
		llms.WithTopK(1),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	); err != nil {
		return fmt.Errorf("chatLLM.GenerateContent: %w", err)
	}

	return nil
}

func buildChatModel(httpCli *http.Client) (llms.Model, error) {
	llm, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithHTTPClient(httpCli),
	)
	if err != nil {
		return nil, fmt.Errorf("ollama.New: %w", err)
	}

	return llm, nil
}

func buildEmbeddingModel(httpCli *http.Client) (embeddings.EmbedderClient, error) {
	llm, err := ollama.New(
		ollama.WithModel("nomic-embed-text:v1.5"),
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithHTTPClient(httpCli),
	)
	if err != nil {
		return nil, fmt.Errorf("ollama.New: %w", err)
	}

	return llm, nil
}

func buildEmbeddingStore(embedder embeddings.Embedder) (vectorstores.VectorStore, error) {

	//docker run -d --name chroma \
	//-p 8000:8000 \
	//--rm \
	//chromadb/chroma:0.5.2

	//return chroma.New(
	//	chroma.WithEmbedder(embedder),
	//	chroma.WithChromaURL("http://localhost:8000"),
	//	chroma.WithNameSpace("Testcontainers"),
	//	chroma.WithDistanceFunction(types.COSINE),
	//)

	//docker run -d --name weaviate \
	//-e QUERY_DEFAULTS=vector \
	//-e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \
	//-e WEAVIATE_DEFAULT_CONTEXTIONARY_LANG=en \
	//-p 8080:8080 \
	//--rm \
	//semitechnologies/weaviate:1.25.33

	return weaviate.New(
		weaviate.WithScheme("http"),
		weaviate.WithHost("localhost:8080"),
		weaviate.WithIndexName("Testcontainers"),
		weaviate.WithEmbedder(embedder),
	)
}

func ingestion(ctx context.Context, store vectorstores.VectorStore) error {
	docs := []schema.Document{
		{
			PageContent: "I like football",
		},
		{
			PageContent: "The weather is good today.",
		},
	}

	if _, err := store.AddDocuments(ctx, docs); err != nil {
		return fmt.Errorf("store.AddDocuments: %w", err)
	}

	return nil
}
