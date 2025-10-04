# RunnerConfig 配置结构体

`RunnerConfig` 是一个配置结构体，用于在 Eino ADK 框架中初始化 `Runner`。
它定义了运行代理的核心组件和行为设置，例如流式模式和检查点持久化。

这个结构体对于设置代理如何执行和恢复至关重要——特别是在长时间运行或可中断的工作流中。

---

### 定义

```go
35:40:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/adk/runner.go
type RunnerConfig struct {
	Agent           Agent
	EnableStreaming bool
	CheckPointStore compose.CheckPointStore
}
```

- **用途**: 保存 `Runner` 的初始化参数，`Runner` 协调 `Agent` 的执行。
- **字段**:
  - `Agent`: 要运行的代理（例如聊天机器人或工作流主管）。
  - `EnableStreaming`: 如果为 `true`，则在执行期间启用事件的实时流式传输。
  - `CheckPointStore`: 可选存储，用于在中断后保存和恢复代理状态。

此配置被传递给 `NewRunner`，后者构造一个用于调用 `Run` 或 `Resume` 操作的 `*Runner` 实例。

```go
42:47:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/adk/runner.go
func NewRunner(_ context.Context, conf RunnerConfig) *Runner {
	return &Runner{
		enableStreaming: conf.EnableStreaming,
		a:               conf.Agent,
		store:           conf.CheckPointStore,
	}
}
```

- **返回值**: `*Runner` — 一个包装代理的有状态执行器。
- **副作用**: 没有直接副作用，但返回的 `Runner` 可能在执行或恢复期间与 `CheckPointStore` 交互。

---

### 使用示例

在实践中，`RunnerConfig` 在各种示例中用于引导具有不同功能的代理执行。

#### 基本用法：启用流式传输
```go
59:62:/home/ruijzhan/data/eino-examples/adk/helloworld/helloworld.go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
	Agent:           agent,
	EnableStreaming: true,
})
```
这里，`RunnerConfig` 设置了一个启用流式传输的简单代理。没有使用检查点——执行运行到完成。

#### 高级用法：使用检查点
```go
37:41:/home/ruijzhan/data/eino-examples/adk/intro/chatmodel/chatmodel.go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
	EnableStreaming: true,
	Agent:           a,
	CheckPointStore: newInMemoryStore(),
})
```
此配置启用了流式传输和检查点功能。代理可以被中断，其状态保存在内存中，以便稍后通过 `Resume` 恢复。

#### 内联构造和查询
```go
46:49:/home/ruijzhan/data/eino-examples/adk/multiagent/layered-supervisor/layered_supervisor.go
iter := adk.NewRunner(ctx, adk.RunnerConfig{
	EnableStreaming: true,
	Agent:           sv,
}).Query(ctx, query)
```
常见模式：内联构造 `Runner` 并立即调用 `Query`。这对于一次性执行很有用。

**总体使用摘要**:
`RunnerConfig` 在整个代码库中被广泛使用——出现在至少 12 个不同的调用者文件中。它是所有代理执行场景的核心，从简单的聊天代理（`helloworld.go`）到复杂的多代理工作流（`layered-supervisor`、`plan-execute-replan`）。一些示例中 `CheckPointStore` 的存在（例如 `integration-project-manager`）表明了它在有状态的、可恢复的 AI 工作流中的作用。

---

### 注意事项

- 尽管在 `NewRunner` 中接受 `context.Context`，但上下文被忽略（下划线参数），表明配置构造是无副作用的，执行上下文在稍后的 `Run` 期间处理。
- `CheckPointStore` 字段对于可恢复执行至关重要：如果为 `nil`，`Resume` 将因错误而失败。
- `RunnerConfig` 不支持超时或重试策略——这些可能通过传递给 `Run` 或 `Query` 的 `AgentRunOption` 来管理。

---

### 另请参见

- `Runner`: 从 `RunnerConfig` 构造的类型。它执行代理并处理流式传输/检查点逻辑。
- `NewRunner`: 使用 `RunnerConfig` 并返回 `*Runner` 的构造函数。
- `compose.CheckPointStore`: 用于保存和检索执行状态的接口。用于可恢复工作流。
- `Agent`: 表示 AI 代理的核心接口。直接传递到 `RunnerConfig` 中。
- `AgentRunOption`: 在 `Run` 或 `Query` 期间使用的选项，可以覆盖或扩展 `RunnerConfig` 中定义的行为。