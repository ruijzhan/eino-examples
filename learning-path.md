# Eino Examples 学习路径

根据项目结构和复杂度分析，制定了以下由浅入深的学习顺序，帮助开发者系统掌握 Eino 框架。

## 🗂️ 目录
- [📚 总览](#-总览)
- [🌱 入门级 - 基础概念和简单使用](#-入门级---基础概念和简单使用)
- [🔧 中级 - 组件深入和功能扩展](#-中级---组件深入和功能扩展)
- [🏗️ 高级 - 编排和工作流](#-高级---编排和工作流)
- [🚀 专家级 - 复杂系统架构](#-专家级---复杂系统架构)
- [💡 学习建议](#-学习建议)
- [🔧 常见问题排查](#-常见问题排查)
- [📌 学习检查清单](#-学习检查清单)
- [📝 进阶建议](#-进阶建议)
- [📚 参考资源](#-参考资源)

## 📚 总览

Eino Examples 项目按学习难度分为四个等级，建议按此顺序学习：

- **🌱 入门级**：基础概念和简单使用
- **🔧 中级**：组件深入和功能扩展
- **🏗️ 高级**：编排和工作流
- **🚀 专家级**：复杂系统架构

---

## 🌱 入门级 - 基础概念和简单使用

### 1. quickstart/chat/ - 基础聊天机器人
- [x] 已完成该模块
**学习目标**：掌握基本的消息创建、LLM 调用和流式响应

**前置知识**：Go 基础语法

**难度**：⭐⭐

**核心内容**：
- 消息模板创建和使用
- OpenAI/Ollama 模型集成
- 同步和流式生成
- 基础的上下文管理

**关键文件**：
- `quickstart/chat/main.go` - 主程序入口
- `quickstart/chat/generate.go` - 同步生成示例
- `quickstart/chat/stream.go` - 流式生成示例
- `quickstart/chat/openai.go` - OpenAI 模型集成
- `quickstart/chat/ollama.go` - Ollama 模型集成
- `quickstart/chat/template.go` - 消息模板示例

### 2. quickstart/todoagent/ - 简单任务管理代理
- [x] 已完成该模块
**学习目标**：了解基本的工具调用和状态管理

**前置知识**：`quickstart/chat`

**难度**：⭐⭐

**核心内容**：
- 工具定义和调用
- 状态持久化
- 任务生命周期管理

**关键文件**：
- `quickstart/todoagent/main.go` - 主程序入口
- `quickstart/todoagent/README.md` - 完整使用说明
- `quickstart/todoagent/Chain.md` - 链式编排详解
- `quickstart/todoagent/NewTool.md` - NewTool 方式说明
- `quickstart/todoagent/InferTool.md` - InferTool 方式说明
- `quickstart/todoagent/InvokableTool.md` - InvokableTool 接口说明
- `quickstart/todoagent/ToolsNode.md` - ToolsNode 使用说明
- `quickstart/todoagent/ToolInfo.md` - ToolInfo 结构说明
- `quickstart/todoagent/Runnable.md` - Runnable 接口说明

### 3. components/tool/ - 工具集成基础
- [x] 已完成该模块
**学习目标**：掌握外部 API 调用和函数执行

**前置知识**：`quickstart/todoagent`

**难度**：⭐⭐

**核心内容**：
- 工具定义和注册
- 参数验证和转换
- 错误处理机制

**关键文件**：
- `components/tool/jsonschema/main.go`
- `components/tool/callback/main.go`

### 4. adk/helloworld/ - ADK 框架入门
- [x] 已完成该模块
**学习目标**：学习 Agent 开发套件的基础用法

**前置知识**：`components/tool`

**难度**：⭐⭐

**核心内容**：
- ADK 框架基本概念
- Agent 创建和配置
- 简单的交互模式

**关键文件**：
- `adk/helloworld/helloworld.go` - 主程序入口
- `adk/helloworld/README.md` - 完整的 ADK 概述和架构说明
- `adk/helloworld/ChatModelAgentConfig.md` - Agent 配置结构详解
- `adk/helloworld/NewChatModelAgent.md` - Agent 构造函数说明
- `adk/helloworld/RunnerConfig.md` - Runner 配置结构详解
- `adk/helloworld/NewRunner.md` - Runner 构造函数说明
- `adk/helloworld/Run.md` - Runner 执行方法详解（包含 Query 方法对比）

---

## 🔧 中级 - 组件深入和功能扩展

### 5. components/prompt/ - 提示工程
- [x] 已完成该模块
**学习目标**：掌握提示模板、变量替换等高级提示技术

**前置知识**：`quickstart/chat`

**难度**：⭐⭐⭐

**核心内容**：
- 提示模板系统
- 动态变量注入
- 条件提示逻辑
- 提示优化策略

**关键文件**：
- `components/prompt/chat_prompt/chat_prompt.go`

### 6. components/lambda/ - 自定义函数组件
- [ ] 已完成该模块
**学习目标**：学习在 Eino 中编排自定义函数

**前置知识**：`components/tool`

**难度**：⭐⭐

**核心内容**：
- Lambda 函数组合
- 事件驱动架构
- 外部服务集成
- 函数链编排

**关键文件**：
- `components/lambda/lambda.go`

### 7. quickstart/eino_assistant/ - 增强型助手
- [ ] 已完成该模块
**学习目标**：构建更复杂的助手应用

**前置知识**：`components/prompt`, `components/tool`

**难度**：⭐⭐⭐

**核心内容**：
- 多模态助手
- 知识库集成
- 工具链管理
- 用户会话管理

**关键文件**：
- `quickstart/eino_assistant/cmd/einoagent/main.go`
- `quickstart/eino_assistant/eino/einoagent/flow.go`

### 8. adk/intro/chatmodel/ - ChatModel Agent 进阶
- [ ] 已完成该模块
**学习目标**：掌握 ChatModelAgent 的高级用法和中断机制

**前置知识**：`adk/helloworld`, `components/tool`

**难度**：⭐⭐

**核心内容**：
- ChatModelAgent 使用
- 中断（Interrupt）机制
- Agent 状态管理
- 高级交互模式

**关键文件**：
- `adk/intro/chatmodel/` 目录下的示例代码

---

## 🏗️ 高级 - 编排和工作流

### 9. compose/chain/ - 链式编排
- [ ] 已完成该模块
**学习目标**：学习线性工作流构建

**前置知识**：`components/prompt`, `components/lambda`

**难度**：⭐⭐⭐

**核心内容**：
- 链式处理模式
- 数据流转设计
- 中间件机制
- 错误传播处理

**关键文件**：
- `compose/chain/main.go`

### 10. compose/workflow/ - 工作流编排
- [ ] 已完成该模块
**学习目标**：掌握复杂业务流程自动化

**前置知识**：`compose/chain`

**难度**：⭐⭐⭐⭐

**核心内容**：
- 工作流定义和执行
- 字段映射（Field Mapping）
- 条件分支（Branch）
- 静态值注入
- 流式字段映射
- 数据流控制

**关键文件**：
- `compose/workflow/1_simple/` - 简单工作流
- `compose/workflow/2_field_mapping/` - 字段映射
- `compose/workflow/3_data_only/` - 纯数据流
- `compose/workflow/4_control_only_branch/` - 条件分支
- `compose/workflow/5_static_values/` - 静态值
- `compose/workflow/6_stream_field_map/` - 流式字段映射

### 11. compose/graph/ - 图形化编排
- [ ] 已完成该模块
**学习目标**：学习非线性流程和分支逻辑

**前置知识**：`compose/workflow`

**难度**：⭐⭐⭐⭐

**核心内容**：
- 图结构定义
- 节点和边的关系
- 状态管理
- ReAct 模式实现
- 工具调用编排
- 中断与恢复机制

**关键文件**：
- `compose/graph/simple/` - 简单图编排
- `compose/graph/state/` - 状态管理
- `compose/graph/react_with_interrupt/` - ReAct 模式与中断
- `compose/graph/tool_call_agent/` - 工具调用 Agent
- `compose/graph/tool_call_once/` - 单次工具调用
- `compose/graph/two_model_chat/` - 双模型对话

### 12. components/retriever/ - 数据检索
- [ ] 已完成该模块
**学习目标**：实现知识库和 RAG 系统

**前置知识**：`components/document`

**难度**：⭐⭐⭐

**核心内容**：
- 向量数据库集成
- 相似度搜索
- 文档索引和检索
- RAG 架构设计

**关键文件**：
- `components/retriever/multiquery/main.go`
- `components/retriever/router/main.go`

### 13. components/document/ - 文档处理
- [ ] 已完成该模块
**学习目标**：学习各种文件格式解析和处理

**前置知识**：`quickstart/chat`

**难度**：⭐⭐⭐

**核心内容**：
- 多格式解析器（PDF、HTML、文本）
- 文档加载器
- 内容提取和转换
- 文档分块策略

**关键文件**：
- `components/document/parser/extparser/` - 扩展格式解析器
- `components/document/parser/textparser/` - 文本解析器
- `components/document/parser/customparser/` - 自定义解析器

---

## 🚀 专家级 - 复杂系统架构

### 14. flow/agent/ - 流式代理系统
- [ ] 已完成该模块
**学习目标**：构建基于流的智能体系统和复杂应用

**前置知识**：`adk/intro`, `compose/graph`

**难度**：⭐⭐⭐⭐⭐

**核心内容**：
- 流式处理架构
- 事件驱动模式
- ReAct 模式实现
- 多智能体协作
- 复杂应用架构

**关键文件**：
- `flow/agent/react/` - ReAct 模式实现
- `flow/agent/deer-go/` - Deer-Go 智能体示例
- `flow/agent/manus/` - Manus 智能体示例
- `flow/agent/multiagent/` - 多智能体系统

### 15. adk/intro/ - ADK 深入教程
- [ ] 已完成该模块
**学习目标**：掌握会话管理、工作流等高级特性

**前置知识**：`adk/helloworld`, `compose/workflow`

**难度**：⭐⭐⭐⭐

**核心内容**：
- 会话状态管理
- 高级工作流模式
- 自定义 Agent 开发
- 性能优化策略

**模块结构**：
- `chatmodel/` - 聊天模型高级用法
- `custom/` - 自定义 Agent
- `session/` - 会话管理
- `transfer/` - 状态转换
- `workflow/` - 工作流模式（包含 sequential、parallel、loop 等子目录）

### 16. adk/multiagent/ - 多智能体系统
- [ ] 已完成该模块
**学习目标**：学习智能体协作和分布式 AI

**前置知识**：`flow/agent`, `adk/intro`

**难度**：⭐⭐⭐⭐⭐

**核心内容**：
- 多智能体架构
- Supervisor 模式
- Plan-Execute-Replan 模式
- 分层 Supervisor 架构
- 项目管理智能体
- 通信协议与协作策略

**关键文件**：
- `adk/multiagent/supervisor/` - 基础 Supervisor 模式
- `adk/multiagent/plan-execute-replan/` - 计划-执行-重新计划模式
- `adk/multiagent/layered-supervisor/` - 分层 Supervisor
- `adk/multiagent/integration-project-manager/` - 项目管理示例

### 17. devops/debug/ - 调试与追踪
- [ ] 已完成该模块
**学习目标**：掌握 Eino 应用的调试、追踪与诊断技术

**前置知识**：`compose/chain`, `compose/graph`, 基础 DevOps 知识

**难度**：⭐⭐⭐

**核心内容**：
- 链式编排调试
- 图编排调试
- 日志与追踪
- 问题定位与诊断
- CozeLoop 集成（可选）

**关键文件**：
- `devops/debug/chain/` - 链式编排调试
- `devops/debug/graph/` - 图编排调试
- `devops/debug/main.go` - 调试工具入口

**注意**：各示例项目中也包含日志和追踪实现，可结合学习

---

## 💡 学习建议

### 学习节奏
- **第1-2周**：完成入门级，掌握基础概念
- **第3-4周**：学习中级，深入各个组件
- **第5-6周**：研究高级，构建复杂应用
- **第7-8周**：挑战专家级，设计架构方案

### 学习方法
1. **理论先行**：每个目录都有 README 文档，先阅读理解概念
2. **代码实践**：运行示例代码，观察实际效果
3. **修改实验**：基于示例进行修改，加深理解
4. **项目应用**：结合实际项目需求进行练习

### 环境准备
- Go 1.24+ 开发环境
- API 密钥配置（OpenAI、Ollama 等）
- Redis 等依赖服务（根据需要）

## 🔧 常见问题排查

### API 连接问题
- 检查 `.env` 配置和 API Key 是否有效
- 确认网络与代理设置允许访问外部模型服务
- 使用 `curl` 或 Postman 验证模型端点可用性

### 依赖问题
- 各子目录可能有独立的 `go.mod`，进入目录后执行 `go mod tidy`
- 如遇缺失依赖，可对照 `vendor/` 目录确认版本

### 运行错误
- 确保使用 Go 1.24+（推荐 go1.24.4 工具链）
- 必要的外部服务（如 Redis）需提前启动
- 若使用 GPU/本地模型，确认驱动与模型文件可用

### 调试技巧
- 使用 `go test -run TestName ./path/...` 定位问题
- 在关键流程添加日志或使用 tracing 工具观测
- 善用 `make` 或脚本查看项目提供的快捷命令

## 📌 学习检查清单

### 入门级
- **配置与运行**：能独立配置并运行 `quickstart/chat`
- **工具调用**：理解并扩展 `quickstart/todoagent` 的工具逻辑
- **ADK 基础**：使用 `adk/helloworld` 创建简单 Agent

### 中级
- **提示工程**：熟练使用 `components/prompt` 的模板与变量
- **自定义函数**：掌握 `components/lambda` 的函数编排
- **增强助手**：扩展 `quickstart/eino_assistant` 实现多工具协同
- **ChatModel Agent**：运行并理解 `adk/intro/chatmodel` 的中断机制

### 高级
- **链式编排**：能设计并实现 `compose/chain` 的自定义链路
- **工作流设计**：掌握 `compose/workflow` 的分支与循环
- **图形编排**：能利用 `compose/graph` 构建复杂路由
- **RAG 系统**：基于 `components/retriever` 完成知识检索

### 专家级
- **流式代理系统**：实现 `flow/agent` 的各种模式（ReAct、多智能体等）
- **复杂图编排**：掌握 `compose/graph` 的所有编排模式
- **多智能体架构**：设计 `adk/multiagent` 的协作策略（Supervisor、Plan-Execute-Replan）
- **调试与追踪**：在 `devops/debug` 中掌握调试技术

---

## 📝 进阶建议

### 实战项目详细方案

#### 1. 智能客服机器人
- **技术栈**：`quickstart/chat` + `components/tool` + `components/retriever`
- **关键能力**：多轮对话、FAQ 检索、业务 API 调用
- **进阶方向**：结合 `flow/agent/react` 实现智能客服协同
- **参考资源**：Eino Examples 的 chat 和 todoagent 项目

#### 2. 文档分析系统
- **技术栈**：`components/document` + `components/retriever` + `compose/workflow`
- **关键能力**：文档解析、向量检索、自动生成报告
- **进阶方向**：接入 `devops/debug` 监控任务执行和调试
- **参考资源**：Eino Examples 的 document 和 retriever 项目

#### 3. 智能助手平台
- **技术栈**：`quickstart/eino_assistant` + `compose/chain` + `adk/intro`
- **关键能力**：多工具协同、流式响应、会话管理
- **进阶方向**：引入 `adk/multiagent` 实现协作助手
- **参考资源**：Eino Examples 的 eino_assistant 和 adk/intro 项目

#### 4. 自动化运营平台
- **技术栈**：`compose/workflow` + `compose/graph` + `flow/agent`
- **关键能力**：复杂流程编排、事件驱动、实时调度
- **进阶方向**：加入 `devops/debug` 进行调试和监控
- **参考资源**：Eino Examples 的 workflow、graph 和 flow/agent 项目

### 深入学习方向
- **性能优化**：使用 Profiling 工具分析链路性能
- **调试技术**：掌握 `devops/debug` 的调试方法和工具
- **可观测性**：集成 CozeLoop、OpenTelemetry 等追踪方案（参考各示例中的实现）
- **扩展开发**：开发自定义组件与插件，贡献给生态
- **生产部署**：学习容器化部署、CI/CD、灰度发布实践

## 📚 参考资源
- [Eino 官方文档](https://github.com/cloudwego/eino)
- [Eino Examples Issue 列表](https://github.com/cloudwego/eino-examples/issues)
- [CloudWeGo 社区](https://www.cloudwego.io/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)