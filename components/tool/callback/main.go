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
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/eino-contrib/jsonschema"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	callbackHelper "github.com/cloudwego/eino/utils/callbacks"
	"github.com/cloudwego/eino/callbacks"
)

func main() {
	fmt.Println("=== Callback 功能演示示例 ===")

	// 演示 callback 功能
	demonstrateCallbackFeature()
}

// 演示 callback 功能的完整示例
func demonstrateCallbackFeature() {
	ctx := context.Background()

	// ========== 步骤 1: 创建 Tool（使用 struct 生成 JSON Schema） ==========
	fmt.Println("【步骤 1】创建天气查询工具")
	fmt.Println("----------------------------------------------")

	// 定义参数结构体（带 tag）
	type WeatherParams struct {
		City string `json:"city" jsonschema:"required,description=城市名称，例如：北京、上海"`
		Unit string `json:"unit,omitempty" jsonschema:"description=温度单位：celsius 或 fahrenheit"`
	}

	// 从 struct 生成 JSON Schema
	weatherSchema := jsonschema.Reflect(WeatherParams{})

	toolInfo := &schema.ToolInfo{
		Name:        "get_weather",
		Desc:        "获取指定城市的天气信息",
		ParamsOneOf: schema.NewParamsOneOfByJSONSchema(weatherSchema),
	}

	// 创建 Tool
	weatherTool := utils.NewTool(toolInfo, func(ctx context.Context, params *WeatherParams) (string, error) {
		// 设置默认单位
		unit := params.Unit
		if unit == "" {
			unit = "celsius"
		}

		// 模拟 HTTP 请求
		apiURL := fmt.Sprintf("https://api.weather.com/v1/current?city=%s&unit=%s",
			strings.ReplaceAll(params.City, " ", "%20"), unit)

		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return "", fmt.Errorf("创建请求失败: %w", err)
		}
		req.Header.Set("Authorization", "Bearer mock-api-key")
		req.Header.Set("Content-Type", "application/json")

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

		return string(responseJSON), nil
	})

	fmt.Printf("✓ 天气查询工具创建成功\n")

	// ========== 步骤 2: 创建 Callback Handler ==========
	fmt.Println("\n【步骤 2】创建 Callback Handler")
	fmt.Println("----------------------------------------------")

	// 创建 callback handler
	handler := &callbackHelper.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			fmt.Printf("🚀 开始执行工具，参数: %s\n", input.ArgumentsInJSON)
			if info != nil {
				fmt.Printf("   组件: %s\n", info.Component)
				fmt.Printf("   类型: %s\n", info.Type)
				fmt.Printf("   名称: %s\n", info.Name)
			}
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			fmt.Printf("✅ 工具执行完成，结果: %s\n", output.Response)
			if info != nil {
				fmt.Printf("   组件: %s\n", info.Component)
				fmt.Printf("   类型: %s\n", info.Type)
				fmt.Printf("   名称: %s\n", info.Name)
			}
			return ctx
		},
		OnEndWithStreamOutput: func(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[*tool.CallbackOutput]) context.Context {
			fmt.Println("🌊 工具开始流式输出")
			go func() {
				defer output.Close()

				for {
					chunk, err := output.Recv()
					if errors.Is(err, io.EOF) {
						fmt.Println("📋 流式输出结束")
						return
					}
					if err != nil {
						fmt.Printf("❌ 流式输出错误: %v\n", err)
						return
					}
					fmt.Printf("📦 收到流式输出: %s\n", chunk.Response)
				}
			}()
			return ctx
		},
	}

	// 创建 handler helper
	helper := callbackHelper.NewHandlerHelper().
		Tool(handler).
		Handler()

	fmt.Printf("✓ Callback Handler 创建成功\n")

	// ========== 步骤 3: 创建 ToolsNode ==========
	fmt.Println("\n【步骤 3】创建 ToolsNode")
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

	// ========== 步骤 4: 构造 Mock Input ==========
	fmt.Println("\n【步骤 4】构造 Mock Input")
	fmt.Println("----------------------------------------------")

	// 构造 mock input
	mockInput := &schema.Message{
		Role: schema.Assistant,
		ToolCalls: []schema.ToolCall{
			{
				ID: "call_weather_001",
				Function: schema.FunctionCall{
					Name:      "get_weather",
					Arguments: `{"city": "北京", "unit": "celsius"}`,
				},
			},
		},
	}

	fmt.Printf("✓ Mock Input 构造完成\n")
	fmt.Printf("   工具调用 ID: %s\n", mockInput.ToolCalls[0].ID)
	fmt.Printf("   工具名称: %s\n", mockInput.ToolCalls[0].Function.Name)
	fmt.Printf("   工具参数: %s\n", mockInput.ToolCalls[0].Function.Arguments)

	// ========== 步骤 5: 创建 Chain 并使用 Callback ==========
	fmt.Println("\n【步骤 5】创建 Chain 并使用 Callback")
	fmt.Println("----------------------------------------------")

	// 创建一个简单的 chain: ToolsNode
	chain := compose.NewChain[*schema.Message, []*schema.Message]()
	chain.AppendToolsNode(toolsNode)

	// 编译 chain
	runnable, err := chain.Compile(ctx)
	if err != nil {
		fmt.Printf("✗ Chain 编译失败: %v\n", err)
		return
	}

	fmt.Println("🔧 开始执行 Chain（带 Callback）...")

	// 调用 chain，使用 callback
	output, err := runnable.Invoke(ctx, mockInput, compose.WithCallbacks(helper))
	if err != nil {
		fmt.Printf("✗ Chain 调用失败: %v\n", err)
		return
	}

	// ========== 步骤 6: 展示调用结果 ==========
	fmt.Println("\n【步骤 6】调用结果展示")
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

	// ========== 步骤 7: 演示多次调用的 Callback ==========
	fmt.Println("\n【步骤 7】演示多次工具调用的 Callback")
	fmt.Println("----------------------------------------------")

	// 构造多个工具调用
	multipleCallsInput := &schema.Message{
		Role: schema.Assistant,
		ToolCalls: []schema.ToolCall{
			{
				ID: "call_weather_002",
				Function: schema.FunctionCall{
					Name:      "get_weather",
					Arguments: `{"city": "上海", "unit": "celsius"}`,
				},
			},
			{
				ID: "call_weather_003",
				Function: schema.FunctionCall{
					Name:      "get_weather",
					Arguments: `{"city": "广州"}`,
				},
			},
		},
	}

	fmt.Println("🔧 执行多个工具调用...")
	_, err = runnable.Invoke(ctx, multipleCallsInput, compose.WithCallbacks(helper))
	if err != nil {
		fmt.Printf("✗ 多次调用失败: %v\n", err)
		return
	}

	fmt.Println("\n==============================================")
	fmt.Println("✓ Callback 功能演示完成！")
	fmt.Println("==============================================")
}