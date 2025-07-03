// Package auth_middleware 提供JWT认证相关中间件
// 创建者：Done-0
// 创建时间：2025-05-10
package auth_middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/internal/global"
	"github.com/Done-0/metaphysics/internal/utils"
)

// JWTConfig 定义了 Token 相关的配置
type JWTConfig struct {
	Authorization string // 认证头名称
	TokenPrefix   string // Token前缀
	RefreshToken  string // 刷新令牌头名称
	UserCache     string // 用户缓存键前缀
}

// DefaultJWTConfig 默认配置
var DefaultJWTConfig = JWTConfig{
	Authorization: "Authorization",
	TokenPrefix:   "Bearer ",
	RefreshToken:  "Refresh-Token",
	UserCache:     "USER_CACHE",
}

// AuthMiddleware 处理 JWT 认证中间件
// 返回值：
//   - gin.HandlerFunc: Gin 框架中间件函数
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中提取 Access-Token
		AuthorizationHeader := c.GetHeader(DefaultJWTConfig.Authorization)
		if AuthorizationHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization 请求头缺失",
			})
			return
		}
		tokenString := strings.TrimPrefix(AuthorizationHeader, DefaultJWTConfig.TokenPrefix)

		// 验证 Access-Token；若 Access-Token 已过期或无效，则尝试使用 Refresh Token 刷新
		_, err := utils.ValidateJWTToken(tokenString, false)
		if err != nil {
			refreshTokenHeader := c.GetHeader(DefaultJWTConfig.RefreshToken)
			if refreshTokenHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Access-Token 已过期，请提供 Refresh-Token",
				})
				return
			}

			newTokens, err := utils.RefreshTokenLogic(refreshTokenHeader)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Refresh-Token 无效或已过期，请重新登录",
				})
				return
			}

			// 设置新的 Access-Token 和 Refresh-Token 到响应头 Authorization 和 Refresh-Token
			c.Header(DefaultJWTConfig.Authorization, DefaultJWTConfig.TokenPrefix+newTokens["Authorization"])
			c.Header(DefaultJWTConfig.RefreshToken, newTokens["Refresh-Token"])
			tokenString = newTokens["Authorization"]
		}

		// 从 access_token 中解析 accountID
		accountID, err := utils.ParseAccountFromJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Access-Token 解析失败，请重新登录",
			})
			return
		}

		// 检验会话有效性
		sessionCacheKey := fmt.Sprintf("%s:%d", DefaultJWTConfig.UserCache, accountID)
		if sessionVal, err := global.RedisClient.Get(c.Request.Context(), sessionCacheKey).Result(); err != nil || sessionVal == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "会话已失效，请重新登录",
			})
			return
		}

		// 将 accountID 存入上下文，方便后续使用
		c.Set("accountID", accountID)

		c.Next()
	}
}
