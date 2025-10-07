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
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
)

func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	// TODO Modify component configuration here.
	requiredEnvVars := []string{
		"OPENAI_EMBEDDING_MODEL",
		"OPENAI_EMBEDDING_API_KEY",
		"OPENAI_EMBEDDING_BASE_URL",
	}

	for _, key := range requiredEnvVars {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("environment variable %s is not set", key)
		}
	}

	config := &openai.EmbeddingConfig{
		Model:   os.Getenv("OPENAI_EMBEDDING_MODEL"),
		APIKey:  os.Getenv("OPENAI_EMBEDDING_API_KEY"),
		BaseURL: os.Getenv("OPENAI_EMBEDDING_BASE_URL"),
	}
	eb, err = openai.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
