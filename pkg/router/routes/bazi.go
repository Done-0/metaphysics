// Package routes 提供八字相关路由
// 创建者：Done-0
// 创建时间：2023-10-18
package routes

import (
	"github.com/gin-gonic/gin"

	auth_middleware "github.com/Done-0/metaphysics/internal/middleware/auth"
	"github.com/Done-0/metaphysics/pkg/serve/controller/bazi"
	baziMapperImpl "github.com/Done-0/metaphysics/pkg/serve/mapper/bazi/impl"
	baziImpl "github.com/Done-0/metaphysics/pkg/serve/service/bazi/impl"
)

// RegisterBaziRoutes 注册八字相关路由
// 参数：
//   - r: Gin 路由组
func RegisterBaziRoutes(r *gin.RouterGroup) {
	mapper := baziMapperImpl.NewBaziMapper()
	service := baziImpl.NewBaziService(mapper)
	controller := bazi.NewBaziController(service)

	// 八字路由组
	baziGroup := r.Group("/bazi")
	{
		baziGroup.POST("/calculate", controller.CalculateOneBazi)
		baziGroup.GET("/record", auth_middleware.AuthMiddleware(), controller.GetOneBazi)
		baziGroup.GET("/records", auth_middleware.AuthMiddleware(), controller.GetBaziList)
	}
}
