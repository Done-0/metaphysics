// Package impl 提供业务逻辑处理，处理账户相关业务
// 创建者：Done-0
// 创建时间：2025-05-10
package impl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/Done-0/metaphysics/internal/global"
	model "github.com/Done-0/metaphysics/internal/model/account"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/account/dto"
	"github.com/Done-0/metaphysics/pkg/serve/mapper"
	accountService "github.com/Done-0/metaphysics/pkg/serve/service/account"
	accountVo "github.com/Done-0/metaphysics/pkg/vo/account"
)

var (
	registerLock      sync.Mutex // 用户注册锁，保护并发用户注册的操作
	passwordResetLock sync.Mutex // 修改密码锁，保护并发修改用户密码的操作
	logoutLock        sync.Mutex // 用户登出锁，保护并发用户登出操作
)

const (
	USER_CACHE             = "USER_CACHE"
	USER_CACHE_EXPIRE_TIME = time.Hour * 2 // Access Token 有效期
)

// AccountServiceImpl 账户服务实现
type AccountServiceImpl struct{}

// NewAccountService 创建账户服务实例
// 返回值：
//   - accountService.AccountService: 账户服务接口
func NewAccountService() accountService.AccountService {
	return &AccountServiceImpl{}
}

// GetOneAccount 获取用户信息逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 获取账户请求
//
// 返回值：
//   - *accountVo.GetAccountVO: 用户账户视图对象
//   - error: 操作过程中的错误
func (a *AccountServiceImpl) GetOneAccount(c *gin.Context, req *dto.GetOneAccountRequest) (*accountVo.GetAccountResponse, error) {
	userInfo, err := mapper.GetOneAccountByEmail(c, req.Email)
	if err != nil {
		utils.BizLogger(c).Errorf("「%s」用户不存在", req.Email)
		return nil, fmt.Errorf("「%s」用户不存在", req.Email)
	}

	getAccountVO, err := utils.MapModelToVO(userInfo, &accountVo.GetAccountResponse{})
	if err != nil {
		utils.BizLogger(c).Errorf("获取用户信息时映射 VO 失败: %v", err)
		return nil, fmt.Errorf("获取用户信息时映射 VO 失败: %w", err)
	}

	return getAccountVO.(*accountVo.GetAccountResponse), nil
}

// RegisterOneAccount 用户注册逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 注册账户请求
//
// 返回值：
//   - *accountVo.RegisterAccountVO: 注册后的账户视图对象
//   - error: 操作过程中的错误
func (a *AccountServiceImpl) RegisterOneAccount(c *gin.Context, req *dto.RegisterOneAccountRequest) (*accountVo.RegisterAccountResponse, error) {
	registerLock.Lock()
	defer registerLock.Unlock()

	var registerVO *accountVo.RegisterAccountResponse

	err := utils.RunDBTransaction(c, func() error {
		existingUser, _ := mapper.GetOneAccountByEmail(c, req.Email)
		if existingUser != nil {
			utils.BizLogger(c).Errorf("「%s」邮箱已被注册", req.Email)
			return fmt.Errorf("「%s」邮箱已被注册", req.Email)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.BizLogger(c).Errorf("哈希加密失败: %v", err)
			return fmt.Errorf("哈希加密失败: %w", err)
		}

		acc := &model.Account{
			Email:    req.Email,
			Password: string(hashedPassword),
			Nickname: req.Nickname,
		}

		if err := mapper.CreateOneAccount(c, acc); err != nil {
			utils.BizLogger(c).Errorf("「%s」用户注册失败: %v", req.Email, err)
			return fmt.Errorf("「%s」用户注册失败: %w", req.Email, err)
		}

		vo, err := utils.MapModelToVO(acc, &accountVo.RegisterAccountResponse{})
		if err != nil {
			utils.BizLogger(c).Errorf("用户注册时映射 VO 失败: %v", err)
			return fmt.Errorf("用户注册时映射 VO 失败: %w", err)
		}

		registerVO = vo.(*accountVo.RegisterAccountResponse)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return registerVO, nil
}

// LoginOneAccount 登录用户逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 登录请求
//
// 返回值：
//   - *accountVo.LoginVO: 登录成功后的令牌视图对象
//   - error: 操作过程中的错误
func (a *AccountServiceImpl) LoginOneAccount(c *gin.Context, req *dto.LoginOneAccountRequest) (*accountVo.LoginResponse, error) {
	acc, err := mapper.GetOneAccountByEmail(c, req.Email)
	if err != nil {
		utils.BizLogger(c).Errorf("「%s」用户不存在: %v", req.Email, err)
		return nil, fmt.Errorf("「%s」用户不存在: %w", req.Email, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password))
	if err != nil {
		utils.BizLogger(c).Errorf("「%s」用户密码输入错误: %v", acc.Email, err)
		return nil, fmt.Errorf("「%s」用户密码输入错误: %w", acc.Email, err)
	}

	accessTokenString, refreshTokenString, err := utils.GenerateJWT(acc.ID)
	if err != nil {
		utils.BizLogger(c).Errorf("access_token 生成失败: %v", err)
		return nil, fmt.Errorf("access_token 生成失败: %w", err)
	}

	cacheKey := fmt.Sprintf("%s:%d", USER_CACHE, acc.ID)

	err = global.RedisClient.Set(context.Background(), cacheKey, accessTokenString, USER_CACHE_EXPIRE_TIME).Err()
	if err != nil {
		utils.BizLogger(c).Errorf("登录时设置缓存失败: %v", err)
		return nil, fmt.Errorf("登录时设置缓存失败: %w", err)
	}

	token := &accountVo.LoginResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	loginVO, err := utils.MapModelToVO(token, &accountVo.LoginResponse{})
	if err != nil {
		utils.BizLogger(c).Errorf("用户登录时映射 VO 失败: %v", err)
		return nil, fmt.Errorf("用户登陆时映射 VO 失败: %v", err)
	}

	return loginVO.(*accountVo.LoginResponse), nil
}

// LogoutOneAccount 处理用户登出逻辑
// 参数：
//   - c: Gin 上下文
//
// 返回值：
//   - error: 操作过程中的错误
func (a *AccountServiceImpl) LogoutOneAccount(c *gin.Context) error {
	logoutLock.Lock()
	defer logoutLock.Unlock()

	accountID, err := utils.ParseAccountFromJWT(c.GetHeader("Authorization"))
	if err != nil {
		utils.BizLogger(c).Errorf("access_token 解析失败: %v", err)
		return fmt.Errorf("access_token 解析失败: %w", err)
	}

	cacheKey := fmt.Sprintf("%s:%d", USER_CACHE, accountID)
	err = global.RedisClient.Del(c.Request.Context(), cacheKey).Err()
	if err != nil {
		utils.BizLogger(c).Errorf("删除 Redis 缓存失败: %v", err)
		return fmt.Errorf("删除 Redis 缓存失败: %w", err)
	}

	return nil
}

// ResetPassword 重置密码逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 重置密码请求
//
// 返回值：
//   - error: 操作过程中的错误
func (a *AccountServiceImpl) ResetPassword(c *gin.Context, req *dto.ResetPwdRequest) error {
	passwordResetLock.Lock()
	defer passwordResetLock.Unlock()

	return utils.RunDBTransaction(c, func() error {
		if req.NewPassword != req.AgainNewPassword {
			utils.BizLogger(c).Errorf("两次密码输入不一致")
			return fmt.Errorf("两次密码输入不一致")
		}

		accountID, err := utils.ParseAccountFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			utils.BizLogger(c).Errorf("access_token 解析失败: %v", err)
			return fmt.Errorf("access_token 解析失败: %w", err)
		}

		acc, err := mapper.GetOneAccountByID(c, accountID)
		if err != nil {
			utils.BizLogger(c).Errorf("「%s」用户不存在: %v", acc.Email, err)
			return fmt.Errorf("「%s」用户不存在: %w", acc.Email, err)
		}

		newPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.BizLogger(c).Errorf("「%s」用户密码加密失败: %v", acc.Email, err)
			return fmt.Errorf("「%s」用户密码加密失败: %w", acc.Email, err)
		}
		acc.Password = string(newPassword)

		if err := mapper.UpdateOneAccountByID(c, acc); err != nil {
			utils.BizLogger(c).Errorf("「%s」用户密码修改失败: %v", acc.Email, err)
			return fmt.Errorf("「%s」用户密码修改失败: %w", acc.Email, err)
		}

		return nil
	})
}

// UpdateOneAccountByID 更新账户信息逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 更新账户请求
//
// 返回值：
//   - *accountVo.UpdateAccountVO: 更新后的账户视图对象
//   - error: 操作过程中的错误
func (a *AccountServiceImpl) UpdateOneAccountByID(c *gin.Context, req *dto.UpdateOneAccountRequest) (*accountVo.UpdateAccountResponse, error) {
	var updateVO *accountVo.UpdateAccountResponse

	err := utils.RunDBTransaction(c, func() error {
		accountID, err := utils.ParseAccountFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			utils.BizLogger(c).Errorf("access_token 解析失败: %v", err)
			return fmt.Errorf("access_token 解析失败: %w", err)
		}

		acc, err := mapper.GetOneAccountByID(c, accountID)
		if err != nil {
			utils.BizLogger(c).Errorf("获取「%s」用户信息失败: %v", acc.Email, err)
			return fmt.Errorf("获取「%s」用户信息失败: %w", acc.Email, err)
		}

		acc.Nickname = req.Nickname
		acc.Avatar = req.Avatar

		if err := mapper.UpdateOneAccountByID(c, acc); err != nil {
			utils.BizLogger(c).Errorf("更新「%s」用户信息失败: %v", acc.Email, err)
			return fmt.Errorf("更新「%s」用户信息失败: %w", acc.Email, err)
		}

		vo, err := utils.MapModelToVO(acc, &accountVo.UpdateAccountResponse{})
		if err != nil {
			utils.BizLogger(c).Errorf("更新「%s」用户信息时映射 VO 失败: %v", acc.Email, err)
			return fmt.Errorf("更新「%s」用户信息时映射 VO 失败: %w", acc.Email, err)
		}

		updateVO = vo.(*accountVo.UpdateAccountResponse)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updateVO, nil
}
