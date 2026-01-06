package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

// Stream 流式问答写法
func Stream() {
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
	// 生成回复，获取流式接收器
	reader, err := model.Stream(ctx, input)
	defer reader.Close()
	if err != nil {
		log.Fatal("Error streaming response, ", err)
	}
	for {
		chunk, err := reader.Recv()
		if err != nil {
			break
		}
		// 注意chunk有可能为空，如果用println会有奇怪表现
		print(chunk.Content)
	}
}
