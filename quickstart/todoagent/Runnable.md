`Runnable` 是 Eino 框架中表示可执行组件的泛型接口。
它标准化了诸如 `Invoke`、`Stream`、`Collect` 和 `Transform` 等数据流操作在诸如 `Graph` 和 `Chain` 等可组合组件中的执行方式。

这个统一的执行模型允许不同的数据处理模式无缝互操作，支持方法之间的自动适配（例如，在仅实现 `Stream` 的组件上调用 `Invoke`）。它是 Eino 能够从更简单的构建块组合复杂 AI 工作流的核心能力。

---

### 定义

`Runnable[I, O]` 接口定义了四种支持不同数据流模式的核心执行方法。任何实现这些方法中一个或多个的组件都可以被组合和自动适配。

```go
28:37:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/compose/runnable.go
// Runnable 是可执行对象的接口。Graph、Chain 可以编译成 Runnable。
// runnable 是 eino 的核心概念，我们为四种数据流模式做降级兼容，
// 并且可以自动连接只实现一个或多个方法的组件。
// 例如，如果组件只实现 Stream() 方法，你仍然可以调用 Invoke() 来将流输出转换为调用输出。
type Runnable[I, O any] interface {
	Invoke(ctx context.Context, input I, opts ...Option) (output O, err error)
	Stream(ctx context.Context, input I, opts ...Option) (output *schema.StreamReader[O], err error)
	Collect(ctx context.Context, input *schema.StreamReader[I], opts ...Option) (output O, err error)
	Transform(ctx context.Context, input *schema.StreamReader[I], opts ...Option) (output *schema.StreamReader[O], err error)
}
```

- **类型参数**:
  - `I`: 输入类型
  - `O`: 输出类型
- **方法**:
  - `Invoke`: 同步执行 — 输入 → 输出
  - `Stream`: 流式输出 — 输入 → 流
  - `Collect`: 流式输入 → 输出（聚合）
  - `Transform`: 流式输入 → 流输出（流处理）
- **目的**: 通过自动方法合成实现具有部分方法实现的组件之间的互操作性。

例如，如果组件只实现 `Stream`，仍然可以调用 `Invoke`，它将消费流以产生单个结果。

---

### 使用示例

在实践中，`Runnable` 用作更高级构造（如 `Graph` 或 `Chain`）的编译输出。一旦编译，它为执行提供统一的 API。

以下是来自 `composeGraph` 的真实用法，它构建有状态的 AI 代理工作流并返回 `Runnable`：

```go
171:230:/home/ruijzhan/data/eino-examples/compose/graph/react_with_interrupt/main.go
func composeGraph[I, O any](
	ctx context.Context,
	tpl prompt.ChatTemplate,
	cm model.ChatModel,
	tn *compose.ToolsNode,
	store compose.CheckPointStore,
) (compose.Runnable[I, O], error) {
	g := compose.NewGraph[I, O](/* initializes local state */)
	_ = g.AddChatTemplateNode("ChatTemplate", tpl)
	_ = g.AddChatModelNode("ChatModel", cm, /* with pre/post state handlers */)
	_ = g.AddToolsNode("ToolsNode", tn, /* state handler */)
	_ = g.AddEdge(compose.START, "ChatTemplate")
	_ = g.AddEdge("ChatTemplate", "ChatModel")
	_ = g.AddBranch("ChatModel", compose.NewGraphBranch(/* routing logic */))

	// 将图编译成 Runnable
	return g.Compile(ctx, compose.WithCheckPointStore(store), compose.WithInterruptBeforeNodes([]string{"ToolsNode"}))
}
```

后来，在 `main` 中，这个 `Runnable` 被直接调用：

```go
62:82:/home/ruijzhan/data/eino-examples/compose/graph/react_with_interrupt/main.go
runner, err := composeGraph[map[string]any, *schema.Message](ctx, /* dependencies */)
if err != nil {
	log.Fatal(err)
}

result, err := runner.Invoke(ctx, map[string]any{"name": "Megumin", "location": "Beijing"}, /* options */)
if err == nil {
	fmt.Printf("final result: %s", result.Content)
}
```

**使用总结**:
- `Runnable` 由 `(*Graph[I, O]).Compile`、`(*Chain[I, O]).Compile` 和 `(*Workflow[I, O]).Compile` 返回
- 广泛用于 AI 代理示例：多代理规划、助手工作流、工具集成
- 在至少 31 个文件中找到，表明它是基础抽象

---

### 注意事项

- **自动方法适配**: 即使组件只实现 `Stream`，你也可以调用 `Invoke` —— 框架将消费流并返回最终值。这是 Eino 灵活性的关键。
- **泛型 + 反射**: `composableRunnable` 结构体使用反射来处理类型擦除，同时通过泛型在编译时保持类型安全。
- **有状态执行**: 当与 `Graph` 一起使用时，`Runnable` 可以在步骤之间保持状态（例如对话历史），特别是当与 `CheckPointStore` 和 `StatePreHandler`/`PostHandler` 结合使用时。

---

### 另请参见

- `Graph[I, O]`: 编译成 `Runnable` 的节点的有向无环图（或循环，有步骤限制）。用于复杂的控制流和状态管理。
- `Chain[I, O]`: `Runnable` 组件的线性序列，也编译成 `Runnable`。比 `Graph` 简单，但支持流和状态。
- `schema.StreamReader[T]`: 在 `Stream` 和 `Transform` 方法中使用的流类型。表示可以增量消费的 `T` 类型值流。
- `(*Graph[I, O]).Compile`: 将 `Graph` 转换为 `Runnable[I, O]` 以启用执行的方法。
- `PlanExecuteMultiAgent`: 具体的代理实现，持有 `runnable compose.Runnable[[]*schema.Message, *schema.Message]` 来执行多代理工作流。