// Package dto 提供对话相关的数据传输对象
// 创建者：Done-0
// 创建时间：2025-07-03
package dto

// ContinueConversationRequest 继续对话请求参数
type ContinueConversationRequest struct {
	Prompt string `json:"prompt" form:"prompt" query:"prompt" binding:"required" validate:"required"` // 用户提示内容
}

// StreamContinueConversationRequest 流式继续对话请求参数
type StreamContinueConversationRequest struct {
	Prompt string `json:"prompt" form:"prompt" query:"prompt" binding:"required" validate:"required"` // 用户提示内容
}
