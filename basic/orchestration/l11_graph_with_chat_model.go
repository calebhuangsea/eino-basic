package orchestration

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func GraphWithChatModel() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	// 注册图
	graph := compose.NewGraph[map[string]string, *schema.Message]()
	// 创建一个lambda，根据输入判断类型
	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (map[string]string, error) {
		if input["role"] == "tsundere" {
			return map[string]string{"role": "傲娇", "content": input["content"]}, nil
		} else if input["role"] == "cute" {
			return map[string]string{"role": "可爱", "content": input["content"]}, nil
		}
		return map[string]string{"role": "user", "content": input["content"]}, nil
	})
	cuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) ([]*schema.Message, error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的小女孩，用可爱的语气回复我",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	tsundereLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) ([]*schema.Message, error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个傲娇的小姐姐，用傲娇的语气回复我",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	defaultLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) ([]*schema.Message, error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "正常回复",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		log.Fatal("Error initializing chat model, ", err)
	}
	// 加入节点
	err = graph.AddLambdaNode("lambda", lambda)
	if err != nil {
		log.Fatal("Error adding lambda node, ", err)
	}
	err = graph.AddLambdaNode("cuteLambda", cuteLambda)
	if err != nil {
		log.Fatal("Error adding cuteLambda node, ", err)
	}
	err = graph.AddLambdaNode("tsundereLambda", tsundereLambda)
	if err != nil {
		log.Fatal("Error adding tsundereLambda node, ", err)
	}
	err = graph.AddLambdaNode("defaultLambda", defaultLambda)
	if err != nil {
		log.Fatal("Error adding defaultLambda node, ", err)
	}
	err = graph.AddChatModelNode("model", model)
	if err != nil {
		log.Fatal("Error adding model node, ", err)
	}
	// 加入分支
	err = graph.AddBranch("lambda", compose.NewGraphBranch(func(ctx context.Context, input map[string]string) (string, error) {
		if input["role"] == "傲娇" {
			return "tsundereLambda", nil
		} else if input["role"] == "可爱" {
			return "cuteLambda", nil
		}
		return "defaultLambda", nil
	}, map[string]bool{"cuteLambda": true, "tsundereLambda": true, "defaultLambda": true}))
	// 连接节点
	err = graph.AddEdge(compose.START, "lambda")
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = graph.AddEdge("cuteLambda", "model")
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = graph.AddEdge("tsundereLambda", "model")
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = graph.AddEdge("defaultLambda", "model")
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = graph.AddEdge("model", compose.END)
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	// 编译
	runnable, err := graph.Compile(ctx)
	if err != nil {
		log.Fatal("Error compiling graph, ", err)
	}
	input := map[string]string{
		"role":    "?",
		"content": "今天天气真好啊",
	}
	answer, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal("Error invoking graph, ", err)
	}
	fmt.Println(answer)
}
