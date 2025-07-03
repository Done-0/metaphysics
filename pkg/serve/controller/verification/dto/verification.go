// Package dto 提供验证码相关的数据传输对象定义
// 创建者：Done-0
// 创建时间：2025-05-10
package dto

// GetOneVerificationCode 获取验证码请求
// @Description 请求获取验证码时所需参数
// @Property    Email body string true "用户邮箱"
type GetOneVerificationCode struct {
	Email string `json:"email" form:"email" query:"email" uri:"email" validate:"required,email"`
}
