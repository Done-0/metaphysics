// Package provider 实现 ollama AI 服务提供者
// 创建者：Done-0
// 创建时间：2024-06-10
package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"

	"github.com/Done-0/metaphysics/configs"
	types "github.com/Done-0/metaphysics/internal/ai/model"
	"github.com/Done-0/metaphysics/internal/ai/prompt"
	"github.com/Done-0/metaphysics/pkg/vo/ai"
)

// ollamaProvider ollama 服务提供者
type ollamaProvider struct {
	config *configs.Config // 配置信息
	llm    *ollama.LLM     // ollama LLM 实例
}

// NewOllamaProvider ollama 服务提供者构造器
// 参数：
//
//	cfg: 配置信息
//
// 返回值：
//
//	types.Service: ollama Provider 实例
//	error: 错误信息
func NewOllamaProvider(cfg *configs.Config) (types.Service, error) {
	return &ollamaProvider{config: cfg}, nil
}

// AnalyzeBazi 分析八字
// 参数：
//
//	ctx: 上下文
//	name: 姓名
//	gender: 性别
//	baziInfo: 八字信息
//
// 返回值：
//
//	string: 分析结果
//	error: 错误信息
func (p *ollamaProvider) AnalyzeBazi(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string) (string, error) {
	llm, err := p.llmInstance()
	if err != nil {
		return "", fmt.Errorf("获取 ollama LLM 实例失败: %w", err)
	}

	promptText := prompt.BuildBaziPrompt(name, gender, birthTime, baziInfo)
	data := map[string]any{"prompt": promptText}
	tpl := prompts.NewPromptTemplate("{{.prompt}}", []string{"prompt"})
	chain := chains.NewLLMChain(llm, tpl)
	result, err := chain.Call(ctx, data)
	if err != nil {
		return "", fmt.Errorf("AI 分析失败: %w", err)
	}
	text, ok := result["text"].(string)
	if !ok {
		return "", fmt.Errorf("AI 结果解析失败")
	}
	return text, nil
}

// AnalyzeBaziWithReasoning 分析八字（带推理过程）
// 参数：
//
//	ctx: 上下文
//	name: 姓名
//	gender: 性别
//	baziInfo: 八字信息
//
// 返回值：
//
//	*ai.DeepseekResponse: 分析结果（包含推理过程）
//	error: 错误信息
func (p *ollamaProvider) AnalyzeBaziWithReasoning(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string) (*ai.DeepseekResponse, error) {
	content, err := p.AnalyzeBazi(ctx, name, gender, birthTime, baziInfo)
	if err != nil {
		return nil, fmt.Errorf("AI 分析失败: %w", err)
	}
	return &ai.DeepseekResponse{Content: content}, nil
}

// StreamAnalyzeBazi 流式分析八字
// 参数：
//
//	ctx: 上下文
//	name: 姓名
//	gender: 性别
//	baziInfo: 八字信息
//	handler: 流式响应处理函数
//
// 返回值：
//
//	error: 错误信息
func (p *ollamaProvider) StreamAnalyzeBazi(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string, handler types.StreamHandler) error {
	llm, err := p.llmInstance()
	if err != nil {
		return fmt.Errorf("获取 ollama LLM 实例失败: %w", err)
	}

	promptText := prompt.BuildBaziPrompt(name, gender, birthTime, baziInfo)
	_, err = llms.GenerateFromSinglePrompt(ctx, llm, promptText, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		return handler(&ai.StreamChunk{Content: string(chunk)})
	}))
	if err != nil {
		return fmt.Errorf("流式分析失败: %w", err)
	}
	return handler(&ai.StreamChunk{Done: true})
}

// DetermineProvider 确定要使用的 AI 提供商
// 返回值：
//
//	types.Provider: AI 服务提供商
func (p *ollamaProvider) DetermineProvider() types.Provider {
	return types.PROVIDER_OLLAMA
}

// llmInstance 获取或初始化 ollama LLM 实例
// 返回值：
//
//	*ollama.LLM: ollama LLM 实例
//	error: 错误信息
func (p *ollamaProvider) llmInstance() (*ollama.LLM, error) {
	if p.llm != nil {
		return p.llm, nil
	}
	llm, err := ollama.New(
		ollama.WithServerURL(p.config.AIConfig.OllamaAPIBase),
		ollama.WithModel(p.config.AIConfig.OllamaModel),
	)
	if err != nil {
		return nil, fmt.Errorf("初始化 ollama LLM 失败: %w", err)
	}
	p.llm = llm
	return llm, nil
}
