`Chain` 是一个用于构建 AI 工作流的可组合组件管道。
它支持组件的顺序、并行或分支执行——例如模型、提示和工具——编译成可运行的图。

---

### 定义

`Chain` 类型是一个泛型结构体，表示 AI 组件（例如聊天模型、提示模板、检索器）的有向序列，可以以构建器(Builder)模式风格组合。它包装了底层的 `Graph[I, O]` 并提供流畅的方法来附加节点，然后编译成 `Runnable[I, O]` 以供执行。

```go
72:82:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/compose/chain.go
type Chain[I, O any] struct {
	err error
	gg *Graph[I, O]               // 底层图
	nodeIdx int                   // 当前节点索引
	preNodeKeys []string          // 前一个节点的键
	hasEnd bool                   // 确保 END 节点只添加一次
}
```

- **类型参数**:
  - `I`: 链的输入类型
  - `O`: 链的输出类型
- **核心字段**:
  - `gg *Graph[I, O]` — 管理节点和边的内部图
- **构建器模式**: 设计为方法链式调用（例如 `chain.AppendX().AppendY().Compile()`）
- **编译要求**: 执行前必须调用 `.Compile()`；编译后修改会因 `ErrChainCompiled` 而失败
- **实现**: 通过 `inputType()`、`outputType()` 和 `compile()` 等方法实现 `AnyGraph` 接口

`Chain` 通过 `NewChain[I, O]()` 初始化，通过附加组件（如 `AppendChatModel`、`AppendLambda` 等）构建，然后编译成可运行形式。

---

### 使用示例

代码库中的一个真实示例展示了使用 `InvokableLambda` 函数创建简单的字符串处理链：

```go
29:49:/home/ruijzhan/data/eino-examples/devops/debug/chain/chain.go
chain := compose.NewChain[string, string]()

c1 := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
	return input + " process by node_1,", nil
})

c2 := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
	return input + " process by node_2,", nil
})

chain.AppendLambda(c1, compose.WithNodeName("c1")).
	AppendLambda(c2, compose.WithNodeName("c2"))

r, err := chain.Compile(ctx)
if err != nil { /* handle error */ }

message, err := r.Invoke(ctx, "eino chain test")
// 输出: "eino chain test process by node_1, process by node_2,"
```

这展示了：
- 输入类型: `string`
- 两个顺序附加的处理步骤（lambda）
- 通过 `Invoke` 获得最终结果

另一个在日志编写代理中的使用示例构建了一个使用聊天模型将消息列表转换为单个日志条目的链：

```go
100:118:/home/ruijzhan/data/eino-examples/flow/agent/multiagent/host/journal/write_journal_specialist.go
chain := compose.NewChain[[]*schema.Message, *schema.Message]()
chain.AppendLambda(/* prepends system message */).
	AppendChatModel(chatModel).
	AppendLambda(/* writes content to journal file */)
```

**使用总结**:
- `Chain` 用于多个领域：调试演示、代理系统和工作流组合。
- 主要调用者包括 `RegisterSimpleChain`、`newWriteJournalSpecialist` 和测试示例。
- 它是 `eino` 框架中构建 AI 管道的核心抽象，通常在 `Graph` 之前使用或通过 `AppendGraph` 在其他链内部使用。

---

### 注意事项

- **非线程安全**: 构建器模式意味着在构建期间会发生变异；不应并发修改。
- **单次编译**: 一旦调用 `Compile()`，不允许进一步修改——由 `ErrChainCompiled` 强制执行。
- **结束节点自动注入**: `addEndIfNeeded()` 方法通过将最后一个节点连接到 `END` 来确保链正确终止，防止悬空执行路径。

---

### 另请参见

- `Graph[I, O]`: `Chain` 包装的底层数据结构；管理节点和边。`Chain` 在其上添加了语法糖和构建器语义。
- `NewChain[I, O]`: 创建新 `Chain` 实例的工厂函数；初始化内部 `Graph`。
- `Runnable[I, O]`: `Chain.Compile()` 的编译输出；提供 `Invoke`、`Stream`、`Collect` 和 `Transform` 方法以供执行。
- `AnyGraph`: `Chain` 实现的接口，使其能够嵌入到其他图或链中。
- `AppendChatModel`、`AppendLambda` 等：使用特定组件类型扩展链的流畅方法。