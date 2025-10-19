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

package einoagent

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildEinoAgent(ctx context.Context) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		InputToQuery   = "InputToQuery"
		ChatTemplate   = "ChatTemplate"
		ReactAgent     = "ReactAgent"
		RedisRetriever = "RedisRetriever"
		InputToHistory = "InputToHistory"
	)
	g := compose.NewGraph[*UserMessage, *schema.Message]()

	// 输入起点：并行构建查询与上下文变量
	_ = g.AddLambdaNode(InputToQuery, compose.InvokableLambdaWithOption(userMessageToQueryLambda), compose.WithNodeName("UserMessageToQuery"))
	_ = g.AddEdge(compose.START, InputToQuery)

	// 上下文变量节点：提取历史消息并与起点连接
	_ = g.AddLambdaNode(InputToHistory, compose.InvokableLambdaWithOption(userMessageToVariablesLambda), compose.WithNodeName("UserMessageToVariables"))
	_ = g.AddEdge(compose.START, InputToHistory)

	// 检索链路：查询 → 检索器 → 模板
	redisRetrieverKeyOfRetriever, err := buildRedisRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(RedisRetriever, redisRetrieverKeyOfRetriever, compose.WithOutputKey("documents"))
	_ = g.AddEdge(InputToQuery, RedisRetriever)

	chatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate, chatTemplateKeyOfChatTemplate)
	_ = g.AddEdge(RedisRetriever, ChatTemplate)
	_ = g.AddEdge(InputToHistory, ChatTemplate)

	// 推理节点：ReAct 代理
	reactAgentKeyOfLambda, err := reactAgentLambda(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(ReactAgent, reactAgentKeyOfLambda, compose.WithNodeName("ReAct Agent"))
	_ = g.AddEdge(ChatTemplate, ReactAgent)
	_ = g.AddEdge(ReactAgent, compose.END)

	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgent"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
