// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package prompt

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bluele/gcache"

	"github.com/coze-dev/cozeloop-go/entity"
	"github.com/coze-dev/cozeloop-go/internal/util"
)

const (
	defaultCacheSize = 100
	cacheKeyPrefix   = "prompt_hub"
	updateInterval   = time.Minute
)

type PromptCache struct {
	workspaceID string
	cache       gcache.Cache
	openAPI     *OpenAPIClient
	once        sync.Once
	stopChan    chan struct{}
	option      CacheOption
}

type CacheOption struct {
	EnableAsyncUpdate bool          // Whether to enable asynchronous updates
	UpdateInterval    time.Duration // Update interval, if 0, use default value
	MaxCacheSize      int
}

type Option func(*CacheOption)

// withAsyncUpdate set whether to enable asynchronous updates
func withAsyncUpdate(enable bool) Option {
	return func(opt *CacheOption) {
		opt.EnableAsyncUpdate = enable
	}
}

// withUpdateInterval set update interval
func withUpdateInterval(interval time.Duration) Option {
	return func(opt *CacheOption) {
		if interval > 0 {
			opt.UpdateInterval = interval
		}
	}
}

// withMaxCacheSize set max cache size
func withMaxCacheSize(size int) Option {
	return func(opt *CacheOption) {
		if size > 0 {
			opt.MaxCacheSize = size
		}
	}
}

func newPromptCache(workspaceID string, openAPI *OpenAPIClient, opts ...Option) *PromptCache {
	// Default configuration
	option := &CacheOption{
		EnableAsyncUpdate: false,
		UpdateInterval:    updateInterval,
		MaxCacheSize:      defaultCacheSize,
	}

	// Apply custom configurations
	for _, opt := range opts {
		opt(option)
	}

	cache := &PromptCache{
		workspaceID: workspaceID,
		cache:       gcache.New(option.MaxCacheSize).LFU().Build(),
		openAPI:     openAPI,
		stopChan:    make(chan struct{}),
		option:      *option,
	}

	// If asynchronous updates are enabled, start the update task
	if option.EnableAsyncUpdate {
		cache.Start()
	}

	return cache
}

func (c *PromptCache) Start() {
	c.once.Do(func() {
		util.GoSafe(context.Background(), c.startAsyncUpdate)
	})
}

func (c *PromptCache) Stop() {
	close(c.stopChan)
}

func (c *PromptCache) startAsyncUpdate() {
	ticker := time.NewTicker(c.option.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.updateAllPrompts()
		case <-c.stopChan:
			return
		}
	}
}

func (c *PromptCache) updateAllPrompts() {
	ctx := context.Background()
	queries := c.GetAllPromptQueries()

	if len(queries) == 0 {
		return
	}

	// Batch update
	promptResults, err := c.openAPI.MPullPrompt(ctx, MPullPromptRequest{
		WorkSpaceID: c.workspaceID,
		Queries:     queries,
	})
	if err != nil {
		return
	}

	// Update cache
	for _, p := range promptResults {
		if p != nil {
			c.Set(p.Query.PromptKey, p.Query.Version, p.Query.Label, toModelPrompt(p.Prompt))
		}
	}
}

func (c *PromptCache) getCacheKey(promptKey, version, label string) string {
	return fmt.Sprintf("%s:%s:%s:%s", cacheKeyPrefix, promptKey, version, label)
}

func (c *PromptCache) Get(promptKey, version, label string) (*entity.Prompt, bool) {
	key := c.getCacheKey(promptKey, version, label)
	if value, err := c.cache.Get(key); err == nil {
		if prompt, ok := value.(*entity.Prompt); ok {
			return prompt, true
		}
	}
	return nil, false
}

func (c *PromptCache) Set(promptKey, version, label string, prompt *entity.Prompt) {
	if prompt == nil {
		return
	}
	key := c.getCacheKey(promptKey, version, label)
	c.cache.Set(key, prompt)
}

// GetAllPromptQueries gets all cached Prompt query conditions
func (c *PromptCache) GetAllPromptQueries() []PromptQuery {
	queries := make([]PromptQuery, 0)
	keys := c.cache.Keys(false)

	for _, key := range keys {
		if strKey, ok := key.(string); ok {
			promptKey, version, label, ok := parseCacheKey(strKey)
			if ok {
				queries = append(queries, PromptQuery{
					PromptKey: promptKey,
					Version:   version,
					Label:     label,
				})
			}
		}
	}
	return queries
}

func parseCacheKey(key string) (promptKey string, version string, label string, ok bool) {
	parts := strings.Split(key, ":")
	if len(parts) == 4 {
		return parts[1], parts[2], parts[3], true
	}
	return "", "", "", false
}
