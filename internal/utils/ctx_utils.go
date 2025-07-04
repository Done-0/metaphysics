// Package utils 提供通用工具函数
// 创建者：Done-0
// 创建时间：2025-07-04
package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext 从上下文中获取用户ID
// 参数：
//   - ctx: Gin上下文
//
// 返回值：
//   - int64: 用户ID
//   - error: 错误信息
func GetUserIDFromContext(ctx *gin.Context) (int64, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("未找到用户ID")
	}
	return userID.(int64), nil
}
