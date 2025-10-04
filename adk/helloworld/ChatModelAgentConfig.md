# ChatModelAgentConfig 配置结构体

`ChatModelAgentConfig` 是一个配置结构体，用于定义基于聊天的 AI 代理的行为和属性。

它作为初始化 `ChatModelAgent` 的基础蓝图，指定了代理的核心属性，如名称、角色、工具和执行限制。此配置会传递给 `NewChatModelAgent`，后者验证并构建实际的代理实例。

---

### 定义

```go
113:147:/home/ruijzhan/data/eino-examples/vendor/github.com/cloudwego/eino/adk/chatmodel.go
type ChatModelAgentConfig struct {
	// Name of the agent. Better be unique across all agents.
	// 代理名称。最好在所有代理中保持唯一。
	Name string

	// Description of the agent's capabilities.
	// Helps other agents determine whether to transfer tasks to this agent.
	// 代理能力的描述。帮助其他代理确定是否向此代理转移任务。
	Description string

	// Instruction used as the system prompt for this agent.
	// Optional. If empty, no system prompt will be used.
	// Supports f-string placeholders for session values.
	// 用于此代理的系统提示指令。可选。如果为空，则不使用系统提示。
	// 支持会话值的 f-string 占位符。
	Instruction string

	// Model that supports tool calling.
	// 支持工具调用的模型。
	Model model.ToolCallingChatModel

	// Configuration for tools available to the agent.
	// 代理可用工具的配置。
	ToolsConfig ToolsConfig

	// Function to format input messages for the model.
	// Optional. Defaults to defaultGenModelInput.
	// 为模型格式化输入消息的函数。可选。默认为 defaultGenModelInput。
	GenModelInput GenModelInput

	// Tool used to terminate the agent process.
	// Optional. If nil, no exit action is generated.
	// 用于终止代理进程的工具。可选。如果为 nil，则不生成退出操作。
	Exit tool.BaseTool

	// Key to store agent output in the session.
	// Optional. If set, output is saved via AddSessionValue.
	// 在会话中存储代理输出的键。可选。如果设置，输出将通过 AddSessionValue 保存。
	OutputKey string

	// Max number of generation cycles.
	// Optional. Defaults to 20.
	// 最大生成周期数。可选。默认为 20。
	MaxIterations int
}
```

- **类型**: 结构体
- **用途**: `NewChatModelAgent` 的配置输入
- **必填字段**: `Name`, `Description`, `Model`
- **可选字段**: `Instruction`, `GenModelInput`, `Exit`, `OutputKey`, `MaxIterations`
- **默认值**:
  - `MaxIterations` 默认为 20
  - `GenModelInput` 默认为 `defaultGenModelInput`，它结合指令和消息，支持 f-string 风格的会话值插值

该结构体在 `NewChatModelAgent` 中进行验证，如果缺少必填字段则返回错误。

---

### 使用示例

`ChatModelAgentConfig` 的典型用法可以在"Hello World"示例中看到，其中创建了一个简单的问候代理：

```go
48:53:/home/ruijzhan/data/eino-examples/adk/helloworld/helloworld.go
agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
	Name:        "hello_agent",
	Description: "A friendly greeting assistant",
	Instruction: "You are a friendly assistant. Please respond to the user in a warm tone.",
	Model:       model,
})
```

此配置创建了一个代理，该代理：
- 标识为 `hello_agent`
- 使用提供的 `model`（可能是 OpenAI 聊天模型）
- 以温暖的系统提示开始
- 没有自定义工具或退出逻辑

在整个代码库中，`ChatModelAgentConfig` 被广泛用于定义各种代理，例如：
- **研究代理**，搜索网络
- **数学代理**，执行计算
- **代码代理**，生成软件
- **审查代理**，批评输出

它在至少 17 个文件中被引用，主要位于 `/adk/` 下的示例中，表明它是 ADK（Agent Development Kit / 代理开发套件）框架的核心部分，用于构建模块化的、使用工具的 AI 代理。

---

### 注意事项

- **`Instruction` 中的 f-string 支持**：`Instruction` 字段通过默认的 `GenModelInput` 函数支持占位符插值（例如 `{Time}`、`{User}`）。这允许使用会话值动态个性化系统提示。
- **工具驱动的工作流**：`ToolsConfig` 使代理能够调用外部工具，`ReturnDirectly` 映射允许某些工具立即终止代理循环——对于路由或最终确定响应很有用。
- **分层代理系统**：`ChatModelAgentConfig` 用于构建分层多代理系统（例如，主管和子代理），代理使用内置的 `transfer_to_agent` 工具转移任务。

---

### 另请参见

- `NewChatModelAgent`: 构造函数，接受 `ChatModelAgentConfig` 并返回 `*ChatModelAgent`。它执行验证并设置默认值。
- `ChatModelAgent`: 从配置创建的运行时代理实例。它封装了执行逻辑和状态。
- `ToolsConfig`: 嵌套配置，定义代理可以使用哪些工具以及调用它们是否应立即返回。
- `ExitTool`: 可以分配给 `Exit` 字段的内置工具实现，允许代理优雅终止。
- `GenModelInput`: 自定义代理输入格式化的函数类型。替换它允许完全控制提示构造。