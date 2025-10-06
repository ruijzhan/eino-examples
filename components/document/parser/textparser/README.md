# TextParser - 文本解析器详解

## 📋 概述

`TextParser` 是 Eino 框架内置的基础文本解析器，提供了最简单直接的纯文本解析功能。虽然代码简洁，但它展示了 Eino 解析器接口的标准实现方式，是理解更复杂解析器的基础。

## 🎯 学习目标

- 理解 Eino 框架内置解析器的使用方法
- 掌握基础的文本解析流程
- 学习解析器接口的标准实现模式
- 了解简单解析器的设计理念

## 📁 文件结构

```
textparser/
├── README.md       # 本文档
└── text_parser.go  # 文本解析器使用示例
```

## 🔧 核心组件分析

### 1. TextParser 结构

```go
textParser := parser.TextParser{}
```

**特点分析**：
- **零配置设计**：不需要任何初始化参数
- **开箱即用**：直接实例化即可使用
- **零依赖**：不依赖外部库或复杂配置
- **线程安全**：无状态设计，支持并发使用

### 2. Parse 方法调用

```go
docs, err := textParser.Parse(ctx, strings.NewReader("hello world"))
```

**方法签名解析**：
```go
func (p TextParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error)
```

**参数说明**：
- `ctx context.Context`：上下文控制，支持取消和超时
- `reader io.Reader`：文本数据源，支持任何实现 Reader 接口的类型
- `opts ...parser.Option`：可变参数列表，用于传递解析选项

**返回值说明**：
- `[]*schema.Document`：解析后的文档数组
- `error`：解析过程中的错误信息

### 3. 输入输出处理

#### 输入源多样性

```go
// 字符串
strings.NewReader("hello world")

// 文件
file, _ := os.Open("text.txt")
textParser.Parse(ctx, file)

// 网络流
resp, _ := http.Get(url)
textParser.Parse(ctx, resp.Body)

// 内存缓冲
buffer := bytes.NewBuffer(data)
textParser.Parse(ctx, buffer)
```

#### 输出结构

```go
type Document struct {
    Content string                 // 文档内容
    Metadata map[string]any       // 元数据信息
}
```

## 🚀 使用示例

### 基础使用

```go
package main

import (
    "context"
    "strings"
    "github.com/cloudwego/eino/components/document/parser"
)

func main() {
    ctx := context.Background()

    // 创建文本解析器
    textParser := parser.TextParser{}

    // 准备文本数据
    textContent := "Hello, this is a sample text for parsing."

    // 执行解析
    docs, err := textParser.Parse(ctx, strings.NewReader(textContent))
    if err != nil {
        panic(err)
    }

    // 处理解析结果
    for i, doc := range docs {
        fmt.Printf("Document %d: %s\n", i+1, doc.Content)
    }
}
```

### 文件解析示例

```go
func parseTextFile(filePath string) error {
    // 打开文件
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    // 创建解析器
    textParser := parser.TextParser{}

    // 解析文件内容
    docs, err := textParser.Parse(context.Background(), file)
    if err != nil {
        return err
    }

    // 输出结果
    fmt.Printf("解析成功，共 %d 个文档\n", len(docs))
    for _, doc := range docs {
        fmt.Printf("内容: %s\n", doc.Content)
    }

    return nil
}
```

### 流式数据处理

```go
func parseStreamingData(dataStream <-chan []byte) {
    textParser := parser.TextParser{}

    for data := range dataStream {
        reader := bytes.NewReader(data)
        docs, err := textParser.Parse(context.Background(), reader)
        if err != nil {
            log.Printf("解析错误: %v", err)
            continue
        }

        // 处理解析结果
        processDocuments(docs)
    }
}
```

## 📚 核心概念解析

### 1. Reader 接口的灵活性

`TextParser` 接受任何实现了 `io.Reader` 接口的数据源：

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**支持的常见类型**：
- `*strings.Reader`：字符串数据
- `*os.File`：文件数据
- `*bytes.Buffer`：字节数据
- `http.Response.Body`：网络响应
- 自定义 Reader：可以实现特殊的读取逻辑

### 2. Context 的作用

```go
docs, err := textParser.Parse(ctx, reader)
```

Context 提供了以下能力：
- **取消控制**：可以提前取消解析操作
- **超时控制**：设置解析超时时间
- **值传递**：在解析过程中传递上下文信息

#### 超时控制示例

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

docs, err := textParser.Parse(ctx, largeReader)
if err == context.DeadlineExceeded {
    fmt.Println("解析超时")
}
```

### 3. Document 结构解析

```go
type Document struct {
    Content string                 // 主要内容
    Metadata map[string]any       // 元数据信息
}
```

**Content 字段**：
- 存储解析后的文本内容
- 对于 TextParser，这是原始文本的直接映射

**Metadata 字段**：
- 包含文档的元信息
- 可能包含文件名、创建时间、编码等信息
- 对于 TextParser，通常是空的或包含基础信息

## 🔍 代码逐行分析

### text_parser.go 关键代码

#### 第31行：创建解析器实例
```go
textParser := parser.TextParser{}
```
- 直接实例化，无需参数
- 使用结构体字面量语法
- 零值即可正常工作

#### 第32行：执行解析
```go
docs, err := textParser.Parse(ctx, strings.NewReader("hello world"))
```
- 调用 Parse 方法
- 传入 context.Background() 作为上下文
- 使用 strings.NewReader 包装字符串数据

#### 第33-36行：错误处理
```go
if err != nil {
    logs.Errorf("TextParser{}.Parse failed, err=%v", err)
    return
}
```
- 标准的错误处理模式
- 使用项目的日志系统记录错误
- 遇到错误时直接返回

#### 第38行：结果输出
```go
logs.Infof("text content: %v", docs[0].Content)
```
- 输出第一个文档的内容
- 使用格式化字符串显示结果
- 假设解析结果至少包含一个文档

## 🛠️ 扩展和变体

### 1. 添加元数据支持

```go
func parseWithMetadata(reader io.Reader, filename string) ([]*schema.Document, error) {
    textParser := parser.TextParser{}
    docs, err := textParser.Parse(context.Background(), reader)
    if err != nil {
        return nil, err
    }

    // 添加元数据
    for _, doc := range docs {
        if doc.Metadata == nil {
            doc.Metadata = make(map[string]any)
        }
        doc.Metadata["filename"] = filename
        doc.Metadata["parser"] = "TextParser"
        doc.Metadata["parsed_at"] = time.Now()
    }

    return docs, nil
}
```

### 2. 文本预处理

```go
func parseWithPreprocessing(text string) ([]*schema.Document, error) {
    // 预处理文本
    cleanedText := preprocessText(text)

    textParser := parser.TextParser{}
    return textParser.Parse(context.Background(), strings.NewReader(cleanedText))
}

func preprocessText(text string) string {
    // 移除多余空白
    text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
    // 去除首尾空格
    text = strings.TrimSpace(text)
    return text
}
```

### 3. 批量处理

```go
func parseMultipleTexts(texts []string) ([]*schema.Document, error) {
    textParser := parser.TextParser{}
    var allDocs []*schema.Document

    for i, text := range texts {
        docs, err := textParser.Parse(context.Background(), strings.NewReader(text))
        if err != nil {
            return nil, fmt.Errorf("解析第 %d 个文本失败: %w", i+1, err)
        }
        allDocs = append(allDocs, docs...)
    }

    return allDocs, nil
}
```

## ⚡ 性能考虑

### 1. 内存使用

```go
// ✅ 好的做法：处理大文件时使用流式读取
file, _ := os.Open("large.txt")
defer file.Close()

docs, err := textParser.Parse(ctx, file)

// ❌ 避免：一次性读取大文件到内存
data, _ := os.ReadFile("large.txt")
reader := bytes.NewReader(data)  // 可能导致内存溢出
```

### 2. 并发处理

```go
func parseConcurrently(texts []string) ([]*schema.Document, error) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var allDocs []*schema.Document
    errs := make(chan error, len(texts))

    for i, text := range texts {
        wg.Add(1)
        go func(idx int, content string) {
            defer wg.Done()

            textParser := parser.TextParser{}
            docs, err := textParser.Parse(context.Background(), strings.NewReader(content))
            if err != nil {
                errs <- fmt.Errorf("处理第 %d 个文本失败: %w", idx+1, err)
                return
            }

            mu.Lock()
            allDocs = append(allDocs, docs...)
            mu.Unlock()
        }(i, text)
    }

    wg.Wait()
    close(errs)

    for err := range errs {
        if err != nil {
            return allDocs, err
        }
    }

    return allDocs, nil
}
```

## 🔗 与其他组件的集成

### 1. 与 ExtParser 的配合

```go
// TextParser 常用作 ExtParser 的默认解析器
extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".html": htmlParser,
        ".pdf":  pdfParser,
    },
    FallbackParser: parser.TextParser{},  // 处理未知格式
})
```

### 2. 与 Retriever 的集成

```go
// 解析文本后添加到检索系统
func addToRetriever(textContent string) error {
    textParser := parser.TextParser{}
    docs, err := textParser.Parse(context.Background(), strings.NewReader(textContent))
    if err != nil {
        return err
    }

    // 添加到检索器
    return retriever.AddDocuments(ctx, docs)
}
```

### 3. 在 Agent 中的使用

```go
// 在文件顶部添加必要的导入
import (
    "context"
    "strings"
    "github.com/cloudwego/eino/components/document/parser"
    "github.com/cloudwego/eino/schema"
    "github.com/cloudwego/eino/adk"
)

// Agent 使用 TextParser 处理用户输入的文本
func processUserInput(agent *adk.ChatModelAgent, input string) (string, error) {
    ctx := context.Background()
    textParser := parser.TextParser{}
    docs, err := textParser.Parse(ctx, strings.NewReader(input))
    if err != nil {
        return "", err
    }

    // 将解析后的文档传递给 Agent
    // 创建 AgentInput
    agentInput := &adk.AgentInput{
        Messages:        []adk.Message{schema.UserMessage(docs[0].Content)},
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
```

## 🧪 测试示例

### 基础单元测试

```go
func TestTextParser_Parse(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"Empty string", "", ""},
        {"Simple text", "hello world", "hello world"},
        {"Multiline text", "line1\nline2", "line1\nline2"},
        {"Unicode text", "你好世界", "你好世界"},
    }

    textParser := parser.TextParser{}
    ctx := context.Background()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            docs, err := textParser.Parse(ctx, strings.NewReader(tt.input))
            assert.NoError(t, err)
            assert.Len(t, docs, 1)
            assert.Equal(t, tt.expected, docs[0].Content)
        })
    }
}
```

### 错误处理测试

```go
func TestTextParser_ErrorHandling(t *testing.T) {
    textParser := parser.TextParser{}

    // 测试 nil reader
    _, err := textParser.Parse(context.Background(), nil)
    assert.Error(t, err)

    // 测试取消的 context
    ctx, cancel := context.WithCancel(context.Background())
    cancel()  // 立即取消

    reader := strings.NewReader("test")
    _, err = textParser.Parse(ctx, reader)
    assert.Equal(t, context.Canceled, err)
}
```

## ⚠️ 常见问题和注意事项

### 1. Reader 管理问题

```go
// ❌ 错误：重复使用已关闭的 Reader
file, _ := os.Open("file.txt")
file.Close()  // 关闭文件
textParser.Parse(ctx, file)  // 错误：文件已关闭

// ✅ 正确：确保 Reader 在解析期间保持打开
func parseFile(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()  // 在函数结束时关闭

    textParser := parser.TextParser{}
    _, err = textParser.Parse(context.Background(), file)
    return err
}
```

### 2. Context 使用问题

```go
// ❌ 错误：使用已取消的 Context
ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
time.Sleep(time.Millisecond)  // 确保超时
textParser.Parse(ctx, reader)  // Context 已取消

// ✅ 正确：及时检查 Context 状态
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

textParser := parser.TextParser{}
docs, err := textParser.Parse(ctx, reader)
if err == context.DeadlineExceeded {
    // 处理超时情况
}
```

### 3. 内存泄漏问题

```go
// ❌ 错误：没有处理大文件的情况
func parseLargeFile(filePath string) error {
    data, err := os.ReadFile(filePath)  // 可能占用大量内存
    if err != nil {
        return err
    }

    reader := bytes.NewReader(data)
    textParser := parser.TextParser{}
    _, err = textParser.Parse(context.Background(), reader)
    return err
}

// ✅ 正确：使用流式处理
func parseLargeFile(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    textParser := parser.TextParser{}
    _, err = textParser.Parse(context.Background(), file)
    return err
}
```

## 🎓 总结

`TextParser` 虽然代码简单，但它完美展示了 Eino 框架解析器的设计哲学：

### 核心优势
1. **简单易用**：零配置，开箱即用
2. **灵活性高**：支持任何 Reader 接口的数据源
3. **标准接口**：遵循统一的解析器接口规范
4. **性能优良**：无状态设计，支持并发使用

### 设计启示
1. **接口抽象**：通过 `io.Reader` 接口实现数据源无关性
2. **上下文支持**：通过 `context.Context` 实现生命周期管理
3. **错误处理**：标准的 Go 错误处理模式
4. **并发安全**：无状态设计天然支持并发

### 学习价值
虽然 `TextParser` 功能简单，但它是理解 Eino 框架解析器系统的最佳起点。掌握了它，你就能更好地理解：
- 更复杂的解析器（如 HTMLParser、PDFParser）
- 自定义解析器的实现方法
- 文档处理在整个 AI 应用中的地位

**下一步学习**：建议继续学习 `customparser` 了解如何实现自定义解析器，或者学习 `extparser` 了解如何组合多个解析器。

---

**实践建议**：尝试基于 `TextParser` 构建一个简单的文本处理工具，比如日志文件分析器或配置文件处理器，这样可以更好地理解解析器在实际应用中的作用。