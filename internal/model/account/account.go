// Package model 提供用户账户数据模型定义
// 创建者：Done-0
// 创建时间：2025-05-10
package account

import "github.com/Done-0/metaphysics/internal/model/base"

// Account 用户账户模型
type Account struct {
	base.Base
	Email    string `gorm:"type:varchar(64);unique;not null" json:"email"` // 邮箱，主登录方式
	Password string `gorm:"type:varchar(255);not null" json:"password"`    // 加密密码
	Nickname string `gorm:"type:varchar(64);not null" json:"nickname"`     // 昵称
	Avatar   string `gorm:"type:varchar(255);default:null" json:"avatar"`  // 用户头像
}

// TableName 指定表名
// 返回值：
//   - string: 表名
func (Account) TableName() string {
	return "accounts"
}
