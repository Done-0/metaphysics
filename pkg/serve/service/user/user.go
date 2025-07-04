// Package user 提供用户相关的服务层接口
// 创建者：Done-0
// 创建时间：2025-05-10
package user

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/pkg/serve/controller/user/dto"
	"github.com/Done-0/metaphysics/pkg/vo/user"
)

// AccountService 用户服务接口
type UserService interface {
	// RegisterOneUser 注册用户账号
	// 参数：
	//   - c: Gin上下文
	//   - req: 注册请求参数
	// 返回值：
	//   - *user.RegisterOneUserResponse: 注册结果视图对象
	//   - error: 错误信息
	RegisterOneUser(c *gin.Context, req *dto.RegisterOneUserRequest) (*user.RegisterOneUserResponse, error)

	// LoginOneUser 用户登录
	// 参数：
	//   - c: Gin上下文
	//   - req: 登录请求参数
	// 返回值：
	//   - *user.LoginOneUserResponse: 登录结果视图对象
	//   - error: 错误信息
	LoginOneUser(c *gin.Context, req *dto.LoginOneUserRequest) (*user.LoginOneUserResponse, error)

	// GetOneUser 获取用户信息
	// 参数：
	//   - c: Gin上下文
	//   - req: 获取用户请求参数
	// 返回值：
	//   - *user.GetOneUserResponse: 用户信息视图对象
	//   - error: 错误信息
	GetOneUser(c *gin.Context, req *dto.GetOneUserRequest) (*user.GetOneUserResponse, error)

	// UpdateOneUserByID 更新用户信息
	// 参数：
	//   - c: Gin上下文
	//   - req: 更新用户请求参数
	// 返回值：
	//   - *user.UpdateOneUserResponse: 更新用户结果视图对象
	//   - error: 错误信息
	UpdateOneUserByID(c *gin.Context, req *dto.UpdateOneUserRequest) (*user.UpdateOneUserResponse, error)

	// LogoutOneAccount 用户登出
	// 参数：
	//   - c: Gin上下文
	// 返回值：
	//   - error: 错误信息
	LogoutOneUser(c *gin.Context) error

	// ResetPassword 重置密码
	// 参数：
	//   - c: Gin上下文
	//   - req: 重置密码请求参数
	// 返回值：
	//   - error: 错误信息
	ResetPassword(c *gin.Context, req *dto.ResetPwdRequest) error
}
