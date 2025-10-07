# Document Parser 组件示例

本目录包含 Eino 框架中各种文档解析器的使用示例和实现指南。

## 📁 目录结构

```
parser/
├── README.md          # 本文档 - 解析器概览
├── textparser/        # 基础文本解析器
├── customparser/      # 自定义解析器实现
└── extparser/         # 扩展解析器（多格式支持）
```

## 🔧 解析器类型

### TextParser - 文本解析器
- **位置**: `textparser/`
- **用途**: 处理纯文本内容
- **特点**: 零配置，开箱即用
- **适用场景**: 日志文件、配置文件、纯文本文档

### CustomParser - 自定义解析器
- **位置**: `customparser/`
- **用途**: 学习如何实现自定义解析器
- **特点**: 展示选项模式、配置管理
- **适用场景**: 需要特殊解析逻辑的场景

### ExtParser - 扩展解析器
- **位置**: `extparser/`
- **用途**: 根据文件扩展名自动选择解析器
- **特点**: 支持多种格式，统一接口
- **适用场景**: 处理多种文档格式的应用

## 🚀 快速开始

### 基础文本解析
```go
import "github.com/cloudwego/eino/components/document/parser"

textParser := parser.TextParser{}
docs, err := textParser.Parse(ctx, strings.NewReader("Hello World"))
```

### 多格式文档解析
```go
extParser, _ := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".html": htmlParser,
        ".pdf":  pdfParser,
    },
    FallbackParser: parser.TextParser{},
})

docs, err := extParser.Parse(ctx, file, parser.WithURI("document.html"))
```

## 📚 学习路径

1. **初学者**: 从 `textparser/` 开始，理解基本解析流程
2. **进阶**: 学习 `customparser/`，掌握自定义解析器实现
3. **高级**: 研究 `extparser/`，构建多格式文档处理系统

## 🎯 核心概念

### Parser 接口
```go
type Parser interface {
    Parse(ctx context.Context, reader io.Reader, opts ...Option) ([]*schema.Document, error)
}
```

### Option 模式
所有解析器都支持选项模式，用于传递配置参数：
- `parser.WithURI(uri string)` - 设置文档URI
- `parser.WithExtraMeta(meta map[string]any)` - 添加额外元数据

### Document 结构
```go
type Document struct {
    Content string            // 文档内容
    Metadata map[string]any   // 元数据信息
}
```

## 🛠️ 扩展开发

### 创建自定义解析器
```go
type MyParser struct{}

func (p *MyParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
    // 实现解析逻辑
    return []*schema.Document{{Content: "parsed content"}}, nil
}
```

### 注册到 ExtParser
```go
extParser, _ := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".myformat": &MyParser{},
    },
    FallbackParser: parser.TextParser{},
})
```

## ⚡ 性能建议

- 使用流式处理大文件，避免一次性加载到内存
- 解析器是无状态的，支持并发使用
- 合理设置 Context 超时时间
- 及时关闭文件 Reader

## 🔗 相关组件

- **Retriever**: 文档检索系统
- **Loader**: 文档加载器
- **Embedding**: 文档向量化
- **Agent**: AI 应用系统

## 📖 更多示例

查看各个子目录中的详细示例：
- [TextParser 使用指南](./textparser/README.md)
- [CustomParser 实现指南](./customparser/README.md)
- [ExtParser 配置指南](./extparser/README.md)