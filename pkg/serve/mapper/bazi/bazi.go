// Package bazi 提供八字相关的数据访问接口
// 创建者：Done-0
// 创建时间：2023-10-18
package bazi

import (
	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/internal/model/bazi"
)

// BaziMapper 八字数据访问接口
type BaziMapper interface {
	// CreateOneBazi 在事务中创建八字
	// 参数：
	//   - ctx: Gin上下文
	//   - bazi: 八字记录
	// 返回值：
	//   - error: 操作过程中的错误
	CreateOneBazi(ctx *gin.Context, bazi *bazi.Bazi) error

	// GetOneBaziByID 根据 ID 获取八字记录
	// 参数：
	//   - ctx: 上下文信息
	//   - id: 八字记录ID
	// 返回值：
	//   - *bazi.Bazi: 八字记录
	//   - error: 错误信息
	GetOneBaziByID(ctx *gin.Context, id int64) (*bazi.Bazi, error)

	// GetBaziRecordList 获取八字记录列表
	// 参数：
	//   - ctx: 上下文信息
	//   - pageNo: 页码
	//   - pageSize: 每页数量
	// 返回值：
	//   - []*bazi.Bazi: 八字记录列表
	//   - int64: 总记录数
	//   - error: 错误信息
	GetBaziList(ctx *gin.Context, pageNo, pageSize int) ([]*bazi.Bazi, int64, error)
}
