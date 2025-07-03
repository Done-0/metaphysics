// Package impl 提供验证码相关的业务逻辑实现
// 创建者：Done-0
// 创建时间：2025-05-10
package impl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	bizErr "github.com/Done-0/metaphysics/internal/error"
	"github.com/Done-0/metaphysics/internal/global"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/verification/dto"
	"github.com/Done-0/metaphysics/pkg/serve/service/verification"
)

const (
	EMAIL_VERIFICATION_CODE_CACHE_EXPIRATION = 3 * time.Minute // 邮箱验证码缓存过期时间
)

// VerificationServiceImpl 验证码服务实现
type VerificationServiceImpl struct{}

// NewVerificationService 创建验证码服务实例
// 返回值：
//   - verification.VerificationService: 验证码服务接口
func NewVerificationService() verification.VerificationService {
	return &VerificationServiceImpl{}
}

// SendEmailVerificationCode 发送邮箱验证码
// 参数：
//   - c: Gin 上下文
//   - req: 获取验证码请求参数
//
// 返回值：
//   - error: 操作过程中的错误
func (v *VerificationServiceImpl) SendEmailVerificationCode(c *gin.Context, req *dto.GetOneVerificationCode) error {
	key := utils.EMAIL_VERIFICATION_CODE_CACHE_KEY_PREFIX + req.Email

	// 检查验证码是否存在
	exists, err := global.RedisClient.Exists(context.Background(), key).Result()
	if err != nil {
		utils.BizLogger(c).Errorf("检查邮箱验证码是否有效失败: %v", err)
		return err
	}
	if exists > 0 {
		return bizErr.New(bizErr.SYSTEM_ERROR, "邮箱验证码已存在")
	}

	// 生成并缓存验证码
	code := utils.NewRand()
	err = global.RedisClient.Set(context.Background(), key, strconv.Itoa(code), EMAIL_VERIFICATION_CODE_CACHE_EXPIRATION).Err()
	if err != nil {
		utils.BizLogger(c).Errorf("邮箱验证码写入缓存失败: %v", err)
		return err
	}

	// 发送验证码邮件
	expirationInMinutes := int(EMAIL_VERIFICATION_CODE_CACHE_EXPIRATION.Round(time.Minute).Minutes())
	emailContent := fmt.Sprintf("您的注册验证码是: %d , 有效期为 %d 分钟。", code, expirationInMinutes)
	success, err := utils.SendEmail(emailContent, []string{req.Email})
	if !success {
		utils.BizLogger(c).Errorf("邮箱验证码发送失败，邮箱地址: %s, 错误: %v", req.Email, err)
		global.RedisClient.Del(context.Background(), key)
		return bizErr.New(bizErr.SYSTEM_ERROR, err.Error())
	}

	return nil
}
