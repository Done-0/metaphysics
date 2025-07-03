// Package account 提供账户相关的服务层接口
// 创建者：Done-0
// 创建时间：2025-05-10
package account

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/pkg/serve/controller/account/dto"
	"github.com/Done-0/metaphysics/pkg/vo/account"
)

// AccountService 账户服务接口
type AccountService interface {
	// RegisterOneAccount 注册用户账号
	// 参数：
	//   - c: Gin上下文
	//   - req: 注册请求参数
	// 返回值：
	//   - *account.RegisterAccountVO: 注册结果视图对象
	//   - error: 错误信息
	RegisterOneAccount(c *gin.Context, req *dto.RegisterOneAccountRequest) (*account.RegisterAccountResponse, error)

	// LoginOneAccount 用户登录
	// 参数：
	//   - c: Gin上下文
	//   - req: 登录请求参数
	// 返回值：
	//   - *account.LoginVO: 登录结果视图对象
	//   - error: 错误信息
	LoginOneAccount(c *gin.Context, req *dto.LoginOneAccountRequest) (*account.LoginResponse, error)

	// GetOneAccount 获取用户信息
	// 参数：
	//   - c: Gin上下文
	//   - req: 获取账户请求参数
	// 返回值：
	//   - *account.GetAccountVO: 账户信息视图对象
	//   - error: 错误信息
	GetOneAccount(c *gin.Context, req *dto.GetOneAccountRequest) (*account.GetAccountResponse, error)

	// UpdateOneAccountByID 更新用户信息
	// 参数：
	//   - c: Gin上下文
	//   - req: 更新账户请求参数
	// 返回值：
	//   - *account.UpdateAccountVO: 更新账户结果视图对象
	//   - error: 错误信息
	UpdateOneAccountByID(c *gin.Context, req *dto.UpdateOneAccountRequest) (*account.UpdateAccountResponse, error)

	// LogoutOneAccount 用户登出
	// 参数：
	//   - c: Gin上下文
	// 返回值：
	//   - error: 错误信息
	LogoutOneAccount(c *gin.Context) error

	// ResetPassword 重置密码
	// 参数：
	//   - c: Gin上下文
	//   - req: 重置密码请求参数
	// 返回值：
	//   - error: 错误信息
	ResetPassword(c *gin.Context, req *dto.ResetPwdRequest) error
}
