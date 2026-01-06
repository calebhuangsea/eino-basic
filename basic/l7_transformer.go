package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

// Transformer 处理传入的Markdown文件
func Transformer() {
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	// 初始化Markdown Header Splitter
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "H1",
			"##":  "H2",
			"###": "H3",
		},
		TrimHeaders: false,
	})
	if err != nil {
		panic(err)
	}
	// 打开文件
	bs, err := os.ReadFile("basic/document.md")
	if err != nil {
		panic(err)
	}
	docs := []*schema.Document{
		{
			ID:      "doc1",
			Content: string(bs),
		},
	}
	result, err := splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}
	for _, doc := range result {
		log.Printf("Document ID: %s\nContent: %s\n", doc.ID, doc.Content)
	}
}
