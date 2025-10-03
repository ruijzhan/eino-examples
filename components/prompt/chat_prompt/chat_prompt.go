/*
 * Copyright 2024 CloudWeGo Authors
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
 */

package main

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"

	"github.com/cloudwego/eino-examples/internal/logs"
)

func main() {
	// Demonstrate FString format type (Python f-string style)
	systemTpl := `你是情绪助手，你的任务是根据用户的输入，生成一段赞美的话，语句优美，韵律强。
用户姓名：{user_name}
用户年龄：{user_age}
用户性别：{user_gender}
用户喜好：{user_hobby}`

	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)

	msgList, err := chatTpl.Format(context.Background(), map[string]any{
		"user_name":   "张三",
		"user_age":    "18",
		"user_gender": "男",
		"user_hobby":  "打篮球、打游戏",
		"message_histories": []*schema.Message{
			schema.UserMessage("我喜欢打羽毛球"),
			schema.AssistantMessage("羽毛球是一项很好的运动，能够锻炼身体的协调性和反应能力。", nil),
		},
		"user_query": "请为我赋诗一首",
	})
	if err != nil {
		logs.Errorf("Format failed, err=%v", err)
		return
	}

	logs.Infof("=== FString Format Example ===")
	for i, msg := range msgList {
		logs.Infof("Message %d: [%s] %s", i+1, msg.Role, msg.Content)
	}

	// Demonstrate GoTemplate format type
	goSystemTpl := `你是专业助手，帮助用户{{.user_query}}。
用户信息：姓名={{.user_name}}，年龄={{.user_age}}，爱好={{.user_hobby}}`

	goChatTpl := prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage(goSystemTpl),
		schema.MessagesPlaceholder("history", false),
		schema.UserMessage("当前请求：{{.current_request}}"),
	)

	goMsgList, err := goChatTpl.Format(context.Background(), map[string]any{
		"user_name":       "李四",
		"user_age":        "25",
		"user_hobby":      "编程、阅读",
		"user_query":      "学习新技术",
		"current_request": "请推荐一些学习 Go 语言的资源",
		"history": []*schema.Message{
			schema.UserMessage("我想学习编程"),
			schema.AssistantMessage("很好！编程是一项很有价值的技能。你想从哪门语言开始？", nil),
		},
	})
	if err != nil {
		logs.Errorf("GoTemplate format failed, err=%v", err)
		return
	}

	logs.Infof("\n=== GoTemplate Format Example ===")
	for i, msg := range goMsgList {
		logs.Infof("Message %d: [%s] %s", i+1, msg.Role, msg.Content)
	}

	// Demonstrate optional placeholder
	optionalTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个有帮助的助手"),
		schema.MessagesPlaceholder("optional_history", true), // optional history messages
		schema.UserMessage("{user_input}"),
	)

	optionalMsgList, err := optionalTpl.Format(context.Background(), map[string]any{
		"user_input": "你好，请介绍一下你自己",
		// Note: optional_history is not provided here, but since it's optional, no error will occur
	})
	if err != nil {
		logs.Errorf("Optional placeholder format failed, err=%v", err)
		return
	}

	logs.Infof("\n=== Optional Placeholder Example ===")
	for i, msg := range optionalMsgList {
		logs.Infof("Message %d: [%s] %s", i+1, msg.Role, msg.Content)
	}

	demonstratePromptAPIs(context.Background())
}

func demonstratePromptAPIs(ctx context.Context) {
	baseTemplate := prompt.FromMessages(schema.FString,
		schema.SystemMessage("嗨 {user_name}，让我们开始学习 {topic} 编程吧！"),
		schema.UserMessage("我想了解更多关于 {topic} 的内容"),
	)

	var tpl prompt.ChatTemplate = baseTemplate
	formattedMessages, err := tpl.Format(ctx, map[string]any{
		"user_name": "王五",
		"topic":     "Go",
	})
	if err != nil {
		logs.Errorf("ChatTemplate example failed, err=%v", err)
		return
	}

	logs.Infof("\n=== ChatTemplate / DefaultChatTemplate / FromMessages Example ===")
	for i, msg := range formattedMessages {
		logs.Infof("Example Message %d: [%s] %s", i+1, msg.Role, msg.Content)
	}

	customTpl := newCustomChatTemplate(baseTemplate)
	uppercaseOpt := prompt.WrapImplSpecificOptFn(func(opt *uppercaseOption) {
		opt.enable = true
	})
	uppercaseMessages, err := customTpl.Format(ctx, map[string]any{
		"user_name": "王五",
		"topic":     "Go",
	}, uppercaseOpt)
	if err != nil {
		logs.Errorf("Option example failed, err=%v", err)
		return
	}

	logs.Infof("\n=== Option / WrapImplSpecificOptFn / GetImplSpecificOptions Example ===")
	for i, msg := range uppercaseMessages {
		logs.Infof("Option Message %d: [%s] %s", i+1, msg.Role, msg.Content)
	}

	manualInput := &prompt.CallbackInput{
		Variables: map[string]any{
			"user_name": "王五",
			"topic":     "Go",
		},
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage("你现在正在学习 {topic}"),
		},
	}
	logs.Infof("\n=== CallbackInput Example ===")
	logs.Infof("Manual CallbackInput variables: %v", manualInput.Variables)
	logs.Infof("Manual CallbackInput templates count: %d", len(manualInput.Templates))

	convInput := prompt.ConvCallbackInput(map[string]any{
		"converted": true,
		"topic":     "Go",
	})
	logs.Infof("Converted CallbackInput variables: %v", convInput.Variables)

	manualOutput := &prompt.CallbackOutput{
		Result: uppercaseMessages,
	}
	logs.Infof("\n=== CallbackOutput Example ===")
	logs.Infof("Manual CallbackOutput result count: %d", len(manualOutput.Result))

	convOutput := prompt.ConvCallbackOutput(formattedMessages)
	if convOutput != nil {
		logs.Infof("Converted CallbackOutput result count: %d", len(convOutput.Result))
	}
}

type uppercaseOption struct {
	enable bool
}

type customChatTemplate struct {
	base prompt.ChatTemplate
}

func newCustomChatTemplate(base prompt.ChatTemplate) *customChatTemplate {
	return &customChatTemplate{base: base}
}

func (c *customChatTemplate) Format(ctx context.Context, vs map[string]any, opts ...prompt.Option) ([]*schema.Message, error) {
	optsCfg := prompt.GetImplSpecificOptions(&uppercaseOption{}, opts...)
	msgs, err := c.base.Format(ctx, vs)
	if err != nil {
		return nil, err
	}

	if !optsCfg.enable {
		return msgs, nil
	}

	formatted := make([]*schema.Message, len(msgs))
	for i, msg := range msgs {
		if msg == nil {
			continue
		}

		copyMsg := *msg
		copyMsg.Content = strings.ToUpper(copyMsg.Content)
		formatted[i] = &copyMsg
	}

	return formatted, nil
}
