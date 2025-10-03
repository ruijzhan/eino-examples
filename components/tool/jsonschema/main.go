/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/eino-contrib/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	fmt.Println("=== JSON Schema 工具创建与 ToolsNode 调用示例 ===")

	// 完整流程演示
	demonstrateCompleteWorkflow()
}

// 完整工作流程演示：JSON Schema -> ToolInfo -> Tool -> ToolsNode -> Invoke
func demonstrateCompleteWorkflow() {
	ctx := context.Background()

	// ========== 步骤 1: 使用 JSON Schema 创建 ToolInfo ==========
	fmt.Println("【步骤 1】使用 JSON Schema 创建 ToolInfo")
	fmt.Println("----------------------------------------------")

	// 定义 JSON Schema（模拟天气查询 API 的参数）
	weatherSchema := &jsonschema.Schema{
		Type:     string(schema.Object),
		Required: []string{"city"},
		Properties: orderedmap.New[string, *jsonschema.Schema](
			orderedmap.WithInitialData[string, *jsonschema.Schema](
				orderedmap.Pair[string, *jsonschema.Schema]{
					Key: "city",
					Value: &jsonschema.Schema{
						Type:        string(schema.String),
						Description: "城市名称，例如：北京、上海",
					},
				},
				orderedmap.Pair[string, *jsonschema.Schema]{
					Key: "unit",
					Value: &jsonschema.Schema{
						Type:        string(schema.String),
						Description: "温度单位：celsius 或 fahrenheit",
					},
				},
			),
		),
	}

	toolInfo := &schema.ToolInfo{
		Name:        "get_weather",
		Desc:        "获取指定城市的天气信息",
		ParamsOneOf: schema.NewParamsOneOfByJSONSchema(weatherSchema),
	}

	fmt.Printf("✓ 创建的 ToolInfo:\n")
	fmt.Printf("  - 名称: %s\n", toolInfo.Name)
	fmt.Printf("  - 描述: %s\n", toolInfo.Desc)
	fmt.Printf("  - 参数 Schema:\n")
	spew.Dump(toolInfo.ParamsOneOf)

	// ========== 步骤 2: 使用 ToolInfo 创建 Tool（包含 HTTP API 模拟） ==========
	fmt.Println("\n【步骤 2】使用 ToolInfo 创建 Tool（模拟 HTTP API 调用）")
	fmt.Println("----------------------------------------------")

	// 定义参数结构体
	type WeatherParams struct {
		City string `json:"city"`
		Unit string `json:"unit,omitempty"`
	}

	// 创建 Tool，在 InvokableFunc 中模拟 HTTP API 调用
	weatherTool := utils.NewTool(toolInfo, func(ctx context.Context, params *WeatherParams) (string, error) {
		// 设置默认单位
		unit := params.Unit
		if unit == "" {
			unit = "celsius"
		}

		fmt.Printf("\n  → 模拟 HTTP API 调用:\n")
		fmt.Printf("    - 方法: GET\n")
		fmt.Printf("    - URL: https://api.weather.com/v1/current\n")
		fmt.Printf("    - 参数: city=%s, unit=%s\n", params.City, unit)

		// 模拟 HTTP 请求
		apiURL := fmt.Sprintf("https://api.weather.com/v1/current?city=%s&unit=%s",
			strings.ReplaceAll(params.City, " ", "%20"), unit)

		// 创建 HTTP 请求（实际不会发送，仅用于演示）
		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return "", fmt.Errorf("创建请求失败: %w", err)
		}
		req.Header.Set("Authorization", "Bearer mock-api-key")
		req.Header.Set("Content-Type", "application/json")

		fmt.Printf("    - Headers: %v\n", req.Header)

		// 模拟网络延迟
		time.Sleep(100 * time.Millisecond)

		// 模拟 API 响应数据
		mockResponse := map[string]interface{}{
			"city":        params.City,
			"temperature": 22,
			"unit":        unit,
			"condition":   "晴朗",
			"humidity":    65,
			"wind_speed":  12,
			"timestamp":   time.Now().Format(time.RFC3339),
		}

		// 将响应转换为 JSON
		responseJSON, err := json.Marshal(mockResponse)
		if err != nil {
			return "", fmt.Errorf("序列化响应失败: %w", err)
		}

		fmt.Printf("  ← API 响应: %s\n", string(responseJSON))

		return string(responseJSON), nil
	})

	fmt.Printf("✓ Tool 创建成功\n")
	// ========== 步骤 3: 使用 Tool 创建 ToolsNode ==========
	fmt.Println("\n【步骤 3】使用 Tool 创建 ToolsNode")
	fmt.Println("----------------------------------------------")

	// 创建 ToolsNode
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{weatherTool},
	})
	if err != nil {
		fmt.Printf("✗ 创建 ToolsNode 失败: %v\n", err)
		return
	}
	fmt.Printf("✓ ToolsNode 创建成功\n")
	fmt.Printf("  - 包含工具数量: 1\n")
	fmt.Printf("  - 工具名称: %s\n", toolInfo.Name)

	// ========== 步骤 4: Mock Input 并调用 ToolsNode ==========
	fmt.Println("\n【步骤 4】Mock Input 并使用 Invoke 调用 ToolsNode")
	fmt.Println("----------------------------------------------")

	// 构造 mock input（模拟 LLM 返回的工具调用请求）
	mockInput := &schema.Message{
		Role: schema.Assistant,
		ToolCalls: []schema.ToolCall{
			{
				ID: "call_001",
				Function: schema.FunctionCall{
					Name:      "get_weather",
					Arguments: `{"city": "北京", "unit": "celsius"}`,
				},
			},
		},
	}

	fmt.Printf("✓ Mock Input 构造完成:\n")
	fmt.Printf("  - 角色: %s\n", mockInput.Role)
	fmt.Printf("  - 工具调用 ID: %s\n", mockInput.ToolCalls[0].ID)
	fmt.Printf("  - 工具名称: %s\n", mockInput.ToolCalls[0].Function.Name)
	fmt.Printf("  - 工具参数: %s\n", mockInput.ToolCalls[0].Function.Arguments)

	// 调用 ToolsNode
	fmt.Println("\n  开始调用 ToolsNode...")
	output, err := toolsNode.Invoke(ctx, mockInput)
	if err != nil {
		fmt.Printf("✗ ToolsNode 调用失败: %v\n", err)
		return
	}

	// ========== 步骤 5: 展示调用结果 ==========
	fmt.Println("\n【步骤 5】ToolsNode 调用结果")
	fmt.Println("----------------------------------------------")

	fmt.Printf("✓ 调用成功！\n")
	fmt.Printf("  - 返回消息数量: %d\n", len(output))

	for i, msg := range output {
		fmt.Printf("\n  消息 #%d:\n", i+1)
		fmt.Printf("    - 角色: %s\n", msg.Role)
		fmt.Printf("    - 工具调用 ID: %s\n", msg.ToolCallID)
		fmt.Printf("    - 工具名称: %s\n", msg.Name)
		fmt.Printf("    - 工具返回内容:\n")

		// 格式化输出 JSON
		var prettyJSON map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Content), &prettyJSON); err == nil {
			prettyBytes, _ := json.MarshalIndent(prettyJSON, "      ", "  ")
			fmt.Printf("      %s\n", string(prettyBytes))
		} else {
			fmt.Printf("      %s\n", msg.Content)
		}
	}

	fmt.Println("\n==============================================")
	fmt.Println("✓ 完整工作流程演示完成！")
	fmt.Println("==============================================")
}
