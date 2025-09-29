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

package adk

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/internal/safe"
	"github.com/cloudwego/eino/schema"
)

type HistoryEntry struct {
	IsUserInput bool
	AgentName   string
	Message     Message
}

type HistoryRewriter func(ctx context.Context, entries []*HistoryEntry) ([]Message, error)

type flowAgent struct {
	Agent

	subAgents   []*flowAgent
	parentAgent *flowAgent

	disallowTransferToParent bool
	historyRewriter          HistoryRewriter

	checkPointStore compose.CheckPointStore
}

func (a *flowAgent) deepCopy() *flowAgent {
	ret := &flowAgent{
		Agent:                    a.Agent,
		subAgents:                make([]*flowAgent, 0, len(a.subAgents)),
		parentAgent:              a.parentAgent,
		disallowTransferToParent: a.disallowTransferToParent,
		historyRewriter:          a.historyRewriter,
		checkPointStore:          a.checkPointStore,
	}

	for _, sa := range a.subAgents {
		ret.subAgents = append(ret.subAgents, sa.deepCopy())
	}
	return ret
}

func SetSubAgents(ctx context.Context, agent Agent, subAgents []Agent) (Agent, error) {
	return setSubAgents(ctx, agent, subAgents)
}

type AgentOption func(options *flowAgent)

func WithDisallowTransferToParent() AgentOption {
	return func(fa *flowAgent) {
		fa.disallowTransferToParent = true
	}
}

func WithHistoryRewriter(h HistoryRewriter) AgentOption {
	return func(fa *flowAgent) {
		fa.historyRewriter = h
	}
}

func toFlowAgent(ctx context.Context, agent Agent, opts ...AgentOption) *flowAgent {
	var fa *flowAgent
	var ok bool
	if fa, ok = agent.(*flowAgent); !ok {
		fa = &flowAgent{Agent: agent}
	} else {
		fa = fa.deepCopy()
	}
	for _, opt := range opts {
		opt(fa)
	}

	if fa.historyRewriter == nil {
		fa.historyRewriter = buildDefaultHistoryRewriter(agent.Name(ctx))
	}

	return fa
}

func AgentWithOptions(ctx context.Context, agent Agent, opts ...AgentOption) Agent {
	return toFlowAgent(ctx, agent, opts...)
}

func setSubAgents(ctx context.Context, agent Agent, subAgents []Agent) (*flowAgent, error) {
	fa := toFlowAgent(ctx, agent)

	if len(fa.subAgents) > 0 {
		return nil, errors.New("agent's sub-agents has already been set")
	}

	if onAgent, ok_ := fa.Agent.(OnSubAgents); ok_ {
		err := onAgent.OnSetSubAgents(ctx, subAgents)
		if err != nil {
			return nil, err
		}
	}

	for _, s := range subAgents {
		fsa := toFlowAgent(ctx, s)

		if fsa.parentAgent != nil {
			return nil, errors.New("agent has already been set as a sub-agent of another agent")
		}

		fsa.parentAgent = fa
		if onAgent, ok__ := fsa.Agent.(OnSubAgents); ok__ {
			err := onAgent.OnSetAsSubAgent(ctx, agent)
			if err != nil {
				return nil, err
			}

			if fsa.disallowTransferToParent {
				err = onAgent.OnDisallowTransferToParent(ctx)
				if err != nil {
					return nil, err
				}
			}
		}

		fa.subAgents = append(fa.subAgents, fsa)
	}

	return fa, nil
}

func (a *flowAgent) getAgent(ctx context.Context, name string) *flowAgent {
	for _, subAgent := range a.subAgents {
		if subAgent.Name(ctx) == name {
			return subAgent
		}
	}

	if a.parentAgent != nil && a.parentAgent.Name(ctx) == name {
		return a.parentAgent
	}

	return nil
}

func belongToRunPath(eventRunPath []RunStep, runPath []RunStep) bool {
	if len(runPath) < len(eventRunPath) {
		return false
	}

	for i, step := range eventRunPath {
		if !runPath[i].Equals(step) {
			return false
		}
	}

	return true
}

func rewriteMessage(msg Message, agentName string) Message {
	var sb strings.Builder
	sb.WriteString("For context:")
	if msg.Role == schema.Assistant {
		if msg.Content != "" {
			sb.WriteString(fmt.Sprintf(" [%s] said: %s.", agentName, msg.Content))
		}
		if len(msg.ToolCalls) > 0 {
			for i := range msg.ToolCalls {
				f := msg.ToolCalls[i].Function
				sb.WriteString(fmt.Sprintf(" [%s] called tool: `%s` with arguments: %s.",
					agentName, f.Name, f.Arguments))
			}
		}
	} else if msg.Role == schema.Tool && msg.Content != "" {
		sb.WriteString(fmt.Sprintf(" [%s] `%s` tool returned result: %s.",
			agentName, msg.ToolName, msg.Content))
	}

	return schema.UserMessage(sb.String())
}

func genMsg(entry *HistoryEntry, agentName string) (Message, error) {
	msg := entry.Message
	if entry.AgentName != agentName {
		msg = rewriteMessage(msg, entry.AgentName)
	}

	return msg, nil
}

func (ai *AgentInput) deepCopy() *AgentInput {
	copied := &AgentInput{
		Messages:        make([]Message, len(ai.Messages)),
		EnableStreaming: ai.EnableStreaming,
	}

	copy(copied.Messages, ai.Messages)

	return copied
}

func (a *flowAgent) genAgentInput(ctx context.Context, runCtx *runContext, skipTransferMessages bool) (*AgentInput, error) {
	input := runCtx.RootInput.deepCopy()
	runPath := runCtx.RunPath

	events := runCtx.Session.getEvents()
	historyEntries := make([]*HistoryEntry, 0)

	for _, m := range input.Messages {
		historyEntries = append(historyEntries, &HistoryEntry{
			IsUserInput: true,
			Message:     m,
		})
	}

	for _, event := range events {
		if !belongToRunPath(event.RunPath, runPath) {
			continue
		}

		if skipTransferMessages && event.Action != nil && event.Action.TransferToAgent != nil {
			// If skipTransferMessages is true and the event contain transfer action, the message in this event won't be appended to history entries.
			if event.Output != nil &&
				event.Output.MessageOutput != nil &&
				event.Output.MessageOutput.Role == schema.Tool &&
				len(historyEntries) > 0 {
				// If the skipped message's role is Tool, remove the previous history entry as it's also a transfer message(from ChatModelAgent and GenTransferMessages).
				historyEntries = historyEntries[:len(historyEntries)-1]
			}
			continue
		}

		msg, err := getMessageFromWrappedEvent(event)
		if err != nil {
			return nil, err
		}

		if msg == nil {
			continue
		}

		historyEntries = append(historyEntries, &HistoryEntry{
			AgentName: event.AgentName,
			Message:   msg,
		})
	}

	messages, err := a.historyRewriter(ctx, historyEntries)
	if err != nil {
		return nil, err
	}
	input.Messages = messages

	return input, nil
}

func buildDefaultHistoryRewriter(agentName string) HistoryRewriter {
	return func(ctx context.Context, entries []*HistoryEntry) ([]Message, error) {
		messages := make([]Message, 0, len(entries))
		var err error
		for _, entry := range entries {
			msg := entry.Message
			if !entry.IsUserInput {
				msg, err = genMsg(entry, agentName)
				if err != nil {
					return nil, fmt.Errorf("gen agent input failed: %w", err)
				}
			}

			if msg != nil {
				messages = append(messages, msg)
			}
		}

		return messages, nil
	}
}

func (a *flowAgent) Run(ctx context.Context, input *AgentInput, opts ...AgentRunOption) *AsyncIterator[*AgentEvent] {
	agentName := a.Name(ctx)

	ctx, runCtx := initRunCtx(ctx, agentName, input)

	o := getCommonOptions(nil, opts...)

	input, err := a.genAgentInput(ctx, runCtx, o.skipTransferMessages)
	if err != nil {
		return genErrorIter(err)
	}

	if wf, ok := a.Agent.(*workflowAgent); ok {
		return wf.Run(ctx, input, opts...)
	}

	aIter := a.Agent.Run(ctx, input, filterOptions(agentName, opts)...)

	iterator, generator := NewAsyncIteratorPair[*AgentEvent]()

	go a.run(ctx, runCtx, aIter, generator, opts...)

	return iterator
}

func (a *flowAgent) Resume(ctx context.Context, info *ResumeInfo, opts ...AgentRunOption) *AsyncIterator[*AgentEvent] {
	runCtx := getRunCtx(ctx)
	if runCtx == nil {
		return genErrorIter(fmt.Errorf("failed to resume agent: run context is empty"))
	}

	agentName := a.Name(ctx)
	targetName := agentName
	if len(runCtx.RunPath) > 0 {
		targetName = runCtx.RunPath[len(runCtx.RunPath)-1].agentName
	}

	if agentName != targetName {
		// go to target flow agent
		targetAgent := recursiveGetAgent(ctx, a, targetName)
		if targetAgent == nil {
			return genErrorIter(fmt.Errorf("failed to resume agent: cannot find agent: %s", agentName))
		}
		return targetAgent.Resume(ctx, info, opts...)
	}

	if wf, ok := a.Agent.(*workflowAgent); ok {
		return wf.Resume(ctx, info, opts...)
	}

	// resume current agent
	ra, ok := a.Agent.(ResumableAgent)
	if !ok {
		return genErrorIter(fmt.Errorf("failed to resume agent: target agent[%s] isn't resumable", agentName))
	}
	iterator, generator := NewAsyncIteratorPair[*AgentEvent]()
	aIter := ra.Resume(ctx, info, opts...)

	go a.run(ctx, runCtx, aIter, generator, opts...)

	return iterator
}

type DeterministicTransferConfig struct {
	Agent        Agent
	ToAgentNames []string
}

func (a *flowAgent) run(
	ctx context.Context,
	runCtx *runContext,
	aIter *AsyncIterator[*AgentEvent],
	generator *AsyncGenerator[*AgentEvent],
	opts ...AgentRunOption) {
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			e := safe.NewPanicErr(panicErr, debug.Stack())
			generator.Send(&AgentEvent{Err: e})
		}

		generator.Close()
	}()

	var lastAction *AgentAction
	for {
		event, ok := aIter.Next()
		if !ok {
			break
		}

		event.AgentName = a.Name(ctx)
		event.RunPath = runCtx.RunPath
		// copy the event so that the copied event's stream is exclusive for any potential consumer
		// copy before adding to session because once added to session it's stream could be consumed by genAgentInput at any time
		copied := copyAgentEvent(event)
		setAutomaticClose(copied)
		setAutomaticClose(event)
		runCtx.Session.addEvent(copied)
		lastAction = event.Action
		generator.Send(event)
	}

	var destName string
	if lastAction != nil {
		if lastAction.Interrupted != nil {
			appendInterruptRunCtx(ctx, runCtx)
			return
		}
		if lastAction.Exit {
			return
		}

		if lastAction.TransferToAgent != nil {
			destName = lastAction.TransferToAgent.DestAgentName
		}
	}

	// handle transferring to another agent
	if destName != "" {
		agentToRun := a.getAgent(ctx, destName)
		if agentToRun == nil {
			e := errors.New(fmt.Sprintf(
				"transfer failed: agent '%s' not found when transferring from '%s'",
				destName, a.Name(ctx)))
			generator.Send(&AgentEvent{Err: e})
			return
		}

		subAIter := agentToRun.Run(ctx, nil /*subagents get input from runCtx*/, opts...)
		for {
			subEvent, ok_ := subAIter.Next()
			if !ok_ {
				break
			}

			setAutomaticClose(subEvent)
			generator.Send(subEvent)
		}
	}
}

func recursiveGetAgent(ctx context.Context, agent *flowAgent, agentName string) *flowAgent {
	if agent == nil {
		return nil
	}
	if agent.Name(ctx) == agentName {
		return agent
	}
	a := agent.getAgent(ctx, agentName)
	if a != nil {
		return a
	}
	for _, sa := range agent.subAgents {
		a = recursiveGetAgent(ctx, sa, agentName)
		if a != nil {
			return a
		}
	}
	return nil
}
