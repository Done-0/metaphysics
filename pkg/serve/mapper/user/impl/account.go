// Package impl 提供账号相关的数据访问实现
// 创建者：Done-0
// 创建时间：2025-05-10
package impl

import (
	"fmt"

	"github.com/gin-gonic/gin"

	userModel "github.com/Done-0/metaphysics/internal/model/user"
	"github.com/Done-0/metaphysics/internal/utils"
	userMapper "github.com/Done-0/metaphysics/pkg/serve/mapper/user"
)

// UserMapperImpl 用户数据访问实现
type UserMapperImpl struct{}

// NewUserMapper 创建用户数据访问实例
// 返回值：
//   - userMapper.UserMapper: 用户数据访问接口
func NewUserMapper() userMapper.UserMapper {
	return &UserMapperImpl{}
}

// GetOneUserByEmail 根据邮箱获取用户账号信息
// 参数：
//   - c: Gin 上下文
//   - email: 用户邮箱
//
// 返回值：
//   - *userModel.User: 用户信息
//   - error: 操作过程中的错误
func (m *UserMapperImpl) GetOneUserByEmail(c *gin.Context, email string) (*userModel.User, error) {
	var user userModel.User
	db := utils.GetDBFromContext(c)
	if err := db.Where("email = ? AND deleted = ?", email, false).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return &user, nil
}

// GetOneAccountByID 根据用户 ID 获取用户账号信息
// 参数：
//   - c: Gin 上下文
//   - userID: 用户 ID
//
// 返回值：
//   - *userModel.User: 用户信息
//   - error: 操作过程中的错误
func (m *UserMapperImpl) GetOneUserByID(c *gin.Context, userID int64) (*userModel.User, error) {
	var user userModel.User
	db := utils.GetDBFromContext(c)
	if err := db.Where("id = ? AND deleted = ?", userID, false).First(&user).Error; err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return &user, nil
}

// CreateOneUser 创建新用户
// 参数：
//   - c: Gin 上下文
//   - user: 用户信息
//
// 返回值：
//   - error: 操作过程中的错误
func (m *UserMapperImpl) CreateOneUser(c *gin.Context, user *userModel.User) error {
	return utils.RunDBTransaction(c, func() error {
		db := utils.GetDBFromContext(c)
		if err := db.Create(user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
		return nil
	})
}

// UpdateOneUserByID 更新用户信息
// 参数：
//   - c: Gin 上下文
//   - user: 用户信息
//
// 返回值：
//   - error: 操作过程中的错误
func (m *UserMapperImpl) UpdateOneUserByID(c *gin.Context, user *userModel.User) error {
	return utils.RunDBTransaction(c, func() error {
		db := utils.GetDBFromContext(c)
		if err := db.Where("id = ? AND deleted = ?", user.ID, false).Updates(user).Error; err != nil {
			return fmt.Errorf("更新用户失败: %w", err)
		}
		return nil
	})
}
