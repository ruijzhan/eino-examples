`InvokableTool` 是 Eino 框架中的一个接口，定义了在意图识别和执行期间可由 `ChatModel` 调用的工具。

它使 AI 代理能够基于用户输入调用外部函数，通过标准化工具的描述和执行方式——特别是通过 JSON 序列化的参数和结构化响应。

---

### 定义

`InvokableTool` 接口扩展了 `BaseTool`，要求实现提供元数据（`Info`）和执行方法（`InvokableRun`）。它主要用于 AI 代理系统，其中语言模型决定调用哪个工具。

```go
31:36:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/components/tool/interface.go
// InvokableTool 用于 ChatModel 意图识别和 ToolsNode 执行的工具
type InvokableTool interface {
	BaseTool

	// InvokableRun 使用 JSON 格式的参数调用函数
	InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
}
```

- **扩展**: `BaseTool` — 必须实现 `Info(ctx) (*schema.ToolInfo, error)` 以向模型暴露工具元数据（名称、描述、参数）。
- **方法**: `InvokableRun` — 使用 JSON 格式的输入参数执行工具。
- **参数**:
  - `ctx context.Context`: 标准 Go 上下文，用于取消和超时。
  - `argumentsInJSON string`: 序列化为 JSON 的输入参数。
  - `opts ...Option`: 工具的可选配置（例如自定义行为、调试）。
- **返回**: `string, error` — 结果为 JSON 字符串，如果执行失败则返回错误。
- **使用上下文**: 用于代理管道，其中 `ChatModel` 通过 `ToolsNode` 选择和执行工具。

此接口对于非流式工具调用至关重要——与返回流的 `StreamableTool` 相对比。

---

### 使用示例

一个常见的模式是通过 `utils.InferOptionableTool` 定义工具，它从类型化函数生成 `InvokableTool`。

```go
42:59:/home/ruijzhan/data/eino-examples/adk/intro/chatmodel/subagents/ask_for_clarification.go
func NewAskForClarificationTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"ask_for_clarification",
		"当用户请求不明确时调用此工具...",
		func(ctx context.Context, input *AskForClarificationInput, opts ...tool.Option) (output string, err error) {
			o := tool.GetImplSpecificOptions[askForClarificationOptions](nil, opts...)
			if o.NewInput == nil {
				return "", compose.NewInterruptAndRerunErr(input.Question)
			}
			output = *o.NewInput
			o.NewInput = nil
			return output, nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return t
}
```

此工具：
- 接受问题输入（`AskForClarificationInput`）。
- 使用 `Option` 在恢复期间注入新输入（例如用户响应后）。
- 返回中断执行并提示用户的错误——显示高级控制流。

在代理管道中，它与其他工具一起注册：

```go
42:43:/home/ruijzhan/data/eino-examples/adk/intro/chatmodel/subagents/agent.go
Tools: []tool.BaseTool{NewBookRecommender(), NewAskForClarificationTool()},
```

**在代码库中的总体使用**:
- 在 26 个文件中找到。
- 在示例中广泛使用：`todoagent`、`deer-go`、`plan-execute-replan` 等。
- 是 `eino` 中基于代理的工作流的核心，特别是在 `React` 和 `ChatModel` 代理中。
- 像 DuckDuckGo 搜索、餐厅/菜品查询和书籍推荐等工具都实现 `InvokableTool`。

---

### 注意事项

- 大多数开发人员不直接实现 `InvokableTool`——而是定义类型化函数并使用 `InferTool` 或 `InferOptionableTool` 生成符合接口的包装器。
- 尽管被称为"可调用"，但它支持通过 `compose.NewInterruptAndRerunErr` 进行**中断和恢复**等高级模式，实现对话来回。
- 输入和输出都使用 `string`（而不是泛型类型）确保与基于 JSON 的模型工具调用 API 兼容。

---

### 另请参见

- `BaseTool`: 提供 `Info()` 的父接口——模型用它来了解可用的工具。
- `StreamableTool`: 用于返回流式结果的工具的兄弟接口（例如长时间运行的查询）。
- `utils.InferTool`: 通过反射其输入/输出类型将 Go 函数转换为 `InvokableTool` 的实用函数。
- `schema.ToolInfo`: 描述工具名称、描述和参数 JSON 模式的结构体——在 `Info()` 和模型提示中使用。
- `compose.ToolsNode`: 在代理工作流中执行 `InvokableTool` 实例的组件。