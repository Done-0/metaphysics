// Package verification 提供验证码相关的服务层接口
// 创建者：Done-0
// 创建时间：2025-05-10
package verification

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/pkg/serve/controller/verification/dto"
)

// VerificationService 验证码服务接口
type VerificationService interface {
	// SendEmailVerificationCode 发送邮箱验证码
	// 参数：
	//   - c: Gin上下文
	//   - req: 获取验证码请求参数
	// 返回值：
	//   - error: 错误信息
	SendEmailVerificationCode(c *gin.Context, req *dto.GetOneVerificationCode) error
}
