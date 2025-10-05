package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// basicLambdaExample 演示基础的 InvokableLambda 使用
func basicLambdaExample() {
	fmt.Println("=== 基础 Lambda 示例 ===")

	// Create a chain to use the lambda properly
	chain := compose.NewChain[string, string]()

	// Add lambda to the chain
	chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return "处理结果: " + strings.ToUpper(input), nil
	}))

	// Compile the chain
	runner, err := chain.Compile(context.Background())
	if err != nil {
		panic(err)
	}

	// Invoke the chain (which contains our lambda)
	result, err := runner.Invoke(context.Background(), "hello lambda")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// streamableLambdaExample 演示 StreamableLambda 的使用
func streamableLambdaExample() {
	fmt.Println("\n=== StreamableLambda 示例 ===")

	// Create a StreamableLambda that splits text into words
	wordStreamer := compose.StreamableLambda(func(ctx context.Context, text string) (*schema.StreamReader[string], error) {
		words := strings.Split(text, " ")
		sr, sw := schema.Pipe[string](len(words))

		go func() {
			defer sw.Close()
			for _, word := range words {
				if ctx.Err() != nil {
					return
				}
				sw.Send(word, nil)
				time.Sleep(100 * time.Millisecond) // 模拟处理延迟
			}
		}()

		return sr, nil
	})

	streamChain := compose.NewChain[string, string]()
	streamChain.AppendLambda(wordStreamer)

	streamRunner, err := streamChain.Compile(context.Background())
	if err != nil {
		panic(err)
	}

	stream, err := streamRunner.Stream(context.Background(), "Go 语言 是 一个 有趣 的 案例")
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	fmt.Println("逐词流式输出:")
	for {
		chunk, chunkErr := stream.Recv()
		if errors.Is(chunkErr, io.EOF) {
			break
		}
		if chunkErr != nil {
			panic(chunkErr)
		}
		fmt.Printf("stream chunk: %s\n", chunk)
	}
}

// collectableLambdaExample 演示 CollectableLambda 的使用
func collectableLambdaExample() {
	fmt.Println("\n=== CollectableLambda 示例 ===")

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

	// 创建测试数据流
	sr, sw := schema.Pipe[int](5)
	go func() {
		defer sw.Close()
		for i := 1; i <= 5; i++ {
			sw.Send(i, nil)
		}
	}()

	collectChain := compose.NewChain[int, int]()
	collectChain.AppendLambda(sumCollector)

	collectRunner, err := collectChain.Compile(context.Background())
	if err != nil {
		panic(err)
	}

	result, err := collectRunner.Collect(context.Background(), sr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("数字流求和结果: %d (1+2+3+4+5 = 15)\n", result)
}

// transformableLambdaExample 演示 TransformableLambda 的使用
func transformableLambdaExample() {
	fmt.Println("\n=== TransformableLambda 示例 ===")

	// 过滤偶数
	evenFilter := compose.TransformableLambda(func(ctx context.Context, numbers *schema.StreamReader[int]) (*schema.StreamReader[int], error) {
		sr, sw := schema.Pipe[int](0) // 动态管道，容量未知

		go func() {
			defer sw.Close()
			for {
				num, err := numbers.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					return
				}
				if num%2 == 0 {
					sw.Send(num, nil)
				}
			}
		}()

		return sr, nil
	})

	// 创建测试数据流 (1-10)
	inputSr, inputSw := schema.Pipe[int](10)
	go func() {
		defer inputSw.Close()
		for i := 1; i <= 10; i++ {
			inputSw.Send(i, nil)
		}
	}()

	transformChain := compose.NewChain[int, int]()
	transformChain.AppendLambda(evenFilter)

	transformRunner, err := transformChain.Compile(context.Background())
	if err != nil {
		panic(err)
	}

	outputSr, err := transformRunner.Transform(context.Background(), inputSr)
	if err != nil {
		panic(err)
	}
	defer outputSr.Close()

	fmt.Println("过滤偶数结果:")
	var evenNumbers []int
	for {
		num, err := outputSr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		evenNumbers = append(evenNumbers, num)
		fmt.Printf("even number: %d\n", num)
	}
	fmt.Printf("偶数列表: %v\n", evenNumbers)
}

// concurrentProcessorExample 演示并发处理
func concurrentProcessorExample() {
	fmt.Println("\n=== 并发处理示例 ===")

	// 模拟任务和结果
	type Task struct {
		ID   int
		Data string
	}
	type Result struct {
		TaskID  int
		Outcome string
	}

	// 模拟任务处理函数
	processTask := func(ctx context.Context, task Task) Result {
		time.Sleep(50 * time.Millisecond) // 模拟处理时间
		return Result{
			TaskID:  task.ID,
			Outcome: fmt.Sprintf("处理完成: %s", task.Data),
		}
	}

	// 并发处理多个任务
	concurrentProcessor := compose.StreamableLambda(func(ctx context.Context, tasks []Task) (*schema.StreamReader[Result], error) {
		sr, sw := schema.Pipe[Result](len(tasks))
		sem := make(chan struct{}, 3) // 限制并发数为3

		go func() {
			defer sw.Close()
			var wg sync.WaitGroup

			for _, task := range tasks {
				wg.Add(1)
				go func(t Task) {
					defer wg.Done()
					sem <- struct{}{}        // 获取信号量
					defer func() { <-sem }() // 释放信号量

					result := processTask(ctx, t)
					sw.Send(result, nil)
				}(task)
			}

			wg.Wait()
		}()

		return sr, nil
	})

	// 创建测试任务
	tasks := []Task{
		{ID: 1, Data: "任务1"},
		{ID: 2, Data: "任务2"},
		{ID: 3, Data: "任务3"},
		{ID: 4, Data: "任务4"},
		{ID: 5, Data: "任务5"},
	}

	streamChain := compose.NewChain[[]Task, Result]()
	streamChain.AppendLambda(concurrentProcessor)

	streamRunner, err := streamChain.Compile(context.Background())
	if err != nil {
		panic(err)
	}

	outputSr, err := streamRunner.Stream(context.Background(), tasks)
	if err != nil {
		panic(err)
	}
	defer outputSr.Close()

	fmt.Println("并发处理结果:")
	var results []Result
	for {
		result, err := outputSr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		results = append(results, result)
		fmt.Printf("任务 %d: %s\n", result.TaskID, result.Outcome)
	}
	fmt.Printf("总共处理了 %d 个任务\n", len(results))
}

func main() {
	fmt.Println("Lambda 组件示例演示")
	fmt.Println("==================")

	// 运行所有示例
	basicLambdaExample()
	streamableLambdaExample()
	collectableLambdaExample()
	transformableLambdaExample()
	concurrentProcessorExample()

	fmt.Println("\n所有示例演示完成!")
}
