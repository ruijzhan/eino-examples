`ToolsNode` 是 Eino 框架中的一个图节点，执行来自 LLM 生成消息的工具调用。
它处理 `AssistantMessage` 中的 `ToolCalls`，调用相应的工具，并返回结构化的 `ToolMessage` 响应。

---

### 定义

`ToolsNode` 结构体在 `tool_node.go` 中定义，表示计算图中能够执行工具的节点——通常是函数或 API——基于 LLM 输出。它支持多个工具调用的顺序和并行执行，并与 Eino 的有状态图系统集成。

```go
70:75:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/compose/tool_node.go
type ToolsNode struct {
	tuple                *toolsTuple               // 已注册工具的索引和元数据
	unknownToolHandler   func(ctx context.Context, name, input string) (string, error)  // 处理对不存在工具的调用
	executeSequentially  bool                      // 如果为 true，工具依次运行
	toolArgumentsHandler func(ctx context.Context, name, input string) (string, error)  // 预处理工具参数
}
```

- **目的**: 执行图工作流中 LLM 生成的 `ToolCalls` 引用的工具。
- **输入**: 包含 `ToolCalls` 列表的 `AssistantMessage` 类型的 `*schema.Message`。
- **输出**: `ToolMessage` 类型的 `*schema.Message` 切片，每个包含工具执行的结果。
- **执行模式**: 可以并行（默认）或通过 `ExecuteSequentially` 配置顺序运行工具。
- **错误处理**: 如果工具不存在且未设置 `UnknownToolsHandler`，则返回错误。否则，处理程序提供回退响应。

`ToolsNode` 使用 `NewToolNode` 构造，它接受 `ToolsNodeConfig` 并初始化内部元数据以进行高效的工具查找和执行。

```go
119:130:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/compose/tool_node.go
func NewToolNode(ctx context.Context, conf *ToolsNodeConfig) (*ToolsNode, error) {
	tuple, err := convTools(ctx, conf.Tools)
	if err != nil {
		return nil, err
	}
	return &ToolsNode{
		tuple:                tuple,
		unknownToolHandler:   conf.UnknownToolsHandler,
		executeSequentially:  conf.ExecuteSequentially,
		toolArgumentsHandler: conf.ToolArgumentsHandler,
	}, nil
}
```

- **参数**:
  - `ctx`: 用于初始化和工具信息检索的上下文。
  - `conf`: 指定工具、处理程序和执行行为的配置。
- **返回**: 成功时返回 `*ToolsNode`，如果工具转换失败则返回错误。
- **副作用**: 直接没有，但工具执行可能根据工具实现产生副作用。

---

### 使用示例

一个常见的用法是在 ReAct 风格的代理循环中，LLM 决定调用工具，`ToolsNode` 执行它们，然后将结果反馈给模型。

在 `react_with_interrupt/main.go` 中，创建 `ToolsNode` 并将其添加到编排工具使用和人工在环中断的对话流图中。

```go
153:160:/home/ruijzhan/data/eino-examples/compose/graph/react_with_interrupt/main.go
func newToolsNode(ctx context.Context) *compose.ToolsNode {
	tools := getTools()
	tn, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{Tools: tools})
	if err != nil {
		log.Fatal(err)
	}
	return tn
}
```

后来，它被集成到图中：

```go
236:241:/home/ruijzhan/data/eino-examples/compose/graph/react_with_interrupt/main.go
err = g.AddToolsNode(
	NodeKeyToolsNode,
	toolsNode,
	compose.WithNodeName(NodeKeyToolsNode),
	compose.WithStatePostHandler(appendNextPrompt(ctx, browserTool)),
)
```

- **主要调用者**: `ToolsNode` 广泛用于 `react`、`multiagent` 和 `manus` 等代理框架中。
- **集成点**: 通过 `AddToolsNode` 添加到图中，通过 `AppendToolsNode` 在链中使用，支持流式和批处理执行。
- **广泛使用**: 在至少 24 个文件中找到，主要用于代理逻辑和图组合层。

---

### 注意事项

- **工具解析是静态的**: 工具在初始化时通过 `convTools` 注册，它构建索引（`map[string]int`）以便在执行期间快速查找。不支持动态工具注册。
- **处理程序灵活性**: `UnknownToolsHandler` 允许优雅处理幻觉工具调用——当 LLM 发明不存在的工具时，这对稳健性至关重要。
- **参数预处理**: `ToolArgumentsHandler` 允许在执行之前清理或转换工具输入，对安全性或兼容性很有用。

---

### 另请参见

- `ToolsNodeConfig`: 用于初始化 `ToolsNode` 的配置结构，定义工具、处理程序和执行模式。
- `tool.BaseTool`: 所有工具必须实现的接口；提供元数据和执行逻辑。
- `graph.AddToolsNode`: 将 `ToolsNode` 集成到给定键下的计算图中的方法。
- `NewToolNode`: 从配置构造 `ToolsNode` 的工厂函数。
- `schema.ToolCall`: 表示来自 LLM 的单个工具调用请求，由 `ToolsNode` 处理。