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
	"os"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino-ext/components/document/parser/html"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document/parser"

	"github.com/cloudwego/eino-examples/components/document/parser/common"
	"github.com/cloudwego/eino-examples/internal/gptr"
	"github.com/cloudwego/eino-examples/internal/logs"
)

func main() {
	// 基础解析示例
	basicParsingExample()

	// Agent 集成示例
	agentIntegrationExample()
}

// basicParsingExample 演示基础的扩展解析功能
func basicParsingExample() {
	ctx := context.Background()

	textParser := parser.TextParser{}

	htmlParser, err := html.NewParser(ctx, &html.Config{
		Selector: gptr.Of("body"),
	})
	if err != nil {
		logs.Errorf("html.NewParser failed, err=%v", err)
		return
	}

	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{})
	if err != nil {
		logs.Errorf("pdf.NewPDFParser failed, err=%v", err)
		return
	}

	// 创建扩展解析器
	extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
		// 注册特定扩展名的解析器
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".pdf":  pdfParser,
		},
		// 设置默认解析器，用于处理未知格式
		FallbackParser: textParser,
	})
	if err != nil {
		logs.Errorf("NewExtParser failed, err=%v", err)
		return
	}

	// 使用解析器
	filePath := "./testdata/test.html"
	file, err := os.Open(filePath)
	if err != nil {
		logs.Errorf("os.Open failed, file=%v, err=%v", filePath, err)
		return
	}
	defer file.Close()

	docs, err := extParser.Parse(ctx, file,
		// 必须提供 URI ExtParser 选择正确的解析器进行解析
		parser.WithURI(filePath),
		parser.WithExtraMeta(map[string]any{
			"source": "local",
		}),
	)
	if err != nil {
		logs.Errorf("extParser.Parse, err=%v", err)
		return
	}

	for idx, doc := range docs {
		logs.Infof("doc_%v content: %v", idx, doc.Content)
	}
}

// agentIntegrationExample 演示如何将扩展解析器与 Agent 集成
func agentIntegrationExample() {
	ctx := context.Background()

	// 定义文档处理函数
	documentProcessor := func(agent *adk.ChatModelAgent, input string) (string, error) {
		// 解析文档文件
		docs, err := createExtParser(ctx).Parse(ctx, nil, parser.WithURI(input))
		if err != nil {
			return "", fmt.Errorf("failed to parse document: %w", err)
		}

		// 使用 ChatModelAgent 处理解析后的内容
		return common.ProcessWithAgent(agent, docs[0].Content)
	}

	// 运行 Agent 集成示例
	common.RunAgentIntegrationExample("DocumentProcessor",
		"./testdata/test.html",
		func(agent *adk.ChatModelAgent, content string) (string, error) {
			return documentProcessor(agent, content)
		})
}

// DocumentAgent 集成了扩展解析器和 Agent 的结构体（向后兼容）
type DocumentAgent struct {
	extParser *parser.ExtParser
	agent     *common.BaseAgent
}

// NewDocumentAgent 创建新的文档处理器
func NewDocumentAgent(agentType string, ctx context.Context) *DocumentAgent {
	return &DocumentAgent{
		extParser: createExtParser(ctx),
		agent:     common.NewBaseAgent(agentType),
	}
}

// ProcessUploadedFile 处理上传的文件并通过 Agent 分析内容（使用共享代码）
func (da *DocumentAgent) ProcessUploadedFile(ctx context.Context, filePath string) (string, error) {
	// 自动解析上传的文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	docs, err := da.extParser.Parse(ctx, file, parser.WithURI(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	logs.Infof("Parsed document content from %s: %s", filePath, docs[0].Content)

	// 使用共享的 Agent 处理函数
	return da.agent.ProcessContent(docs[0].Content)
}

// createExtParser 创建扩展解析器
func createExtParser(ctx context.Context) *parser.ExtParser {
	textParser := parser.TextParser{}

	htmlParser, err := html.NewParser(ctx, &html.Config{
		Selector: gptr.Of("body"),
	})
	if err != nil {
		logs.Errorf("html.NewParser failed, err=%v", err)
		return nil
	}

	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{})
	if err != nil {
		logs.Errorf("pdf.NewPDFParser failed, err=%v", err)
		return nil
	}

	// 创建扩展解析器
	extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
		// 注册特定扩展名的解析器
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".pdf":  pdfParser,
		},
		// 设置默认解析器，用于处理未知格式
		FallbackParser: textParser,
	})
	if err != nil {
		logs.Errorf("NewExtParser failed, err=%v", err)
		return nil
	}

	return extParser
}
