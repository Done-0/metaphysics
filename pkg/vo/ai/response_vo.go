// Package ai 提供 AI 服务相关的视图对象
// 创建者：Done-0
// 创建时间：2024-07-09
package ai

// DeepseekResponse Deepseek 响应结构
type DeepseekResponse struct {
	ReasoningContent string `json:"reasoning_content"` // 思维链内容
	Content          string `json:"content"`           // 最终回答内容
}

// StreamChunk 流式响应数据块
type StreamChunk struct {
	ReasoningContent string `json:"reasoning_content,omitempty"` // 思维链内容片段
	Content          string `json:"content,omitempty"`           // 回答内容片段
	Done             bool   `json:"done"`                        // 是否完成
}
