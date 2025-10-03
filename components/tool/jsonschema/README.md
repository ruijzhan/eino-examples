# JSON Schema 工具示例

本示例展示了如何使用 JSON Schema 在 Eino 框架中创建工具，并演示完整的工具调用流程。

## 主要内容

示例实现了一个天气查询工具，完整演示了以下过程：

1. **使用 JSON Schema 定义工具参数** - 定义城市和温度单位参数
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

## 读者掌握要点

学习本示例后，读者应该能够：

- **独立创建** 基于 JSON Schema 的 Eino 工具
- **理解并应用** 工具调用的完整流程
- **集成外部 API** 到 Eino 框架中
- **处理参数转换** 和错误处理机制
- **构建复杂的工具组合** 和调用链

## 运行示例

在 `components/tool/jsonschema/` 目录下执行：

```bash
go run main.go
```

程序将演示完整的工具创建和调用流程。

