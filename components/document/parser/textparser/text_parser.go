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
	"github.com/cloudwego/eino/components/document/parser"

	"github.com/cloudwego/eino-examples/components/document/parser/common"
	"github.com/cloudwego/eino-examples/internal/logs"
)

func main() {
	// 基础解析示例
	basicParsingExample()

	// Agent 集成示例
	agentIntegrationExample()
}

// basicParsingExample 演示基础的文本解析功能
func basicParsingExample() {
	ctx := context.Background()

	textParser := parser.TextParser{}
	docs, err := textParser.Parse(ctx, strings.NewReader("hello world"))
	if err != nil {
		logs.Errorf("TextParser{}.Parse failed, err=%v", err)
		return
	}

	logs.Infof("text content: %v", docs[0].Content)
}

// agentIntegrationExample 演示如何将解析结果传递给 Agent 处理
func agentIntegrationExample() {
	// 定义文本处理函数
	textProcessor := func(agent *adk.ChatModelAgent, input string) (string, error) {
		ctx := context.Background()

		// 使用 TextParser 解析文本
		textParser := parser.TextParser{}
		docs, err := textParser.Parse(ctx, strings.NewReader(input))
		if err != nil {
			return "", fmt.Errorf("failed to parse text: %w", err)
		}

		logs.Infof("Parsed document content: %s", docs[0].Content)

		// 直接使用 ChatModelAgent 处理解析后的内容
		return common.ProcessWithAgent(agent, docs[0].Content)
	}

	// 运行 Agent 集成示例
	common.RunAgentIntegrationExample("TextProcessor",
		"请分析这段文本内容：人工智能正在改变我们的生活方式。",
		func(agent *adk.ChatModelAgent, content string) (string, error) {
			return textProcessor(agent, content)
		})
}
