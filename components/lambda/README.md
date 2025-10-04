# Lambda 组件 - 自定义函数编排

## 📚 目录

- [📋 概述](#-概述)
- [🎯 四种交互模式](#-四种交互模式)
  - [1. Invoke 模式](#1-invoke-模式)
  - [2. Stream 模式](#2-stream-模式)
  - [3. Collect 模式](#3-collect-模式)
  - [4. Transform 模式](#4-transform-模式)
- [🛠️ 构建方法](#️-构建方法)
  - [单一交互模式构建](#单一交互模式构建)
    - [1. 不带自定义 Option](#1-不带自定义-option)
    - [2. 使用自定义 Option](#2-使用自定义-option)
  - [AnyLambda - 多模式组合](#anylambda---多模式组合)
- [🔧 内置 Lambda 组件](#-内置-lambda-组件)
  - [1. ToList Lambda](#1-tolist-lambda)
  - [2. MessageParser Lambda](#2-messageparser-lambda)
- [🏗️ 编排集成](#️-编排集成)
  - [Graph 中使用](#graph-中使用)
  - [Chain 中使用](#chain-中使用)
- [🎯 主要使用场景](#-主要使用场景)
  - [1. 数据预处理与后处理](#1-数据预处理与后处理)
  - [2. 外部服务集成](#2-外部服务集成)
  - [3. 业务逻辑封装](#3-业务逻辑封装)
  - [4. 流程编排增强](#4-流程编排增强)
- [💡 典型应用示例](#-典型应用示例)
- [📚 学习要点](#-学习要点)
- [📋 Lambda API 完整参考](#-lambda-api-完整参考)
  - [🎯 核心创建函数（8个）](#-核心创建函数8个)
    - [单一模式创建 - 带 Option](#单一模式创建---带-option)
    - [单一模式创建 - 无 Option](#单一模式创建---无-option)
    - [多模式组合](#多模式组合)
  - [🔧 配置选项（2个）](#-配置选项2个)
  - [🏗️ 内置组件（2个）](#-内置组件2个)
  - [📊 编排集成方法（6个）](#-编排集成方法6个)
  - [📝 函数类型定义](#-函数类型定义)
- [🔗 相关资源](#-相关资源)

## 📋 概述

Lambda 组件是 Eino 框架中的核心功能，允许开发者将自定义函数无缝集成到 AI 编排流程中。通过 Lambda，你可以：

- **灵活转换数据**：在链路中进行任意类型的数据转换
- **集成外部服务**：调用第三方 API 或执行复杂的业务逻辑
- **流式处理**：支持流式数据的实时转换和处理
- **类型安全**：保持强类型检查，确保数据流转的安全性

## 🎯 四种交互模式

Lambda 底层由输入输出是否为流所形成的 4 种运行函数组成：

### 1. Invoke 模式
```go
func(ctx context.Context, input I, opts ...TOption) (output O, err error)
```
- **特点**：同步调用，一次输入一次输出
- **适用场景**：简单的数据转换、计算任务、API 调用

### 2. Stream 模式
```go
func(ctx context.Context, input I, opts ...TOption) (output *schema.StreamReader[O], err error)
```
- **特点**：同步输入，流式输出
- **适用场景**：生成连续数据流，如文本生成、数据流处理

### 3. Collect 模式
```go
func(ctx context.Context, input *schema.StreamReader[I], opts ...TOption) (output O, err error)
```
- **特点**：流式输入，同步输出
- **适用场景**：聚合流数据为单个结果，如数据汇总、流数据收集

### 4. Transform 模式
```go
func(ctx context.Context, input *schema.StreamReader[I], opts ...TOption) (output *schema.StreamReader[O], err error)
```
- **特点**：流式输入，流式输出
- **适用场景**：实时流数据处理和转换

## 🛠️ 构建方法

### 单一交互模式构建

#### 1. 不带自定义 Option

**InvokableLambda**
```go
lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
    // some logic
})
```

**StreamableLambda**
```go
lambda := compose.StreamableLambda(func(ctx context.Context, input string) (output *schema.StreamReader[string], err error) {
    // some logic
})
```

**CollectableLambda**
```go
lambda := compose.CollectableLambda(func(ctx context.Context, input *schema.StreamReader[string]) (output string, err error) {
    // some logic
})
```

**TransformableLambda**
```go
lambda := compose.TransformableLambda(func(ctx context.Context, input *schema.StreamReader[string]) (output *schema.StreamReader[string], err error) {
    // some logic
})
```

#### 2. 使用自定义 Option

```go
type Options struct {
    Field1 string
}
type MyOption func(*Options)

lambda := compose.InvokableLambdaWithOption(
    func(ctx context.Context, input string, opts ...MyOption) (output string, err error) {
        // 处理 opts
        // some logic
    }
)
```

### AnyLambda - 多模式组合

最灵活的创建方式，允许同时实现多种交互模式：

```go
lambda, err := compose.AnyLambda(
    // Invoke 函数
    func(ctx context.Context, input string, opts ...MyOption) (output string, err error) {
        // some logic
    },
    // Stream 函数
    func(ctx context.Context, input string, opts ...MyOption) (output *schema.StreamReader[string], err error) {
        // some logic
    },
    // Collect 函数
    func(ctx context.Context, input *schema.StreamReader[string], opts ...MyOption) (output string, err error) {
        // some logic
    },
    // Transform 函数
    func(ctx context.Context, input *schema.StreamReader[string], opts ...MyOption) (output *schema.StreamReader[string], err error) {
        // some logic
    },
)
```

## 🔧 内置 Lambda 组件

### 1. ToList Lambda
将单个输入元素转换为包含该元素的切片（数组）：

```go
// 创建一个 ToList Lambda
lambda := compose.ToList[*schema.Message]()

// 在 Chain 中使用
chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
chain.AppendChatModel(chatModel)  // chatModel 返回 *schema.Message
chain.AppendLambda(lambda)        // 将 *schema.Message 转换为 []*schema.Message
```

### 2. MessageParser Lambda
将 JSON 消息（通常由 LLM 生成）解析为指定的结构体：

```go
// 定义解析目标结构体
type MyStruct struct {
    ID int `json:"id"`
}

// 创建解析器
parser := schema.NewMessageJSONParser[*MyStruct](&schema.MessageJSONParseConfig{
    ParseFrom:    schema.MessageParseFromContent,
    ParseKeyPath: "", // 如果仅需要 parse 子字段，可用 "key.sub.grandsub"
})

// 创建解析 Lambda
parserLambda := compose.MessageParser(parser)

// 在 Chain 中使用
chain := compose.NewChain[*schema.Message, *MyStruct]()
chain.AppendLambda(parserLambda)

// 使用示例
runner, err := chain.Compile(context.Background())
parsed, err := runner.Invoke(context.Background(), &schema.Message{
    Content: `{"id": 1}`,
})
// parsed.ID == 1
```

**支持的解析来源：**
- `schema.MessageParseFromContent` - 从消息内容解析
- `schema.MessageParseFromToolCall` - 从工具调用结果解析

## 🏗️ 编排集成

### Graph 中使用
```go
graph := compose.NewGraph[string, *MyStruct]()
graph.AddLambdaNode(
    "node1",
    compose.InvokableLambda(func(ctx context.Context, input string) (*MyStruct, error) {
        // some logic
        return &MyStruct{ID: 1}, nil
    }),
)
```

### Chain 中使用
```go
chain := compose.NewChain[string, string]()
chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
    // some logic
    return "", nil
}))
```

## 🎯 主要使用场景

### 1. 数据预处理与后处理
- 格式转换（如 JSON 解析、数据清洗）
- 类型映射（将消息转换为结构化数据）
- 内容过滤和验证

### 2. 外部服务集成
- API 调用和数据获取
- 数据库查询和缓存操作
- 文件处理和存储操作

### 3. 业务逻辑封装
- 复杂计算和算法实现
- 决策逻辑和规则引擎
- 状态管理和持久化

### 4. 流程编排增强
- 在 Chain 和 Graph 中插入自定义处理节点
- 实现条件分支和循环逻辑
- 数据聚合和拆分操作

## 💡 典型应用示例

### 消息解析
将 AI 生成的文本内容解析为结构化数据，便于后续处理和存储。

### API 调用封装
将外部 API 调用封装为 Lambda 节点，统一错误处理和重试逻辑。

### 数据格式转换
在不同系统组件之间进行数据格式转换，确保数据兼容性。

### 意图识别
使用 MessageParser 从 LLM 的工具调用结果中解析意图信息。

## 📚 学习要点

- 理解 Lambda 的四种操作模式及其适用场景
- 掌握不同构建方法的选择和使用
- 学会使用内置 Lambda 简化常见任务
- 了解流式处理的实现方式
- 掌握在 Chain 和 Graph 中的集成方法

## 📋 Lambda API 完整参考

### 🎯 核心创建函数（8个）

#### 单一模式创建 - 带 Option
```go
// 创建支持自定义选项的 Invoke Lambda
InvokableLambdaWithOption[I, O, TOption](i Invoke[I, O, TOption], opts ...LambdaOpt) *Lambda

// 创建支持自定义选项的 Stream Lambda
StreamableLambdaWithOption[I, O, TOption](s Stream[I, O, TOption], opts ...LambdaOpt) *Lambda

// 创建支持自定义选项的 Collect Lambda
CollectableLambdaWithOption[I, O, TOption](c Collect[I, O, TOption], opts ...LambdaOpt) *Lambda

// 创建支持自定义选项的 Transform Lambda
TransformableLambdaWithOption[I, O, TOption](t Transform[I, O, TOption], opts ...LambdaOpt) *Lambda
```

#### 单一模式创建 - 无 Option
```go
// 创建不带选项的 Invoke Lambda
InvokableLambda[I, O](i InvokeWOOpt[I, O], opts ...LambdaOpt) *Lambda

// 创建不带选项的 Stream Lambda
StreamableLambda[I, O](s StreamWOOpt[I, O], opts ...LambdaOpt) *Lambda

// 创建不带选项的 Collect Lambda
CollectableLambda[I, O](c CollectWOOpt[I, O], opts ...LambdaOpt) *Lambda

// 创建不带选项的 Transform Lambda
TransformableLambda[I, O](t TransformWOOpts[I, O], opts ...LambdaOpt) *Lambda
```

#### 多模式组合
```go
// 最灵活的创建方式，可同时实现多种交互模式
AnyLambda[I, O, TOption](i Invoke, s Stream, c Collect, t Transform, opts ...LambdaOpt) (*Lambda, error)
```

### 🔧 配置选项（2个）
```go
// 启用 Lambda 函数的回调功能
WithLambdaCallbackEnable(y bool) LambdaOpt

// 设置 Lambda 函数的类型标识
WithLambdaType(t string) LambdaOpt
```

### 🏗️ 内置组件（2个）
```go
// 将单个输入转换为包含该元素的切片
ToList[I any](opts ...LambdaOpt) *Lambda

// 将消息解析为指定结构体
MessageParser[T any](p schema.MessageParser[T], opts ...LambdaOpt) *Lambda
```

### 📊 编排集成方法（6个）
```go
// Graph 中添加 Lambda 节点
AddLambdaNode(key string, node *Lambda, opts ...GraphAddNodeOpt) error

// Chain 中追加 Lambda 节点
AppendLambda(node *Lambda, opts ...GraphAddNodeOpt) *Chain[I, O]

// Parallel 中添加 Lambda 节点
AddLambda(outputKey string, node *Lambda, opts ...GraphAddNodeOpt) *Parallel

// ChainBranch 中添加 Lambda 节点
AddLambda(key string, node *Lambda, opts ...GraphAddNodeOpt) *ChainBranch

// Workflow 中添加 Lambda 节点
AddLambdaNode(key string, lambda *Lambda, opts ...GraphAddNodeOpt) *WorkflowNode

// Graph 调用时传递 Lambda 选项
WithLambdaOption(opts ...any) Option
```

### 📝 函数类型定义
```go
// 四种核心交互模式的函数类型
type Invoke[I, O, TOption any] func(ctx context.Context, input I, opts ...TOption) (output O, err error)
type Stream[I, O, TOption any] func(ctx context.Context, input I, opts ...TOption) (output *schema.StreamReader[O], err error)
type Collect[I, O, TOption any] func(ctx context.Context, input *schema.StreamReader[I], opts ...TOption) (output O, err error)
type Transform[I, O, TOption any] func(ctx context.Context, input *schema.StreamReader[I], opts ...TOption) (output *schema.StreamReader[O], err error)

// 无选项版本的函数类型
type InvokeWOOpt[I, O any] func(ctx context.Context, input I) (output O, err error)
type StreamWOOpt[I, O any] func(ctx context.Context, input I) (output *schema.StreamReader[O], err error)
type CollectWOOpt[I, O any] func(ctx context.Context, input *schema.StreamReader[I]) (output O, err error)
type TransformWOOpts[I, O any] func(ctx context.Context, input *schema.StreamReader[I]) (output *schema.StreamReader[O], err error)
```

## 🔗 相关资源

- [官方文档](https://www.cloudwego.io/zh/docs/eino/core_modules/components/lambda_guide/)
- [示例代码](https://github.com/cloudwego/eino-examples/blob/main/components/lambda)
- 源码位置：`eino/compose/types_lambda.go`
- API 包：`github.com/cloudwego/eino/compose`
- 前置知识：`components/tool`、`quickstart/todoagent`

---

**提示**：Lambda 组件是构建复杂 AI 应用的基础工具，建议先从简单的 InvokableLambda 数据转换场景开始练习，逐步掌握流式处理和多模式组合。总计 20 个 Lambda 相关的 API 函数，覆盖了创建、配置、集成等完整的使用流程。