// Package model 提供应用程序的数据模型定义和聚合
// 创建者：Done-0
// 创建时间：2025-07-01
package model

import (
	"github.com/Done-0/metaphysics/internal/model/bazi"
	"github.com/Done-0/metaphysics/internal/model/user"
)

// GetAllModels 获取并注册所有模型
// 返回值：
//   - []any: 所有需要注册到数据库的模型列表
func GetAllModels() []any {
	return []any{
		&bazi.Bazi{},       // 八字模型
		&user.User{},       // 用户模型
	}
}
