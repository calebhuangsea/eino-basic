package rag

import (
	"context"
	"log"

	"github.com/joho/godotenv"
)

// Rag 用组件实现整个Rag流程
func Rag() {
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	embedder := NewArkEmbedder(ctx)
	InitClient()
	// 先存储文档的向量
	//indexer := NewArkIndexer(ctx, embedder)
	//splitter := NewTransformer(ctx)
	//// 打开文件
	//bs, err := os.ReadFile("basic/document.md")
	//if err != nil {
	//	panic(err)
	//}
	//docs := []*schema.Document{
	//	{
	//		ID:      "doc1",
	//		Content: string(bs),
	//	},
	//}
	//results, err := splitter.Transform(ctx, docs)
	//if err != nil {
	//	panic(err)
	//}
	//for i, doc := range results {
	//	doc.ID = docs[0].ID + "_" + strconv.Itoa(i)
	//	fmt.Println(doc.ID)
	//}
	//ids, err := indexer.Store(ctx, results)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(ids)
	// 查询
	retriever := NewArkRetriever(ctx, embedder)
	results, err := retriever.Retrieve(ctx, "1996 年 1 月 11 日是什么日子")
	if err != nil {
		panic(err)
	}
	for _, doc := range results {
		println(doc.ID, doc.Content)
	}
}
