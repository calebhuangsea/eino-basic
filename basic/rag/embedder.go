package rag

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
)

func NewArkEmbedder(ctx context.Context) *ark.Embedder {
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("EMBEDDING_MODEL"),
	})
	if err != nil {
		panic(err)
	}
	return embedder
}
