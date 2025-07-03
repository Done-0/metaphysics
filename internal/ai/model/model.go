// Package types 提供 AI 服务相关类型定义
// 创建者：Done-0
// 创建时间：2024-06-10
package types

import (
	"context"
	"time"

	"github.com/Done-0/metaphysics/pkg/vo/ai"
)

// Provider AI 服务提供商类型
type Provider string

const (
	PROVIDER_OLLAMA   Provider = "ollama"   // ollama 本地模型服务
	PROVIDER_DEEPSEEK Provider = "deepseek" // deepseek 推理模型服务
)

// StreamHandler 流式响应处理器
type StreamHandler func(chunk *ai.StreamChunk) error

// Service AI 服务接口，所有 Provider 必须实现
type Service interface {
	// AnalyzeBazi 分析八字
	// 参数：
	//   ctx: 上下文
	//   name: 姓名
	//   gender: 性别
	//   baziInfo: 八字信息
	// 返回值：
	//   string: 分析结果
	//   error: 错误信息
	AnalyzeBazi(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string) (string, error)

	// StreamAnalyzeBazi 流式分析八字
	// 参数：
	//   ctx: 上下文
	//   name: 姓名
	//   gender: 性别
	//   baziInfo: 八字信息
	//   handler: 流式响应处理函数
	// 返回值：
	//   error: 错误信息
	StreamAnalyzeBazi(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string, handler StreamHandler) error

	// DetermineProvider 确定使用的 AI 提供商
	// 返回值：
	//   Provider: AI 服务提供商
	DetermineProvider() Provider
}
