# JSON Schema 工具示例

本示例展示了如何在 Eino 框架中手动构建 JSON Schema 来创建工具，帮助理解框架的底层工作机制。

## 示例定位

> **注意**: 对于日常开发，推荐使用 [`utils.InferTool()`](../callback/README.md#5-infertool-自动推断机制) 来自动生成 JSON Schema。本示例主要用于学习目的和复杂场景。

## 主要内容

示例实现了一个天气查询工具，完整演示了以下过程：

1. **手动构建 JSON Schema** - 使用 `eino-contrib/jsonschema` 包定义复杂参数结构
2. **创建工具信息对象** - 将 JSON Schema 转换为 Eino 可识别的 ToolInfo
3. **实现工具逻辑** - 模拟 HTTP API 调用获取天气数据
4. **创建工具节点** - 将工具包装为可调用的节点
5. **执行工具调用** - 通过 Invoke 方法执行工具并获取结果

## 核心知识点

### 1. JSON Schema 与 Eino 的集成
- 学会使用 `eino-contrib/jsonschema` 包定义标准 JSON Schema
- 掌握 `schema.NewParamsOneOfByJSONSchema()` 将 JSON Schema 转换为 Eino 参数定义
- 理解 JSON Schema 中的类型、必需字段、属性描述等概念

### 2. 工具创建模式
- 掌握 `utils.NewTool()` 创建工具实例的标准模式
- 学会定义与 JSON Schema 匹配的 Go 结构体参数
- 理解工具函数的签名：`func(ctx context.Context, params *Struct) (string, error)`

### 3. 工具节点组合
- 学会使用 `compose.NewToolNode()` 创建工具节点
- 理解 ToolsNodeConfig 的配置方式
- 掌握多个工具组合到同一个节点的方法

### 4. 工具调用机制
- 理解 `schema.Message` 和 `schema.ToolCall` 的结构
- 学会构造工具调用请求，包括函数名和参数 JSON
- 掌握 `toolsNode.Invoke()` 执行工具调用的方法
- 了解工具调用结果的格式和处理方式

### 5. HTTP API 集成模式
- 学会在工具中模拟 HTTP 请求的完整流程
- 掌握请求头设置、URL 构建、参数编码等技术
- 理解响应数据的序列化和格式化输出

## 使用场景

### 何时需要手动构建 JSON Schema？

1. **复杂 Schema 结构** - 需要嵌套对象、数组、联合类型等复杂结构
2. **自定义验证规则** - 需要超出标准 tag 支持的验证逻辑
3. **多 Schema 组合** - 需要使用 OneOf、AllOf、AnyOf 等复杂关系
4. **动态 Schema 生成** - 运行时根据条件生成不同 Schema
5. **框架学习理解** - 了解 `InferTool` 的内部工作原理

### 日常开发建议

- **简单场景**: 使用 `utils.InferTool()` 自动生成
- **复杂场景**: 参考本示例手动构建
- **调试问题**: 理解底层机制有助于排查问题

## 读者掌握要点

学习本示例后，读者应该能够：

- **理解底层机制** - 掌握 JSON Schema 在 Eino 中的工作原理
- **处理复杂场景** - 在需要时手动构建复杂 Schema 结构
- **调试 Schema 问题** - 理解 `InferTool` 内部工作，便于问题排查
- **扩展框架功能** - 在需要时自定义或扩展框架的 Schema 处理能力
- **选择合适方案** - 根据场景选择使用 `InferTool` 或手动构建

## 运行示例

在 `components/tool/jsonschema/` 目录下执行：

```bash
go run main.go
```

程序将演示完整的工具创建和调用流程。

## 与 InferTool 的对比

```go
// 本示例：手动构建 JSON Schema
weatherSchema := &jsonschema.Schema{
    Type:     string(schema.Object),
    Required: []string{"city"},
    Properties: orderedmap.New[string, *jsonschema.Schema]...,
}
toolInfo := &schema.ToolInfo{
    Name:        "get_weather",
    Desc:        "获取指定城市的天气信息",
    ParamsOneOf: schema.NewParamsOneOfByJSONSchema(weatherSchema),
}
weatherTool := utils.NewTool(toolInfo, weatherFunc)

// 推荐：使用 InferTool 自动生成
type WeatherParams struct {
    City string `json:"city" jsonschema:"required,description=城市名称"`
    Unit string `json:"unit,omitempty" jsonschema:"description=温度单位"`
}
weatherTool, err := utils.InferTool("get_weather", "获取天气", weatherFunc)
```

选择合适的方式取决于具体需求和场景复杂度。

