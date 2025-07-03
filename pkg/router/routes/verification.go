// Package routes 提供路由注册功能
// 创建者：Done-0
// 创建时间：2025-05-10
package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/pkg/serve/controller/verification"
	verificationImpl "github.com/Done-0/metaphysics/pkg/serve/service/verification/impl"
)

// RegisterVerificationRoutes 注册验证码相关路由
// 参数：
//   - r: Gin 路由组
func RegisterVerificationRoutes(r *gin.RouterGroup) {
	// 创建服务实例
	service := verificationImpl.NewVerificationService()
	// 创建控制器
	controller := verification.NewVerificationController(service)

	// 验证码路由组
	verificationGroup := r.Group("/verification")
	{
		verificationGroup.GET("/sendEmailCode", controller.SendEmailVerificationCode) // 发送邮箱验证码
	}
}
