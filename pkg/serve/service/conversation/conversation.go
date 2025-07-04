// Package conversation 提供对话相关的服务接口
// 创建者：Done-0
// 创建时间：2025-07-03
package conversation

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/pkg/serve/controller/conversation/dto"
	"github.com/Done-0/metaphysics/pkg/vo/conversation"
)

// ConversationService 对话服务接口
type ConversationService interface {
	// AnalyzeBaziByUserID 根据用户ID分析八字
	// 参数：
	//   - ctx: 上下文信息
	//
	// 返回值：
	//   - *conversation.BaziAnalysisResponse: 八字分析结果
	//   - error: 错误信息
	AnalyzeBaziByUserID(ctx *gin.Context) (*conversation.BaziAnalysisResponse, error)

	// StreamAnalyzeBaziByUserID 流式分析用户八字
	// 参数：
	//   - ctx: 上下文信息
	//   - handler: 流式响应处理函数
	//
	// 返回值：
	//   - error: 错误信息
	StreamAnalyzeBaziByUserID(ctx *gin.Context, handler func(content string, done bool) error) error

	// ContinueConversation 继续与AI的对话
	// 参数：
	//   - ctx: 上下文信息
	//   - req: 请求参数
	//
	// 返回值：
	//   - *conversation.ConversationResponse: 对话响应
	//   - error: 错误信息
	ContinueConversation(ctx *gin.Context, req *dto.ContinueConversationRequest) (*conversation.ConversationResponse, error)

	// StreamContinueConversation 流式继续对话
	// 参数：
	//   - ctx: 上下文信息
	//   - req: 请求参数
	//   - handler: 流式响应处理函数
	//
	// 返回值：
	//   - error: 错误信息
	StreamContinueConversation(ctx *gin.Context, req *dto.StreamContinueConversationRequest, handler func(content string, done bool) error) error

	// GetMessageIDs 获取当前用户的消息ID
	// 参数：
	//   - ctx: 上下文信息
	//
	// 返回值：
	//   - int: 请求消息ID
	//   - int: 响应消息ID
	//   - error: 错误信息
	GetMessageIDs(ctx *gin.Context) (int, int, error)
}
