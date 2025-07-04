// Package user 提供用户相关的视图对象定义
// 创建者：Done-0
// 创建时间：2025-05-10
package user

// GetOneUserResponse     获取用户信息请求体
// @Description	请求获取用户信息时所需参数
// @Property			Nickname	body	string	true	"用户昵称"
// @Property			Email	    body	string	true	"用户邮箱"
// @Property			Avatar	    body	string	true	"用户头像"
type GetOneUserResponse struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// LoginOneUserResponse           返回给前端的登录信息
// @Description	登录成功后返回的访问令牌和刷新令牌
// @Property			AccessToken 	body	string	true	"访问令牌"
// @Property			RefreshToken	body	string	true	"刷新令牌"
type LoginOneUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RegisterUserResponse     获取用户信息请求体
// @Description	请求获取用户信息时所需参数
// @Property			Nickname	body	string	true	"用户昵称"
// @Property			Email	    body	string	true	"用户邮箱"
type RegisterOneUserResponse struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// UpdateUserResponse     更新用户信息响应体
// @Description	更新用户信息后返回的参数
// @Property			Nickname	body	string	true	"用户昵称"
// @Property			Email	    body	string	true	"用户邮箱"
// @Property			Avatar	    body	string	true	"用户头像"
type UpdateOneUserResponse struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
