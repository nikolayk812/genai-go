package pgvector

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

// NewStore creates a new PgVector store. It will use a Postgres container with the pgvector module to store the data.
func NewStore(ctx context.Context, embedder embeddings.Embedder) (pgvector.Store, error) {
	conn := mustGetConnection(ctx)

	return pgvector.New(
		ctx,
		pgvector.WithConnectionURL(conn),
		pgvector.WithEmbedder(embedder),
		pgvector.WithVectorDimensions(384),
		pgvector.WithCollectionName(`Testcontainers`),
		pgvector.WithCollectionTableName("tctable"),
	)
}

func mustGetConnection(ctx context.Context) string {
	c, err := tcpostgres.Run(ctx, "pgvector/pgvector:pg16",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name: "pgvector-db",
			},
			Reuse: true,
		},
		))
	if err != nil {
		panic(err)
	}

	conn, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	return conn
}
