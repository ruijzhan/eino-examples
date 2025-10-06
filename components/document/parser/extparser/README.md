# ExtParser - 扩展解析器详解

## 📋 概述

`ExtParser` 是 Eino 框架中的扩展解析器，它能够根据文件扩展名自动选择合适的解析器来处理不同格式的文档。这个组件展示了如何构建一个统一的文档解析入口，支持多种文档格式的无缝处理。

## 🎯 学习目标

- 理解扩展解析器的自动格式识别机制
- 掌握多解析器的注册和管理方法
- 学会配置默认解析器和回退策略
- 了解 URI 在格式识别中的重要作用

## 📁 文件结构

```
extparser/
├── README.md       # 本文档
├── ext_parser.go   # 扩展解析器使用示例
└── testdata/       # 测试数据
    └── test.html   # HTML 测试文件
```

## 🔧 核心组件分析

### 1. ExtParser 配置结构

```go
type ExtParserConfig struct {
    Parsers        map[string]parser.Parser  // 按扩展名映射的解析器
    FallbackParser parser.Parser            // 默认解析器
}
```

**设计要点**：
- **解析器注册表**：通过 `map[string]parser.Parser` 管理不同格式的解析器
- **回退机制**：`FallbackParser` 处理未知格式或注册失败的情况
- **扩展名映射**：以文件扩展名为键，解析器实例为值

### 2. 解析器注册示例

```go
extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".html": htmlParser,  // HTML 文件使用 HTML 解析器
        ".pdf":  pdfParser,   // PDF 文件使用 PDF 解析器
    },
    FallbackParser: textParser,  // 其他格式使用文本解析器
})
```

**解析器分工**：
- **HTML 解析器**：提取 HTML 内容，支持 CSS 选择器
- **PDF 解析器**：解析 PDF 文档，提取文本内容
- **文本解析器**：处理纯文本和其他未知格式

### 3. URI 的重要性

```go
docs, err := extParser.Parse(ctx, file,
    parser.WithURI(filePath),  // 必须提供 URI
    parser.WithExtraMeta(map[string]any{
        "source": "local",
    }),
)
```

**URI 的作用**：
- **格式识别**：通过文件扩展名确定使用哪个解析器
- **元数据来源**：为解析结果提供文件路径信息
- **解析器选择**：ExtParser 根据扩展名在注册表中查找对应的解析器

## 🚀 使用示例

### 基础多格式解析

```go
package main

import (
    "context"
    "github.com/cloudwego/eino/components/document/parser"
    "github.com/cloudwego/eino-ext/components/document/parser/html"
    "github.com/cloudwego/eino-ext/components/document/parser/pdf"
)

func main() {
    ctx := context.Background()

    // 1. 创建各种解析器
    textParser := parser.TextParser{}

    htmlParser, _ := html.NewParser(ctx, &html.Config{
        Selector: gptr.Of("body"),  // 只提取 body 内容
    })

    pdfParser, _ := pdf.NewPDFParser(ctx, &pdf.Config{})

    // 2. 创建扩展解析器
    extParser, _ := parser.NewExtParser(ctx, &parser.ExtParserConfig{
        Parsers: map[string]parser.Parser{
            ".html": htmlParser,
            ".pdf":  pdfParser,
        },
        FallbackParser: textParser,
    })

    // 3. 解析不同格式的文件
    files := []string{
        "document.html",
        "report.pdf",
        "readme.txt",
    }

    for _, file := range files {
        docs, err := parseFile(extParser, file)
        if err != nil {
            fmt.Printf("解析 %s 失败: %v\n", file, err)
            continue
        }
        fmt.Printf("成功解析 %s: %s\n", file, docs[0].Content[:50])
    }
}

func parseFile(extParser *parser.ExtParser, filePath string) ([]*schema.Document, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    return extParser.Parse(context.Background(), file,
        parser.WithURI(filePath),
        parser.WithExtraMeta(map[string]any{
            "source": "local_file",
            "format": filepath.Ext(filePath),
        }),
    )
}
```

### 动态解析器注册

```go
func createDynamicExtParser(ctx context.Context) (*parser.ExtParser, error) {
    // 基础解析器
    parsers := map[string]parser.Parser{
        ".txt":  parser.TextParser{},
        ".md":   createMarkdownParser(),
    }

    // 条件性添加解析器
    if isHTMLSupportEnabled() {
        htmlParser, _ := html.NewParser(ctx, &html.Config{
            Selector: gptr.Of("main"),
        })
        parsers[".html"] = htmlParser
        parsers[".htm"] = htmlParser
    }

    if isPDFSupportEnabled() {
        pdfParser, _ := pdf.NewPDFParser(ctx, &pdf.Config{})
        parsers[".pdf"] = pdfParser
    }

    return parser.NewExtParser(ctx, &parser.ExtParserConfig{
        Parsers:        parsers,
        FallbackParser: parser.TextParser{},
    })
}
```

### 批量文档处理

```go
func batchProcessDocuments(extParser *parser.ExtParser, dirPath string) ([]*schema.Document, error) {
    var allDocs []*schema.Document

    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        docs, err := parseFile(extParser, path)
        if err != nil {
            fmt.Printf("跳过文件 %s: %v\n", path, err)
            return nil
        }

        allDocs = append(allDocs, docs...)
        return nil
    })

    return allDocs, err
}
```

## 🔍 代码逐行分析

### ext_parser.go 关键代码

#### 第34-48行：创建各种解析器
```go
// 基础文本解析器
textParser := parser.TextParser{}

// HTML 解析器，配置只提取 body 内容
htmlParser, err := html.NewParser(ctx, &html.Config{
    Selector: gptr.Of("body"),
})

// PDF 解析器，使用默认配置
pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{})
```

**解析器配置说明**：
- **HTML 解析器**：使用 CSS 选择器 `body` 只提取页面主体内容
- **PDF 解析器**：使用默认配置，提取所有文本内容
- **文本解析器**：无需配置，直接处理纯文本

#### 第51-59行：创建 ExtParser 配置
```go
extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".html": htmlParser,
        ".pdf":  pdfParser,
    },
    FallbackParser: textParser,
})
```

**配置详解**：
- **Parsers 映射**：定义了扩展名到解析器的映射关系
- **FallbackParser**：当找不到对应解析器时使用的默认解析器
- **错误处理**：如果创建失败，err 会被设置

#### 第72-78行：解析选项配置
```go
docs, err := extParser.Parse(ctx, file,
    // 必须提供 URI ExtParser 选择正确的解析器进行解析
    parser.WithURI(filePath),
    parser.WithExtraMeta(map[string]any{
        "source": "local",
    }),
)
```

**选项说明**：
- **WithURI**：必需选项，提供文件路径用于格式识别
- **WithExtraMeta**：可选选项，添加额外的元数据信息

## 📚 核心概念解析

### 1. 格式识别机制

ExtParser 的格式识别流程：

```
输入文件路径 → 提取扩展名 → 查找解析器映射 → 选择解析器 → 执行解析
     ↓              ↓              ↓              ↓           ↓
"./test.html"  →  ".html"   →  Parsers[".html"] → htmlParser → 解析HTML
```

**识别失败处理**：
- 如果扩展名不在映射表中，使用 `FallbackParser`
- 如果 `FallbackParser` 也失败，返回错误

### 2. 解析器选择策略

```go
func selectParser(extension string, config *ExtParserConfig) parser.Parser {
    // 1. 精确匹配
    if parser, exists := config.Parsers[extension]; exists {
        return parser
    }

    // 2. 大小写不敏感匹配
    for ext, parser := range config.Parsers {
        if strings.EqualFold(ext, extension) {
            return parser
        }
    }

    // 3. 使用回退解析器
    return config.FallbackParser
}
```

### 3. 元数据传递

ExtParser 会合并来自不同源的元数据：

```go
// 最终的 Document 元数据包含：
{
    "uri": "./test.html",           // 来自 WithURI
    "source": "local",             // 来自 WithExtraMeta
    "parser_type": "html",         // 来自 ExtParser 自动添加
    "original_format": ".html",    // 来自 ExtParser 自动添加
}
```

## 🛠️ 扩展和配置

### 1. 添加新的文档格式

```go
func addCustomParser(extParser *parser.ExtParser) error {
    // 创建自定义解析器
    jsonParser := createJSONParser()
    xmlParser := createXMLParser()

    // 重新配置 ExtParser
    newConfig := &parser.ExtParserConfig{
        Parsers: map[string]parser.Parser{
            ".html": htmlParser,
            ".pdf":  pdfParser,
            ".json": jsonParser,  // 新增 JSON 支持
            ".xml":  xmlParser,   // 新增 XML 支持
        },
        FallbackParser: textParser,
    }

    // 重新创建 ExtParser
    return parser.NewExtParser(ctx, newConfig)
}
```

### 2. 配置文件驱动的解析器

```go
type ParserConfig struct {
    Extension string                 `json:"extension"`
    Type      string                 `json:"type"`
    Options   map[string]interface{} `json:"options"`
}

func loadExtParserFromConfig(ctx context.Context, configFile string) (*parser.ExtParser, error) {
    // 读取配置文件
    configData, err := os.ReadFile(configFile)
    if err != nil {
        return nil, err
    }

    var parserConfigs []ParserConfig
    json.Unmarshal(configData, &parserConfigs)

    // 创建解析器映射
    parsers := make(map[string]parser.Parser)
    for _, pc := range parserConfigs {
        parser, err := createParserFromConfig(pc)
        if err != nil {
            continue  // 跳过创建失败的解析器
        }
        parsers[pc.Extension] = parser
    }

    return parser.NewExtParser(ctx, &parser.ExtParserConfig{
        Parsers:        parsers,
        FallbackParser: parser.TextParser{},
    })
}
```

### 3. 动态解析器加载

```go
func loadPluginParsers(pluginDir string) (map[string]parser.Parser, error) {
    parsers := make(map[string]parser.Parser)

    // 扫描插件目录
    files, err := os.ReadDir(pluginDir)
    if err != nil {
        return nil, err
    }

    for _, file := range files {
        if !strings.HasSuffix(file.Name(), ".so") {
            continue
        }

        // 加载插件
        plugin, err := plugin.Open(filepath.Join(pluginDir, file.Name()))
        if err != nil {
            continue
        }

        // 查找解析器创建函数
        createParserSymbol, err := plugin.Lookup("CreateParser")
        if err != nil {
            continue
        }

        createParser := createParserSymbol.(func(context.Context) (parser.Parser, error))
        p, err := createParser(context.Background())
        if err != nil {
            continue
        }

        // 根据插件名确定扩展名
        ext := "." + strings.TrimSuffix(file.Name(), ".so")
        parsers[ext] = p
    }

    return parsers, nil
}
```

## 🔗 与其他组件的集成

### 1. 与文档检索系统集成

```go
type DocumentIndexer struct {
    extParser *parser.ExtParser
    retriever retriever.Retriever
}

func (di *DocumentIndexer) IndexDirectory(dirPath string) error {
    // 批量处理目录中的所有文档
    docs, err := batchProcessDocuments(di.extParser, dirPath)
    if err != nil {
        return err
    }

    // 添加到检索系统
    return di.retriever.AddDocuments(context.Background(), docs)
}
```

### 2. 在 Agent 中的使用

```go
import (
    "context"
    "os"
    "github.com/cloudwego/eino/components/document/parser"
    "github.com/cloudwego/eino/schema"
    "github.com/cloudwego/eino/adk"
)

type DocumentAgent struct {
    extParser *parser.ExtParser
    agent     *adk.ChatModelAgent
}

func (da *DocumentAgent) ProcessUploadedFile(ctx context.Context, filePath string) (string, error) {
    // 自动解析上传的文件
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    docs, err := da.extParser.Parse(ctx, file, parser.WithURI(filePath))
    if err != nil {
        return "", err
    }

    // 将解析结果传递给 Agent 处理
    // 创建 AgentInput
    agentInput := &adk.AgentInput{
        Messages:        []adk.Message{schema.UserMessage(docs[0].Content)},
        EnableStreaming: false,
    }

    // 运行 Agent 并处理结果
    iterator := da.agent.Run(ctx, agentInput)

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
```

### 3. Web API 集成

```go
func handleDocumentUpload(w http.ResponseWriter, r *http.Request) {
    // 解析上传的文件
    file, header, err := r.FormFile("document")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // 使用 ExtParser 解析
    docs, err := extParser.Parse(r.Context(), file,
        parser.WithURI(header.Filename),
        parser.WithExtraMeta(map[string]any{
            "uploaded_by": r.Header.Get("User-ID"),
            "content_type": header.Header.Get("Content-Type"),
        }),
    )

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 返回解析结果
    json.NewEncoder(w).Encode(map[string]interface{}{
        "content": docs[0].Content,
        "metadata": docs[0].Metadata,
    })
}
```

## ⚡ 性能优化

### 1. 解析器池化

```go
type ParserPool struct {
    parsers map[string]chan parser.Parser
    factory map[string]func() (parser.Parser, error)
}

func NewParserPool() *ParserPool {
    return &ParserPool{
        parsers: make(map[string]chan parser.Parser),
        factory: make(map[string]func() (parser.Parser, error)),
    }
}

func (pp *ParserPool) GetParser(extension string) (parser.Parser, error) {
    if ch, exists := pp.parsers[extension]; exists {
        select {
        case parser := <-ch:
            return parser, nil
        default:
            // 池为空，创建新的解析器
            return pp.factory[extension]()
        }
    }

    return nil, fmt.Errorf("unknown parser for extension: %s", extension)
}

func (pp *ParserPool) ReturnParser(extension string, p parser.Parser) {
    if ch, exists := pp.parsers[extension]; exists {
        select {
        case ch <- p:
            // 成功归还到池中
        default:
            // 池满了，丢弃这个解析器
        }
    }
}
```

### 2. 并发解析

```go
func parseConcurrently(extParser *parser.ExtParser, files []string) ([]*schema.Document, error) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var allDocs []*schema.Document
    errChan := make(chan error, len(files))

    // 限制并发数量
    semaphore := make(chan struct{}, 10)

    for _, file := range files {
        wg.Add(1)
        go func(filePath string) {
            defer wg.Done()

            semaphore <- struct{}{}  // 获取信号量
            defer func() { <-semaphore }()  // 释放信号量

            docs, err := parseFile(extParser, filePath)
            if err != nil {
                errChan <- fmt.Errorf("解析 %s 失败: %w", filePath, err)
                return
            }

            mu.Lock()
            allDocs = append(allDocs, docs...)
            mu.Unlock()
        }(file)
    }

    wg.Wait()
    close(errChan)

    // 检查是否有错误
    for err := range errChan {
        return allDocs, err  // 返回第一个错误
    }

    return allDocs, nil
}
```

## 🧪 测试策略

### 1. 单元测试

```go
func TestExtParser_FormatSelection(t *testing.T) {
    tests := []struct {
        name     string
        filename string
        expected string
    }{
        {"HTML file", "test.html", "Hello World"},
        {"PDF file", "test.pdf", "PDF content"},
        {"Text file", "test.txt", "plain text"},
        {"Unknown file", "test.xyz", "fallback text"},
    }

    extParser := createTestExtParser(t)
    ctx := context.Background()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            reader := strings.NewReader(tt.expected)

            docs, err := extParser.Parse(ctx, reader, parser.WithURI(tt.filename))
            assert.NoError(t, err)
            assert.Len(t, docs, 1)
            assert.Contains(t, docs[0].Content, tt.expected)
        })
    }
}
```

### 2. 集成测试

```go
func TestExtParser_RealFiles(t *testing.T) {
    extParser := createTestExtParser(t)

    testFiles := []string{
        "testdata/sample.html",
        "testdata/sample.pdf",
        "testdata/sample.txt",
    }

    for _, file := range testFiles {
        t.Run(file, func(t *testing.T) {
            fileReader, err := os.Open(file)
            require.NoError(t, err)
            defer fileReader.Close()

            docs, err := extParser.Parse(context.Background(), fileReader,
                parser.WithURI(file))
            assert.NoError(t, err)
            assert.NotEmpty(t, docs[0].Content)
        })
    }
}
```

## ⚠️ 常见问题和注意事项

### 1. URI 必须提供

```go
// ❌ 错误：没有提供 URI
docs, err := extParser.Parse(ctx, file)
// ExtParser 无法确定使用哪个解析器

// ✅ 正确：提供 URI
docs, err := extParser.Parse(ctx, file, parser.WithURI("document.html"))
```

### 2. 解析器映射冲突

```go
// ❌ 错误：多个扩展名映射到同一个解析器但配置不同
parsers := map[string]parser.Parser{
    ".html": htmlParser1,
    ".htm":  htmlParser2,  // 配置不同可能导致行为不一致
}

// ✅ 正确：统一相同类型的解析器配置
htmlParser := createHTMLParser()
parsers := map[string]parser.Parser{
    ".html": htmlParser,
    ".htm":  htmlParser,  // 使用同一个解析器
}
```

### 3. FallbackParser 的选择

```go
// ❌ 错误：FallbackParser 设置为 nil
config := &parser.ExtParserConfig{
    Parsers:        parsers,
    FallbackParser: nil,  // 未知格式会失败
}

// ✅ 正确：设置合适的默认解析器
config := &parser.ExtParserConfig{
    Parsers:        parsers,
    FallbackParser: parser.TextParser{},  // 处理未知格式
}
```

## 🎓 总结

ExtParser 展示了 Eino 框架在文档处理方面的强大能力：

### 核心优势
1. **自动格式识别**：根据文件扩展名自动选择合适的解析器
2. **统一接口**：为不同格式的文档提供一致的解析接口
3. **可扩展性**：易于添加新的文档格式支持
4. **容错机制**：通过 FallbackParser 处理未知格式

### 设计模式
1. **策略模式**：根据文件类型选择不同的解析策略
2. **注册表模式**：通过映射表管理解析器
3. **适配器模式**：统一不同解析器的接口
4. **模板方法模式**：定义统一的解析流程

### 实际应用价值
- **文档管理系统**：处理各种格式的文档
- **内容管理系统**：自动解析上传的文档
- **知识库系统**：从多种文档源提取知识
- **Agent 系统**：让 AI 能够理解各种格式的文档

**下一步学习**：建议学习如何实现自定义解析器，或者了解如何将 ExtParser 集成到更大的应用系统中。

---

**实践建议**：尝试基于 ExtParser 构建一个完整的文档处理管道，包括格式识别、内容提取、元数据提取和后续处理等功能。这样可以更好地理解多格式文档处理的完整流程。