# Callback 功能演示示例

本示例展示了如何在 Eino 框架中使用 Callback 机制来监控和追踪工具执行过程。

## 主要内容

示例演示了完整的 Callback 使用流程，包括：

1. **创建 Callback Handler** - 定义工具执行的生命周期回调函数
2. **使用 jsonschema.Reflect() 优化** - 从 struct 自动生成 JSON Schema，避免重复定义
3. **集成到 Chain 中** - 将 Callback 与 Eino 的 Chain 组合使用
4. **工具生命周期监控** - 追踪工具执行的开始和结束状态
5. **参数和结果监控** - 捕获工具调用的输入参数和返回结果
6. **并发 vs 串行执行** - 演示多个工具调用的执行模式

## 核心知识点

### 1. Callback Handler 结构
CallbackHandler 提供了三个主要的生命周期钩子：
- **OnStart**: 工具开始执行时触发，可以访问传入参数和运行信息
- **OnEnd**: 工具执行完成时触发，可以访问执行结果和运行信息
- **OnEndWithStreamOutput**: 处理流式输出的工具调用

### 2. 运行信息获取
通过 `callbacks.RunInfo` 可以获取：
- **Component**: 组件类型（如 Tool）
- **Type**: 具体实现类型（如 GetWeather）
- **Name**: 工具名称（如 get_weather）

### 3. Callback 集成模式
- 创建 `ToolCallbackHandler` 实例
- 通过 `callbackHelper.NewHandlerHelper()` 包装 handler
- 使用 `compose.WithCallbacks()` 将 callback 注入到 Chain 执行中

### 4. 工具调用模式
- **并发执行**（默认）: 多个工具调用同时执行，提高效率
- **串行执行**（可选）: 通过 `ExecuteSequentially: true` 按顺序执行

### 5. InferTool 自动推断机制
- **传统方式缺陷**: 手动定义 JSON Schema 和 Go struct 存在重复
- **最佳实践**: 使用 `utils.InferTool()` 自动从 struct 生成 JSON Schema 并创建工具
- **内部机制**: `InferTool` 内部使用 `jsonschema.Reflect()` 自动推断参数结构
- **一站式创建**: 一次调用同时完成 ToolInfo 创建和工具实例化
- **推荐用法**: 这是 Eino 框架中创建工具的推荐方式

### 6. Chain 组合模式
- 使用 `compose.NewChain[InputType, OutputType]()` 创建 Chain
- 通过 `chain.AppendToolsNode()` 添加工具节点
- 使用 `chain.Compile(ctx)` 编译为可执行对象

## Callback 机制

### 工具执行流程
```
输入 Message → Callback OnStart → 工具执行 → Callback OnEnd → 输出 Message
     ↓              ↓                ↓             ↓           ↓
  工具调用参数    监控开始状态      业务逻辑处理    监控完成状态   工具执行结果
```

### 并发执行特点
- 多个工具调用同时启动
- 执行完成顺序不确定
- 适合独立工具调用场景

### 串行执行特点
- 按输入顺序逐个执行
- 执行完成顺序确定
- 适合有依赖关系的工具调用

## 读者掌握要点

学习本示例后，读者应该能够：

- **理解 Callback 机制** - 掌握 Eino 框架中的生命周期回调概念
- **使用 InferTool()** - 掌握 Eino 框架推荐的工具创建方式
- **创建自定义 Callback** - 能够根据业务需求定制监控逻辑
- **集成 Callback 到应用** - 学会在 Chain 和 Graph 中使用 Callback
- **监控工具执行** - 追踪工具调用的参数、结果和执行状态
- **控制执行模式** - 根据需求选择并发或串行执行策略

## 运行示例

在 `components/tool/callback/` 目录下执行：

```bash
go run main.go
```

程序将演示完整的 Callback 功能，包括：
- 单次工具调用的 Callback 监控
- 多次并发工具调用的 Callback 监控
- 工具执行过程的详细追踪信息

## 进阶用法

### 串行执行设置
如果需要按顺序执行工具调用，可以在创建 ToolsNode 时设置：

```go
toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
    Tools:              []tool.BaseTool{weatherTool},
    ExecuteSequentially: true,  // 强制串行执行
})
```

### InferTool 使用示例
```go
// 推荐：使用 InferTool 创建工具
type WeatherParams struct {
    City string `json:"city" jsonschema:"required,description=城市名称"`
    Unit string `json:"unit,omitempty" jsonschema:"description=温度单位"`
}

weatherTool, err := utils.InferTool(
    "get_weather",
    "获取指定城市的天气信息",
    func(ctx context.Context, params *WeatherParams) (string, error) {
        // 工具逻辑实现
        return "天气数据", nil
    },
)
```

### 自定义 Callback 扩展
可以根据业务需求扩展 Callback 功能：
- 添加日志记录
- 集成监控系统
- 实现执行时间统计
- 添加错误处理逻辑