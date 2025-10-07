# CustomParser - 自定义解析器实现

展示如何在 Eino 框架中实现符合标准的自定义解析器，重点演示选项模式和配置管理。

## 📁 文件结构

```
customparser/
├── README.md           # 本文档
├── custom_parser.go    # 解析器核心实现
└── parse.go           # 使用示例
```

## 🔧 核心实现

### 解析器结构
```go
type CustomParser struct {
    defaultEncoding string
    defaultMaxSize  int64
}
```

### 配置管理
```go
// 解析器级配置
type Config struct {
    DefaultEncoding string
    DefaultMaxSize  int64
}

// 运行时选项
type options struct {
    Encoding string
    MaxSize  int64
}
```

### 选项函数
```go
func WithEncoding(encoding string) parser.Option {
    return parser.WrapImplSpecificOptFn(func(o *options) {
        o.Encoding = encoding
    })
}

func WithMaxSize(size int64) parser.Option {
    return parser.WrapImplSpecificOptFn(func(o *options) {
        o.MaxSize = size
    })
}
```

## 🚀 使用示例

### 基本使用
```go
// 创建解析器
customParser, err := NewCustomParser(&Config{
    DefaultEncoding: "utf-8",
    DefaultMaxSize:  1024,
})

// 执行解析
docs, err := customParser.Parse(ctx, reader)
```

### 使用自定义选项
```go
docs, err := customParser.Parse(ctx, reader,
    WithMaxSize(2048),        // 临时修改最大大小
    WithEncoding("gbk"),      // 临时修改编码
)
```

## 📚 关键模式

### 选项模式
- **Config**: 解析器创建时的配置，创建后不可更改
- **options**: 每次解析时可覆盖的运行时配置
- **优先级**: 运行时选项 > 解析器默认值 > 框架默认值

### 解析流程
1. 处理通用选项 (`parser.GetCommonOptions`)
2. 处理特定选项 (`parser.GetImplSpecificOptions`)
3. 执行解析逻辑

## 🛠️ 扩展指南

### 添加新选项
```go
// 1. 在 options 中添加字段
type options struct {
    Encoding    string
    MaxSize     int64
    Timeout     time.Duration  // 新增
}

// 2. 创建 Option 函数
func WithTimeout(timeout time.Duration) parser.Option {
    return parser.WrapImplSpecificOptFn(func(o *options) {
        o.Timeout = timeout
    })
}
```

### 实现真实解析逻辑
```go
func (p *CustomParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
    // 处理选项...

    // 读取内容
    content, err := io.ReadAll(reader)
    if err != nil {
        return nil, err
    }

    // 检查大小限制
    if int64(len(content)) > myOpts.MaxSize {
        return nil, fmt.Errorf("content too large")
    }

    return []*schema.Document{{
        Content: string(content),
        Metadata: map[string]any{
            "encoding": myOpts.Encoding,
            "parser":   "CustomParser",
        },
    }}, nil
}
```

## ⚠️ 注意事项

- **资源管理**: 解析器不应关闭传入的 Reader
- **并发安全**: 无状态设计，天然支持并发
- **错误处理**: 及时处理和返回错误

## 🎓 学习价值

通过 CustomParser 学习：
- Eino 解析器接口的标准实现
- 选项模式的设计和运用
- 配置管理的最佳实践
- 可扩展解析器的设计方法

**下一步**: 学习 [ExtParser](../extparser/) 了解如何组合多个解析器。