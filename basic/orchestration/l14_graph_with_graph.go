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

// GraphWithGraph 展示Graph作为特殊lambda节点被另一个Graph调用的例子，展示了Graph的嵌套调用能力
func GraphWithGraph() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	// 注册图
	graph := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(GetStateFunc))
	// 创建一个lambda，根据输入判断类型
	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (map[string]string, error) {
		// 处理State
		compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
			state.History["tsundere_action"] = "我喜欢你"
			state.History["cute_action"] = "摸摸头"
			return nil
		})

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

	cutePreHandler := func(ctx context.Context, input map[string]string, state *State) (map[string]string, error) {
		input["content"] = input["content"] + state.History["cute_action"].(string)
		return input, nil
	}

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
	err = graph.AddLambdaNode("cuteLambda", cuteLambda, compose.WithStatePreHandler(cutePreHandler))
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
	input := map[string]string{
		"role":    "cute",
		"content": "今天天气真好啊",
	}

	outsideGraph := compose.NewGraph[map[string]string, string]()
	outerLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (map[string]string, error) {
		// 直接返回输入
		return input, nil
	})
	writeLambda := compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output string, err error) {
		f, err := os.OpenFile("orc_graph_withgraph.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := f.WriteString(input.Content + "\n---\n"); err != nil {
			return "", err
		}
		return "已经写入文件，请前往文件内查看内容", nil
	})
	err = outsideGraph.AddLambdaNode("outerLambda", outerLambda)
	if err != nil {
		panic(err)
	}
	err = outsideGraph.AddGraphNode("inside", graph)
	if err != nil {
		panic(err)
	}
	err = outsideGraph.AddLambdaNode("write", writeLambda)
	if err != nil {
		panic(err)
	}
	err = outsideGraph.AddEdge(compose.START, "outerLambda")
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = outsideGraph.AddEdge("outerLambda", "inside")
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = outsideGraph.AddEdge("inside", "write")
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = outsideGraph.AddEdge("write", compose.END)
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	runnable, err := outsideGraph.Compile(ctx)
	if err != nil {
		log.Fatal("Error compiling graph, ", err)
	}
	answer, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal("Error invoking graph, ", err)
	}
	fmt.Println(answer)
}
