/*
 * Copyright 2025 CloudWeGo Authors
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

package common

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino/adk"

	"github.com/cloudwego/eino-examples/internal/logs"
)

// ProcessWithAgent 通用的 Agent 处理函数
// 接受解析后的内容，通过 Agent 处理并返回结果
func ProcessWithAgent(agent *adk.ChatModelAgent, content string) (string, error) {
	if agent == nil {
		return "", fmt.Errorf("agent is nil")
	}

	ctx := context.Background()

	// 创建 AgentInput
	agentInput := &adk.AgentInput{
		Messages:        []adk.Message{schema.UserMessage(content)},
		EnableStreaming: false,
	}

	// 运行 Agent 并处理结果
	iterator := agent.Run(ctx, agentInput)

	// 获取最终响应
	var result string
	for {
		event, ok := iterator.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			return "", event.Err
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			msg, err := event.Output.MessageOutput.GetMessage()
			if err != nil {
				return "", err
			}
			result = msg.Content
		}
	}

	return result, nil
}

// CreateMockAgentForExample 创建用于示例的模拟 Agent
// 在实际使用中需要配置真实的模型
func CreateMockAgentForExample(agentType string) *adk.ChatModelAgent {
	// 这里应该配置真实的 ChatModelAgent
	// 由于示例代码无法访问真实的 API 密钥，这里返回 nil
	// 在实际项目中，您需要：
	// 1. 配置真实的模型（如 OpenAI GPT、DeepSeek 等）
	// 2. 设置正确的 API 密钥和端点
	// 3. 使用 adk.NewChatModelAgent 创建 Agent

	logs.Infof("Note: In production, configure a real model here:")
	logs.Infof("agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{")

	// 根据不同的 Agent 类型提供相应的配置示例
	switch agentType {
	case "TextProcessor":
		logs.Infof("    Name:        \"TextProcessor\",")
		logs.Infof("    Description: \"Processes text content\",")
		logs.Infof("    Model:       yourChatModel,")
		logs.Infof("    Instruction: \"请分析并总结用户提供的文本内容\",")
	case "DocumentProcessor":
		logs.Infof("    Name:        \"DocumentProcessor\",")
		logs.Infof("    Description: \"Processes and analyzes document content\",")
		logs.Infof("    Model:       yourChatModel,")
		logs.Infof("    Instruction: \"请分析并总结文档内容，提取关键信息\",")
	case "CustomContentProcessor":
		logs.Infof("    Name:        \"CustomContentProcessor\",")
		logs.Infof("    Description: \"Processes custom parsed content\",")
		logs.Infof("    Model:       yourChatModel,")
		logs.Infof("    Instruction: \"请分析自定义解析器处理后的内容，提供有意义的见解\",")
	default:
		logs.Infof("    Name:        \"ContentProcessor\",")
		logs.Infof("    Description: \"Processes content\",")
		logs.Infof("    Model:       yourChatModel,")
		logs.Infof("    Instruction: \"请分析并处理提供的内容\",")
	}

	logs.Infof("})")

	return nil
}

// RunAgentIntegrationExample 运行 Agent 集成示例的通用函数
func RunAgentIntegrationExample(agentType string, content string, processor func(*adk.ChatModelAgent, string) (string, error)) {
	logs.Infof("=== Agent Integration Example ===")

	// 创建模拟的 Agent（在实际使用中需要真实配置）
	agent := CreateMockAgentForExample(agentType)
	if agent == nil {
		logs.Infof("Skip agent integration: no valid model configuration found")
		return
	}

	// 处理内容
	result, err := processor(agent, content)
	if err != nil {
		logs.Errorf("Agent processing failed, err=%v", err)
		return
	}

	logs.Infof("Agent 处理结果: %s", result)
}

// BaseAgent 基础 Agent 结构体，可以被其他解析器嵌入
type BaseAgent struct {
	agent *adk.ChatModelAgent
}

// NewBaseAgent 创建基础 Agent
func NewBaseAgent(agentType string) *BaseAgent {
	return &BaseAgent{
		agent: CreateMockAgentForExample(agentType),
	}
}

// ProcessContent 处理内容的通用方法
func (ba *BaseAgent) ProcessContent(content string) (string, error) {
	if ba.agent == nil {
		return "", fmt.Errorf("agent is not initialized")
	}
	return ProcessWithAgent(ba.agent, content)
}

// IsAvailable 检查 Agent 是否可用
func (ba *BaseAgent) IsAvailable() bool {
	return ba.agent != nil
}