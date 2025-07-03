// Package mapper 提供数据模型与数据库交互的映射层，处理用户相关数据操作
// 创建者：Done-0
// 创建时间：2025-05-10
package mapper

import (
	"fmt"

	"github.com/gin-gonic/gin"

	account "github.com/Done-0/metaphysics/internal/model/account"
	"github.com/Done-0/metaphysics/internal/utils"
)

// GetOneAccountByEmail 根据邮箱获取用户账号信息
// 参数：
//   - c: Gin 上下文
//   - email: 用户邮箱
//
// 返回值：
//   - *account.Account: 用户信息
//   - error: 操作过程中的错误
func GetOneAccountByEmail(c *gin.Context, email string) (*account.Account, error) {
	var user account.Account
	db := utils.GetDBFromContext(c)
	if err := db.Where("email = ? AND deleted = ?", email, false).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return &user, nil
}

// GetOneAccountByID 根据用户 ID 获取用户账号信息
// 参数：
//   - c: Gin 上下文
//   - accountID: 用户 ID
//
// 返回值：
//   - *account.Account: 账户信息
//   - error: 操作过程中的错误
func GetOneAccountByID(c *gin.Context, accountID int64) (*account.Account, error) {
	var user account.Account
	db := utils.GetDBFromContext(c)
	if err := db.Where("id = ? AND deleted = ?", accountID, false).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return &user, nil
}

// CreateOneAccount 创建新用户
// 参数：
//   - c: Gin 上下文
//   - acc: 账户信息
//
// 返回值：
//   - error: 操作过程中的错误
func CreateOneAccount(c *gin.Context, account *account.Account) error {
	db := utils.GetDBFromContext(c)
	if err := db.Create(account).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

// UpdateOneAccountByID 更新账户信息
// 参数：
//   - c: Gin 上下文
//   - acc: 账户信息
//
// 返回值：
//   - error: 操作过程中的错误
func UpdateOneAccountByID(c *gin.Context, account *account.Account) error {
	db := utils.GetDBFromContext(c)
	if err := db.Where("id = ? AND deleted = ?", account.ID, false).Updates(account).Error; err != nil {
		return fmt.Errorf("更新账户失败: %w", err)
	}
	return nil
}
