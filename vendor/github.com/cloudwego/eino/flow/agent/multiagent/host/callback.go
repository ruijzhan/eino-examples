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
 * limitations under the License.
 */

package host

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/internal/safe"
	"github.com/cloudwego/eino/schema"
	template "github.com/cloudwego/eino/utils/callbacks"
)

// MultiAgentCallback is the callback interface for host multi-agent.
type MultiAgentCallback interface {
	OnHandOff(ctx context.Context, info *HandOffInfo) context.Context
}

// HandOffInfo is the info which will be passed to MultiAgentCallback.OnHandOff, representing a hand off event.
type HandOffInfo struct {
	ToAgentName string
	Argument    string
}

// ConvertCallbackHandlers converts []host.MultiAgentCallback to callbacks.Handler.
func ConvertCallbackHandlers(handlers ...MultiAgentCallback) callbacks.Handler {
	onChatModelEnd := func(ctx context.Context, info *callbacks.RunInfo, output *model.CallbackOutput) context.Context {
		msg := output.Message
		if msg == nil || msg.Role != schema.Assistant || len(msg.ToolCalls) == 0 {
			return ctx
		}

		for _, cb := range handlers {
			for _, toolCall := range msg.ToolCalls {
				ctx = cb.OnHandOff(ctx, &HandOffInfo{
					ToAgentName: toolCall.Function.Name,
					Argument:    toolCall.Function.Arguments,
				})
			}
		}

		return ctx
	}

	onChatModelEndWithStreamOutput := func(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[*model.CallbackOutput]) context.Context {
		go func() {
			defer func() {
				panicInfo := recover()
				if panicInfo != nil {
					fmt.Println(safe.NewPanicErr(panicInfo, debug.Stack()))
				}
				output.Close()
			}()

			handOffs := make(map[string]string)
			var handOffOrder []string
			for {
				oneOutput, err := output.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					return
				}

				for _, toolCall := range oneOutput.Message.ToolCalls {
					if len(toolCall.Function.Name) > 0 {
						if existing, ok := handOffs[toolCall.Function.Name]; !ok {
							handOffOrder = append(handOffOrder, toolCall.Function.Name)
							handOffs[toolCall.Function.Name] = toolCall.Function.Arguments
						} else {
							handOffs[toolCall.Function.Name] = existing + toolCall.Function.Arguments
						}
					}
				}
			}

			for _, cb := range handlers {
				for _, name := range handOffOrder {
					args := handOffs[name]
					_ = cb.OnHandOff(ctx, &HandOffInfo{
						ToAgentName: name,
						Argument:    args,
					})
				}
			}
		}()

		return ctx
	}

	return template.NewHandlerHelper().ChatModel(&template.ModelCallbackHandler{
		OnEnd:                 onChatModelEnd,
		OnEndWithStreamOutput: onChatModelEndWithStreamOutput,
	}).Handler()
}

// convertCallbacks reads graph call options, extract host.MultiAgentCallback and convert it to callbacks.Handler.
func convertCallbacks(opts ...agent.AgentOption) callbacks.Handler {
	agentOptions := agent.GetImplSpecificOptions(&options{}, opts...)
	if len(agentOptions.agentCallbacks) == 0 {
		return nil
	}

	handlers := agentOptions.agentCallbacks
	return ConvertCallbackHandlers(handlers...)
}
