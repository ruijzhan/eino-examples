# Eino ADK Examples / Eino ADK 示例

## About Eino ADK / 关于 Eino ADK

**Eino ADK (Agent Development Kit / 代理开发套件)** 是一个用于构建智能代理和多代理系统的综合框架。它为开发者提供了创建复杂 AI 应用程序所需的工具、模式和抽象，同时最大限度地降低了复杂性。

### Key Features / 主要特性

- **Agent Lifecycle Management** / 代理生命周期管理：从初始化到执行的完整代理生命周期控制
- **Multi-Agent Coordination** / 多代理协调：支持代理间的协作、竞争和层次化组织
- **State Persistence** / 状态持久化：通过会话管理实现跨代理的数据和状态传递
- **Tool Integration** / 工具集成：无缝集成外部工具和 API 调用
- **Interrupt Handling** / 中断处理：支持代理执行过程中的暂停、恢复和用户交互
- **Workflow Orchestration** / 工作流编排：提供循环、并行和顺序执行模式

### Use Cases / 使用场景

- **Chatbots & Virtual Assistants** / 聊天机器人和虚拟助手
- **Process Automation** / 流程自动化
- **Knowledge Management Systems** / 知识管理系统
- **Multi-Team Coordination** / 多团队协作
- **Complex Decision Making** / 复杂决策系统

This directory provides examples for Eino ADK:
本目录提供 Eino ADK 的示例：

- Agent（代理）
  - `helloworld`: simple hello-world chat agent / 简单的问候聊天代理。
  - `intro`（入门）
    - `chatmodel`: example about using `ChatModelAgent` with interrupt / 使用 `ChatModelAgent` 和中断功能的示例。
    - `custom`: shows how to implement an agent which meets the definition of ADK / 展示如何实现符合 ADK 定制的代理。
    - `workflow`: examples about using `Loop` / `Parallel` / `Sequential` agent / 关于使用 `Loop` / `Parallel` / `Sequential` 代理的示例。
    - `session`: shows how to pass data and state across agents by using session / 展示如何通过 session 在代理间传递数据和状态。
    - `transfer`: shows transfer ability by using ChatModelAgent / 通过使用 ChatModelAgent 展示传输能力。
  - `multiagent`（多代理）
    - `plan-execute-replan`: basic example of plan-execute-replan agent / 计划-执行-再计划代理的基础示例。
    - `supervisor`: basic example of supervisor agent / 主管代理的基础示例。
    - `layered-supervisor`: another example of supervisor agent, which set a supervisor agent as sub-agent of another supervisor agent / 主管代理的另一个示例，将一个主管代理设为另一个主管代理的子代理。
    - `integration-project-manager`: another example of using supervisor agent / 使用主管代理的另一个示例。
  - `common`: utils / 工具类。


Additionally, you can enable [coze-loop](https://github.com/coze-dev/coze-loop) trace for examples, see .example.env for keys.
另外，你可以为示例启用 [coze-loop](https://github.com/coze-dev/coze-loop) 跟踪，具体密钥请查看 .example_env。 