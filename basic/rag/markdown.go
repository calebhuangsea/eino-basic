package rag

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
)

func NewTransformer(ctx context.Context) document.Transformer {
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
	return splitter
}
