# TextParser - 文本解析器

Eino 框架内置的基础文本解析器，提供简单直接的纯文本解析功能。是理解复杂解析器的基础。

## 📁 文件结构

```
textparser/
├── README.md       # 本文档
└── text_parser.go  # 使用示例
```

## 🔧 基本使用

### 创建解析器
```go
textParser := parser.TextParser{}
```

### 解析文本
```go
docs, err := textParser.Parse(ctx, strings.NewReader("Hello World"))
if err != nil {
    return err
}
fmt.Println(docs[0].Content) // 输出: Hello World
```

## 🚀 使用示例

### 解析字符串
```go
textContent := "Hello, this is a sample text."
docs, err := textParser.Parse(ctx, strings.NewReader(textContent))
```

### 解析文件
```go
file, err := os.Open("document.txt")
if err != nil {
    return err
}
defer file.Close()

docs, err := textParser.Parse(ctx, file)
```

### 超时控制
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

docs, err := textParser.Parse(ctx, largeReader)
if err == context.DeadlineExceeded {
    fmt.Println("解析超时")
}
```

## 📚 核心特性

### Reader 接口支持
TextParser 支持任何实现了 `io.Reader` 接口的数据源：
- `*strings.Reader` - 字符串数据
- `*os.File` - 文件数据
- `*bytes.Buffer` - 字节数据
- `http.Response.Body` - 网络响应

### 零配置设计
- 无需初始化参数
- 开箱即用
- 线程安全

### 标准接口
```go
func (p TextParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error)
```

## 🛠️ 扩展用法

### 添加元数据
```go
docs, err := textParser.Parse(ctx, reader)
if err == nil {
    for _, doc := range docs {
        if doc.Metadata == nil {
            doc.Metadata = make(map[string]any)
        }
        doc.Metadata["parser"] = "TextParser"
        doc.Metadata["timestamp"] = time.Now()
    }
}
```

### 文本预处理
```go
func preprocessAndParse(text string) ([]*schema.Document, error) {
    // 清理文本
    cleaned := strings.TrimSpace(text)
    cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")

    textParser := parser.TextParser{}
    return textParser.Parse(ctx, strings.NewReader(cleaned))
}
```

### 批量处理
```go
func parseMultipleTexts(texts []string) ([]*schema.Document, error) {
    textParser := parser.TextParser{}
    var allDocs []*schema.Document

    for _, text := range texts {
        docs, err := textParser.Parse(ctx, strings.NewReader(text))
        if err != nil {
            return nil, err
        }
        allDocs = append(allDocs, docs...)
    }

    return allDocs, nil
}
```

## 🔗 组件集成

### 作为 ExtParser 的回退解析器
```go
extParser, _ := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".html": htmlParser,
        ".pdf":  pdfParser,
    },
    FallbackParser: parser.TextParser{}, // 处理未知格式
})
```

### 在 Agent 中使用
```go
// 处理用户输入
textParser := parser.TextParser{}
docs, err := textParser.Parse(ctx, strings.NewReader(userInput))

// 传递给 Agent
agentInput := &adk.AgentInput{
    Messages: []adk.Message{schema.UserMessage(docs[0].Content)},
}
```

## ⚠️ 注意事项

- **资源管理**: 由调用者负责 Reader 的生命周期
- **内存使用**: 大文件使用流式处理，避免一次性读取
- **错误处理**: 及时检查和处理返回的错误
- **Context 使用**: 合理设置超时和取消控制

## 🎓 学习价值

TextParser 展示了 Eino 解析器的核心概念：
- 统一的解析器接口设计
- Context 生命周期管理
- Reader 接口的数据源无关性
- 无状态设计的并发安全性

**下一步**: 学习 [CustomParser](../customparser/) 了解如何实现自定义解析器。

**适用场景**: 日志处理、配置文件解析、简单文本分析等。