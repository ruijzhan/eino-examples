# ExtParser - 多格式解析器

根据文件扩展名自动选择合适的解析器，为多种文档格式提供统一的解析接口。

## 📁 文件结构

```
extparser/
├── README.md       # 本文档
├── ext_parser.go   # 使用示例
└── testdata/       # 测试数据
    └── test.html   # HTML 测试文件
```

## 🔧 核心配置

### ExtParserConfig
```go
type ExtParserConfig struct {
    Parsers        map[string]parser.Parser  // 扩展名 -> 解析器映射
    FallbackParser parser.Parser            // 默认回退解析器
}
```

### 创建 ExtParser
```go
extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".html": htmlParser,  // HTML 文件
        ".pdf":  pdfParser,   // PDF 文件
        ".md":   markdownParser, // Markdown 文件
    },
    FallbackParser: parser.TextParser{},  // 处理未知格式
})
```

## 🚀 使用示例

### 基本使用
```go
// 必须提供 URI 用于格式识别
docs, err := extParser.Parse(ctx, file,
    parser.WithURI("document.html"),
    parser.WithExtraMeta(map[string]any{
        "source": "local",
    }),
)
```

### 批量处理
```go
files := []string{"doc.html", "report.pdf", "readme.txt"}
for _, file := range files {
    docs, err := parseFile(extParser, file)
    // 处理解析结果...
}
```

### 动态注册
```go
// 运行时动态添加解析器
parsers := map[string]parser.Parser{
    ".txt": parser.TextParser{},
}

if htmlEnabled {
    htmlParser, _ := html.NewParser(ctx, &html.Config{
        Selector: gptr.Of("body"),
    })
    parsers[".html"] = htmlParser
    parsers[".htm"] = htmlParser
}
```

## 📚 工作机制

### 格式识别流程
```
文件路径 → 提取扩展名 → 查找解析器 → 执行解析 → 返回结果
   ↓           ↓           ↓          ↓         ↓
"test.html" → ".html" → htmlParser → 解析HTML → 文档内容
```

### 解析器选择策略
1. **精确匹配**: 根据扩展名在 Parsers 映射中查找
2. **回退处理**: 使用 FallbackParser 处理未知格式
3. **错误处理**: 所有解析器都失败时返回错误

### 元数据合并
ExtParser 自动合并元数据：
```go
{
    "uri": "document.html",           // 来自 WithURI
    "source": "local",               // 来自 WithExtraMeta
    "parser_type": "html",           // 自动添加
    "original_format": ".html",      // 自动添加
}
```

## 🛠️ 扩展方式

### 添加新格式支持
```go
// 1. 实现自定义解析器
jsonParser := &JSONParser{}

// 2. 注册到 ExtParser
config.Parsers[".json"] = jsonParser
config.Parsers[".jsonl"] = jsonParser
```

### 配置文件驱动
```go
// 从配置文件加载解析器配置
type ParserConfig struct {
    Extension string                 `json:"extension"`
    Type      string                 `json:"type"`
    Options   map[string]interface{} `json:"options"`
}
```

## 🔗 组件集成

### 与检索系统集成
```go
// 解析后直接添加到检索系统
docs, _ := extParser.Parse(ctx, file, parser.WithURI(path))
retriever.AddDocuments(ctx, docs)
```

### 在 Agent 中使用
```go
// 自动解析用户上传的文件
docs, _ := extParser.Parse(ctx, uploadedFile, parser.WithURI(filename))
agentInput := &adk.AgentInput{
    Messages: []adk.Message{schema.UserMessage(docs[0].Content)},
}
```

### Web API 集成
```go
// 处理文件上传
docs, _ := extParser.Parse(r.Context(), file,
    parser.WithURI(header.Filename),
    parser.WithExtraMeta(map[string]any{
        "uploaded_by": userID,
    }),
)
```

## ⚡ 性能优化

### 并发处理
```go
// 限制并发数量的批量解析
semaphore := make(chan struct{}, 10)
for _, file := range files {
    go func(filePath string) {
        semaphore <- struct{}{}
        defer func() { <-semaphore }()

        docs, _ := parseFile(extParser, filePath)
        // 处理结果...
    }(file)
}
```

### 解析器池化
```go
// 复用解析器实例以减少创建开销
type ParserPool struct {
    parsers map[string]chan parser.Parser
    // 实现获取和归还方法...
}
```

## ⚠️ 注意事项

- **URI 必需**: 必须使用 `parser.WithURI()` 提供文件路径
- **回退解析器**: 始终设置合适的 FallbackParser
- **解析器一致性**: 相同格式应使用相同解析器配置
- **资源管理**: 由调用者负责 Reader 的生命周期

## 🎓 应用场景

- **文档管理系统**: 处理各种格式的用户文档
- **内容聚合平台**: 从多源提取和处理内容
- **知识库构建**: 统一处理不同格式的知识文档
- **AI 助手**: 让 AI 理解多种格式的用户输入

**下一步**: 学习 [CustomParser](../customparser/) 了解如何实现自定义解析器。