package weaviate

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	tcweaviate "github.com/testcontainers/testcontainers-go/modules/weaviate"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
)

func NewStore(ctx context.Context, embedder embeddings.Embedder) (weaviate.Store, error) {
	schema, host := mustGetAddress(ctx)

	return weaviate.New(
		weaviate.WithScheme(schema),
		weaviate.WithHost(host),
		weaviate.WithIndexName("Testcontainers"),
		weaviate.WithEmbedder(embedder),
	)
}

func mustGetAddress(ctx context.Context) (string, string) {
	c, err := tcweaviate.Run(ctx, "semitechnologies/weaviate:1.27.2", testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name: "weaviate-db",
		},
		Reuse: true,
	}))
	if err != nil {
		panic(err)
	}

	schema, host, err := c.HttpHostAddress(ctx)
	if err != nil {
		panic(err)
	}

	return schema, host
}
