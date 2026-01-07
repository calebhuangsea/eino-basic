package rag

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

func NewArkModel(ctx context.Context) *ark.ChatModel {
	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		log.Fatal("Error initializing chat model, ", err)
	}
	return model
}
