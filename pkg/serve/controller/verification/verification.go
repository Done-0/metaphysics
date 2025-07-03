// Package verification 提供验证码相关的HTTP接口处理
// 创建者：Done-0
// 创建时间：2025-05-10
package verification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bizErr "github.com/Done-0/metaphysics/internal/error"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/verification/dto"
	verificationService "github.com/Done-0/metaphysics/pkg/serve/service/verification"
	"github.com/Done-0/metaphysics/pkg/vo"
)

// VerificationController 验证码控制器
type VerificationController struct {
	verificationService verificationService.VerificationService
}

// NewVerificationController 创建验证码控制器
// 参数：
//   - verificationService: 验证码服务
//
// 返回值：
//   - *VerificationController: 验证码控制器
func NewVerificationController(verificationService verificationService.VerificationService) *VerificationController {
	return &VerificationController{
		verificationService: verificationService,
	}
}

// SendEmailVerificationCode godoc
// @Summary 	发送邮箱验证码
// @Description 向指定邮箱发送验证码，验证码有效期为3分钟
// @Tags 		用户
// @Accept 		json
// @Produce 	json
// @Param 		email  query	string  true  "邮箱地址，用于发送验证码"
// @Success 	200    {object} vo.Result "邮箱验证码发送成功, 请注意查收邮件"
// @Failure 	400    {object} vo.Result "请求参数错误，邮箱地址为空"
// @Failure 	500    {object} vo.Result "服务器错误，邮箱验证码发送失败"
// @Router /api/v1/verification/sendEmailCode [get]
func (c *VerificationController) SendEmailVerificationCode(ctx *gin.Context) {
	req := new(dto.GetOneVerificationCode)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	err := c.verificationService.SendEmailVerificationCode(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, "邮箱验证码发送成功, 请注意查收！"))
}
