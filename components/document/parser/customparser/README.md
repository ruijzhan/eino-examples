# CustomParser - 自定义解析器详解

## 📋 概述

`CustomParser` 是 Eino 框架中自定义解析器的完整示例，展示了如何实现一个符合 Eino 标准的文档解析器。这个例子详细演示了解析器的核心设计模式、选项处理机制和扩展方法。

## 🎯 学习目标

- 理解 Eino 解析器接口的设计理念
- 掌握选项模式的实现方法
- 学会处理通用选项和特定选项
- 了解解析器的配置和初始化过程

## 📁 文件结构

```
customparser/
├── README.md           # 本文档
├── custom_parser.go    # 解析器核心实现
└── parse.go           # 使用示例和测试代码
```

## 🔧 核心组件分析

### 1. 解析器结构体 (CustomParser)

```go
type CustomParser struct {
    defaultEncoding string  // 默认编码格式
    defaultMaxSize  int64   // 默认最大文件大小
}
```

**设计要点**：
- 保持简单的状态管理
- 只存储配置信息，不维护解析状态
- 解析是无状态的，每次调用都是独立的

### 2. 配置结构 (Config & options)

#### 初始化配置
```go
type Config struct {
    DefaultEncoding string  // 解析器创建时的默认编码
    DefaultMaxSize  int64   // 解析器创建时的默认最大大小
}
```

#### 运行时选项
```go
type options struct {
    Encoding string  // 每次解析时可以覆盖的编码
    MaxSize  int64   // 每次解析时可以覆盖的最大大小
}
```

**双层配置的设计理念**：
- **Config**：解析器级别的配置，创建后不可更改
- **options**：调用级别的配置，每次解析可以不同

### 3. 选项函数 (Option Functions)

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

**关键设计模式**：
- 使用 `parser.WrapImplSpecificOptFn` 包装特定选项
- 返回 `parser.Option` 接口，保持统一的调用方式
- 通过闭包修改配置结构体

### 4. 核心解析方法 (Parse)

```go
func (p *CustomParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error)
```

**解析流程详解**：

#### 步骤1：处理通用选项
```go
commonOpts := parser.GetCommonOptions(&parser.Options{}, opts...)
```
- 提取所有解析器都支持的通用选项
- 如文件路径、元数据等

#### 步骤2：处理特定选项
```go
myOpts := &options{
    Encoding: p.defaultEncoding,  // 使用解析器默认值
    MaxSize:  p.defaultMaxSize,
}
myOpts = parser.GetImplSpecificOptions(myOpts, opts...)
```
- 初始化特定选项的默认值
- 应用用户传入的特定选项
- 实现配置的动态覆盖

#### 步骤3：实现解析逻辑
```go
return []*schema.Document{{
    Content: "Hello World",
}}, nil
```
- 当前示例只是返回固定内容
- 实际应用中这里应该是复杂的解析逻辑

## 🚀 使用示例

### 基本使用

```go
// 1. 创建解析器
customParser, err := NewCustomParser(&Config{
    DefaultEncoding: "utf-8",
    DefaultMaxSize:  1024,
})

// 2. 执行解析
docs, err := customParser.Parse(ctx, reader)
```

### 使用自定义选项

```go
// 覆盖默认配置
docs, err := customParser.Parse(ctx, reader,
    WithMaxSize(2048),        // 临时修改最大大小
    WithEncoding("gbk"),      // 临时修改编码
)
```

## 📚 关键概念解析

### 1. 选项模式 (Option Pattern)

选项模式是 Go 语言中处理可选参数的常用模式：

**传统方式的问题**：
```go
// 参数过多，调用复杂
func NewParser(encoding string, maxSize int64, timeout time.Duration, retries int) // ...
```

**选项模式的解决方案**：
```go
func NewParser(opts ...Option) *Parser // 简洁的函数签名
parser := NewParser(WithEncoding("utf-8"), WithMaxSize(1024)) // 灵活的调用
```

### 2. 接口隔离

Eino 使用 `parser.Option` 接口来统一选项类型：

```go
type Option interface {
    // 内部实现，用户无需关心
}
```

这样的设计允许：
- 统一的选项传递方式
- 不同解析器可以有不同类型的选项
- 类型安全的选项处理

### 3. 配置优先级

配置的优先级从高到低：
1. **运行时选项**：`WithEncoding("gbk")`
2. **解析器默认值**：`Config.DefaultEncoding`
3. **框架默认值**：空字符串或零值

## 🔍 代码逐行分析

### custom_parser.go 关键代码

#### 第27-32行：options 结构体
```go
type options struct {
    Encoding string
    MaxSize  int64
}
```
- 定义了解析器的特定选项
- 这些选项只对 CustomParser 有效

#### 第36-40行：WithEncoding 函数
```go
func WithEncoding(encoding string) parser.Option {
    return parser.WrapImplSpecificOptFn(func(o *options) {
        o.Encoding = encoding
    })
}
```
- 创建一个修改编码的选项
- 使用 `WrapImplSpecificOptFn` 包装
- 通过闭包修改 options 结构体

#### 第65-82行：Parse 方法核心逻辑
```go
func (p *CustomParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
    // 1. 处理通用选项
    commonOpts := parser.GetCommonOptions(&parser.Options{}, opts...)
    _ = commonOpts

    // 2. 处理特定选项
    myOpts := &options{
        Encoding: p.defaultEncoding,
        MaxSize:  p.defaultMaxSize,
    }
    myOpts = parser.GetImplSpecificOptions(myOpts, opts...)
    _ = myOpts

    // 3. 实现解析逻辑
    return []*schema.Document{{
        Content: "Hello World",
    }}, nil
}
```

### parse.go 关键代码

#### 第28-31行：解析器初始化
```go
customParser, err := NewCustomParser(&Config{
    DefaultEncoding: "default",
    DefaultMaxSize:  1024,
})
```
- 创建解析器实例
- 设置默认配置

#### 第37-39行：使用自定义选项
```go
docs, err := customParser.Parse(ctx, nil,
    WithMaxSize(2048),
)
```
- 临时修改最大大小配置
- 不会影响解析器的默认配置

## 🛠️ 扩展指南

### 1. 添加新的选项

```go
// 1. 在 options 结构体中添加字段
type options struct {
    Encoding    string
    MaxSize     int64
    Timeout     time.Duration  // 新增
}

// 2. 创建对应的 Option 函数
func WithTimeout(timeout time.Duration) parser.Option {
    return parser.WrapImplSpecificOptFn(func(o *options) {
        o.Timeout = timeout
    })
}
```

### 2. 实现真实的解析逻辑

```go
func (p *CustomParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
    // 处理选项（如前所述）

    // 实际的解析逻辑
    content, err := io.ReadAll(reader)
    if err != nil {
        return nil, err
    }

    // 检查大小限制
    if int64(len(content)) > myOpts.MaxSize {
        return nil, fmt.Errorf("content too large: %d > %d", len(content), myOpts.MaxSize)
    }

    // 编码转换
    if myOpts.Encoding != "utf-8" {
        // 执行编码转换逻辑
    }

    return []*schema.Document{{
        Content: string(content),
    }}, nil
}
```

### 3. 添加元数据支持

```go
return []*schema.Document{{
    Content: string(content),
    Metadata: map[string]any{
        "encoding": myOpts.Encoding,
        "size":     len(content),
        "parser":   "CustomParser",
    },
}}, nil
```

## 🔗 与其他组件的关系

### 1. 与 ExtParser 的集成

```go
// CustomParser 可以注册到 ExtParser 中
extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".custom": customParser,  // 注册自定义解析器
    },
    FallbackParser: textParser,
})
```

### 2. 与 Retriever 的配合

解析后的文档可以直接用于检索：
```go
docs, _ := customParser.Parse(ctx, reader)
retriever.AddDocuments(docs)  // 添加到检索系统
```

## ⚠️ 常见问题和注意事项

### 1. 错误处理

```go
// ❌ 错误：忽略错误
_ = commonOpts

// ✅ 正确：处理错误
commonOpts, err := parser.GetCommonOptions(&parser.Options{}, opts...)
if err != nil {
    return nil, err
}
```

### 2. 资源管理

```go
// ❌ 错误：没有关闭 reader
func (p *CustomParser) Parse(ctx context.Context, reader io.Reader, ...) {
    // 使用 reader 但不关闭
}

// ✅ 正确：由调用者负责资源管理
// 解析器不应该关闭传入的 reader
```

### 3. 并发安全

```go
// CustomParser 是并发安全的，因为：
// 1. 没有可变状态
// 2. 每次调用都有独立的 options
// 3. 不存储解析过程中的数据
```

## 🧪 测试建议

### 1. 单元测试

```go
func TestCustomParser_Parse(t *testing.T) {
    parser := NewCustomParser(&Config{
        DefaultEncoding: "utf-8",
        DefaultMaxSize:  1024,
    })

    docs, err := parser.Parse(context.Background(), strings.NewReader("test"))
    assert.NoError(t, err)
    assert.Len(t, docs, 1)
    assert.Equal(t, "Hello World", docs[0].Content)
}
```

### 2. 选项测试

```go
func TestCustomParser_WithOptions(t *testing.T) {
    parser := NewCustomParser(&Config{DefaultMaxSize: 100})

    // 测试选项覆盖
    docs, err := parser.Parse(context.Background(), largeReader, WithMaxSize(1000))
    assert.NoError(t, err)
}
```

## 📈 性能优化建议

1. **避免不必要的字符串拷贝**
2. **使用缓冲读取大文件**
3. **实现流式解析**（支持大文件）
4. **添加并发控制**（如果解析过程很耗时）

## 🎓 总结

`CustomParser` 展示了 Eino 框架中解析器的标准实现模式：

1. **接口设计**：遵循 `parser.Parser` 接口
2. **配置管理**：双层配置系统（Config + options）
3. **选项模式**：类型安全的参数传递
4. **无状态设计**：并发安全的解析器实现
5. **扩展性**：易于添加新功能和新选项

掌握了这个示例，你就具备了实现任意自定义解析器的基础能力！

---

**下一步学习**：建议继续学习 `extparser` 示例，了解如何将多个解析器组合使用。