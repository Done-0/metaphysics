// Package ai 提供 AI 服务能力
// 创建者：Done-0
// 创建时间：2024-06-10
package ai

import (
	"fmt"
	"sync"

	"github.com/Done-0/metaphysics/configs"
	types "github.com/Done-0/metaphysics/internal/ai/model"
	"github.com/Done-0/metaphysics/internal/ai/provider"
)

var (
	instance types.Service
	once     sync.Once
)

// New 返回 AI 服务实例
// 返回值：
//
//	types.Service: AI 服务接口
func New() types.Service {
	once.Do(func() {
		cfg, err := configs.GetConfig()
		if err != nil {
			panic(fmt.Errorf("配置加载失败: %w", err))
		}

		// 按优先级选择 Provider
		switch {
		case cfg.AIConfig.DeepseekEnabled:
			provider, err := provider.NewDeepseekProvider(cfg)
			if err == nil {
				instance = provider
				return
			}
			fmt.Printf("[AI] deepseek 初始化失败: %v，降级使用 ollama\n", err)
			fallthrough
		case cfg.AIConfig.OllamaEnabled:
			provider, err := provider.NewOllamaProvider(cfg)
			if err == nil {
				instance = provider
				return
			}
			panic(fmt.Errorf("ollama 初始化失败: %w", err))
		default:
			panic("未启用任何 AI Provider")
		}
	})
	return instance
}
