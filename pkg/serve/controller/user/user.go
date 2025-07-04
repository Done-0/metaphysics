// Package user 提供用户相关的HTTP接口处理
// 创建者：Done-0
// 创建时间：2025-05-10
package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bizErr "github.com/Done-0/metaphysics/internal/error"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/user/dto"
	userSrv "github.com/Done-0/metaphysics/pkg/serve/service/user"
	"github.com/Done-0/metaphysics/pkg/vo"
)

// UserController 用户控制器
type UserController struct {
	userService userSrv.UserService
}

// NewUserController 创建用户控制器
// 参数：
//   - userService: 用户服务
//
// 返回值：
//   - *UserController: 用户控制器
func NewUserController(userService userSrv.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// RegisterOneUser godoc
// @Summary      用户注册
// @Description  注册新用户账号，支持邮箱验证码校验
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterOneUserRequest  true  "注册信息"
// @Param        EmailVerificationCode  query   string  true  "邮箱验证码"
// @Success      200     {object}   vo.Result{data=dto.RegisterOneUserRequest}  "注册成功"
// @Failure      400     {object}   vo.Result         "参数错误，验证码校验失败"
// @Failure      500     {object}   vo.Result         "服务器错误"
// @Router       /api/v1/user/register [post]
func (c *UserController) RegisterOneUser(ctx *gin.Context) {
	req := new(dto.RegisterOneUserRequest)
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

	acc, err := c.userService.RegisterOneUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, acc))
}

// LoginOneUser godoc
// @Summary      用户登录
// @Description  用户登录并获取访问令牌
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginOneUserRequest  true  "登录信息"
// @Success      200     {object}   vo.Result{data=user.LoginOneUserResponse}  "登录成功，返回访问令牌"
// @Failure      400     {object}   vo.Result         "参数错误"
// @Failure      401     {object}   vo.Result         "登录失败，凭证无效"
// @Router       /api/v1/user/login [post]
func (c *UserController) LoginOneUser(ctx *gin.Context) {
	req := new(dto.LoginOneUserRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	response, err := c.userService.LoginOneUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// GetOneUser godoc
// @Summary      获取用户信息
// @Description  根据提供的邮箱获取对应用户的详细信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  query      dto.GetOneUserRequest  true  "获取用户请求参数"
// @Success      200     {object}   vo.Result{data=user.GetOneUserResponse}  "获取成功"
// @Failure      400     {object}   vo.Result              "请求参数错误"
// @Failure      404     {object}   vo.Result              "用户不存在"
// @Router       /api/v1/user/info [get]
func (c *UserController) GetOneUser(ctx *gin.Context) {
	req := new(dto.GetOneUserRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	response, err := c.userService.GetOneUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// UpdateOneUser godoc
// @Summary      更新用户信息
// @Description  更新当前登录用户的用户信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UpdateOneUserRequest  true  "更新用户信息"
// @Success      200     {object}   vo.Result{data=user.UpdateOneUserResponse}  "更新成功"
// @Failure      400     {object}   vo.Result              "请求参数错误"
// @Failure      401     {object}   vo.Result              "未授权"
// @Failure      500     {object}   vo.Result              "服务器错误"
// @Security     BearerAuth
// @Router       /api/v1/user/update [post]
func (c *UserController) UpdateOneUser(ctx *gin.Context) {
	req := new(dto.UpdateOneUserRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	errors := utils.Validator(req)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, errors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	response, err := c.userService.UpdateOneUserByID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// LogoutOneUser godoc
// @Summary      用户登出
// @Description  退出当前用户登录状态
// @Tags         用户
// @Produce      json
// @Success      200  {object}  vo.Result{data=string}  "登出成功"
// @Failure      401  {object}  vo.Result  "未授权"
// @Failure      500  {object}  vo.Result  "服务器错误"
// @Security     BearerAuth
// @Router       /api/v1/user/logout [post]
func (c *UserController) LogoutOneUser(ctx *gin.Context) {
	if err := c.userService.LogoutOneUser(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, "用户注销成功"))
}

// ResetPassword godoc
// @Summary      重置密码
// @Description  重置用户用户密码，支持邮箱验证码校验
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ResetPasswordRequest  true  "重置密码信息"
// @Success      200     {object}   vo.Result{data=string}  "密码重置成功"
// @Failure      400     {object}   vo.Result         "参数错误，验证码校验失败"
// @Failure      401     {object}   vo.Result         "未授权，用户未登录"
// @Failure      500     {object}   vo.Result         "服务器错误"
// @Security     BearerAuth
// @Router       /api/v1/user/resetPassword [post]
func (c *UserController) ResetPassword(ctx *gin.Context) {
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

	err := c.userService.ResetPassword(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, "密码重置成功"))
}
