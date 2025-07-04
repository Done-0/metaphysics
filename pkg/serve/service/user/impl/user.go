// Package impl 提供业务逻辑处理，处理用户相关业务
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
	model "github.com/Done-0/metaphysics/internal/model/user"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/user/dto"
	userMapper "github.com/Done-0/metaphysics/pkg/serve/mapper/user"
	userService "github.com/Done-0/metaphysics/pkg/serve/service/user"
	userVo "github.com/Done-0/metaphysics/pkg/vo/user"
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

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userMapper userMapper.UserMapper
}

// NewUserService 创建用户服务实例
// 参数：
//   - userMapperImpl: 用户数据访问接口
//
// 返回值：
//   - userService.UserService: 用户服务接口
func NewUserService(userMapperImpl userMapper.UserMapper) userService.UserService {
	return &UserServiceImpl{
		userMapper: userMapperImpl,
	}
}

// GetOneUser 获取用户信息逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 获取用户请求
//
// 返回值：
//   - *userVo.GetOneUserResponse: 用户用户视图对象
//   - error: 操作过程中的错误
func (a *UserServiceImpl) GetOneUser(c *gin.Context, req *dto.GetOneUserRequest) (*userVo.GetOneUserResponse, error) {
	userInfo, err := a.userMapper.GetOneUserByEmail(c, req.Email)
	if err != nil {
		utils.BizLogger(c).Errorf("「%s」用户不存在", req.Email)
		return nil, fmt.Errorf("「%s」用户不存在", req.Email)
	}

	getUserVO, err := utils.MapModelToVO(userInfo, &userVo.GetOneUserResponse{})
	if err != nil {
		utils.BizLogger(c).Errorf("获取用户信息时映射 VO 失败: %v", err)
		return nil, fmt.Errorf("获取用户信息时映射 VO 失败: %w", err)
	}

	return getUserVO.(*userVo.GetOneUserResponse), nil
}

// RegisterOneUser 用户注册逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 注册用户请求
//
// 返回值：
//   - *userVo.RegisterOneUserResponse: 注册后的用户视图对象
//   - error: 操作过程中的错误
func (a *UserServiceImpl) RegisterOneUser(c *gin.Context, req *dto.RegisterOneUserRequest) (*userVo.RegisterOneUserResponse, error) {
	registerLock.Lock()
	defer registerLock.Unlock()

	var registerVO *userVo.RegisterOneUserResponse

	err := utils.RunDBTransaction(c, func() error {
		existingUser, _ := a.userMapper.GetOneUserByEmail(c, req.Email)
		if existingUser != nil {
			utils.BizLogger(c).Errorf("「%s」邮箱已被注册", req.Email)
			return fmt.Errorf("「%s」邮箱已被注册", req.Email)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.BizLogger(c).Errorf("哈希加密失败: %v", err)
			return fmt.Errorf("哈希加密失败: %w", err)
		}

		user := &model.User{
			Email:    req.Email,
			Password: string(hashedPassword),
			Nickname: req.Nickname,
		}

		if err := a.userMapper.CreateOneUser(c, user); err != nil {
			utils.BizLogger(c).Errorf("「%s」用户注册失败: %v", req.Email, err)
			return fmt.Errorf("「%s」用户注册失败: %w", req.Email, err)
		}

		vo, err := utils.MapModelToVO(user, &userVo.RegisterOneUserResponse{})
		if err != nil {
			utils.BizLogger(c).Errorf("用户注册时映射 VO 失败: %v", err)
			return fmt.Errorf("用户注册时映射 VO 失败: %w", err)
		}

		registerVO = vo.(*userVo.RegisterOneUserResponse)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return registerVO, nil
}

// LoginOneUser 登录用户逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 登录请求
//
// 返回值：
//   - *userVo.LoginOneUserResponse: 登录成功后的令牌视图对象
//   - error: 操作过程中的错误
func (a *UserServiceImpl) LoginOneUser(c *gin.Context, req *dto.LoginOneUserRequest) (*userVo.LoginOneUserResponse, error) {
	user, err := a.userMapper.GetOneUserByEmail(c, req.Email)
	if err != nil {
		utils.BizLogger(c).Errorf("「%s」用户不存在: %v", req.Email, err)
		return nil, fmt.Errorf("「%s」用户不存在: %w", req.Email, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		utils.BizLogger(c).Errorf("「%s」用户密码输入错误: %v", user.Email, err)
		return nil, fmt.Errorf("「%s」用户密码输入错误: %w", user.Email, err)
	}

	accessTokenString, refreshTokenString, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.BizLogger(c).Errorf("access_token 生成失败: %v", err)
		return nil, fmt.Errorf("access_token 生成失败: %w", err)
	}

	cacheKey := fmt.Sprintf("%s:%d", USER_CACHE, user.ID)

	err = global.RedisClient.Set(context.Background(), cacheKey, accessTokenString, USER_CACHE_EXPIRE_TIME).Err()
	if err != nil {
		utils.BizLogger(c).Errorf("登录时设置缓存失败: %v", err)
		return nil, fmt.Errorf("登录时设置缓存失败: %w", err)
	}

	token := &userVo.LoginOneUserResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	loginVO, err := utils.MapModelToVO(token, &userVo.LoginOneUserResponse{})
	if err != nil {
		utils.BizLogger(c).Errorf("用户登录时映射 VO 失败: %v", err)
		return nil, fmt.Errorf("用户登陆时映射 VO 失败: %v", err)
	}

	return loginVO.(*userVo.LoginOneUserResponse), nil
}

// LogoutOneUser 处理用户登出逻辑
// 参数：
//   - c: Gin 上下文
//
// 返回值：
//   - error: 操作过程中的错误
func (a *UserServiceImpl) LogoutOneUser(c *gin.Context) error {
	logoutLock.Lock()
	defer logoutLock.Unlock()

	userID, err := utils.ParseAccountFromJWT(c.GetHeader("Authorization"))
	if err != nil {
		utils.BizLogger(c).Errorf("access_token 解析失败: %v", err)
		return fmt.Errorf("access_token 解析失败: %w", err)
	}

	cacheKey := fmt.Sprintf("%s:%d", USER_CACHE, userID)
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
func (a *UserServiceImpl) ResetPassword(c *gin.Context, req *dto.ResetPwdRequest) error {
	passwordResetLock.Lock()
	defer passwordResetLock.Unlock()

	return utils.RunDBTransaction(c, func() error {
		if req.NewPassword != req.AgainNewPassword {
			utils.BizLogger(c).Errorf("两次密码输入不一致")
			return fmt.Errorf("两次密码输入不一致")
		}

		userID, err := utils.ParseAccountFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			utils.BizLogger(c).Errorf("access_token 解析失败: %v", err)
			return fmt.Errorf("access_token 解析失败: %w", err)
		}

		user, err := a.userMapper.GetOneUserByID(c, userID)
		if err != nil {
			utils.BizLogger(c).Errorf("「%s」用户不存在: %v", user.Email, err)
			return fmt.Errorf("「%s」用户不存在: %w", user.Email, err)
		}

		newPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.BizLogger(c).Errorf("「%s」用户密码加密失败: %v", user.Email, err)
			return fmt.Errorf("「%s」用户密码加密失败: %w", user.Email, err)
		}
		user.Password = string(newPassword)

		if err := a.userMapper.UpdateOneUserByID(c, user); err != nil {
			utils.BizLogger(c).Errorf("「%s」用户密码修改失败: %v", user.Email, err)
			return fmt.Errorf("「%s」用户密码修改失败: %w", user.Email, err)
		}

		return nil
	})
}

// UpdateOneUserByID 更新用户信息逻辑
// 参数：
//   - c: Gin 上下文
//   - req: 更新用户请求
//
// 返回值：
//   - *userVo.UpdateOneUserResponse: 更新后的用户视图对象
//   - error: 操作过程中的错误
func (a *UserServiceImpl) UpdateOneUserByID(c *gin.Context, req *dto.UpdateOneUserRequest) (*userVo.UpdateOneUserResponse, error) {
	var updateVO *userVo.UpdateOneUserResponse

	err := utils.RunDBTransaction(c, func() error {
		userID, err := utils.ParseAccountFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			utils.BizLogger(c).Errorf("access_token 解析失败: %v", err)
			return fmt.Errorf("access_token 解析失败: %w", err)
		}

		user, err := a.userMapper.GetOneUserByID(c, userID)
		if err != nil {
			utils.BizLogger(c).Errorf("获取「%s」用户信息失败: %v", user.Email, err)
			return fmt.Errorf("获取「%s」用户信息失败: %w", user.Email, err)
		}

		user.Nickname = req.Nickname
		user.Avatar = req.Avatar

		if err := a.userMapper.UpdateOneUserByID(c, user); err != nil {
			utils.BizLogger(c).Errorf("更新「%s」用户信息失败: %v", user.Email, err)
			return fmt.Errorf("更新「%s」用户信息失败: %w", user.Email, err)
		}

		vo, err := utils.MapModelToVO(user, &userVo.UpdateOneUserResponse{})
		if err != nil {
			utils.BizLogger(c).Errorf("更新「%s」用户信息时映射 VO 失败: %v", user.Email, err)
			return fmt.Errorf("更新「%s」用户信息时映射 VO 失败: %w", user.Email, err)
		}

		updateVO = vo.(*userVo.UpdateOneUserResponse)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updateVO, nil
}
