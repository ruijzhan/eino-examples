# HelloWorld 代理示例

## 概述

本示例通过创建一个简单的对话代理，展示了 Eino ADK 的基本构建模块。它演示了初始化组件、配置代理以及通过运行器执行对话的核心工作流程。

## 演示的核心概念

- **ChatModel Integration** / 聊天模型集成：如何将 OpenAI 模型集成到 ADK 中
- **Agent Configuration** / 代理配置：设置代理的基本属性和行为指令
- **Streaming Execution** / 流式执行：处理实时的流式响应
- **Event-Driven Architecture** / 事件驱动架构：通过事件流处理代理响应
- **Message Handling** / 消息处理：管理用户输入和代理输出

## 执行流程

```
环境设置
         ↓
    OpenAI 模型初始化
         ↓
   聊天模型代理创建
         ↓
      运行器配置
         ↓
   用户输入消息
         ↓
   流式事件处理
         ↓
   响应输出
```

## 组件流程

```
OpenAI ChatModel
       ↓ (注入到)
ChatModelAgent {
  Name: "hello_agent"
  Description: "A friendly greeting assistant"
  Instruction: "You are a friendly assistant..."
}
       ↓ (配置到)
Runner {
  Agent: ChatModelAgent
  EnableStreaming: true
}
       ↓ (执行)
Events Stream {
  MessageInput → Agent Processing → MessageOutput
}
```

## 架构模式

本示例遵循 **模型-代理-运行器** 模式，这是 Eino ADK 中的基础模式：

1. **模型层**：提供 AI 推理能力
2. **代理层**：封装模型并添加业务逻辑
3. **运行器层**：管理执行流程和事件处理

## 说明的关键特性

- **模块化设计**：清晰的组件分离和职责划分
- **配置驱动**：通过配置而非代码控制行为
- **流式处理**：实时响应处理能力
- **错误处理**：优雅的错误处理机制
- **可扩展性**：易于添加新功能和组件

## 使用方法

1. 设置环境变量：
   ```bash
   export OPENAI_API_KEY=your-api-key
   export OPENAI_MODEL=your-model-name
   export OPENAI_BASE_URL=your-base-url
   ```

2. 运行示例：
   ```bash
   go run helloworld.go
   ```

代理将以友好的方式回应问候语 "Hello, please introduce yourself"。