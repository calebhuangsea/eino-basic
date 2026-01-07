package orchestration

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)
import "github.com/cloudwego/eino-ext/components/tool/browseruse"

// Tool 允许大模型调用工具
func Tool() {
	ctx := context.Background()
	tool, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		panic(err)
	}
	url := "https://www.bilibili.com"
	result, err := tool.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	time.Sleep(10 * time.Second)
	tool.Cleanup()
}

// 创建一个自定义Tool
type Game struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type InputParam struct {
	Name string `json:"name" jsonschema:"description=name of the game"`
}

func GetGame(_ context.Context, param *InputParam) (string, error) {
	GameSet := []Game{
		{"王者荣耀", "https://pvp.qq.com/"},
		{"英雄联盟", "https://www.leagueoflegends.com/"},
		{"和平精英", "https://gp.qq.com/"},
	}
	for _, game := range GameSet {
		if game.Name == param.Name {
			return game.Url, nil
		}
	}
	return "", nil
}

func CreateTool() tool.InvokableTool {
	GetGameTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "get_game",
			Desc: "get the url of the game",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"name": &schema.ParameterInfo{
						Type:     schema.String,
						Desc:     "game's name",
						Required: true,
					},
				},
			),
		}, GetGame)
	return GetGameTool
}
