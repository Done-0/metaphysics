// Package conversation 提供对话相关模型定义
// 创建者：Done-0
// 创建时间：2025-07-03
package conversation

import (
	"time"

	"gorm.io/gorm"
)

// Conversation 对话记录
type Conversation struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`     // 主键ID
	UserID      int64          `gorm:"index:idx_user_id" json:"user_id"`       // 用户ID
	Title       string         `gorm:"size:255" json:"title"`                  // 对话标题
	SessionID   string         `gorm:"size:64;uniqueIndex" json:"session_id"`  // 会话ID
	FirstPrompt string         `gorm:"type:text" json:"first_prompt"`          // 首次提示语
	GmtCreate   time.Time      `gorm:"autoCreateTime" json:"gmt_create"`       // 创建时间
	GmtModified time.Time      `gorm:"autoUpdateTime" json:"gmt_modified"`     // 修改时间
	Deleted     bool           `gorm:"default:false" json:"deleted"`           // 是否删除
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"deleted_at"` // 删除时间
}

// TableName 表名
func (Conversation) TableName() string {
	return "t_conversation"
}

// Message 消息记录
type Message struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`               // 主键ID
	ConversationID int64          `gorm:"index:idx_conversation_id" json:"conversation_id"` // 对话ID
	UserID         int64          `gorm:"index:idx_user_id" json:"user_id"`                 // 用户ID
	SessionID      string         `gorm:"size:64;index:idx_session_id" json:"session_id"`   // 会话ID
	Role           string         `gorm:"size:20" json:"role"`                              // 角色（USER/ASSISTANT）
	Content        string         `gorm:"type:text" json:"content"`                         // 消息内容
	RequestID      int            `gorm:"index:idx_request_id" json:"request_id"`           // 请求消息ID
	ResponseID     int            `gorm:"index:idx_response_id" json:"response_id"`         // 响应消息ID
	ParentID       int            `json:"parent_id"`                                        // 父消息ID
	TokenUsage     int            `json:"token_usage"`                                      // Token使用量
	GmtCreate      time.Time      `gorm:"autoCreateTime" json:"gmt_create"`                 // 创建时间
	GmtModified    time.Time      `gorm:"autoUpdateTime" json:"gmt_modified"`               // 修改时间
	Deleted        bool           `gorm:"default:false" json:"deleted"`                     // 是否删除
	DeletedAt      gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"deleted_at"`           // 删除时间
}

// TableName 表名
func (Message) TableName() string {
	return "t_message"
}

// MessageCounter 消息计数器
type MessageCounter struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`     // 主键ID
	UserID      int64     `gorm:"uniqueIndex:idx_user_id" json:"user_id"` // 用户ID
	NextID      int       `json:"next_id"`                                // 下一个消息ID
	GmtCreate   time.Time `gorm:"autoCreateTime" json:"gmt_create"`       // 创建时间
	GmtModified time.Time `gorm:"autoUpdateTime" json:"gmt_modified"`     // 修改时间
}

// TableName 表名
func (MessageCounter) TableName() string {
	return "t_message_counter"
}
