/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/cloudwego/eino-ext/callbacks/apmplus"
	clc "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino-ext/callbacks/langfuse"
	"github.com/cloudwego/eino-ext/devops"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/cozeloop-go"

	"github.com/cloudwego/eino-examples/quickstart/eino_assistant/eino/einoagent"
	"github.com/cloudwego/eino-examples/quickstart/eino_assistant/pkg/mem"
)

// 命令行参数：会话ID，用于标识不同的对话会话
var id = flag.String("id", "", "conversation id")

// 内存管理器：用于存储和管理对话历史记录
// GetDefaultMemory() 返回默认配置的内存实例
var memory = mem.GetDefaultMemory()

// 回调处理器：用于处理Eino框架执行过程中的各种回调事件
// 包括开始、结束、错误等生命周期事件的记录和处理
var cbHandler callbacks.Handler

// main 主函数：程序入口点
// 负责初始化系统组件并启动交互式对话循环
func main() {
	// 解析命令行参数，包括 -id 参数
	flag.Parse()

	// 开启 Eino 的可视化调试能力
	// devops.Init 初始化监控和调试功能，包括链路追踪、性能监控等
	err := devops.Init(context.Background())
	if err != nil {
		log.Printf("[eino dev] init failed, err=%v", err)
		return
	}

	// 如果没有提供会话ID，则生成一个随机的6位数字ID
	// 这确保每次运行都有唯一的会话标识
	if *id == "" {
		*id = strconv.Itoa(rand.Intn(1000000))
	}

	// 创建背景上下文，用于控制整个应用程序的生命周期
	ctx := context.Background()

	// 执行系统初始化，包括回调处理器、日志系统等
	err = Init()
	if err != nil {
		log.Printf("[eino agent] init failed, err=%v", err)
		return
	}

	// 启动交互式对话循环
	// 使用缓冲读取器从标准输入读取用户输入
	reader := bufio.NewReader(os.Stdin)
	for {
		// 显示用户提示符并等待输入
		fmt.Printf("🧑‍ : ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}

		// 清理输入，移除首尾空白字符
		input = strings.TrimSpace(input)
		// 检查退出命令
		if input == "" || input == "exit" || input == "quit" {
			return
		}

		// 调用 RunAgent 函数处理用户输入
		// 返回流式响应，允许实时显示AI回复
		sr, err := RunAgent(ctx, *id, input)
		if err != nil {
			fmt.Printf("Error from RunAgent: %v\n", err)
			continue
		}

		// 打印AI响应
		// 使用流式读取，实时显示AI的回复内容
		fmt.Print("🤖 : ")
		for {
			msg, err := sr.Recv()
			if err != nil {
				// 当收到 EOF 时表示流结束，正常退出循环
				if err == io.EOF {
					break
				}
				fmt.Printf("Error receiving message: %v\n", err)
				break
			}
			// 实时打印消息内容，实现流式显示效果
			fmt.Print(msg.Content)
		}
		// 在每轮对话结束后添加空行，提高可读性
		fmt.Println()
		fmt.Println()
	}
}

// Init 系统初始化函数
// 负责初始化日志系统、回调处理器和各种监控组件
func Init() error {
	// 创建初始化上下文
	ctx := context.Background()

	// 创建日志目录，确保日志文件有地方存储
	os.MkdirAll("log", 0755)

	// 打开日志文件，以追加模式写入
	// O_CREATE: 文件不存在时创建
	// O_WRONLY: 只写模式
	// O_APPEND: 追加模式，不清空原有内容
	var f *os.File
	f, err := os.OpenFile("log/eino.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// 配置日志回调处理器
	cbConfig := &LogCallbackConfig{
		Detail: true,  // 启用详细日志记录
		Writer: f,     // 指定日志写入到文件
	}
	// 如果环境变量 DEBUG=true，则启用调试模式
	if os.Getenv("DEBUG") == "true" {
		cbConfig.Debug = true
	}
	// 创建日志回调处理器，用于 WithCallback 调用选项
	cbHandler = LogCallback(cbConfig)

	// 初始化全局回调处理器集合，用于链路追踪和指标收集
	callbackHandlers := make([]callbacks.Handler, 0)

	// 配置 APM+ 监控（字节跳动云监控平台）
	// 当设置了 APMPLUS_APP_KEY 环境变量时启用
	if os.Getenv("APMPLUS_APP_KEY") != "" {
		// 获取监控区域，默认为北京
		region := os.Getenv("APMPLUS_REGION")
		if region == "" {
			region = "cn-beijing"
		}
		fmt.Println("[eino agent] INFO: use apmplus as callback, watch at: https://console.volcengine.com/apmplus-server")

		// 创建 APM+ 处理器，配置监控参数
		cbh, _, err := apmplus.NewApmplusHandler(&apmplus.Config{
			Host:        fmt.Sprintf("apmplus-%s.volces.com:4317", region),  // 监控数据上报地址
			AppKey:      os.Getenv("APMPLUS_APP_KEY"),                       // 应用密钥
			ServiceName: "eino-assistant",                                   // 服务名称
			Release:     "release/v0.0.1",                                   // 版本号
		})
		if err != nil {
			log.Fatal(err)
		}
		// 将 APM+ 处理器添加到回调处理器列表
		callbackHandlers = append(callbackHandlers, cbh)
	}

	// 配置 Langfuse 监控（LLM 应用可观测性平台）
	// 当同时设置了 LANGFUSE_PUBLIC_KEY 和 LANGFUSE_SECRET_KEY 环境变量时启用
	if os.Getenv("LANGFUSE_PUBLIC_KEY") != "" && os.Getenv("LANGFUSE_SECRET_KEY") != "" {
		fmt.Println("[eino agent] INFO: use langfuse as callback, watch at: https://cloud.langfuse.com")

		// 创建 Langfuse 处理器，专门用于 LLM 应用的追踪和分析
		cbh, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
			Host:      "https://cloud.langfuse.com",                    // Langfuse 服务器地址
			PublicKey: os.Getenv("LANGFUSE_PUBLIC_KEY"),                // 公钥
			SecretKey: os.Getenv("LANGFUSE_SECRET_KEY"),                // 私钥
			Name:      "Eino Assistant",                                // 应用名称
			Public:    true,                                            // 公开模式
			Release:   "release/v0.0.1",                               // 版本号
			UserID:    "eino_god",                                      // 用户标识
			Tags:      []string{"eino", "assistant"},                  // 标签，用于分类和过滤
		})
		// 将 Langfuse 处理器添加到回调处理器列表
		callbackHandlers = append(callbackHandlers, cbh)
	}

	// 配置 Coze Loop 监控（字节跳动 AI 开发平台）
	// 参考文档: https://loop.coze.cn/open/docs/cozeloop/go-sdk#4a8c980e
	cozeloopApiToken := os.Getenv("COZELOOP_API_TOKEN")
	cozeloopWorkspaceID := os.Getenv("COZELOOP_WORKSPACE_ID")
	if cozeloopApiToken != "" && cozeloopWorkspaceID != "" {
		// 创建 Coze Loop 客户端，用于连接 Coze AI 平台
		client, err := cozeloop.NewClient(
			cozeloop.WithAPIToken(cozeloopApiToken),           // API 令牌
			cozeloop.WithWorkspaceID(cozeloopWorkspaceID),    // 工作空间 ID
		)
		if err != nil {
			panic(err)
		}
		// 确保在函数退出时关闭客户端连接
		defer client.Close(ctx)
		// 创建 Coze Loop 处理器并添加到回调处理器列表
		callbackHandlers = append(callbackHandlers, clc.NewLoopHandler(client))
	}

	// 如果配置了任何回调处理器，则初始化全局回调处理器
	// 这样所有 Eino 组件的执行都会被追踪和记录
	if len(callbackHandlers) > 0 {
		callbacks.InitCallbackHandlers(callbackHandlers)
	}

	return nil
}

// RunAgent 运行AI代理处理用户消息的核心函数
// 返回流式响应，支持实时显示AI回复
func RunAgent(ctx context.Context, id string, msg string) (*schema.StreamReader[*schema.Message], error) {
	// 构建 Eino 代理执行器
	// BuildEinoAgent 创建一个包含所有必要组件的AI代理图
	runner, err := einoagent.BuildEinoAgent(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build agent graph: %w", err)
	}

	// 从内存中获取或创建指定ID的对话会话
	// 第二个参数 true 表示如果不存在则创建新会话
	conversation := memory.GetConversation(id, true)

	// 构造用户消息结构体，包含当前输入和历史对话
	userMessage := &einoagent.UserMessage{
		ID:      id,                              // 会话ID
		Query:   msg,                             // 用户输入的查询内容
		History: conversation.GetMessages(),      // 获取该会话的历史消息列表
	}

	// 使用流式方式运行代理
	// compose.WithCallbacks 指定使用自定义回调处理器进行追踪
	sr, err := runner.Stream(ctx, userMessage, compose.WithCallbacks(cbHandler))
	if err != nil {
		return nil, fmt.Errorf("failed to stream: %w", err)
	}

	// 将流复制成两路
	// srs[0] 用于返回给调用者显示给用户
	// srs[1] 用于在后台goroutine中保存到内存
	srs := sr.Copy(2)

	// 启动后台goroutine处理消息流保存到内存
	go func() {
		// 用于收集完整的AI回复消息片段
		fullMsgs := make([]*schema.Message, 0)

		defer func() {
			// 关闭第二路流，释放资源
			srs[1].Close()

			// 将用户输入添加到对话历史
			conversation.Append(schema.UserMessage(msg))

			// 将收集到的消息片段合并成完整消息
			fullMsg, err := schema.ConcatMessages(fullMsgs)
			if err != nil {
				fmt.Println("error concatenating messages: ", err.Error())
			}
			// 将AI的完整回复添加到对话历史
			conversation.Append(fullMsg)
		}()

		// 循环接收流式消息直到结束
	outer:
		for {
			select {
			// 检查上下文是否被取消
			case <-ctx.Done():
				fmt.Println("context done", ctx.Err())
				return
			default:
				// 接收消息片段
				chunk, err := srs[1].Recv()
				if err != nil {
					// 如果收到EOF，表示流结束，正常退出
					if errors.Is(err, io.EOF) {
						break outer
					}
				}

				// 收集消息片段，用于后续合并保存
				fullMsgs = append(fullMsgs, chunk)
			}
		}
	}()

	// 返回第一路流给调用者，用于实时显示给用户
	return srs[0], nil
}

// LogCallbackConfig 日志回调配置结构体
// 控制日志记录的详细程度和输出目标
type LogCallbackConfig struct {
	Detail bool        // 是否记录详细信息
	Debug  bool        // 是否启用调试模式（包含格式化JSON）
	Writer io.Writer   // 日志输出目标（文件或标准输出）
}

// LogCallback 创建日志回调处理器
// 用于记录 Eino 组件的执行过程和结果
func LogCallback(config *LogCallbackConfig) callbacks.Handler {
	// 如果配置为空，使用默认配置
	if config == nil {
		config = &LogCallbackConfig{
			Detail: true,          // 默认启用详细信息
			Writer: os.Stdout,     // 默认输出到标准输出
		}
	}
	// 确保输出目标不为空
	if config.Writer == nil {
		config.Writer = os.Stdout
	}

	// 使用构建器模式创建回调处理器
	builder := callbacks.NewHandlerBuilder()

	// 注册组件开始执行的回调函数
	builder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		// 记录组件开始执行的信息：组件类型:组件类别:组件名称
		fmt.Fprintf(config.Writer, "[view]: start [%s:%s:%s]\n", info.Component, info.Type, info.Name)

		// 如果启用了详细信息记录，则记录输入数据
		if config.Detail {
			var b []byte
			if config.Debug {
				// 调试模式：格式化JSON输出，便于阅读
				b, _ = json.MarshalIndent(input, "", "  ")
			} else {
				// 普通模式：压缩JSON输出，节省空间
				b, _ = json.Marshal(input)
			}
			fmt.Fprintf(config.Writer, "%s\n", string(b))
		}
		return ctx
	})

	// 注册组件结束执行的回调函数
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		// 记录组件执行结束的信息
		fmt.Fprintf(config.Writer, "[view]: end [%s:%s:%s]\n", info.Component, info.Type, info.Name)
		return ctx
	})

	// 构建并返回回调处理器
	return builder.Build()
}
