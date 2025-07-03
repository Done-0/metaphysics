// Package routes 提供路由注册功能
// 创建者：Done-0
// 创建时间：2025-05-10
package routes

import (
	"github.com/gin-gonic/gin"

	auth_middleware "github.com/Done-0/metaphysics/internal/middleware/auth"
	"github.com/Done-0/metaphysics/pkg/serve/controller/account"
	accountImpl "github.com/Done-0/metaphysics/pkg/serve/service/account/impl"
)

// RegisterAccountRoutes 注册账户相关路由
// 参数：
//   - r: Gin 路由组
func RegisterAccountRoutes(r *gin.RouterGroup) {
	// 创建服务实例
	service := accountImpl.NewAccountService()
	// 创建控制器
	controller := account.NewAccountController(service)

	// 账户路由组
	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/register", controller.RegisterAcc)
		accountGroup.POST("/login", controller.LoginAccount)
		accountGroup.GET("/info", auth_middleware.AuthMiddleware(), controller.GetAccount)
		accountGroup.POST("/update", auth_middleware.AuthMiddleware(), controller.UpdateAccount)
		accountGroup.POST("/logout", auth_middleware.AuthMiddleware(), controller.LogoutAccount)
		accountGroup.POST("/resetPassword", auth_middleware.AuthMiddleware(), controller.ResetPassword)
	}
}
