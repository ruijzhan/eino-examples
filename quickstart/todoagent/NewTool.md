`NewTool` 是一个函数，用于根据函数及其元数据创建可调用工具。
它将一个简单函数包装成标准化的 `tool.InvokableTool` 接口，实现与 AI 模型的工具调用工作流集成。

---

### 定义

```go
148:152:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/components/tool/utils/invokable_func.go
// NewTool 创建一个工具，输入和输出均为 JSON 格式
func NewTool[T, D any](desc *schema.ToolInfo, i InvokeFunc[T, D], opts ...Option) tool.InvokableTool {
	return newOptionableTool(desc, func(ctx context.Context, input T, _ ...tool.Option) (D, error) {
		return i(ctx, input)
	}, opts...)
}
```

- **参数**:
  - `desc`: 工具的元数据（名称、描述、参数），类型为 `*schema.ToolInfo`
  - `i`: 实际的实现函数，类型为 `InvokeFunc[T, D]`，接受输入 `T` 并返回输出 `D`
  - `opts`: 可选配置，通过 `Option` 函数提供（如自定义编组/解组）
- **副作用**: 无 — 这是一个纯构造函数
- **返回值**: `tool.InvokableTool` 类型的值，可以绑定到聊天模型以启用结构化函数调用

该函数通过抽象实现中支持 `tool.Option` 参数的需求来简化工具创建。它通过包装将基本的 `InvokeFunc` 转换为 `OptionableInvokeFunc`，然后委托给 `newOptionableTool`，后者构建具体的 `invokableTool[T, D]` 结构体。

---

### 使用示例

一个常见的用例是定义一个根据姓名和邮箱查询用户信息的工具：

```go
83:106:/home/ruijzhan/data/eino-examples/compose/graph/tool_call_agent/tool_call_agent.go
userInfoTool := utils.NewTool(
	&schema.ToolInfo{
		Name: "user_info",
		Desc: "根据用户的姓名和邮箱，查询用户的公司、职位、薪酬信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"name":  {Type: "string", Desc: "用户的姓名"},
			"email": {Type: "string", Desc: "用户的邮箱"},
		}),
	},
	func(ctx context.Context, input *userInfoRequest) (output *userInfoResponse, err error) {
		return &userInfoResponse{
			Name:     input.Name,
			Email:    input.Email,
			Company:  "Awesome company",
			Position: "CEO",
			Salary:   "9999",
		}, nil
	})
```

此示例展示了 `NewTool` 的使用方式：
- 定义一个名为 `user_info` 的工具，具有两个字符串参数
- 提供实现逻辑的 Go 函数
- 自动处理参数和结果的 JSON 序列化/反序列化

在代码库中，`NewTool` 用于多个示例，包括：
- `todoagent/main.go`: 创建具有显式模式的 `add_todo` 工具
- `duckduckgo/v2/search.go`: 将搜索客户端方法包装为工具
- `tool_call_once.go`: 用于简单的一次性工具调用演示

它是构建工具的基础工具，不依赖基于反射的推理，在需要对 `ToolInfo` 进行完全控制时非常理想。

---

### 注意事项

- `NewTool` **不会**从结构体推断参数模式 — 与 `InferTool` 不同，它需要手动指定 `ToolInfo`
- 它在内部将不支持选项的函数转换为支持选项的函数，通过丢弃运行时传递的 `opts`，确保与 `InvokableTool` 接口的兼容性
- 尽管名称如此，`NewTool` 不是新类型的构造函数 — 它是一个包装器，用于标准化函数接口以便在 AI 编排系统中使用

---

### 相关参考

- `InferTool`: 使用反射从函数的输入结构体自动推断 `ToolInfo`；在推断后内部调用 `NewTool`
- `newOptionableTool`: 构建 `invokableTool` 结构体的实际构造函数；被 `NewTool` 和 `InferOptionableTool` 使用
- `tool.InvokableTool`: 定义 `Info()` 和 `InvokableRun()` 的接口，用于将工具绑定到聊天模型
- `schema.ToolInfo`: 包含元数据（名称、描述、参数）的结构体，用于向 LLM 描述工具