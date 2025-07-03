// Package dto 提供用户相关的数据传输对象定义
// 创建者：Done-0
// 创建时间：2025-05-10
package dto

// RegisterOneAccountRequest 注册账户请求体
// @Description 请求注册账户时所需参数
// @Property    Email      body string true "用户邮箱"
// @Property    Password   body string true "用户密码"
// @Property    Nickname   body string true "用户昵称"
// @Property    EmailVerificationCode body string true "邮箱验证码"
type RegisterOneAccountRequest struct {
	Email                 string `json:"email" form:"email" validate:"required,email"`
	Password              string `json:"password" form:"password" validate:"required,min=6,max=20"`
	Nickname              string `json:"nickname" form:"nickname" validate:"required,min=2,max=20"`
	EmailVerificationCode string `json:"email_verification_code" form:"verification_code" validate:"required,len=6"`
}

// LoginOneAccountRequest 登录账户请求体
// @Description 请求登录账户时所需参数
// @Property    Email    body string true "用户邮箱"
// @Property    Password body string true "用户密码"
type LoginOneAccountRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=20"`
}

// GetOneAccountRequest 获取账户信息请求体
// @Description 请求获取账户信息时所需参数
// @Property    Email body string true "用户邮箱"
type GetOneAccountRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

// UpdateOneAccountRequest 更新账户信息请求体
// @Description 请求更新账户信息时所需参数
// @Property    Nickname body string true "用户昵称"
// @Property    Avatar   body string true "用户头像"
type UpdateOneAccountRequest struct {
	Nickname string `json:"nickname" form:"nickname" validate:"required,min=2,max=20"`
	Avatar   string `json:"avatar" form:"avatar" validate:"omitempty,url"`
}

// ResetPwdRequest 重置密码请求体
// @Description 请求重置密码时所需参数
// @Property    Email             body string true "用户邮箱"
// @Property    NewPassword       body string true "新密码"
// @Property    AgainNewPassword  body string true "确认新密码"
// @Property    EmailVerificationCode body string true "邮箱验证码"
type ResetPwdRequest struct {
	Email                 string `json:"email" form:"email" validate:"required,email"`
	NewPassword           string `json:"new_password" form:"new_password" validate:"required,min=6,max=20"`
	AgainNewPassword      string `json:"again_new_password" form:"again_new_password" validate:"required,eqfield=NewPassword"`
	EmailVerificationCode string `json:"email_verification_code" form:"email_verification_code" validate:"required,len=6"`
}
