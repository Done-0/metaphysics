// Package mapper 提供数据模型与数据库交互的映射层，处理八字相关数据操作
// 创建者：Done-0
// 创建时间：2023-10-18
package mapper

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Done-0/metaphysics/internal/model/bazi"
	"github.com/Done-0/metaphysics/internal/utils"
)

// CreateBaziRecord 在事务中创建八字记录
// 参数：
//   - ctx: Gin上下文
//   - record: 八字记录
//
// 返回值：
//   - error: 操作过程中的错误
func CreateOneBaziRecord(ctx *gin.Context, record *bazi.BaziRecord) error {
	return utils.RunDBTransaction(ctx, func() error {
		db := utils.GetDBFromContext(ctx)
		if err := db.Create(record).Error; err != nil {
			return fmt.Errorf("保存八字记录失败: %w", err)
		}

		return nil
	})
}

// GetOneBaziRecordByRequestID 根据请求ID获取八字记录
// 参数：
//   - ctx: 上下文信息
//   - requestID: 请求ID
//
// 返回值：
//   - *bazi.BaziRecord: 八字记录
//   - error: 错误信息
func GetOneBaziRecordByRequestID(ctx *gin.Context, requestID string) (*bazi.BaziRecord, error) {
	var record bazi.BaziRecord
	db := utils.GetDBFromContext(ctx)
	err := db.Where("request_id = ? AND deleted = ?", requestID, false).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("八字记录不存在")
		}
		return nil, fmt.Errorf("查询八字记录失败: %w", err)
	}

	return &record, nil
}
