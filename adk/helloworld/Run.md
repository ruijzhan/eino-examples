# Run 方法

`(*Runner).Run` 是一个方法，使用输入消息执行代理并返回事件的异步迭代器。
它协调代理的执行流程，处理流式传输，并集成检查点功能以支持可恢复运行。

---

### 定义

`(*Runner).Run` 方法是 `adk` 包中 `Runner` 类型的核心执行函数。它接受上下文、消息列表和可选的运行时配置，然后触发代理的执行。根据是否配置了检查点存储，它要么直接返回代理的输出迭代器，要么包装它以支持容错、可恢复的执行。

```go
50:74:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/adk/runner.go
func (r *Runner) Run(ctx context.Context, messages []Message,
    opts ...AgentRunOption) *AsyncIterator[*AgentEvent] {
    // 提取常用选项，如会话值和检查点 ID
    o := getCommonOptions(nil, opts...)

    // 将代理转换为支持流的代理
    fa := toFlowAgent(ctx, r.a)

    // 用消息和流式设置准备输入
    input := &AgentInput{
        Messages:        messages,
        EnableStreaming: r.enableStreaming,
    }

    // 在上下文中创建新的运行上下文
    ctx = ctxWithNewRunCtx(ctx)

    // 将会话值注入上下文
    AddSessionValues(ctx, o.sessionValues)

    // 运行代理并获取事件流
    iter := fa.Run(ctx, input, opts...)

    // 如果没有检查点存储，返回原始迭代器
    if r.store == nil {
        return iter
    }

    // 否则，创建新的迭代器对
    niter, gen := NewAsyncIteratorPair[*AgentEvent]()

    // 启动 goroutine 处理事件和检查点
    go r.handleIter(ctx, iter, gen, o.checkPointID)

    // 返回包装的迭代器
    return niter
}
```

- **参数**:
  - `ctx context.Context`: 执行上下文，可能携带截止时间和取消信号。
  - `messages []Message`: 代理的初始输入消息（例如，用户查询）。
  - `opts ...AgentRunOption`: 可选配置，如检查点 ID、会话数据或跟踪选项。

- **副作用**:
  - 通过注入新的 `runContext` 修改上下文。
  - 如果启用检查点，则生成 goroutine（`r.handleIter`）。
  - 如果执行被中断且配置了存储，可能会保存检查点。

- **返回值**:
  - `*AsyncIterator[*AgentEvent]`: 可以异步消费的代理事件流（例如，消息、工具调用、错误）。

---

### 使用示例

一个常见的用法是在 `helloworld` 示例中，其中运行器用于发送问候消息并处理代理的响应流。

```go
65:73:/home/ruijzhan/data/eino-examples/adk/helloworld/helloworld.go
	input := []adk.Message{
		schema.UserMessage("Hello, please introduce yourself."),
	}

	events := runner.Run(ctx, input)  // 开始执行
	for {
		event, ok := events.Next()  // 消费事件流
		if !ok {
			break
		}
		// 处理事件...
	}
```

另一个内部用途是在 `agent_tool.go` 中，其中创建了一个运行器来调用代理工具，启用模拟输入和检查点：

```go
143:151:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/adk/agent_tool.go
	iter = newInvokableAgentToolRunner(at.agent, ms).Run(
		ctx,
		input,
		append(getOptionsByAgentName(at.agent.Name(ctx), opts),
			WithCheckPointID(mockCheckPointID))...,
	)
```

**使用摘要**:
`(*Runner).Run` 用于：
- 最终用户代理执行（例如 `helloworld`）。
- 代理调用的内部工具（`agent_tool.go`）。
- 像 `Query` 这样的高级方法，简化了使用单个字符串调用 `Run`。

它是一个中心执行点，在示例和内部逻辑中直接或间接调用。尽管通过静态分析没有找到直接调用者，但它通过方法引用被调用，是代理执行模型的基础。

---

### 注意事项

- **检查点是可选的但有影响**：如果 `r.store` 为 `nil`，`Run` 返回原始迭代器。否则，它用 `handleIter` 包装它，启用中断后的可恢复执行。
- **panic 安全的事件处理**：`handleIter` goroutine 从 panic 中恢复并将它们作为 `AgentEvent{Err: ...}` 发送，确保流不会静默失败。
- **上下文变异**：该方法将新的 `runContext` 注入 `ctx`，稍后 `getInterruptRunCtx` 使用它在中断时保存检查点。

---

### Query 方法对比

`(*Runner).Query` 是 `(*Runner).Run` 的一个便利包装方法，两者有以下关键区别：

#### 定义对比

**Run 方法**:
```go
func (r *Runner) Run(ctx context.Context, messages []Message,
    opts ...AgentRunOption) *AsyncIterator[*AgentEvent]
```

**Query 方法**:
```go
func (r *Runner) Query(ctx context.Context,
    query string, opts ...AgentRunOption) *AsyncIterator[*AgentEvent] {

    return r.Run(ctx, []Message{schema.UserMessage(query)}, opts...)
}
```

#### 主要区别

| 特性 | Run 方法 | Query 方法 |
|------|----------|------------|
| **输入参数** | `[]Message` 消息数组 | `string` 查询字符串 |
| **灵活性** | 高，支持复杂消息历史 | 低，仅支持单条用户消息 |
| **使用复杂度** | 需要手动构造消息结构 | 简单，直接传入字符串 |
| **适用场景** | 多轮对话、复杂上下文 | 简单查询、一次性交互 |

#### 使用示例对比

**Run 方法使用**:
```go
// 需要手动构造消息
input := []adk.Message{
    schema.UserMessage("Hello, please introduce yourself."),
    schema.AssistantMessage("I'm a helpful assistant."),
    schema.UserMessage("What can you help me with?"),
}
events := runner.Run(ctx, input)
```

**Query 方法使用**:
```go
// 直接传入查询字符串
events := runner.Query(ctx, "Hello, please introduce yourself.")
```

#### 内部实现

`Query` 方法内部实际是对 `Run` 方法的简单包装：
1. 将传入的字符串参数转换为 `schema.UserMessage`
2. 将其包装在单元素的消息数组中
3. 调用 `Run` 方法并返回结果

这意味着 `Query` 方法在功能上是 `Run` 方法的一个子集，所有 `Query` 能做的事情都可以通过 `Run` 实现，但 `Query` 提供了更简洁的 API 用于常见的查询场景。

#### 使用场景建议

**选择 Query 方法的场景**:
- **简单交互**: 一轮对话，用户直接提问
- **快速测试**: 原型开发或功能验证
- **示例代码**: 教学或文档中的简单演示
- **API 简化**: 不需要复杂消息历史的场景

```go
// 适合使用 Query 的场景
events := runner.Query(ctx, "今天天气怎么样？")
events := runner.Query(ctx, "帮我推荐一本书")
events := runner.Query(ctx, "解释一下机器学习的概念")
```

**选择 Run 方法的场景**:
- **多轮对话**: 需要维护对话历史和上下文
- **复杂交互**: 涉及多种角色消息（用户、助手、系统）
- **消息预处理**: 需要对消息进行特殊处理或格式化
- **上下文管理**: 需要精确控制消息顺序和内容

```go
// 适合使用 Run 的场景
input := []adk.Message{
    schema.SystemMessage("你是一个专业的技术顾问"),
    schema.UserMessage("我想了解微服务架构"),
    schema.AssistantMessage("微服务架构是一种将应用程序构建为一组小型服务的方法"),
    schema.UserMessage("微服务架构有什么优缺点？"),
}
events := runner.Run(ctx, input)
```

#### 设计模式分析

这种设计模式在软件工程中很常见，类似于：

1. **Builder 模式**: `Run` 提供完整的构建能力，`Query` 提供简化的构建路径
2. **Facade 模式**: `Query` 作为简化门面，隐藏了消息构造的复杂性
3. **Template Method 模式**: 两者共享相同的执行模板，只是在输入处理上有所不同

**优势**:
- **渐进式复杂度**: 初学者可以使用 `Query`，高级用户可以使用 `Run`
- **向后兼容**: 可以从 `Query` 开始，需要时无缝迁移到 `Run`
- **代码复用**: 两个方法共享相同的底层执行逻辑
- **API 一致性**: 相同的选项和返回类型，学习成本低

---

### 另请参见

- `(*Runner).Query`: `Run` 方法的便利包装，用于简单的基于字符串的查询交互。
- `(*Runner).Resume`: 通过从保存的检查点重新启动执行来补充 `Run`，依赖相同的 `handleIter` 逻辑。
- `handleIter`: 管理事件转发、检查点保存和 panic 恢复的 goroutine 辅助程序。对于可恢复执行至关重要。
- `AgentRunOption`: 传递给 `Run` 和 `Query` 的函数选项，配置检查点 ID 和会话值等行为。