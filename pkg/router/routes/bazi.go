// Package routes 提供八字相关路由
// 创建者：Done-0
// 创建时间：2023-10-18
package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/pkg/serve/controller/bazi"
	baziImpl "github.com/Done-0/metaphysics/pkg/serve/service/bazi/impl"
)

// RegisterBaziRoutes 注册八字相关路由
func RegisterBaziRoutes(r *gin.Engine) {
	// 创建服务实例
	service := baziImpl.NewBaziService()
	// 创建控制器
	controller := bazi.NewBaziController(service)

	v1 := r.Group("/api/v1")
	{
		baziGroup := v1.Group("/bazi")
		{
			baziGroup.POST("/analyze", controller.AnalyzeOneBazi)              // 普通分析八字
			baziGroup.POST("/analyze/stream", controller.StreamAnalyzeOneBazi) // 流式分析八字
			baziGroup.GET("/record", controller.GetBaziRecord)                 // 获取八字记录
		}
	}
}
