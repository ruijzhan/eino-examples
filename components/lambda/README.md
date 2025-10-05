# Lambda 组件 - 自定义函数编排

> **难度**：⭐⭐ | **前置知识**：`quickstart/chat`、`components/tool` | **预计时间**：2-3天

## 📋 学习导航

- [🎯 快速开始](#-快速开始) - 5分钟上手 Lambda
- [🔄 四种交互模式](#-四种交互模式) - 核心概念理解
- [🛠️ 实践演练](#️-实践演练) - 动手编写 Lambda 函数
- [🏗️ 编排集成](#️-编排集成) - 在 Chain 和 Graph 中使用
- [📚 进阶技巧](#-进阶技巧) - 最佳实践和性能优化
- [🔧 API 参考](#-api-参考) - 完整 API 文档

---

## 🎯 快速开始

### 核心概念

Lambda 组件是 Eino 框架中的**万能适配器**，让你能够在 AI 编排流程中插入自定义 Go 函数。通过 Lambda，你可以：

- **🔄 数据转换**：在链路中进行任意类型的数据转换
- **🌐 服务集成**：调用第三方 API 或执行复杂业务逻辑
- **⚡ 流式处理**：支持流式数据的实时转换和处理
- **🛡️ 类型安全**：保持强类型检查，确保数据流转安全

### 你的第一个 Lambda

```go
package main

import (
    "context"
    "fmt"
    "strings"

    "github.com/cloudwego/eino/components/compose"
)

func main() {
    // 创建一个简单的文本处理 Lambda
    lambda := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
        // 转换为大写并添加前缀
        return "处理结果: " + strings.ToUpper(input), nil
    })

    // 使用 Lambda
    result, err := lambda.Invoke(context.Background(), "hello lambda")
    if err != nil {
        panic(err)
    }

    fmt.Println(result) // 输出: 处理结果: HELLO LAMBDA
}
```

**💡 运行这个例子：**
```bash
cd components/lambda
go run main.go
```

---

## 🔄 四种交互模式

Lambda 基于输入/输出是否为流形成 4 种模式。理解这 4 种模式是掌握 Lambda 的关键。

### 模式选择矩阵

| 输入类型 | 输出类型 | 模式名称 | 适用场景 | 难度 |
|---------|---------|---------|---------|------|
| 单个值 | 单个值 | **Invoke** | 简单转换、API调用 | ⭐ |
| 单个值 | 流式 | **Stream** | 文本生成、数据流 | ⭐⭐ |
| 流式 | 单个值 | **Collect** | 数据聚合、总结 | ⭐⭐ |
| 流式 | 流式 | **Transform** | 实时处理、过滤 | ⭐⭐⭐ |

### 1. Invoke 模式 ⭐

**同步处理，一次输入一次输出**

```go
func(ctx context.Context, input I, opts ...TOption) (output O, err error)
```

**适用场景**：
- ✅ 数据格式转换（JSON解析、类型转换）
- ✅ 简单计算和字符串处理
- ✅ HTTP API 调用
- ✅ 数据验证和过滤

**实战示例**：
```go
// 温度转换 Lambda
tempConverter := compose.InvokableLambda(func(ctx context.Context, celsius float64) (fahrenheit float64, err error) {
    return celsius*9/5 + 32, nil
})

// 用户信息验证
userValidator := compose.InvokableLambda(func(ctx context.Context, user User) (bool, error) {
    return user.Age >= 18 && user.Email != "", nil
})
```

### 2. Stream 模式 ⭐⭐

**单个输入，流式输出**

```go
func(ctx context.Context, input I, opts ...TOption) (output *schema.StreamReader[O], err error)
```

**适用场景**：
- ✅ 文本生成（逐字输出）
- ✅ 数据分批处理
- ✅ 长时间计算的结果流式返回

**实战示例**：
```go
// 文本分词流式输出
wordStreamer := compose.StreamableLambda(func(ctx context.Context, text string) (*schema.StreamReader[string], error) {
    words := strings.Split(text, " ")
    reader := schema.NewStreamReader[string]()

    go func() {
        defer reader.Close()
        for _, word := range words {
            reader.Send(ctx, word)
            time.Sleep(100 * time.Millisecond) // 模拟处理延迟
        }
    }()

    return reader.Recv(), nil
})

// 数据库查询结果流式返回
dataStreamer := compose.StreamableLambda(func(ctx context.Context, query string) (*schema.StreamReader[Record], error) {
    // 模拟数据库查询
    return streamQueryResults(ctx, query), nil
})
```

### 3. Collect 模式 ⭐⭐

**流式输入，单个输出**

```go
func(ctx context.Context, input *schema.StreamReader[I], opts ...TOption) (output O, err error)
```

**适用场景**：
- ✅ 数据汇总和统计
- ✅ 流数据的聚合计算
- ✅ 批量处理结果收集

**实战示例**：
```go
// 数字流求和
sumCollector := compose.CollectableLambda(func(ctx context.Context, numbers *schema.StreamReader[int]) (int, error) {
    sum := 0
    for {
        num, err := numbers.Recv()
        if err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            return 0, err
        }
        sum += num
    }
    return sum, nil
})

// 文本片段合并
textMerger := compose.CollectableLambda(func(ctx context.Context, fragments *schema.StreamReader[string]) (string, error) {
    var builder strings.Builder
    for {
        fragment, err := fragments.Recv()
        if err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            return "", err
        }
        builder.WriteString(fragment)
    }
    return builder.String(), nil
})
```

### 4. Transform 模式 ⭐⭐⭐

**流式输入，流式输出**

```go
func(ctx context.Context, input *schema.StreamReader[I], opts ...TOption) (output *schema.StreamReader[O], err error)
```

**适用场景**：
- ✅ 实时数据过滤
- ✅ 流式数据转换
- ✅ 数据清洗和预处理

**实战示例**：
```go
// 过滤偶数
evenFilter := compose.TransformableLambda(func(ctx context.Context, numbers *schema.StreamReader[int]) (*schema.StreamReader[int], error) {
    output := schema.NewStreamReader[int]()

    go func() {
        defer output.Close()
        for {
            num, err := numbers.Recv()
            if err != nil {
                if errors.Is(err, io.EOF) {
                    break
                }
                return
            }
            if num%2 == 0 {
                output.Send(ctx, num)
            }
        }
    }()

    return output.Recv(), nil
})

// 实时翻译
translator := compose.TransformableLambda(func(ctx context.Context, texts *schema.StreamReader[string]) (*schema.StreamReader[string], error) {
    output := schema.NewStreamReader[string]()

    go func() {
        defer output.Close()
        for {
            text, err := texts.Recv()
            if err != nil {
                if errors.Is(err, io.EOF) {
                    break
                }
                return
            }
            // 调用翻译API
            translated, _ := translateText(ctx, text, "en", "zh")
            output.Send(ctx, translated)
        }
    }()

    return output.Recv(), nil
})
```

---

## 🛠️ 实践演练

### 创建方法对比

| 方法 | 灵活性 | 复杂度 | 推荐场景 |
|------|--------|--------|----------|
| `InvokableLambda` | 低 | ⭐ | 简单转换 |
| `InvokableLambdaWithOption` | 中 | ⭐⭐ | 需要配置选项 |
| `AnyLambda` | 高 | ⭐⭐⭐ | 多模式支持 |

### 练习1：简单的数据处理 Lambda

```go
// 任务：创建一个Lambda，将字符串转换为JSON格式
type Person struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}

stringToPerson := compose.InvokableLambda(func(ctx context.Context, input string) (*Person, error) {
    parts := strings.Split(input, ",")
    if len(parts) != 3 {
        return nil, fmt.Errorf("输入格式错误，应为：姓名,年龄,邮箱")
    }

    age, err := strconv.Atoi(strings.TrimSpace(parts[1]))
    if err != nil {
        return nil, fmt.Errorf("年龄转换失败: %v", err)
    }

    return &Person{
        Name:  strings.TrimSpace(parts[0]),
        Age:   age,
        Email: strings.TrimSpace(parts[2]),
    }, nil
})

// 使用示例
person, _ := stringToPerson.Invoke(ctx, "张三,25,zhangsan@example.com")
fmt.Printf("%+v\n", person)
```

### 练习2：带自定义选项的 Lambda

```go
// 任务：创建一个可配置的文本格式化 Lambda
type FormatOptions struct {
    Prefix string
    Suffix string
    Upper  bool
}

type FormatOption func(*FormatOptions)

func WithPrefix(prefix string) FormatOption {
    return func(opts *FormatOptions) {
        opts.Prefix = prefix
    }
}

func WithSuffix(suffix string) FormatOption {
    return func(opts *FormatOptions) {
        opts.Suffix = suffix
    }
}

func WithUpper() FormatOption {
    return func(opts *FormatOptions) {
        opts.Upper = true
    }
}

formatter := compose.InvokableLambdaWithOption(
    func(ctx context.Context, input string, formatOpts ...FormatOption) (string, error) {
        opts := &FormatOptions{
            Prefix: "",
            Suffix: "",
            Upper:  false,
        }

        for _, opt := range formatOpts {
            opt(opts)
        }

        result := input
        if opts.Upper {
            result = strings.ToUpper(result)
        }

        return opts.Prefix + result + opts.Suffix, nil
    },
)

// 使用示例
result1, _ := formatter.Invoke(ctx, "hello", WithPrefix(">>> "), WithSuffix(" <<<"))
fmt.Println(result1) // >>> hello <<<

result2, _ := formatter.Invoke(ctx, "world", WithUpper(), WithPrefix("[INFO] "))
fmt.Println(result2) // [INFO] WORLD
```

### 练习3：多模式组合 Lambda

```go
// 任务：创建一个既支持同步又支持异步的文本处理 Lambda
textProcessor, err := compose.AnyLambda(
    // Invoke 模式：一次性处理
    func(ctx context.Context, input string, opts ...processOption) (string, error) {
        return processText(input), nil
    },
    // Stream 模式：逐词处理
    func(ctx context.Context, input string, opts ...processOption) (*schema.StreamReader[string], error) {
        words := strings.Split(input, " ")
        reader := schema.NewStreamReader[string]()

        go func() {
            defer reader.Close()
            for _, word := range words {
                processed := processText(word)
                reader.Send(ctx, processed)
                time.Sleep(50 * time.Millisecond)
            }
        }()

        return reader.Recv(), nil
    },
    // Collect 模式：合并多个文本
    func(ctx context.Context, texts *schema.StreamReader[string], opts ...processOption) (string, error) {
        var result strings.Builder
        for {
            text, err := texts.Recv()
            if err != nil {
                if errors.Is(err, io.EOF) {
                    break
                }
                return "", err
            }
            processed := processText(text)
            builder.WriteString(processed + " ")
        }
        return strings.TrimSpace(result.String()), nil
    },
)

func processText(text string) string {
    return strings.ToUpper(strings.TrimSpace(text))
}
```

---

## 🏗️ 编排集成

### 在 Chain 中使用

```go
// 构建一个文本处理流水线
chain := compose.NewChain[string, string]()

// 1. 文本清洗
chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
    return strings.TrimSpace(input), nil
}))

// 2. 语言检测
chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
    if containsChinese(input) {
        return "zh-CN", nil
    }
    return "en-US", nil
}))

// 3. 格式化输出
chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, lang string) (string, error) {
    return fmt.Sprintf("检测到语言: %s", lang), nil
}))

// 编译并运行
runner, _ := chain.Compile(ctx)
result, _ := runner.Invoke(ctx, "  你好世界  ")
fmt.Println(result) // 检测到语言: zh-CN
```

### 在 Graph 中使用

```go
// 构建一个复杂的数据处理图
graph := compose.NewGraph[string, ProcessedData]()

// 添加节点
graph.AddLambdaNode("parse_input", parseInputLambda)      // 解析输入
graph.AddLambdaNode("validate_data", validateDataLambda) // 验证数据
graph.AddLambdaNode("process_data", processDataLambda)   // 处理数据
graph.AddLambdaNode("format_output", formatOutputLambda) // 格式化输出

// 添加边（连接节点）
graph.AddEdge("parse_input", "validate_data")
graph.AddEdge("validate_data", "process_data")
graph.AddEdge("process_data", "format_output")

// 编译并运行
runner, _ := graph.Compile(ctx)
result, _ := runner.Invoke(ctx, "input data")
```

### 内置 Lambda 组件

#### ToList Lambda - 类型转换神器

```go
// 将单个消息转换为消息列表
msgToList := compose.ToList[*schema.Message]()

// 在Chain中的常见用法
chain := compose.NewChain[string, []*schema.Message]()
chain.AppendChatModel(chatModel)  // 返回 *schema.Message
chain.AppendLambda(msgToList)     // 转换为 []*schema.Message
chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, messages []*schema.Message) (int, error) {
    return len(messages), nil
}))
```

#### MessageParser Lambda - JSON解析利器

```go
// 定义要解析的结构体
type WeatherInfo struct {
    City        string  `json:"city"`
    Temperature float64 `json:"temperature"`
    Humidity    int     `json:"humidity"`
    Description string  `json:"description"`
}

// 创建解析器
weatherParser := schema.NewMessageJSONParser[*WeatherInfo](&schema.MessageJSONParseConfig{
    ParseFrom:    schema.MessageParseFromContent,
    ParseKeyPath: "", // 如果只需要解析子字段，可以用 "weather.data"
})

// 创建解析 Lambda
parseWeatherLambda := compose.MessageParser(weatherParser)

// 使用示例
chain := compose.NewChain[*schema.Message, *WeatherInfo]()
chain.AppendLambda(parseWeatherLambda)

// 运行
runner, _ := chain.Compile(ctx)
weather, _ := runner.Invoke(ctx, &schema.Message{
    Content: `{"city": "北京", "temperature": 25.5, "humidity": 60, "description": "晴天"}`,
})
fmt.Printf("城市: %s, 温度: %.1f°C\n", weather.City, weather.Temperature)
```

---

## 📚 进阶技巧

### 性能优化

#### 1. 避免不必要的内存分配

```go
// ❌ 不好的做法：每次都创建新的slice
badLambda := compose.InvokableLambda(func(ctx context.Context, items []string) ([]string, error) {
    result := make([]string, len(items)) // 新分配内存
    for i, item := range items {
        result[i] = strings.ToUpper(item)
    }
    return result, nil
})

// ✅ 好的做法：预分配内存或重用缓冲区
goodLambda := compose.InvokableLambda(func(ctx context.Context, items []string) ([]string, error) {
    result := make([]string, 0, len(items)) // 预分配容量
    for _, item := range items {
        result = append(result, strings.ToUpper(item))
    }
    return result, nil
})
```

#### 2. 并发处理

```go
// 并发处理多个任务
concurrentProcessor := compose.StreamableLambda(func(ctx context.Context, tasks []Task) (*schema.StreamReader[Result], error) {
    output := schema.NewStreamReader[Result]()
    sem := make(chan struct{}, 10) // 限制并发数

    go func() {
        defer output.Close()
        var wg sync.WaitGroup

        for _, task := range tasks {
            wg.Add(1)
            go func(t Task) {
                defer wg.Done()
                sem <- struct{}{}     // 获取信号量
                defer func() { <-sem }() // 释放信号量

                result := processTask(ctx, t)
                output.Send(ctx, result)
            }(task)
        }

        wg.Wait()
    }()

    return output.Recv(), nil
})
```

### 错误处理最佳实践

```go
// 自定义错误类型
type LambdaError struct {
    Code    int
    Message string
    Cause   error
}

func (e *LambdaError) Error() string {
    return fmt.Sprintf("Lambda错误 [%d]: %s (原因: %v)", e.Code, e.Message, e.Cause)
}

// 带重试机制的 Lambda
retryableLambda := compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
    const maxRetries = 3
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        result, err := callExternalAPI(ctx, input)
        if err == nil {
            return result, nil
        }

        lastErr = err
        if i < maxRetries-1 {
            // 指数退避
            backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
            select {
            case <-ctx.Done():
                return "", ctx.Err()
            case <-time.After(backoff):
                continue
            }
        }
    }

    return "", &LambdaError{
        Code:    500,
        Message: "API调用失败",
        Cause:   lastErr,
    }
})
```

### 配置和选项管理

```go
// 配置结构体
type LambdaConfig struct {
    Timeout     time.Duration
    MaxRetries  int
    EnableDebug bool
}

type LambdaOption func(*LambdaConfig)

func WithTimeout(timeout time.Duration) LambdaOption {
    return func(cfg *LambdaConfig) {
        cfg.Timeout = timeout
    }
}

func WithMaxRetries(retries int) LambdaOption {
    return func(cfg *LambdaConfig) {
        cfg.MaxRetries = retries
    }
}

func WithDebug() LambdaOption {
    return func(cfg *LambdaConfig) {
        cfg.EnableDebug = true
    }
}

// 工厂函数创建可配置的Lambda
func NewConfigurableLambda(opts ...LambdaOption) *compose.Lambda {
    config := &LambdaConfig{
        Timeout:    30 * time.Second,
        MaxRetries: 3,
        EnableDebug: false,
    }

    for _, opt := range opts {
        opt(config)
    }

    return compose.InvokableLambdaWithOption(
        func(ctx context.Context, input string, processOpts ...ProcessOption) (string, error) {
            // 应用配置
            if config.Timeout > 0 {
                var cancel context.CancelFunc
                ctx, cancel = context.WithTimeout(ctx, config.Timeout)
                defer cancel()
            }

            if config.EnableDebug {
                log.Printf("处理输入: %s", input)
            }

            result, err := processWithConfig(ctx, input, config, processOpts...)

            if config.EnableDebug {
                log.Printf("处理结果: %s, 错误: %v", result, err)
            }

            return result, err
        },
    )
}
```

---

## 🔧 API 参考

### 快速查找

| 功能 | 函数名 | 难度 | 常用度 |
|------|--------|------|--------|
| 基础Lambda | `InvokableLambda` | ⭐ | ⭐⭐⭐⭐⭐ |
| 流式输出 | `StreamableLambda` | ⭐⭐ | ⭐⭐⭐⭐ |
| 流式输入 | `CollectableLambda` | ⭐⭐ | ⭐⭐⭐ |
| 流式转换 | `TransformableLambda` | ⭐⭐⭐ | ⭐⭐ |
| 多模式 | `AnyLambda` | ⭐⭐⭐ | ⭐⭐⭐ |
| 带选项 | `InvokableLambdaWithOption` | ⭐⭐ | ⭐⭐⭐⭐ |
| 类型转换 | `ToList` | ⭐ | ⭐⭐⭐⭐ |
| JSON解析 | `MessageParser` | ⭐⭐ | ⭐⭐⭐⭐⭐ |

### 核心创建函数

#### 基础创建函数（8个）

```go
// 无选项版本 - 简单场景
InvokableLambda[I, O](i InvokeWOOpt[I, O], opts ...LambdaOpt) *Lambda
StreamableLambda[I, O](s StreamWOOpt[I, O], opts ...LambdaOpt) *Lambda
CollectableLambda[I, O](c CollectWOOpt[I, O], opts ...LambdaOpt) *Lambda
TransformableLambda[I, O](t TransformWOOpts[I, O], opts ...LambdaOpt) *Lambda

// 带选项版本 - 复杂场景
InvokableLambdaWithOption[I, O, TOption](i Invoke[I, O, TOption], opts ...LambdaOpt) *Lambda
StreamableLambdaWithOption[I, O, TOption](s Stream[I, O, TOption], opts ...LambdaOpt) *Lambda
CollectableLambdaWithOption[I, O, TOption](c Collect[I, O, TOption], opts ...LambdaOpt) *Lambda
TransformableLambdaWithOption[I, O, TOption](t Transform[I, O, TOption], opts ...LambdaOpt) *Lambda

// 多模式组合 - 最灵活
AnyLambda[I, O, TOption](i Invoke, s Stream, c Collect, t Transform, opts ...LambdaOpt) (*Lambda, error)
```

#### 配置选项（2个）

```go
// 启用Lambda函数的回调功能
WithLambdaCallbackEnable(y bool) LambdaOpt

// 设置Lambda函数的类型标识
WithLambdaType(t string) LambdaOpt
```

#### 内置组件（2个）

```go
// 单个输入转列表
ToList[I any](opts ...LambdaOpt) *Lambda

// 消息JSON解析器
MessageParser[T any](p schema.MessageParser[T], opts ...LambdaOpt) *Lambda
```

#### 编排集成方法（6个）

```go
// Graph相关
AddLambdaNode(key string, node *Lambda, opts ...GraphAddNodeOpt) error

// Chain相关
AppendLambda(node *Lambda, opts ...GraphAddNodeOpt) *Chain[I, O]

// Parallel相关
AddLambda(outputKey string, node *Lambda, opts ...GraphAddNodeOpt) *Parallel

// ChainBranch相关
AddLambda(key string, node *Lambda, opts ...GraphAddNodeOpt) *ChainBranch

// Workflow相关
AddLambdaNode(key string, lambda *Lambda, opts ...GraphAddNodeOpt) *WorkflowNode

// 调用时传递选项
WithLambdaOption(opts ...any) Option
```

#### 函数类型定义

```go
// 核心交互模式
type Invoke[I, O, TOption any] func(ctx context.Context, input I, opts ...TOption) (output O, err error)
type Stream[I, O, TOption any] func(ctx context.Context, input I, opts ...TOption) (output *schema.StreamReader[O], err error)
type Collect[I, O, TOption any] func(ctx context.Context, input *schema.StreamReader[I], opts ...TOption) (output O, err error)
type Transform[I, O, TOption any] func(ctx context.Context, input *schema.StreamReader[I], opts ...TOption) (output *schema.StreamReader[O], err error)

// 简化版本（无选项）
type InvokeWOOpt[I, O any] func(ctx context.Context, input I) (output O, err error)
type StreamWOOpt[I, O any] func(ctx context.Context, input I) (output *schema.StreamReader[O], err error)
type CollectWOOpt[I, O any] func(ctx context.Context, input *schema.StreamReader[I]) (output O, err error)
type TransformWOOpts[I, O any] func(ctx context.Context, input *schema.StreamReader[I]) (output *schema.StreamReader[O], err error)
```

---

## 🎯 学习检查点

### 基础达标 ✅
- [ ] 能够创建简单的 InvokableLambda 进行数据转换
- [ ] 理解四种交互模式的区别和适用场景
- [ ] 掌握 ToList 和 MessageParser 两个内置组件的使用
- [ ] 能够在 Chain 和 Graph 中正确集成 Lambda

### 进阶达标 🚀
- [ ] 能够使用 AnyLambda 创建多模式 Lambda
- [ ] 掌握自定义选项的设计和使用
- [ ] 理解流式处理的实现方式
- [ ] 能够优化 Lambda 的性能和错误处理

### 实战项目 🏆
- [ ] 创建一个文本预处理管道（清洗→验证→转换）
- [ ] 实现一个带重试机制的 API 调用 Lambda
- [ ] 构建一个实时数据处理流（过滤→转换→聚合）
- [ ] 开发一个配置化的 Lambda 工厂

---

## 🔗 相关资源

### 📚 学习资源
- [官方文档](https://www.cloudwego.io/zh/docs/eino/core_modules/components/lambda_guide/)
- [示例代码](https://github.com/cloudwego/eino-examples/blob/main/components/lambda)
- [源码位置：`eino/compose/types_lambda.go`](https://github.com/cloudwego/eino/blob/main/compose/types_lambda.go)

### 🎯 前置知识
- [`quickstart/chat`](../quickstart/chat/) - 基础聊天应用
- [`components/tool`](../tool/) - 工具集成基础

### 🚀 后续学习
- [`compose/chain`](../../compose/chain/) - 链式编排
- [`compose/graph`](../../compose/graph/) - 图形编排
- [`adk/helloworld`](../../adk/helloworld/) - Agent开发套件

### 💡 最佳实践
- **错误处理**：始终考虑错误情况，提供有意义的错误信息
- **性能优化**：避免不必要的内存分配，合理使用并发
- **可测试性**：将 Lambda 函数设计为纯函数，便于单元测试
- **配置管理**：使用选项模式管理 Lambda 的配置参数

---

## 🎓 总结

Lambda 组件是 Eino 框架中最灵活、最强大的组件之一。通过掌握 Lambda，你可以：

- 🔄 **无缝集成**：将任何 Go 函数集成到 AI 编排流程中
- ⚡ **性能优化**：通过流式处理提升用户体验
- 🛡️ **类型安全**：享受 Go 语言的类型检查优势
- 🏗️ **架构灵活**：构建任意复杂的数据处理管道

**学习建议**：
1. 从简单的 InvokableLambda 开始练习
2. 逐步掌握流式处理和多模式组合
3. 在实际项目中应用最佳实践
4. 关注性能优化和错误处理

**记住**：Lambda 是连接 AI 能力和业务逻辑的桥梁，掌握它就能构建出真正强大的 AI 应用！