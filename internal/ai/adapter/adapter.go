// Package adapter 提供AI模块与对话模块之间的适配器
// 创建者：Done-0
// 创建时间：2025-07-03
package adapter

import (
	"github.com/Done-0/metaphysics/internal/ai/types"
	"github.com/Done-0/metaphysics/pkg/vo/conversation"
)

// AIToConversationStreamHandler 将对话模块的StreamChunk处理函数转换为AI模块的StreamHandler
// 参数：
//   - handler: 对话模块的StreamChunk处理函数
//
// 返回值：
//   - types.StreamHandler: AI模块的StreamHandler
func AIToConversationStreamHandler(handler func(chunk *conversation.StreamChunk) error) types.StreamHandler {
	return func(chunk *conversation.StreamChunk) error {
		return handler(chunk)
	}
}

// AIToConversationBaziAnalysisResponse 将AI模块的BaziAnalysisResponse转换为对话模块的BaziAnalysisResponse
// 参数：
//   - aiResp: AI模块的BaziAnalysisResponse
//
// 返回值：
//   - *conversation.BaziAnalysisResponse: 对话模块的BaziAnalysisResponse
func AIToConversationBaziAnalysisResponse(aiResp *conversation.BaziAnalysisResponse) *conversation.BaziAnalysisResponse {
	return aiResp
}
