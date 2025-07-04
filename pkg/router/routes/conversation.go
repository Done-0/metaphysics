// Package routes 提供对话相关路由
// 创建者：Done-0
// 创建时间：2025-07-03
package routes

import (
	"github.com/gin-gonic/gin"

	auth_middleware "github.com/Done-0/metaphysics/internal/middleware/auth"
	conversationController "github.com/Done-0/metaphysics/pkg/serve/controller/conversation"
	conversationMapperImpl "github.com/Done-0/metaphysics/pkg/serve/mapper/conversation/impl"
	conversationServiceImpl "github.com/Done-0/metaphysics/pkg/serve/service/conversation/impl"
)

// RegisterConversationRoutes 注册对话相关路由
// 参数：
//   - r: Gin 路由组
func RegisterConversationRoutes(r *gin.RouterGroup) {
	mapper := conversationMapperImpl.NewConversationMapper()
	service := conversationServiceImpl.NewConversationService(mapper)
	controller := conversationController.NewConversationController(service)

	// 对话路由组
	conversationGroup := r.Group("/conversation")
	conversationGroup.Use(auth_middleware.AuthMiddleware())
	{
		// 八字分析
		conversationGroup.GET("/bazi/analyze", controller.AnalyzeBazi)
		conversationGroup.GET("/bazi/analyze/stream", controller.StreamAnalyzeBazi)

		// 对话
		conversationGroup.POST("/continue", controller.ContinueConversation)
		conversationGroup.POST("/continue/stream", controller.StreamContinueConversation)
	}
}
