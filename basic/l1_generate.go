package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

// Generate 单次问答写法
func Generate() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	timeout := 30 * time.Second
	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("MODEL"),
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatal("Error initializing chat model, ", err)
	}
	// 准备信息
	input := []*schema.Message{
		schema.SystemMessage("你是一个Go语言大师"),
		schema.UserMessage("你好"),
	}
	// 生成回复
	response, err := model.Generate(ctx, input)
	if err != nil {
		log.Fatal("Error generating response, ", err)
	}
	fmt.Println(response.Content)
	// 获取 Token 使用情况
	if usage := response.ResponseMeta.Usage; usage != nil {
		fmt.Println("提示 Tokens:", usage.PromptTokens)
		fmt.Println("生成 Tokens:", usage.CompletionTokens)
		fmt.Println("总 Tokens:", usage.TotalTokens)
	}
}
