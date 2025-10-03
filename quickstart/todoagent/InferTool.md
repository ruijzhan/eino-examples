`InferTool` 是一个实用函数，通过从 Go 函数的输入类型推断模式信息来创建 `InvokableTool`。
它使用反射和 JSON Schema 推断自动生成工具元数据（如参数结构），减少了定义 AI 代理工具时的样板代码。

---

### 定义

`InferTool` 是一个泛型函数，它接受一个具有结构化输入的函数，并自动从该输入的 Go 结构体标签推导出工具的接口模式（例如参数名称、类型、描述）。这允许开发人员完全通过惯用的 Go 代码来定义工具，而无需手动编写 JSON Schema 或 `ToolInfo`。

```go
40:47:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/components/tool/utils/invokable_func.go
// InferTool 通过从函数的请求参数推断 ToolInfo 来为给定函数创建 InvokableTool
// 最终用户可以在 opts 中传递 SchemaCustomizerFn 来自定义 go 结构体标签解析过程
func InferTool[T, D any](
    toolName, toolDesc string,
    i InvokeFunc[T, D],
    opts ...Option,
) (tool.InvokableTool, error) {
    ti, err := goStruct2ToolInfo[T](toolName, toolDesc, opts...)
    if err != nil {
        return nil, err
    }
    return NewTool(ti, i, opts...), nil
}
```

- **参数**:
  - `toolName`: 工具名称（例如 `"search_book"`）。
  - `toolDesc`: 人类可读的描述，供 LLM 理解工具的用途。
  - `i`: 实际的实现函数（`func(ctx context.Context, input T) (D, error)`）。
  - `opts`: 可选配置，例如自定义结构体标签的解析方式。
- **副作用**: 无；它是一个纯粹的构造函数。
- **返回**: 一个完全配置的 `tool.InvokableTool`，可以绑定到聊天模型，或者如果模式推断失败则返回错误。

在底层，`InferTool` 使用 `github.com/eino-contrib/jsonschema` 来反射输入类型 `T` 并生成 JSON Schema，然后将其转换为 `schema.ToolInfo`。这包括处理 `json` 和 `jsonschema` 结构体标签以进行字段命名、描述和约束。

---

### 使用示例

一个常见的用例是定义书籍搜索工具，其中输入结构包含用户偏好：

```go
38:48:/home/ruijzhan/data/eino-examples/adk/intro/chatmodel/subagents/booksearch.go
bookSearchTool, err := utils.InferTool(
    "search_book",
    "根据用户偏好搜索书籍",
    func(ctx context.Context, input *BookSearchInput) (output *BookSearchOutput, err error) {
        // 模拟实现返回固定的推荐书籍列表
        return &BookSearchOutput{Books: []string{"为这美好的世界献上祝福！"}}, nil
    },
)
if err != nil {
    log.Fatalf("创建搜索书籍工具失败: %v", err)
}
```

这里，`BookSearchInput` 可能是这样的结构：

```go
type BookSearchInput struct {
    Genre string `json:"genre" jsonschema:"description=要搜索的书籍类型"`
    Mood  string `json:"mood" jsonschema:"description=用户偏好的情绪或主题"`
}
```

`InferTool` 调用会自动生成一个模式，告诉 LLM 在调用工具时应该提供 `genre` 和 `mood`。

**使用总结**:
`InferTool` 在整个代码库中被广泛使用——在至少 12 个不同的调用者中出现——用于定义天气查询、航班搜索、酒店预订等工具。它是 Eino 框架"函数优先"工具定义风格的核心，特别是在 `adk/`、`quickstart/` 和 `compose/` 下的示例中。它的普遍性凸显了对最小样板代码和通过 Go 泛型和反射实现最大类型安全的设计偏好。

---

### 注意事项

- `InferTool` 仅适用于具有正确 `json` 标签的**结构体输入**；使用基本类型（如 `string`）将导致模式生成失败。
- 它依赖 `sonic` 进行 JSON 编组，依赖 `jsonschema.Reflect` 进行模式推断，因此结构体字段的可见性和标记至关重要。
- 该函数假设输入在运行时作为 JSON 字符串传递，然后使用高性能 JSON 解析将其反序列化为 Go 结构体。

---

### 另请参见

- `InferOptionableTool`: 类似于 `InferTool`，但在函数签名中接受额外的 `tool.Option` 参数，允许运行时自定义工具行为。
- `NewTool`: 底层构造函数，需要手动指定 `ToolInfo`，当不希望进行模式推断时使用。
- `goStruct2ToolInfo`: 执行反射和模式生成的内部函数；由 `InferTool` 调用以从输入结构体提取元数据。
- `SchemaCustomizerFn`: 可以通过 `opts` 传递的函数选项，用于自定义模式生成期间结构体标签的解析方式。