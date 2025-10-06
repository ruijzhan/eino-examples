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
	"fmt"
	"strings"

	"github.com/cloudwego/eino/adk"

	"github.com/cloudwego/eino-examples/components/document/parser/common"
	"github.com/cloudwego/eino-examples/internal/logs"
)

func main() {
	// 基础解析示例
	basicParsingExample()

	// Agent 集成示例
	agentIntegrationExample()
}

// basicParsingExample 演示基础的自定义解析功能
func basicParsingExample() {
	ctx := context.Background()

	customParser, err := NewCustomParser(&Config{
		DefaultEncoding: "default",
		DefaultMaxSize:  1024,
	})
	if err != nil {
		logs.Errorf("NewCustomParser failed, err=%v", err)
		return
	}

	docs, err := customParser.Parse(ctx, nil,
		WithMaxSize(2048),
	)
	if err != nil {
		logs.Errorf("customParser.Parse, err=%v", err)
		return
	}

	for idx, doc := range docs {
		logs.Infof("doc_%v content: %v", idx, doc.Content)
	}
}

// agentIntegrationExample 演示如何将自定义解析器与 Agent 集成
func agentIntegrationExample() {
	// 定义自定义内容处理函数
	customContentProcessor := func(agent *adk.ChatModelAgent, input string) (string, error) {
		// 使用自定义解析器解析内容
		customParser := createCustomParser()
		docs, err := customParser.Parse(context.Background(), strings.NewReader(input),
			WithEncoding("utf-8"),
			WithMaxSize(4096),
		)
		if err != nil {
			return "", fmt.Errorf("failed to parse custom content: %w", err)
		}

		logs.Infof("Custom parsed content: %s", docs[0].Content)

		// 使用 ChatModelAgent 处理解析后的内容
		return common.ProcessWithAgent(agent, docs[0].Content)
	}

	// 运行 Agent 集成示例
	common.RunAgentIntegrationExample("CustomContentProcessor",
		"这是一个测试内容，用于演示自定义解析器和 Agent 的集成。",
		func(agent *adk.ChatModelAgent, content string) (string, error) {
			return customContentProcessor(agent, content)
		})
}

// CustomAgent 集成了自定义解析器和 Agent 的结构体（向后兼容）
type CustomAgent struct {
	customParser *CustomParser
	agent        *common.BaseAgent
}

// NewCustomAgent 创建新的自定义内容处理器
func NewCustomAgent(agentType string) *CustomAgent {
	return &CustomAgent{
		customParser: createCustomParser(),
		agent:        common.NewBaseAgent(agentType),
	}
}

// ProcessCustomContent 使用自定义解析器处理内容并通过 Agent 分析（使用共享代码）
func (ca *CustomAgent) ProcessCustomContent(ctx context.Context, input string) (string, error) {
	// 使用自定义解析器解析内容
	docs, err := ca.customParser.Parse(ctx, strings.NewReader(input),
		WithEncoding("utf-8"),
		WithMaxSize(4096),
	)
	if err != nil {
		return "", fmt.Errorf("failed to parse custom content: %w", err)
	}

	logs.Infof("Custom parsed content: %s", docs[0].Content)

	// 使用共享的 Agent 处理函数
	return ca.agent.ProcessContent(docs[0].Content)
}

// createCustomParser 创建自定义解析器
func createCustomParser() *CustomParser {
	customParser, err := NewCustomParser(&Config{
		DefaultEncoding: "utf-8",
		DefaultMaxSize:  4096,
	})
	if err != nil {
		logs.Errorf("NewCustomParser failed, err=%v", err)
		return nil
	}

	return customParser
}
