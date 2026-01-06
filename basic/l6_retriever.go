package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/joho/godotenv"
)

// Retriever 根据输入从向量数据库检索Document
func Retriever() {
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	// 初始化向量模型
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("EMBEDDING_MODEL"),
	})
	// 初始化milvus client
	InitClient()
	// 定义Collection
	var collection = "AwesomeEino"
	retriever, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      MilvusCli,
		Collection:  collection,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:      1,
		Embedding: embedder,
	})
	if err != nil {
		panic(err)
	}
	results, err := retriever.Retrieve(ctx, "Caleb是干什么的？")
	if err != nil {
		panic(err)
	}
	fmt.Println(results)
	fmt.Println(results[0].ID)
}
