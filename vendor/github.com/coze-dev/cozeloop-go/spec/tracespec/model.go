// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package tracespec

import "encoding/json"

// ModelInput is the input for model span, for tag key: input
type ModelInput struct {
	Messages        []*ModelMessage  `json:"messages,omitempty"`
	Tools           []*ModelTool     `json:"tools,omitempty"`
	ModelToolChoice *ModelToolChoice `json:"tool_choice,omitempty"`
}

// ModelOutput is the output for model span, for tag key: output
type ModelOutput struct {
	Choices []*ModelChoice `json:"choices"`
}

// ModelCallOption is the option for model span, for tag key: call_options
type ModelCallOption struct {
	Temperature      float32  `json:"temperature"`
	MaxTokens        int64    `json:"max_tokens,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	TopP             float32  `json:"top_p,omitempty"`
	N                int64    `json:"n,omitempty"`
	TopK             *int64   `json:"top_k,omitempty"`
	PresencePenalty  *float32 `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float32 `json:"frequency_penalty,omitempty"`
	ReasoningEffort  string   `json:"reasoning_effort,omitempty"`
}

type ModelMessage struct {
	Role             string              `json:"role"`                        // from enum VRole in span_value
	Content          string              `json:"content,omitempty"`           // single content
	ReasoningContent string              `json:"reasoning_content,omitempty"` // only for output
	Parts            []*ModelMessagePart `json:"parts,omitempty"`             // multi-modality content
	Name             string              `json:"name,omitempty"`
	ToolCalls        []*ModelToolCall    `json:"tool_calls,omitempty"`
	ToolCallID       string              `json:"tool_call_id,omitempty"`
	Metadata         map[string]string   `json:"metadata,omitempty"`
}

type ModelMessagePart struct {
	Type     ModelMessagePartType `json:"type"` // Required. The type of the content.
	Text     string               `json:"text,omitempty"`
	ImageURL *ModelImageURL       `json:"image_url,omitempty"`
	FileURL  *ModelFileURL        `json:"file_url,omitempty"`
}

type ModelMessagePartType string

var (
	ModelMessagePartTypeText  ModelMessagePartType = "text"
	ModelMessagePartTypeImage ModelMessagePartType = "image_url"
	ModelMessagePartTypeFile  ModelMessagePartType = "file_url"
)

type ModelImageURL struct {
	Name string `json:"name,omitempty"`
	// Required. You can enter a valid image URL or MDN Base64 data of image.
	// MDN: https://developer.mozilla.org/en-US/docs/Web/URI/Reference/Schemes/data#syntax
	URL    string `json:"url,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type ModelFileURL struct {
	Name string `json:"name,omitempty"`
	// Required. You can enter a valid file URL or MDN Base64 data of file.
	// MDN: https://developer.mozilla.org/en-US/docs/Web/URI/Reference/Schemes/data#syntax
	URL    string `json:"url,omitempty"`
	Detail string `json:"detail,omitempty"`
	Suffix string `json:"suffix,omitempty"`
}

type ModelToolCall struct {
	ID       string                 `json:"id,omitempty"`
	Type     string                 `json:"type,omitempty"` // Always be: "function"
	Function *ModelToolCallFunction `json:"function"`
}

type ModelToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments,omitempty"`
}

type ModelTool struct {
	Type     string             `json:"type"` // Always be: "function"
	Function *ModelToolFunction `json:"function"`
}

type ModelToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type ModelChoice struct {
	FinishReason string        `json:"finish_reason"`
	Index        int64         `json:"index"`
	Message      *ModelMessage `json:"message"`
}

type ModelToolChoice struct {
	Type     string                 `json:"type"`               // from enum VToolChoice in span_value
	Function *ModelToolCallFunction `json:"function,omitempty"` // field name only.
}
