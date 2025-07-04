// Package dto 提供八字分析相关的数据传输对象
// 创建者：Done-0
// 创建时间：2023-10-18
package dto

import (
	"time"
)

// CalculateBaziRequest 八字计算请求参数
type CalculateBaziRequest struct {
	Name      string    `json:"name" form:"name" query:"name" binding:"required"`                               // 姓名
	Gender    string    `json:"gender" form:"gender" query:"gender" binding:"required"`                         // 性别
	BirthTime time.Time `json:"birth_time" form:"birth_time" query:"birth_time" binding:"required"`             // 出生时间
	Calendar  string    `json:"calendar" form:"calendar" query:"calendar" binding:"required,oneof=lunar solar"` // 日历类型 (lunar/solar)
}

// GetOneBaziRequest 获取八字请求参数
type GetOneBaziRequest struct {
	ID int64 `json:"id,string" form:"id" query:"id" binding:"required"` // 八字 ID
}

// GetBaziListRequest 获取八字列表请求参数
type GetBaziListRequest struct {
	PageNo   int `json:"page_no" form:"page_no" query:"page_no"`       // 页码
	PageSize int `json:"page_size" form:"page_size" query:"page_size"` // 每页数量
}
