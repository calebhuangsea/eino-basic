package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/joho/godotenv"
)

// Retrieval-Augmented Generation 将input strings使用向量模型转换成向量
func Rag() {
	// 加载环境变量
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
	if err != nil {
		panic(err)
	}
	input := []string{
		"你好",
		"测试",
		"测试2",
	}
	embeddings, err := embedder.EmbedStrings(ctx, input)
	if err != nil {
		panic(err)
	}
	for i, embedding := range embeddings {
		// 维度都是2560，因为这个向量模型的唯独是2560
		fmt.Println("文本", i+1, "的向量维度：", len(embedding))
	}
}
