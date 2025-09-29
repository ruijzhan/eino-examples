// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package cozeloop

import (
	"context"

	"github.com/coze-dev/cozeloop-go/entity"
	"github.com/coze-dev/cozeloop-go/internal/prompt"
)

// PromptClient interface of prompt client.
type PromptClient interface {
	// GetPrompt get prompt by prompt key and version.
	// if version is not set,  the latest version of the corresponding prompt will be obtained.
	GetPrompt(ctx context.Context, param GetPromptParam, options ...GetPromptOption) (*entity.Prompt, error)
	// PromptFormat format prompt with variables
	PromptFormat(ctx context.Context, prompt *entity.Prompt, variables map[string]any, options ...PromptFormatOption) (messages []*entity.Message, err error)
}

type GetPromptParam = prompt.GetPromptParam

type GetPromptOption func(option *prompt.GetPromptOptions)

type PromptFormatOption func(option *prompt.PromptFormatOptions)
