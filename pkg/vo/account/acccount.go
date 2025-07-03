// Package account 提供账户相关的视图对象定义
// 创建者：Done-0
// 创建时间：2025-05-10
package account

// GetAccountResponse     获取账户信息请求体
// @Description	请求获取账户信息时所需参数
// @Property			Nickname	body	string	true	"用户昵称"
// @Property			Email	    body	string	true	"用户邮箱"
// @Property			Avatar	    body	string	true	"用户头像"
type GetAccountResponse struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// LoginResponse           返回给前端的登录信息
// @Description	登录成功后返回的访问令牌和刷新令牌
// @Property			AccessToken 	body	string	true	"访问令牌"
// @Property			RefreshToken	body	string	true	"刷新令牌"
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RegisterAccountResponse     获取账户信息请求体
// @Description	请求获取账户信息时所需参数
// @Property			Nickname	body	string	true	"用户昵称"
// @Property			Email	    body	string	true	"用户邮箱"
type RegisterAccountResponse struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// UpdateAccountResponse     更新账户信息响应体
// @Description	更新账户信息后返回的参数
// @Property			Nickname	body	string	true	"用户昵称"
// @Property			Email	    body	string	true	"用户邮箱"
// @Property			Avatar	    body	string	true	"用户头像"
type UpdateAccountResponse struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
