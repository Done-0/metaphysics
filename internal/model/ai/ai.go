// Package ai 提供AI相关模型定义
// 创建者：Done-0
// 创建时间：2025-07-03
package ai

import (
	"github.com/Done-0/metaphysics/internal/model/base"
)

// Conversation 对话模型，存储用户与AI的对话内容
type Conversation struct {
	base.Base

	UserID      int64  `json:"user_id" gorm:"index"`            // 用户ID
	Title       string `json:"title" gorm:"size:255"`           // 对话标题
	SessionID   string `json:"session_id" gorm:"size:64;index"` // 会话ID
	FirstPrompt string `json:"first_prompt" gorm:"type:text"`   // 首次提问内容
}

// TableName 指定表名
// 返回值：
//   - string: 表名
func (Conversation) TableName() string {
	return "ai_conversations"
}

// Message 消息模型，存储对话中的单条消息
type Message struct {
	base.Base

	ConversationID int64  `json:"conversation_id" gorm:"index"`                                      // 对话ID
	UserID         int64  `json:"user_id" gorm:"index"`                                              // 用户ID
	SessionID      string `json:"session_id" gorm:"size:64;index"`                                   // 会话ID
	Role           string `json:"role" gorm:"size:20;check:role IN ('USER', 'ASSISTANT', 'SYSTEM')"` // 角色（USER/ASSISTANT/SYSTEM）
	Content        string `json:"content" gorm:"type:text"`                                          // 消息内容
	RequestID      int    `json:"request_id"`                                                        // 请求消息ID
	ResponseID     int    `json:"response_id"`                                                       // 响应消息ID
	ParentID       int    `json:"parent_id"`                                                         // 父消息ID
	TokenUsage     int    `json:"token_usage"`                                                       // Token使用量
}

// TableName 指定表名
// 返回值：
//   - string: 表名
func (Message) TableName() string {
	return "ai_messages"
}

// MessageCounter 消息计数器，存储用户的消息ID计数
type MessageCounter struct {
	base.Base

	UserID int64 `json:"user_id" gorm:"uniqueIndex"` // 用户ID
	NextID int   `json:"next_id"`                    // 下一个消息ID
}

// TableName 指定表名
// 返回值：
//   - string: 表名
func (MessageCounter) TableName() string {
	return "ai_message_counters"
}
