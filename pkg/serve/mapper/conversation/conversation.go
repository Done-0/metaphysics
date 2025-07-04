// Package conversation 提供对话相关的数据访问接口
// 创建者：Done-0
// 创建时间：2025-07-03
package conversation

import (
	"github.com/gin-gonic/gin"

	conversationModel "github.com/Done-0/metaphysics/internal/model/conversation"
)

// ConversationMapper 对话数据访问接口
type ConversationMapper interface {
	// SaveConversation 保存对话
	// 参数：
	//   - ctx: 上下文信息
	//   - conversation: 对话模型
	//
	// 返回值：
	//   - error: 错误信息
	SaveConversation(ctx *gin.Context, conversation *conversationModel.Conversation) error

	// GetConversationByID 根据ID获取对话
	// 参数：
	//   - ctx: 上下文信息
	//   - id: 对话ID
	//
	// 返回值：
	//   - *conversationModel.Conversation: 对话模型
	//   - error: 错误信息
	GetConversationByID(ctx *gin.Context, id int64) (*conversationModel.Conversation, error)

	// GetConversationsByUserID 获取用户的所有对话
	// 参数：
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//   - pageNo: 页码
	//   - pageSize: 每页数量
	//
	// 返回值：
	//   - []*conversationModel.Conversation: 对话列表
	//   - int64: 总记录数
	//   - error: 错误信息
	GetConversationsByUserID(ctx *gin.Context, userID int64, pageNo, pageSize int) ([]*conversationModel.Conversation, int64, error)

	// UpdateConversation 更新对话
	// 参数：
	//   - ctx: 上下文信息
	//   - conversation: 对话模型
	//
	// 返回值：
	//   - error: 错误信息
	UpdateConversation(ctx *gin.Context, conversation *conversationModel.Conversation) error

	// DeleteConversation 删除对话
	// 参数：
	//   - ctx: 上下文信息
	//   - id: 对话ID
	//
	// 返回值：
	//   - error: 错误信息
	DeleteConversation(ctx *gin.Context, id int64) error

	// SaveMessage 保存消息
	// 参数：
	//   - ctx: 上下文信息
	//   - message: 消息模型
	//
	// 返回值：
	//   - error: 错误信息
	SaveMessage(ctx *gin.Context, message *conversationModel.Message) error

	// GetMessagesByConversationID 获取对话的所有消息
	// 参数：
	//   - ctx: 上下文信息
	//   - conversationID: 对话ID
	//
	// 返回值：
	//   - []*conversationModel.Message: 消息列表
	//   - error: 错误信息
	GetMessagesByConversationID(ctx *gin.Context, conversationID int64) ([]*conversationModel.Message, error)

	// GetLatestConversationHistory 获取最近的对话历史
	// 参数：
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//
	// 返回值：
	//   - string: 会话ID
	//   - string: 对话历史
	//   - error: 错误信息
	GetLatestConversationHistory(ctx *gin.Context, userID int64) (string, string, error)

	// SaveConversationHistory 保存对话历史
	// 参数：
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//   - sessionID: 会话ID
	//   - history: 对话历史
	//
	// 返回值：
	//   - error: 错误信息
	SaveConversationHistory(ctx *gin.Context, userID int64, sessionID string, history string) error

	// GetNextMessageIDs 获取下一个消息ID
	// 参数：
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//
	// 返回值：
	//   - int: 请求消息ID
	//   - int: 响应消息ID
	//   - error: 错误信息
	GetNextMessageIDs(ctx *gin.Context, userID int64) (int, int, error)
}
