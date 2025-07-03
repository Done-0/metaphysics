// Package provider 实现 Deepseek AI 服务提供者
// 创建者：Done-0
// 创建时间：2024-06-10
package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Done-0/metaphysics/configs"
	types "github.com/Done-0/metaphysics/internal/ai/model"
	"github.com/Done-0/metaphysics/internal/ai/prompt"
	"github.com/Done-0/metaphysics/pkg/vo/ai"
)

const (
	DEEPSEEK_TIMEOUT    = 120 * time.Second // Deepseek 请求超时时间
	DEEPSEEK_MAX_TOKENS = 4096              // Deepseek 最大 token 数
)

// deepseekProvider Deepseek 服务提供者
type deepseekProvider struct {
	config *configs.Config // 配置信息
	client *http.Client    // HTTP 客户端
}

// chatRequest Deepseek 聊天请求结构
type chatRequest struct {
	Model     string        `json:"model"`                // 模型名称
	Messages  []chatMessage `json:"messages"`             // 消息内容
	MaxTokens int           `json:"max_tokens,omitempty"` // 最大 token 数
	Stream    bool          `json:"stream,omitempty"`     // 是否流式
}

// chatMessage Deepseek 聊天消息结构
type chatMessage struct {
	Role    string `json:"role"`    // 角色
	Content string `json:"content"` // 内容
}

// chatResponse Deepseek 聊天响应结构
type chatResponse struct {
	Choices []struct {
		Message struct {
			ReasoningContent string `json:"reasoning_content"` // 思维链内容
			Content          string `json:"content"`           // 回答内容
		} `json:"message"`
		FinishReason string `json:"finish_reason"` // 结束原因
	} `json:"choices"`
}

// streamChunk Deepseek 流式响应数据块
type streamChunk struct {
	Choices []struct {
		Delta struct {
			ReasoningContent string `json:"reasoning_content,omitempty"` // 思维链内容片段
			Content          string `json:"content,omitempty"`           // 回答内容片段
		} `json:"delta"`
		FinishReason string `json:"finish_reason"` // 结束原因
	} `json:"choices"`
}

// NewDeepseekProvider Deepseek 服务提供者构造器
// 参数：
//
//	cfg: 配置信息
//
// 返回值：
//
//	types.Service: Deepseek Provider 实例
//	error: 错误信息
func NewDeepseekProvider(cfg *configs.Config) (types.Service, error) {
	return &deepseekProvider{
		config: cfg,
		client: &http.Client{Timeout: DEEPSEEK_TIMEOUT},
	}, nil
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
func (p *deepseekProvider) AnalyzeBazi(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string) (string, error) {
	resp, err := p.analyze(ctx, name, gender, birthTime, baziInfo)
	if err != nil {
		return "", fmt.Errorf("AI 分析失败: %w", err)
	}
	return resp.Content, nil
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
func (p *deepseekProvider) StreamAnalyzeBazi(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string, handler types.StreamHandler) error {
	promptText := prompt.BuildBaziPrompt(name, gender, birthTime, baziInfo)

	req := chatRequest{
		Model:     p.config.AIConfig.DeepseekModel,
		Messages:  []chatMessage{{Role: "user", Content: promptText}},
		MaxTokens: DEEPSEEK_MAX_TOKENS,
		Stream:    true,
	}
	return p.stream(ctx, req, handler)
}

// DetermineProvider 确定要使用的 AI 提供商
// 返回值：
//
//	types.Provider: AI 服务提供商
func (p *deepseekProvider) DetermineProvider() types.Provider {
	return types.PROVIDER_DEEPSEEK
}

// GetProviderName 获取提供商名称
// 返回值：
//
//	string: 提供商名称
func (p *deepseekProvider) GetProviderName() string {
	return string(types.PROVIDER_DEEPSEEK)
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
//	*ai.DeepseekResponse: 分析结果
//	error: 错误信息
func (p *deepseekProvider) AnalyzeBaziWithReasoning(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string) (*ai.DeepseekResponse, error) {
	return p.analyze(ctx, name, gender, birthTime, baziInfo)
}

// analyze 内部八字分析方法
// 参数：
//
//	ctx: 上下文
//	name: 姓名
//	gender: 性别
//	baziInfo: 八字信息
//
// 返回值：
//
//	*ai.DeepseekResponse: 分析结果
//	error: 错误信息
func (p *deepseekProvider) analyze(ctx context.Context, name, gender string, birthTime time.Time, baziInfo map[string]string) (*ai.DeepseekResponse, error) {
	promptText := prompt.BuildBaziPrompt(name, gender, birthTime, baziInfo)

	req := chatRequest{
		Model:     p.config.AIConfig.DeepseekModel,
		Messages:  []chatMessage{{Role: "user", Content: promptText}},
		MaxTokens: DEEPSEEK_MAX_TOKENS,
	}
	data, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.config.AIConfig.DeepseekAPIBase+"/v1/chat/completions", bytes.NewBuffer(data))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.AIConfig.DeepseekAPIKey)
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("deepseek API 请求失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("deepseek API 错误: %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var cr chatResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return nil, fmt.Errorf("AI 结果解析失败: %w", err)
	}
	if len(cr.Choices) == 0 {
		return nil, fmt.Errorf("无 AI 结果")
	}
	return &ai.DeepseekResponse{
		ReasoningContent: cr.Choices[0].Message.ReasoningContent,
		Content:          cr.Choices[0].Message.Content,
	}, nil
}

// stream 内部流式处理方法
// 参数：
//
//	ctx: 上下文
//	req: 请求体
//	handler: 处理函数
//
// 返回值：
//
//	error: 错误信息
func (p *deepseekProvider) stream(ctx context.Context, req chatRequest, handler types.StreamHandler) error {
	data, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.config.AIConfig.DeepseekAPIBase+"/v1/chat/completions", bytes.NewBuffer(data))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.AIConfig.DeepseekAPIKey)
	httpReq.Header.Set("Accept", "text/event-stream")
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("deepseek API 请求失败: %w", err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		if !bytes.HasPrefix(line, []byte("data: ")) {
			continue
		}
		data := bytes.TrimPrefix(line, []byte("data: "))
		if string(data) == "[DONE]" {
			handler(&ai.StreamChunk{Done: true})
			break
		}
		var chunk streamChunk
		if err := json.Unmarshal(data, &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		c := chunk.Choices[0]
		_ = handler(&ai.StreamChunk{
			ReasoningContent: c.Delta.ReasoningContent,
			Content:          c.Delta.Content,
			Done:             c.FinishReason == "stop",
		})
		if c.FinishReason == "stop" {
			break
		}
	}
	return nil
}
