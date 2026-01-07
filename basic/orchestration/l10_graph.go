package orchestration

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/compose"
	"github.com/joho/godotenv"
)

// Graph 展示如何构建一个简单的图，包含一个起始分支，以及三个lambda节点，根据输入不同，进入不同的分支，最后输出结果
func Graph() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		// 处理环境变量加载异常
		log.Fatal("Error loading .env file, ", err)
	}
	ctx := context.Background()
	// 注册图
	graph := compose.NewGraph[string, string]()
	// 创建一个lambda，根据输入判断类型
	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		if input == "1" {
			return "豪猫", nil
		} else if input == "2" {
			return "耄耋", nil
		} else if input == "3" {
			return "device", nil
		}
		return "", nil
	})
	// 三个分支的lambda
	lambda0 := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return "喵～", nil
	})
	lambda1 := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return "哈！", nil
	})
	lambda2 := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return "没有人类了～", nil
	})
	// 加入节点
	err = graph.AddLambdaNode("lambda", lambda)
	if err != nil {
		log.Fatal("Error adding lambda node, ", err)
	}
	err = graph.AddLambdaNode("lambda0", lambda0)
	if err != nil {
		log.Fatal("Error adding lambda0 node, ", err)
	}
	err = graph.AddLambdaNode("lambda1", lambda1)
	if err != nil {
		log.Fatal("Error adding lambda1 node, ", err)
	}
	err = graph.AddLambdaNode("lambda2", lambda2)
	if err != nil {
		log.Fatal("Error adding lambda2 node, ", err)
	}
	// 加入分支
	err = graph.AddBranch("lambda", compose.NewGraphBranch(func(ctx context.Context, input string) (string, error) {
		if input == "豪猫" {
			return "lambda0", nil
		} else if input == "耄耋" {
			return "lambda1", nil
		} else if input == "device" {
			return "lambda2", nil
		}
		return compose.END, nil
	}, map[string]bool{"lambda0": true, "lambda1": true, "lambda2": true}))
	if err != nil {
		log.Fatal("Error adding lambda2 node, ", err)
	}
	// 连接节点
	err = graph.AddEdge(compose.START, "lambda0")
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = graph.AddEdge("lambda0", compose.END)
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = graph.AddEdge("lambda1", compose.END)
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	err = graph.AddEdge("lambda2", compose.END)
	if err != nil {
		log.Fatal("Error adding edge, ", err)
	}
	// 编译
	runnable, err := graph.Compile(ctx)
	if err != nil {
		log.Fatal("Error compiling graph, ", err)
	}
	answer, err := runnable.Invoke(ctx, "1")
	if err != nil {
		log.Fatal("Error invoking graph, ", err)
	}
	fmt.Println(answer)
}
