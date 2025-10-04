# NewRunner 构造函数

`NewRunner` 是一个构造函数，用于创建执行 AI 代理的新 `*Runner` 实例。

它使用代理、流式设置和可选的检查点存储初始化运行器，启用受控执行和可恢复的工作流。

---

### 定义

```go
42:48:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/adk/runner.go
func NewRunner(_ context.Context, conf RunnerConfig) *Runner {
	return &Runner{
		enableStreaming: conf.EnableStreaming,
		a:               conf.Agent,
		store:           conf.CheckPointStore,
	}
}
```

- **参数**:
  - `conf RunnerConfig`: 包含代理、流式标志和检查点存储的配置结构体。
  - `_ context.Context`: 未使用的参数（通过下划线忽略），表明初始化没有与上下文相关的副作用。
- **副作用**: 无。该函数是纯函数——只构造并返回 `*Runner`。
- **返回值**: `*Runner` — 一个包装 `Agent` 并支持流式输出和基于检查点恢复的有状态执行器。

`Runner` 结构体本身包含：
- `a Agent`: 要执行的底层 AI 代理。
- `enableStreaming bool`: 控制事件是否增量流式传输。
- `store compose.CheckPointStore`: 用于保存/恢复执行状态的可选持久化层。

此构造函数不启动执行——它为稍后调用 `Run` 或 `Resume` 准备运行器。

---

### 使用示例

常见模式是创建一个 `Runner` 来执行具有实时流式传输和可选检查点功能的代理。

#### 基本用法：流式执行

```go
59:62:/home/ruijzhan/data/eino-examples/adk/helloworld/helloworld.go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
	Agent:           agent,
	EnableStreaming: true,
})
```

这里，`NewRunner` 设置了一个启用流式传输的简单代理运行器。代理可以通过 `Run` 或 `Query` 调用，产生 `AgentEvent` 流。

#### 高级用法：使用检查点

```go
37:41:/home/ruijzhan/data/eino-examples/adk/intro/chatmodel/chatmodel.go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
	EnableStreaming: true,
	Agent:           a,
	CheckPointStore: newInMemoryStore(),
})
```

这启用了可恢复执行。如果代理被中断，其状态可以被保存，稍后使用 `Resume(checkPointID)` 恢复。

#### 内联构造和立即查询

```go
46:49:/home/ruijzhan/data/eino-examples/adk/multiagent/layered-supervisor/layered_supervisor.go
iter := adk.NewRunner(ctx, adk.RunnerConfig{
	EnableStreaming: true,
	Agent:           sv,
}).Query(ctx, query)
```

常见习语：内联构造 `Runner` 并立即调用 `Query`。对于不需要存储运行器的一次性代理调用很有用。

**总体使用摘要**:
`NewRunner` 在整个代码库中至少 11 个不同的文件中被调用，从像 `helloworld.go` 这样的简单演示到复杂的多代理工作流（`layered-supervisor`、`integration-project-manager`）。它是 ADK 框架中所有代理执行模式的核心。

---

### 注意事项

- `context.Context` 参数被忽略（下划线），表明 `NewRunner` 是无副作用的，不会启动任何异步操作或资源获取。
- 如果 `CheckPointStore` 为 `nil`，`Resume` 将因错误而失败：`"failed to resume: store is nil"`。这使得检查点成为可选的，但在使用时严格执行。
- `NewRunner` 不验证 `Agent` 或 `CheckPointStore`——验证可能在 `Run` 或 `Resume` 期间进行。

---

### 另请参见

- `Runner`: `NewRunner` 返回的类型。它提供像 `Run` 和 `Resume` 这样的方法来执行代理。
- `RunnerConfig`: 传递给 `NewRunner` 的输入结构体，定义代理、流式传输和检查点行为。
- `compose.CheckPointStore`: 用于保存和检索执行状态的接口。可恢复工作流所需。
- `Agent`: 表示 AI 代理的核心接口。直接传递到 `RunnerConfig` 中并在执行期间使用。
- `AgentRunOption`: 传递给 `Run` 的函数选项，配置检查点 ID 和会话值等行为。