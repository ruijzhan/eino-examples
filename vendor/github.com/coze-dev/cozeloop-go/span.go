// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package cozeloop

import (
	"context"
	"time"

	"github.com/coze-dev/cozeloop-go/entity"
	"github.com/coze-dev/cozeloop-go/spec/tracespec"
)

// Span is the interface for span.
type Span interface {
	SpanContext
	commonSpanSetter

	// SetTags sets business custom tags.
	SetTags(ctx context.Context, tagKVs map[string]interface{})

	// SetBaggage sets tags and also passes these tags to other downstream spans (assuming
	//the user uses ToHeader and FromHeader to handle header passing between services).
	SetBaggage(ctx context.Context, baggageItems map[string]string)

	// Finish The span will be reported only after an explicit call to Finish.
	// Under the hood, it is actually placed in an asynchronous queue waiting to be reported.
	Finish(ctx context.Context)

	// GetStartTime returns the start time of the Span.
	GetStartTime() time.Time

	// ToHeader Convert the span to headers. Used for cross-process correlation.
	ToHeader() (map[string]string, error)
}

// Set system-defined fields
type commonSpanSetter interface {
	// SetInput key: `input`
	// Input information. The input will be serialized into a JSON string.
	// You can find recommended specification in https://github.com/coze-dev/cozeloop-go/tree/main/spec/tracespec
	// Or you can use any struct you like.
	SetInput(ctx context.Context, input interface{})

	// SetOutput key: `output`
	// Output information. The output will be serialized into a JSON string.
	// You can find recommended specification in https://github.com/coze-dev/cozeloop-go/tree/main/spec/tracespec
	// Or you can use any struct you like.
	SetOutput(ctx context.Context, output interface{})

	// SetError key: `error`
	// Set error message.
	SetError(ctx context.Context, err error)

	// SetStatusCode key: `status_code`
	// Set status code. A non-zero code is considered an exception.
	SetStatusCode(ctx context.Context, code int)

	// SetUserID key: `user_id`
	// Set user id.
	SetUserID(ctx context.Context, userID string)
	SetUserIDBaggage(ctx context.Context, userID string)

	// SetMessageID key: `message_id`
	// Set message id.
	SetMessageID(ctx context.Context, messageID string)
	SetMessageIDBaggage(ctx context.Context, messageID string)

	// SetThreadID key: `thread_id`
	// Set thread id.
	// thread_id is used to correlate multiple requests and is passed in by the business.
	SetThreadID(ctx context.Context, threadID string)
	SetThreadIDBaggage(ctx context.Context, threadID string)

	// SetPrompt key: `prompt
	// Associated with PromptKey and PromptVersion, it will write two tags: prompt_key and prompt_version.
	// SetPrompt is used to set the PromptKey and PromptVersion to tag.
	SetPrompt(ctx context.Context, prompt entity.Prompt)

	// SetModelProvider key: `model_provider`
	// The provider of the LLM, such as OpenAI, etc.
	SetModelProvider(ctx context.Context, modelProvider string)

	// SetModelName key: `model_name`
	// The Name of the LLM model, such as: gpt-4-1106-preview.
	SetModelName(ctx context.Context, modelName string)

	// SetModelCallOptions key: `call_options`
	// The call options of the LLM, such as temperature, max_tokens, etc.
	// The recommended standard format is CallOption of spec package
	SetModelCallOptions(ctx context.Context, callOptions interface{})

	// SetInputTokens key: `input_tokens`
	// The usage of input tokens. When the value of input_tokens is set,
	// It will be automatically summed with output_tokens to calculate the tokens tag.
	SetInputTokens(ctx context.Context, inputTokens int)

	// SetOutputTokens key: `output_tokens`
	// The usage of output tokens. When the value of output_tokens is set,
	// It will be automatically summed with input_tokens to calculate the tokens tag.
	SetOutputTokens(ctx context.Context, outputTokens int)

	// SetStartTimeFirstResp key: `start_time_first_resp`
	// Timestamp of the first packet return from LLM, unit: microseconds.
	// When `start_time_first_resp` is set, a tag named `latency_first_resp` calculated
	// based on the span's StartTime will be added, meaning the latency for the first packet.
	SetStartTimeFirstResp(ctx context.Context, startTimeFirstResp int64)

	// SetRuntime key: `runtime`
	// The runtime of the LLM, such as language, library, scene, etc.
	// The recommended standard format is Runtime of spec package
	SetRuntime(ctx context.Context, runtime tracespec.Runtime)

	// SetServiceName
	// set the custom service name, identify different services.
	SetServiceName(ctx context.Context, serviceName string)

	// SetLogID
	// set the custom log id, identify different query.
	SetLogID(ctx context.Context, logID string)

	// SetFinishTime
	// Default is time.Now() when span Finish(). DO NOT set unless you do not use default time.
	SetFinishTime(finishTime time.Time)

	// SetSystemTags
	// set the system tags. DO NOT set unless you know what you are doing.
	SetSystemTags(ctx context.Context, systemTags map[string]interface{})

	// SetDeploymentEnv
	// set the deployment env, identify custom env.
	SetDeploymentEnv(ctx context.Context, deploymentEnv string)
}

// SpanContext is the interface for span Baggage transfer.
type SpanContext interface {
	GetSpanID() string
	GetTraceID() string
	GetBaggage() map[string]string
}
