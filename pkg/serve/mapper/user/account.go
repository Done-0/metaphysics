// Package user 提供用户相关的数据访问接口
// 创建者：Done-0
// 创建时间：2025-05-10
package user

import (
	"github.com/gin-gonic/gin"

	userModel "github.com/Done-0/metaphysics/internal/model/user"
)

// UserMapper 用户数据访问接口
type UserMapper interface {
	// GetOneUserByEmail 根据邮箱获取用户账号信息
	// 参数：
	//   - c: Gin 上下文
	//   - email: 用户邮箱
	//
	// 返回值：
	//   - *userModel.User: 用户信息
	//   - error: 操作过程中的错误
	GetOneUserByEmail(c *gin.Context, email string) (*userModel.User, error)

	// GetOneUserByID 根据用户 ID 获取用户账号信息
	// 参数：
	//   - c: Gin 上下文
	//   - userID: 用户 ID
	//
	// 返回值：
	//   - *userModel.User: 用户信息
	//   - error: 操作过程中的错误
	GetOneUserByID(c *gin.Context, userID int64) (*userModel.User, error)

	// CreateOneUser 创建新用户
	// 参数：
	//   - c: Gin 上下文
	//   - user: 用户信息
	//
	// 返回值：
	//   - error: 操作过程中的错误
	CreateOneUser(c *gin.Context, user *userModel.User) error

	// UpdateOneUserByID 更新用户信息
	// 参数：
	//   - c: Gin 上下文
	//   - user: 用户信息
	//
	// 返回值：
	//   - error: 操作过程中的错误
	UpdateOneUserByID(c *gin.Context, user *userModel.User) error
}
