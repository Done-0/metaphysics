// Package routes 提供路由注册功能
// 创建者：Done-0
// 创建时间：2025-05-10
package routes

import (
	"github.com/gin-gonic/gin"

	auth_middleware "github.com/Done-0/metaphysics/internal/middleware/auth"
	userController "github.com/Done-0/metaphysics/pkg/serve/controller/user"
	userMapperImpl "github.com/Done-0/metaphysics/pkg/serve/mapper/user/impl"
	userImpl "github.com/Done-0/metaphysics/pkg/serve/service/user/impl"
)

// RegisterUserRoutes 注册用户相关路由
// 参数：
//   - r: Gin 路由组
func RegisterUserRoutes(r *gin.RouterGroup) {
	mapper := userMapperImpl.NewUserMapper()
	service := userImpl.NewUserService(mapper)
	controller := userController.NewUserController(service)

	// 用户路由组
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", controller.RegisterOneUser)
		userGroup.POST("/login", controller.LoginOneUser)
		userGroup.GET("/info", auth_middleware.AuthMiddleware(), controller.GetOneUser)
		userGroup.POST("/update", auth_middleware.AuthMiddleware(), controller.UpdateOneUser)
		userGroup.POST("/logout", auth_middleware.AuthMiddleware(), controller.LogoutOneUser)
		userGroup.POST("/resetPassword", auth_middleware.AuthMiddleware(), controller.ResetPassword)
	}
}
