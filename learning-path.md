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
- [ ] 已完成该模块
**学习目标**：掌握基本的消息创建、LLM 调用和流式响应

**前置知识**：Go 基础语法

**预计时间**：1-2 天

**核心内容**：
- 消息模板创建和使用
- OpenAI/Ollama 模型集成
- 同步和流式生成
- 基础的上下文管理

**关键文件**：
- `quickstart/chat/main.go`
- `quickstart/chat/generate.go`
- `quickstart/chat/stream.go`

### 2. quickstart/todoagent/ - 简单任务管理代理
- [ ] 已完成该模块
**学习目标**：了解基本的工具调用和状态管理

**前置知识**：`quickstart/chat`

**预计时间**：1 天

**核心内容**：
- 工具定义和调用
- 状态持久化
- 任务生命周期管理

**关键文件**：
- `quickstart/todoagent/main.go`
- `quickstart/todoagent/tool.go`

### 3. components/tool/ - 工具集成基础
- [ ] 已完成该模块
**学习目标**：掌握外部 API 调用和函数执行

**前置知识**：`quickstart/todoagent`

**预计时间**：1 天

**核心内容**：
- 工具定义和注册
- 参数验证和转换
- 错误处理机制

**关键文件**：
- `components/tool/basic/main.go`
- `components/tool/advanced/main.go`

### 4. adk/helloworld/ - ADK 框架入门
- [ ] 已完成该模块
**学习目标**：学习 Agent 开发套件的基础用法

**前置知识**：`components/tool`

**预计时间**：1-2 天

**核心内容**：
- ADK 框架基本概念
- Agent 创建和配置
- 简单的交互模式

**关键文件**：
- `adk/helloworld/main.go`
- `adk/helloworld/agent.go`

---

## 🔧 中级 - 组件深入和功能扩展

### 5. components/prompt/ - 提示工程
- [ ] 已完成该模块
**学习目标**：掌握提示模板、变量替换等高级提示技术

**前置知识**：`quickstart/chat`

**预计时间**：2 天

**核心内容**：
- 提示模板系统
- 动态变量注入
- 条件提示逻辑
- 提示优化策略

**关键文件**：
- `components/prompt/template/main.go`
- `components/prompt/variables/main.go`

### 6. components/lambda/ - 自定义函数组件
- [ ] 已完成该模块
**学习目标**：学习在 Eino 中编排自定义函数

**前置知识**：`components/tool`

**预计时间**：1-2 天

**核心内容**：
- Lambda 函数组合
- 事件驱动架构
- 外部服务集成
- 函数链编排

**关键文件**：
- `components/lambda/basic/main.go`
- `components/lambda/chain/main.go`

### 7. quickstart/eino_assistant/ - 增强型助手
- [ ] 已完成该模块
**学习目标**：构建更复杂的助手应用

**前置知识**：`components/prompt`, `components/tool`

**预计时间**：2-3 天

**核心内容**：
- 多模态助手
- 知识库集成
- 工具链管理
- 用户会话管理

**关键文件**：
- `quickstart/eino_assistant/main.go`
- `quickstart/eino_assistant/router.go`

### 8. flow/todoagent/ - 流式任务代理
- [ ] 已完成该模块
**学习目标**：理解流式处理下的任务管理

**前置知识**：`quickstart/todoagent`, `components/tool`

**预计时间**：1-2 天

**核心内容**：
- 流式事件处理
- 状态同步
- 实时任务更新

**关键文件**：
- `flow/todoagent/main.go`
- `flow/todoagent/handler.go`

---

## 🏗️ 高级 - 编排和工作流

### 9. compose/chain/ - 链式编排
- [ ] 已完成该模块
**学习目标**：学习线性工作流构建

**前置知识**：`components/prompt`, `components/lambda`

**预计时间**：2-3 天

**核心内容**：
- 链式处理模式
- 数据流转设计
- 中间件机制
- 错误传播处理

**关键文件**：
- `compose/chain/basic/main.go`
- `compose/chain/with_memory/main.go`

### 10. compose/workflow/ - 工作流编排
- [ ] 已完成该模块
**学习目标**：掌握复杂业务流程自动化

**前置知识**：`compose/chain`

**预计时间**：3-4 天

**核心内容**：
- 工作流定义和执行
- 条件分支和循环
- 并行处理
- 状态机模式

**示例文件**：
- `1_simple/` - 简单工作流
- `2_field_mapping/` - 字段映射
- `3_data_only/` - 数据流处理
- `4_control_only_branch/` - 控制分支
- `5_static_values/` - 静态值处理
- `6_stream_field_map/` - 流式字段映射

### 11. compose/graph/ - 图形化编排
- [ ] 已完成该模块
**学习目标**：学习非线性流程和分支逻辑

**前置知识**：`compose/workflow`

**预计时间**：3 天

**核心内容**：
- 图结构定义
- 节点和边的关系
- 循环依赖处理
- 复杂路由逻辑

### 12. components/retriever/ - 数据检索
- [ ] 已完成该模块
**学习目标**：实现知识库和 RAG 系统

**前置知识**：`components/document`

**预计时间**：2-3 天

**核心内容**：
- 向量数据库集成
- 相似度搜索
- 文档索引和检索
- RAG 架构设计

**关键文件**：
- `components/retriever/vector/main.go`
- `components/retriever/rag/main.go`

### 13. components/document/ - 文档处理
- [ ] 已完成该模块
**学习目标**：学习各种文件格式解析和处理

**前置知识**：`quickstart/chat`

**预计时间**：2 天

**核心内容**：
- 多格式解析器（PDF、HTML、文本）
- 文档加载器
- 内容提取和转换
- 文档分块策略

**关键文件**：
- `components/document/loader/main.go`
- `components/document/parser/main.go`

---

## 🚀 专家级 - 复杂系统架构

### 14. flow/chat/ - 流式聊天
- [ ] 已完成该模块
**学习目标**：掌握流式消息处理模式

**前置知识**：`quickstart/chat`, `compose/chain`

**预计时间**：2 天

**核心内容**：
- 双向流式对话
- 上下文同步
- 实时事件广播

**关键文件**：
- `flow/chat/main.go`
- `flow/chat/session.go`

### 15. flow/eino_assistant/ - 流式助手
- [ ] 已完成该模块
**学习目标**：构建流式多模态助手

**前置知识**：`quickstart/eino_assistant`, `flow/chat`

**预计时间**：3 天

**核心内容**：
- 多模态流式处理
- 知识库实时检索
- 用户会话同步

**关键文件**：
- `flow/eino_assistant/main.go`
- `flow/eino_assistant/router.go`

### 16. flow/agent/ - 流式代理
- [ ] 已完成该模块
**学习目标**：构建基于流的智能体系统

**前置知识**：`flow/chat`

**预计时间**：3-4 天

**核心内容**：
- 流式处理架构
- 事件驱动模式
- 实时响应系统
- 分布式代理

**子项目**：
- `deer-go/` - Deer 智能体实现
- `manus/` - Manus 智能体框架
- `multiagent/` - 多智能体协作
- `react/` - React 模式智能体

### 17. adk/intro/ - ADK 深入教程
- [ ] 已完成该模块
**学习目标**：掌握会话管理、工作流等高级特性

**前置知识**：`adk/helloworld`, `compose/workflow`

**预计时间**：4 天

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
- `workflow/` - 工作流模式

### 18. adk/multiagent/ - 多智能体系统
- [ ] 已完成该模块
**学习目标**：学习智能体协作和分布式 AI

**前置知识**：`flow/agent`, `adk/intro`

**预计时间**：4-5 天

**核心内容**：
- 多智能体架构
- 通信协议
- 协作策略
- 任务分配机制

**关键文件**：
- `adk/multiagent/collaboration/main.go`
- `adk/multiagent/coordination/main.go`

### 19. devops/observability/ - 监控与可观测性
- [ ] 已完成该模块
**学习目标**：掌握 Eino 应用的监控、追踪与诊断

**前置知识**：`compose/workflow`, 基础 DevOps 知识

**预计时间**：2-3 天

**核心内容**：
- OpenTelemetry 集成
- 日志与指标收集
- Langfuse & APM+ 监控
- 性能分析与告警

**关键文件**：
- `devops/observability/main.go`
- `devops/observability/config.go`

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
- **流式代理入门**：运行并理解 `flow/todoagent`

### 高级
- **链式编排**：能设计并实现 `compose/chain` 的自定义链路
- **工作流设计**：掌握 `compose/workflow` 的分支与循环
- **图形编排**：能利用 `compose/graph` 构建复杂路由
- **RAG 系统**：基于 `components/retriever` 完成知识检索

### 专家级
- **流式聊天**：实现 `flow/chat` 的实时协同会话
- **多模态助手**：整合 `flow/eino_assistant` 与知识库
- **多智能体架构**：设计 `adk/multiagent` 的协作策略
- **可观测性**：在 `devops/observability` 中引入监控与告警

---

## 📝 进阶建议

### 实战项目详细方案

#### 1. 智能客服机器人
- **技术栈**：`quickstart/chat` + `components/tool` + `components/retriever`
- **关键能力**：多轮对话、FAQ 检索、业务 API 调用
- **进阶方向**：结合 `flow/chat` 实现实时客服协同
- **参考资源**：Eino Examples 的 chat 项目

#### 2. 文档分析系统
- **技术栈**：`components/document` + `components/retriever` + `compose/workflow`
- **关键能力**：文档解析、向量检索、自动生成报告
- **进阶方向**：接入 `devops/observability` 监控任务执行
- **参考资源**：Eino Examples 的 document 项目

#### 3. 智能助手平台
- **技术栈**：`quickstart/eino_assistant` + `compose/chain` + `flow/eino_assistant`
- **关键能力**：多工具协同、流式响应、会话管理
- **进阶方向**：引入 `adk/multiagent` 实现协作助手
- **参考资源**：Eino Examples 的 eino_assistant 项目

#### 4. 自动化运营平台
- **技术栈**：`compose/workflow` + `compose/graph` + `flow/agent`
- **关键能力**：复杂流程编排、事件驱动、实时调度
- **进阶方向**：加入 `devops/observability` 进行全链路监控
- **参考资源**：Eino Examples 的 agent 项目

### 深入学习方向
- **性能优化**：使用 Profiling 工具分析链路性能
- **监控运维**：集成 OpenTelemetry、Langfuse 等 observability 方案
- **扩展开发**：开发自定义组件与插件，贡献给生态
- **生产部署**：学习容器化部署、CI/CD、灰度发布实践

## 📚 参考资源
- [Eino 官方文档](https://github.com/cloudwego/eino)
- [Eino Examples Issue 列表](https://github.com/cloudwego/eino-examples/issues)
- [CloudWeGo 社区](https://www.cloudwego.io/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)