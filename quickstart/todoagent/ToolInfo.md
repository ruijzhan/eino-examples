`ToolInfo` 是一个结构体，定义了 AI 模型可用工具的元数据和参数模式。
它通过描述工具的目的、输入和行为，使模型能够理解如何调用外部函数。

---

### 定义

`ToolInfo` 是在 `tool.go` 中定义的 Go 结构体，封装了语言模型与工具（即外部函数或 API）交互所需的所有信息。它包括工具的名称、描述、可选的额外数据，以及关键的输入参数——通过高级 `ParameterInfo` 映射或 JSON Schema、OpenAPI 等正式模式。

```go
61:76:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/schema/tool.go
type ToolInfo struct {
	// 清晰传达工具目的的唯一名称。
	Name string
	// 用于告诉模型如何/何时/为什么使用工具。
	// 你可以提供少量示例作为描述的一部分。
	Desc string
	// Extra 是工具的额外信息。
	Extra map[string]any

	// 函数接受的参数。
	// 可以通过以下方式描述：
	//   - schema.NewParamsOneOfByParams(params)
	//   - schema.NewParamsOneOfByJSONSchema(js)
	// 如果为 nil，则工具不需要输入。
	*ParamsOneOf
}
```

- **字段**:
  - `Name`: 工具的唯一标识符（例如 `"todo_manager"`）。
  - `Desc`: 模型的人类可读解释，用于决定何时以及如何使用工具。可能包括使用示例。
  - `Extra`: 用于自定义处理或路由的可选元数据。
  - `*ParamsOneOf`: 参数模式的嵌入指针；支持多种描述格式。

该结构体使用 `*ParamsOneOf` 的**结构体嵌入**，允许 `ToolInfo` 直接访问 `ToJSONSchema()` 等方法，而无需显式委托。这种设计实现了灵活的参数规范，同时保持了清晰的接口。

---

### 使用示例

一个常见的模式是使用 JSON Schema 创建 `ToolInfo` 来定义结构化输入。例如，待办事项管理工具使用指定必填字段（如 `title`）和可选字段（如 `completed`）的 JSON Schema 对象来定义。

```go
55:59:/home/ruijzhan/data/eino-examples/components/tool/jsonschema/main.go
toolInfo := schema.ToolInfo{
	Name:        "todo_manager",
	Desc:        "管理待办事项列表",
	ParamsOneOf: schema.NewParamsOneOfByJSONSchema(js), // js 定义对象模式
}
```

在代理系统中，`ToolInfo` 实例从 `InvokableTool` 实现中收集，并通过 `BindTools` 或 `WithTools` 传递给聊天模型，使模型能够在推理期间调用它们。

```go
314:323:/home/ruijzhan/data/eino-examples/flow/agent/manus/manus.go
infos := make([]*schema.ToolInfo, 0, len(tools))
for _, t := range tools {
	info, err := t.Info(ctx) // 从每个工具检索 ToolInfo
	if err != nil { log.Fatal(err) }
	infos = append(infos, info)
}
ncm, err := cm.WithTools(infos) // 将工具绑定到模型
```

`ToolInfo` 在整个代码库中广泛使用，特别是在：
- 代理框架（`plan_execute`、`react`、`manus`）
- 工具绑定逻辑（`BindTools`、`WithTools`）
- 模拟和调试模型

它是在 LLM 驱动的代理中启用**函数调用**的核心，至少有 30 个文件直接引用它。

---

### 注意事项

- `ToolInfo` 通过 `ParamsOneOf` 支持**多种模式格式**，但 `NewParamsOneOfByOpenAPIV3` 已被弃用，取而代之的是 JSON Schema，如[讨论 #397](https://github.com/cloudwego/eino/discussions/397) 中所述。
- `Desc` 字段不仅是文档——它是**提示材料**，提供给模型以指导工具选择和使用。
- 由于 `ParamsOneOf` 被嵌入，`ToolInfo` 可以直接在模式转换管道中使用（例如 `ToJSONSchema()`），使其与模型后端可互操作。

---

### 另请参见

- `ParamsOneOf`: 持有不同参数模式表示的联合类型；对 `ToolInfo` 的灵活性至关重要。
- `ParameterInfo`: 在 `NewParamsOneOfByParams` 中使用的高级参数描述结构体。
- `NewParamsOneOfByJSONSchema`: 从 JSON Schema 构建 `ParamsOneOf` 的推荐构造函数。
- `ToolChoice`: 控制是否可以调用工具的枚举（`allowed`、`forbidden`、`forced`），通常与 `ToolInfo` 一起使用。
- `InvokableTool.Info()`: 返回 `*ToolInfo` 的接口方法，用于将工具集成到代理工作流中。