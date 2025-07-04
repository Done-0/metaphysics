// Package impl 提供八字相关的数据访问实现
// 创建者：Done-0
// 创建时间：2023-10-18
package impl

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Done-0/metaphysics/internal/model/bazi"
	"github.com/Done-0/metaphysics/internal/utils"
	baziMapper "github.com/Done-0/metaphysics/pkg/serve/mapper/bazi"
)

// BaziMapperImpl 八字数据访问实现
type BaziMapperImpl struct{}

// NewBaziMapper 创建八字数据访问实例
// 返回值：
//   - baziMapper.BaziMapper: 八字数据访问接口
func NewBaziMapper() baziMapper.BaziMapper {
	return &BaziMapperImpl{}
}

// CreateOneBazi 在事务中创建八字
// 参数：
//   - ctx: Gin上下文
//   - bazi: 八字记录
//
// 返回值：
//   - error: 操作过程中的错误
func (m *BaziMapperImpl) CreateOneBazi(ctx *gin.Context, bazi *bazi.Bazi) error {
	return utils.RunDBTransaction(ctx, func() error {
		db := utils.GetDBFromContext(ctx)
		if err := db.Create(bazi).Error; err != nil {
			return fmt.Errorf("保存八字失败: %w", err)
		}

		return nil
	})
}

// GetOneBaziByID 根据 ID 获取八字记录
// 参数：
//   - ctx: 上下文信息
//   - id: 八字记录ID
//
// 返回值：
//   - *bazi.Bazi: 八字记录
//   - error: 错误信息
func (m *BaziMapperImpl) GetOneBaziByID(ctx *gin.Context, id int64) (*bazi.Bazi, error) {
	var bazi bazi.Bazi
	db := utils.GetDBFromContext(ctx)
	err := db.Where("id = ? AND deleted = ?", id, false).First(&bazi).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("八字不存在")
		}
		return nil, fmt.Errorf("查询八字失败: %w", err)
	}

	return &bazi, nil
}

// GetBaziList 获取八字列表
// 参数：
//   - ctx: 上下文信息
//   - pageNo: 页码
//   - pageSize: 每页数量
//
// 返回值：
//   - []*bazi.Bazi: 八字记录列表
//   - int64: 总记录数
//   - error: 错误信息
func (m *BaziMapperImpl) GetBaziList(ctx *gin.Context, pageNo, pageSize int) ([]*bazi.Bazi, int64, error) {
	var bazis []*bazi.Bazi
	var total int64

	db := utils.GetDBFromContext(ctx)

	query := db.Model(&bazi.Bazi{}).Where("deleted = ?", false)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询八字总数失败: %w", err)
	}

	// 分页查询 - 倒序排序
	offset := (pageNo - 1) * pageSize
	if err := query.
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bazis).Error; err != nil {
		return nil, 0, fmt.Errorf("查询八字列表失败: %w", err)
	}

	return bazis, total, nil
}
