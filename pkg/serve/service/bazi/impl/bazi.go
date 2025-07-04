// Package impl 提供八字分析相关的服务层实现
// 创建者：Done-0
// 创建时间：2023-10-18
package impl

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/internal/model/bazi"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/bazi/dto"
	baziMapper "github.com/Done-0/metaphysics/pkg/serve/mapper/bazi"
	baziSrv "github.com/Done-0/metaphysics/pkg/serve/service/bazi"
	baziVO "github.com/Done-0/metaphysics/pkg/vo/bazi"
)

// BaziServiceImpl 八字服务实现
type BaziServiceImpl struct {
	baziMapper baziMapper.BaziMapper
}

// NewBaziService 创建八字服务实例
// 参数：
//   - mapper: 八字数据访问接口
//
// 返回值：
//   - baziSrv.BaziService: 八字服务接口
func NewBaziService(mapper baziMapper.BaziMapper) baziSrv.BaziService {
	return &BaziServiceImpl{
		baziMapper: mapper,
	}
}

// CalculateOneBazi 生成八字分析
// 参数：
//
//	ctx: 上下文信息
//	req: 请求参数
//
// 返回值：
//
//	*baziVO.BaziResponse: 八字分析结果
//	error: 错误信息
func (b *BaziServiceImpl) CalculateOneBazi(ctx *gin.Context, req *dto.CalculateBaziRequest) (*baziVO.BaziResponse, error) {
	baziInfo := utils.CalculateBazi(req.BirthTime, req.Calendar)

	bazi := &bazi.Bazi{
		UserID:      0,
		Name:        req.Name,
		Gender:      req.Gender,
		BirthTime:   req.BirthTime,
		Calendar:    req.Calendar,
		YearPillar:  baziInfo["year"],
		MonthPillar: baziInfo["month"],
		DayPillar:   baziInfo["day"],
		HourPillar:  baziInfo["hour"],
		YearGan:     baziInfo["year_gan"],
		YearZhi:     baziInfo["year_zhi"],
		MonthGan:    baziInfo["month_gan"],
		MonthZhi:    baziInfo["month_zhi"],
		DayGan:      baziInfo["day_gan"],
		DayZhi:      baziInfo["day_zhi"],
		HourGan:     baziInfo["hour_gan"],
		HourZhi:     baziInfo["hour_zhi"],
	}

	if err := b.baziMapper.CreateOneBazi(ctx, bazi); err != nil {
		utils.BizLogger(ctx).Errorf("存储八字失败: %v", err)
		return nil, fmt.Errorf("存储八字失败: %w", err)
	}

	vo, err := utils.MapModelToVO(bazi, &baziVO.BaziResponse{})
	if err != nil {
		utils.BizLogger(ctx).Errorf("分析八字结果时映射 VO 失败: %v", err)
		return nil, fmt.Errorf("分析八字结果时映射 VO 失败: %w", err)
	}

	return vo.(*baziVO.BaziResponse), nil
}

// GetOneBazi 获取八字
// 参数：
//
//	ctx: 上下文信息
//	req: 请求参数
//
// 返回值：
//
//	*baziVO.BaziResponse: 八字视图对象
//	error: 错误信息
func (b *BaziServiceImpl) GetOneBazi(ctx *gin.Context, req *dto.GetOneBaziRequest) (*baziVO.BaziResponse, error) {
	bazi, err := b.baziMapper.GetOneBaziByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	getBaziVO, err := utils.MapModelToVO(bazi, &baziVO.BaziResponse{})
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取八字时映射 VO 失败: %v", err)
		return nil, fmt.Errorf("获取八字时映射 VO 失败: %w", err)
	}

	return getBaziVO.(*baziVO.BaziResponse), nil
}

// GetBaziList 获取八字列表
// 参数：
//
//	ctx: 上下文信息
//	req: 请求参数
//
// 返回值：
//
//	*baziVO.BaziListResponse: 八字列表视图对象
//	error: 错误信息
func (b *BaziServiceImpl) GetBaziList(ctx *gin.Context, req *dto.GetBaziListRequest) (*baziVO.BaziListResponse, error) {
	bazis, total, err := b.baziMapper.GetBaziList(ctx, req.PageNo, req.PageSize)
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取八字列表失败: %v", err)
		return nil, fmt.Errorf("获取八字列表失败: %w", err)
	}

	baziVOList := make([]*baziVO.BaziResponse, 0, len(bazis))
	for _, item := range bazis {
		vo, err := utils.MapModelToVO(item, &baziVO.BaziResponse{})
		if err != nil {
			utils.BizLogger(ctx).Errorf("获取八字列表时映射单个VO失败: %v", err)
			continue
		}
		baziVOList = append(baziVOList, vo.(*baziVO.BaziResponse))
	}

	// 构造分页响应
	result := &baziVO.BaziListResponse{
		Total:    total,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		List:     baziVOList,
	}

	return result, nil
}
