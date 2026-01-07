package orchestration

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	callbackHelpers "github.com/cloudwego/eino/utils/callbacks"
	"github.com/joho/godotenv"
)

func Chain() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	// 初始化模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		log.Fatal("Error initializing chat model, ", err)
	}
	// 创建一个lambda，增加固定结尾
	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		suffix := input + "【我是固定结尾】"
		output = []*schema.Message{
			{
				Role:    schema.User,
				Content: suffix,
			},
		}
		return
	})
	// 注册chain
	chain := compose.NewChain[string, *schema.Message]()
	chain.AppendLambda(lambda)
	chain.AppendChatModel(model)
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatal("Error compiling chain, ", err)
	}
	// 运行chain
	response, err := runnable.Invoke(ctx, "这是输入")
	if err != nil {
		log.Fatal("Error running chain, ", err)
	}
	log.Println(response.Content)
}

func SimpleAgent() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	getGameTool := CreateTool()
	//大模型回调函数
	modelHandler := &callbackHelpers.ModelCallbackHandler{
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *model.CallbackOutput) context.Context {
			// 1. output.Result 类型是 string
			fmt.Println("模型思考过程为：")
			fmt.Println(output.Message.ReasoningContent)
			return ctx
		},
	}
	//工具回调函数
	toolHandler := &callbackHelpers.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			fmt.Printf("开始执行工具，参数: %s\n", input.ArgumentsInJSON)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			fmt.Printf("工具执行完成，结果: %s\n", output.Response)
			return ctx
		},
	}
	//构建实际回调函数Handler
	handler := callbackHelpers.NewHandlerHelper().
		ChatModel(modelHandler).
		Tool(toolHandler).
		Handler()
	// 初始化模型
	timeout := 30 * time.Second
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("MODEL"),
		Timeout: &timeout,
	})
	if err != nil {
		log.Fatal("Error initializing chat model, ", err)
	}
	// 绑定工具
	info, err := getGameTool.Info(ctx)
	if err != nil {
		panic(err)
	}
	infos := []*schema.ToolInfo{
		info,
	}
	err = model.BindTools(infos)
	if err != nil {
		panic(err)
	}
	// 创建tools节点
	//创建tools节点
	ToolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			getGameTool,
		},
	})
	if err != nil {
		panic(err)
	}
	//创建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(model, compose.WithNodeName("chat_model")).
		AppendToolsNode(ToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//运行Agent
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "请告诉我王者荣耀的URL是什么",
		},
	}, compose.WithCallbacks(handler))
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	fmt.Println("打印结果：")
	for _, msg := range resp {
		fmt.Println(msg.Content)
	}
}
