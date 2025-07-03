// Package account 提供用户相关的HTTP接口处理
// 创建者：Done-0
// 创建时间：2025-05-10
package account

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bizErr "github.com/Done-0/metaphysics/internal/error"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/account/dto"
	accountService "github.com/Done-0/metaphysics/pkg/serve/service/account"
	"github.com/Done-0/metaphysics/pkg/vo"
)

// AccountController 用户控制器
type AccountController struct {
	accountService accountService.AccountService
}

// NewAccountController 创建用户控制器
// 参数：
//   - accountService: 账户服务
//
// 返回值：
//   - *AccountController: 用户控制器
func NewAccountController(accountService accountService.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

// RegisterAcc godoc
// @Summary      用户注册
// @Description  注册新用户账号，支持邮箱验证码校验
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterOneAccountRequest  true  "注册信息"
// @Param        EmailVerificationCode  query   string  true  "邮箱验证码"
// @Success      200     {object}   vo.Result{data=dto.RegisterOneAccountRequest}  "注册成功"
// @Failure      400     {object}   vo.Result         "参数错误，验证码校验失败"
// @Failure      500     {object}   vo.Result         "服务器错误"
// @Router       /api/v1/account/register [post]
func (c *AccountController) RegisterAcc(ctx *gin.Context) {
	req := new(dto.RegisterOneAccountRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	if !utils.VerifyEmailCode(ctx, req.EmailVerificationCode, req.Email) {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "邮箱验证码校验失败")))
		return
	}

	acc, err := c.accountService.RegisterOneAccount(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, acc))
}

// LoginAccount godoc
// @Summary      用户登录
// @Description  用户登录并获取访问令牌
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginOneAccountRequest  true  "登录信息"
// @Success      200     {object}   vo.Result{data=account.LoginVO}  "登录成功，返回访问令牌"
// @Failure      400     {object}   vo.Result         "参数错误"
// @Failure      401     {object}   vo.Result         "登录失败，凭证无效"
// @Router       /api/v1/account/login [post]
func (c *AccountController) LoginAccount(ctx *gin.Context) {
	req := new(dto.LoginOneAccountRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	response, err := c.accountService.LoginOneAccount(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// GetAccount godoc
// @Summary      获取用户信息
// @Description  根据提供的邮箱获取对应用户的详细信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  query      dto.GetOneAccountRequest  true  "获取账户请求参数"
// @Success      200     {object}   vo.Result{data=account.GetAccountVO}  "获取成功"
// @Failure      400     {object}   vo.Result              "请求参数错误"
// @Failure      404     {object}   vo.Result              "用户不存在"
// @Router       /api/v1/account/info [get]
func (c *AccountController) GetAccount(ctx *gin.Context) {
	req := new(dto.GetOneAccountRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	response, err := c.accountService.GetOneAccount(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// UpdateAccount godoc
// @Summary      更新用户信息
// @Description  更新当前登录用户的账户信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UpdateOneAccountRequest  true  "更新账户信息"
// @Success      200     {object}   vo.Result{data=account.UpdateAccountVO}  "更新成功"
// @Failure      400     {object}   vo.Result              "请求参数错误"
// @Failure      401     {object}   vo.Result              "未授权"
// @Failure      500     {object}   vo.Result              "服务器错误"
// @Security     BearerAuth
// @Router       /api/v1/account/update [post]
func (c *AccountController) UpdateAccount(ctx *gin.Context) {
	req := new(dto.UpdateOneAccountRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	response, err := c.accountService.UpdateOneAccountByID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// LogoutAccount godoc
// @Summary      用户登出
// @Description  退出当前用户登录状态
// @Tags         用户
// @Produce      json
// @Success      200  {object}  vo.Result{data=string}  "登出成功"
// @Failure      401  {object}  vo.Result  "未授权"
// @Failure      500  {object}  vo.Result  "服务器错误"
// @Security     BearerAuth
// @Router       /api/v1/account/logout [post]
func (c *AccountController) LogoutAccount(ctx *gin.Context) {
	if err := c.accountService.LogoutOneAccount(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, "用户注销成功"))
}

// ResetPassword godoc
// @Summary      重置密码
// @Description  重置用户账户密码，支持邮箱验证码校验
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ResetPwdRequest  true  "重置密码信息"
// @Success      200     {object}   vo.Result{data=string}  "密码重置成功"
// @Failure      400     {object}   vo.Result         "参数错误，验证码校验失败"
// @Failure      401     {object}   vo.Result         "未授权，用户未登录"
// @Failure      500     {object}   vo.Result         "服务器错误"
// @Security     BearerAuth
// @Router       /api/v1/account/resetPassword [post]
func (c *AccountController) ResetPassword(ctx *gin.Context) {
	req := new(dto.ResetPwdRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	if !utils.VerifyEmailCode(ctx, req.EmailVerificationCode, req.Email) {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "邮箱验证码校验失败")))
		return
	}

	err := c.accountService.ResetPassword(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, "密码重置成功"))
}
