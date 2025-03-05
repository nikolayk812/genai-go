package weaviate

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
)

// NewStore creates a new Weaviate store. It will use a weaviate container to store the data.
func NewStore(embedder embeddings.Embedder) (weaviate.Store, error) {
	/*
		docker run -d --name weaviate \
		-e QUERY_DEFAULTS=vector \
		-e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \
		-e WEAVIATE_DEFAULT_CONTEXTIONARY_LANG=en \
		-p 8080:8080 \
		--rm \
		semitechnologies/weaviate:1.25.33
	*/

	return weaviate.New(
		weaviate.WithScheme("http"),
		weaviate.WithHost("localhost:8080"),
		weaviate.WithIndexName("Testcontainers"),
		weaviate.WithEmbedder(embedder),
	)
}
