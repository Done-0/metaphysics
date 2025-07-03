// Package utils 提供各种工具函数，包括数据库事务管理
// 创建者：Done-0
// 创建时间：2025-07-01
package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Done-0/metaphysics/internal/global"
)

// DB_TRANSACTION_CONTEXT_KEY 事务相关常量
const DB_TRANSACTION_CONTEXT_KEY = "tx" // 存储在 Gin 上下文中的数据库事务键名

// GetDBFromContext 从上下文中获取数据库连接
// 参数：
//   - c: Gin 上下文
//
// 返回值：
//   - *gorm.DB: 数据库连接（事务优先，无事务则返回全局连接）
func GetDBFromContext(c *gin.Context) *gorm.DB {
	if tx, exists := c.Get(DB_TRANSACTION_CONTEXT_KEY); exists {
		if db, ok := tx.(*gorm.DB); ok && db != nil {
			return db
		}
	}
	return global.DB
}

// RunDBTransaction 在事务中执行函数
// 参数：
//   - c: Gin 上下文
//   - fn: 事务内执行的函数
//
// 返回值：
//   - error: 执行过程中的错误
func RunDBTransaction(c *gin.Context, fn func() error) error {
	tx := global.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %w", tx.Error)
	}

	c.Set(DB_TRANSACTION_CONTEXT_KEY, tx)
	defer c.Set(DB_TRANSACTION_CONTEXT_KEY, nil)

	// panic 处理
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 执行业务逻辑
	if err := fn(); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}
