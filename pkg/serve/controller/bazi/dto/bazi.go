// Package dto 提供八字分析相关的数据传输对象
// 创建者：Done-0
// 创建时间：2023-10-18
package dto

import (
	"time"
)

// CreateOneBaziAnalysisRequest 八字分析请求参数
type CreateOneBaziAnalysisRequest struct {
	Name      string    `json:"name" form:"name" query:"name" binding:"required"`                               // 姓名
	Gender    string    `json:"gender" form:"gender" query:"gender" binding:"required"`                         // 性别
	BirthTime time.Time `json:"birth_time" form:"birth_time" query:"birth_time" binding:"required"`             // 出生时间
	Calendar  string    `json:"calendar" form:"calendar" query:"calendar" binding:"required,oneof=lunar solar"` // 日历类型 (lunar/solar)
}

// GetOneBaziRecordRequest 八字记录请求参数
type GetOneBaziRecordRequest struct {
	RequestID string `json:"request_id" form:"request_id" query:"request_id" binding:"required"` // 请求 ID
}
