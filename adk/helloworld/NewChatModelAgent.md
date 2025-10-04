# NewChatModelAgent 构造函数

`NewChatModelAgent` 是一个构造函数，用于创建基于聊天的 AI 代理的新实例。
它使用配置结构体初始化 `ChatModelAgent`，验证必填字段并设置默认值。

---

### 定义

`NewChatModelAgent` 是一个工厂函数，从 `*ChatModelAgentConfig` 构造 `*ChatModelAgent`。它对关键字段（`Name`、`Description`、`Model`）进行验证，并设置默认行为，比如除非被覆盖，否则使用 `defaultGenModelInput` 进行消息格式化。

```go
177:204:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/adk/chatmodel.go
func NewChatModelAgent(_ context.Context, config *ChatModelAgentConfig) (*ChatModelAgent, error) {
	if config.Name == "" {
		return nil, errors.New("agent 'Name' is required")
	}
	if config.Description == "" {
		return nil, errors.New("agent 'Description' is required")
	}
	if config.Model == nil {
		return nil, errors.New("agent 'Model' is required")
	}

	// 如果提供了自定义 GenModelInput 则使用，否则使用默认值
	genInput := defaultGenModelInput
	if config.GenModelInput != nil {
		genInput = config.GenModelInput
	}

	// 构造并返回代理
	return &ChatModelAgent{
		name:          config.Name,
		description:   config.Description,
		instruction:   config.Instruction,
		model:         config.Model,
		toolsConfig:   config.ToolsConfig,
		genModelInput: genInput,
		exit:          config.Exit,
		outputKey:     config.OutputKey,
		maxIterations: config.MaxIterations,
	}, nil
}
```

- **参数**:
  - `context.Context`: 未使用（下划线），表明该函数当前不使用 context 进行取消或获取值。
  - `*ChatModelAgentConfig`: 定义代理属性（名称、模型、工具等）的配置结构体。
- **副作用**: 无 — 纯粹创建并返回代理实例。
- **返回值**:
  - `*ChatModelAgent`: 完全初始化的代理，准备好在多代理系统中使用。
  - `error`: 如果缺少必填字段（`Name`、`Description`、`Model`）。

`ChatModelAgent` 结构体保存核心代理状态，包括 LLM 模型、工具、输入生成逻辑和会话输出设置。它支持工具调用、消息格式化和工作流中的迭代执行。

---

### 使用示例

`NewChatModelAgent` 的典型用法可以在"Hello World"示例中看到，它使用 OpenAI 模型创建一个简单的问候代理。

```go
48:53:/home/ruijzhan/data/eino-examples/adk/helloworld/helloworld.go
agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
	Name:        "hello_agent",
	Description: "A friendly greeting assistant",
	Instruction: "You are a friendly assistant. Please respond to the user in a warm tone.",
	Model:       model,
})
if err != nil {
	log.Fatal(err)
}
```

然后该代理与 `Runner` 一起使用来处理用户输入并生成响应。

**使用摘要**:
`NewChatModelAgent` 在整个代码库中被广泛用于实例化各种类型的代理：
- **通用聊天代理**（例如 `ChatAgent`）
- **专业化代理**（例如 `WeatherAgent`、`CodeAgent`、`ResearchAgent`）
- **主管/路由代理**（例如 `RouterAgent`、`ProjectManagerAgent`）
- **数学和数据处理代理**（例如 `multiply_agent`、`StockDataCollectionAgent`）

它出现在至少 **17 个文件**中，主要在 `/adk/` 下的示例和集成项目中，表明它是 EINO 框架中构建代理工作流的核心入口点。

---

### 注意事项

- 尽管接受 `context.Context`，但该函数忽略了它——这可能表明未来支持上下文感知初始化的可扩展性（例如跟踪、超时）。
- 该函数**不**验证 `ToolsConfig` 或 `Exit` 工具，意味着无效工具可能只在运行时失败。
- `GenModelInput` 是可插拔的，允许开发者在发送到模型之前自定义指令和消息的格式化方式。

---

### 另请参见

- `ChatModelAgentConfig`: 传递给 `NewChatModelAgent` 的配置结构体，定义所有代理属性。
- `defaultGenModelInput`: 当没有提供自定义 `GenModelInput` 时使用的默认消息格式化器。
- `ChatModelAgent`: 生成的代理类型，实现了 `Agent` 接口并支持工具调用和消息生成。
- `ExitTool`: 可以分配给配置中的 `Exit` 字段的内置工具，允许代理终止。
- `Runner`: 通常与 `NewChatModelAgent` 创建的代理一起使用来执行和流式传输代理响应。