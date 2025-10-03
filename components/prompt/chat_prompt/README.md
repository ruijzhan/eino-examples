# ChatPrompt 聊天提示词模板

Eino 的 ChatPrompt 组件提供了结构化的聊天提示词模板功能，支持多种格式化类型和动态变量替换。

## 核心概念

### 格式化类型 (FormatType)
Eino 支持三种模板格式化类型：

- **FString**：Python f-string 风格，使用 `{variable}` 语法
- **GoTemplate**：Go 标准模板语法，使用 `{{.variable}}` 语法
- **Jinja2**：Python Jinja2 模板语法

### 消息类型
通过 `schema` 包创建不同角色的消息：
- `schema.SystemMessage(content)` - 系统角色定义
- `schema.UserMessage(content)` - 用户消息模板
- `schema.AssistantMessage(content, toolCalls)` - 助手回复消息

### 变量替换机制
在模板中使用占位符，运行时通过 `Format()` 方法注入实际数据。

**FString 语法：**
```go
systemTpl := `你是情绪助手，你的任务是根据用户的输入，生成一段赞美的话。
用户姓名：{user_name}
用户年龄：{user_age}`
```

**GoTemplate 语法：**
```go
goSystemTpl := `你是专业助手，帮助用户{{.user_query}}。
用户信息：姓名={{.user_name}}，年龄={{.user_age}}`
```

字符串模板的处理流程如下：
```
原始字符串模板 "嗨 {user_name}..."
    │ prompt.FromMessages
    ▼
构建出的 ChatTemplate
    │ 调用 Format(ctx, variables)
    ▼
渲染后的 []*schema.Message
```

### 对话历史管理
使用 `MessagesPlaceholder` 动态插入历史对话记录：
```go
schema.MessagesPlaceholder("history", optional)
```

- 第一个参数：变量名
- 第二个参数：是否必需（`true` 必需，`false` 可选）

## 关键 API

### prompt.FromMessages / DefaultChatTemplate / ChatTemplate
- **`prompt.FromMessages`**：通过格式化类型和若干 `schema.MessagesTemplate` 构造 `*prompt.DefaultChatTemplate`
- **`prompt.DefaultChatTemplate`**：默认实现了 `prompt.ChatTemplate` 接口
- **`prompt.ChatTemplate`**：统一的聊天模板接口，暴露 `Format(ctx, map[string]any, opts ...Option)` 方法

```go
tpl := prompt.FromMessages(schema.FString,
    schema.SystemMessage("嗨 {user_name}，让我们开始学习 {topic} 编程吧！"),
    schema.UserMessage("我想了解更多关于 {topic} 的内容"),
)

var chat prompt.ChatTemplate = tpl
msgs, err := chat.Format(ctx, map[string]any{
    "user_name": "王五",
    "topic":     "Go",
})
```

### Format 方法
格式化模板并生成消息列表：
```go
msgList, err := chatTpl.Format(ctx, variables)
```

### Option / WrapImplSpecificOptFn / GetImplSpecificOptions
- **`prompt.Option`**：向模板实现传递自定义配置的载体
- **`prompt.WrapImplSpecificOptFn`**：将实现内部的配置函数封装成 `Option`
- **`prompt.GetImplSpecificOptions`**：在自定义实现中提取配置，并允许提供默认值

```go
type uppercaseOption struct { enable bool }

opt := prompt.WrapImplSpecificOptFn(func(cfg *uppercaseOption) {
    cfg.enable = true
})

cfg := prompt.GetImplSpecificOptions(&uppercaseOption{}, opt)
```

### CallbackInput / CallbackOutput / ConvCallbackInput / ConvCallbackOutput
回调辅助结构体及其转换方法：
- **`prompt.CallbackInput`**：传递变量、模板等上下文信息
- **`prompt.CallbackOutput`**：传递格式化结果
- **`prompt.ConvCallbackInput`**：将 `map[string]any` 或 `*CallbackInput` 转为 `*CallbackInput`
- **`prompt.ConvCallbackOutput`**：将 `[]*schema.Message` 或 `*CallbackOutput` 转为 `*CallbackOutput`

```go
input := prompt.ConvCallbackInput(map[string]any{"topic": "Go"})
output := prompt.ConvCallbackOutput(msgs)
```

## 高级用法示例

```go
base := prompt.FromMessages(schema.FString,
    schema.SystemMessage("嗨 {user_name}，让我们开始学习 {topic} 编程吧！"),
    schema.UserMessage("我想了解更多关于 {topic} 的内容"),
)

customTpl := newCustomChatTemplate(base)
uppercaseOpt := prompt.WrapImplSpecificOptFn(func(opt *uppercaseOption) {
    opt.enable = true
})

msgs, _ := customTpl.Format(ctx, map[string]any{
    "user_name": "王五",
    "topic":     "Go",
}, uppercaseOpt)

input := &prompt.CallbackInput{Variables: map[string]any{"topic": "Go"}}
output := prompt.ConvCallbackOutput(msgs)
```

> **提示**：示例中 `newCustomChatTemplate` 使用 `prompt.GetImplSpecificOptions` 解析自定义 `Option`，并在需要时将消息内容转换为大写，展示了如何扩展默认模板。

## 最佳实践

- **必需 vs 可选**：根据业务需求合理设置占位符的必需性
- **格式选择**：FString 适合简单替换，GoTemplate 适合复杂逻辑
- **变量命名**：使用有意义的变量名，提高代码可读性
- **错误处理**：始终检查 Format 方法的返回错误

## 应用场景

- **角色定制**：为 AI 助手定义特定角色和个性特征
- **个性化服务**：基于用户信息提供定制化响应
- **多轮对话**：维护对话历史，支持上下文感知的交互
- **模板复用**：一次定义，多处使用，提高开发效率