# Components Document - 文档处理组件

## 📋 概述

`components/document/` 模块是 Eino 框架中用于文档处理的组件集合，提供了多种文档格式的解析、处理和转换功能。该模块是学习路径中的第7个模块，帮助开发者掌握文档处理的基础知识和实际应用。

## 🎯 学习目标

- 掌握多种文档格式的解析技术
- 理解文档加载器和分块策略
- 学习自定义解析器的开发
- 掌握内容提取和转换技术

## 📚 模块定位

根据学习路径规划，本模块属于 **🔧 中级 - 组件深入和功能扩展** 阶段：

- **前置知识**：`quickstart/chat` - 基础聊天功能
- **后续关联**：`components/retriever` - 数据检索系统
- **难度等级**：⭐⭐⭐
- **预计学习时间**：2-3天

## 🏗️ 目录结构

```
components/document/
└── parser/                   # 解析器集合
    ├── customparser/         # 自定义解析器示例
    │   ├── custom_parser.go  # 自定义解析器实现
    │   └── parse.go          # 使用示例
    ├── textparser/           # 文本解析器
    │   └── text_parser.go    # 纯文本解析示例
    ├── extparser/            # 扩展解析器
    │   ├── ext_parser.go     # 多格式解析器示例
    │   └── testdata/         # 测试数据
    │       └── test.html     # HTML 测试文件
    └── README.md             # 本文档
```

## 🔧 核心组件

### 1. 自定义解析器 (CustomParser)
**位置**：`parser/customparser/`

展示如何实现自定义的文档解析器，包括：
- 解析器配置和选项设计
- 通用选项和特定选项的处理
- 解析逻辑的实现

```go
type CustomParser struct {
    defaultEncoding string
    defaultMaxSize  int64
}

func (p *CustomParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error)
```

### 2. 文本解析器 (TextParser)
**位置**：`parser/textparser/`

演示最基础的文本解析功能：
- 使用 Eino 框架内置的文本解析器
- 处理纯文本内容
- 基础的文档模式理解

```go
textParser := parser.TextParser{}
docs, err := textParser.Parse(ctx, strings.NewReader("hello world"))
```

### 3. 扩展解析器 (ExtParser)
**位置**：`parser/extparser/`

展示多格式文档解析能力：
- HTML 解析器集成
- PDF 解析器集成
- 默认解析器配置
- 文件类型自动识别

```go
extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
    Parsers: map[string]parser.Parser{
        ".html": htmlParser,
        ".pdf":  pdfParser,
    },
    FallbackParser: textParser,
})
```

## 🚀 快速开始

### 运行自定义解析器示例

```bash
cd components/document/parser/customparser
go run parse.go
```

### 运行文本解析器示例

```bash
cd components/document/parser/textparser
go run text_parser.go
```

### 运行扩展解析器示例

```bash
cd components/document/parser/extparser
go run ext_parser.go
```

## 📖 学习路径

### 第一阶段：理解基础概念 (0.5天)
1. 阅读 `textparser/text_parser.go` 理解基础解析
2. 学习 `schema.Document` 结构和字段含义
3. 理解解析器的接口设计

### 第二阶段：自定义解析器开发 (1天)
1. 深入研究 `customparser/custom_parser.go`
2. 理解选项模式的实现
3. 掌握解析器配置和扩展方法

### 第三阶段：多格式处理 (1天)
1. 学习 `extparser/ext_parser.go` 的设计
2. 理解扩展解析器的注册机制
3. 掌握文件类型识别和路由

### 第四阶段：实践应用 (0.5天)
1. 尝试实现自己的解析器
2. 测试不同文档格式的处理
3. 集成到实际项目中

## 🔑 关键概念

### Parser 接口
所有解析器都需要实现 `parser.Parser` 接口：
```go
type Parser interface {
    Parse(ctx context.Context, reader io.Reader, opts ...Option) ([]*schema.Document, error)
}
```

### 选项模式
解析器使用选项模式进行配置：
- **通用选项**：所有解析器都支持的选项
- **特定选项**：特定解析器独有的选项
- **配置传递**：通过 `Option` 函数传递参数

### ExtParser 机制
扩展解析器提供了格式自动识别能力：
- 根据文件扩展名选择解析器
- 支持默认解析器处理未知格式
- 可以注册多种解析器

## 📊 学习检查点

完成本模块学习后，你应该能够：

- [ ] 理解 Eino 文档处理组件的架构设计
- [ ] 能够使用内置的文本解析器处理文本内容
- [ ] 掌握自定义解析器的开发方法
- [ ] 理解扩展解析器的工作原理
- [ ] 能够实现自己的文档解析器
- [ ] 掌握解析器选项和配置的使用

## 🔗 相关模块

- **`components/retriever`** - 文档检索系统，依赖本模块的解析能力
- **`quickstart/chat`** - 基础聊天功能，展示了文档内容的简单应用
- **`components/lambda`** - 函数组件，可与文档解析结合使用

## 🛠️ 扩展练习

1. **实现 Markdown 解析器**
   - 添加 Markdown 格式支持
   - 处理标题、列表、代码块等结构

2. **批量文档处理**
   - 实现文件夹遍历和批量解析
   - 添加进度监控和错误处理

3. **文档元数据提取**
   - 提取文档标题、作者、创建时间等元信息
   - 实现文档摘要生成

4. **集成到实际应用**
   - 将解析器集成到聊天机器人中
   - 实现文档问答功能

## 📝 注意事项

1. **依赖管理**：本模块使用了 `eino-ext` 扩展包，确保正确安装依赖
2. **测试数据**：`extparser/testdata/` 包含测试用的 HTML 文件
3. **错误处理**：注意解析过程中的错误处理和资源管理
4. **性能考虑**：大文件解析时注意内存使用和处理效率

## 🔗 参考资料

- [Eino 官方文档 - Document 组件](https://github.com/cloudwego/eino)
- [Eino 扩展组件库](https://github.com/cloudwego/eino-ext)
- [Go 文件处理最佳实践](https://golang.org/pkg/io/)

---

**提示**：建议按照学习路径逐步深入，先理解基础概念，再进行实践操作。遇到问题时可以参考其他模块的实现方式。