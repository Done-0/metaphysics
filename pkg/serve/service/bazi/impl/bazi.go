// Package impl 提供八字分析相关的服务层实现
// 创建者：Done-0
// 创建时间：2023-10-18
package impl

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/metaphysics/internal/ai"
	"github.com/Done-0/metaphysics/internal/ai/types"
	bizErr "github.com/Done-0/metaphysics/internal/error"
	"github.com/Done-0/metaphysics/internal/model/bazi"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/bazi/dto"
	"github.com/Done-0/metaphysics/pkg/serve/mapper"
	baziService "github.com/Done-0/metaphysics/pkg/serve/service/bazi"
	"github.com/Done-0/metaphysics/pkg/vo"
	aiVO "github.com/Done-0/metaphysics/pkg/vo/ai"
	baziVO "github.com/Done-0/metaphysics/pkg/vo/bazi"
)

// BaziServiceImpl 八字服务实现
type BaziServiceImpl struct {
	aiService types.Service
}

// NewBaziService 创建八字服务实例
// 返回值：
//
//	baziService.BaziService: 八字服务接口
func NewBaziService() baziService.BaziService {
	return &BaziServiceImpl{
		aiService: ai.New(),
	}
}

// CreateOneBaziAnalysis 生成八字分析
// 参数：
//
//	ctx: 上下文信息
//	req: 请求参数
//
// 返回值：
//
//	*baziVO.BaziResponse: 八字分析结果
//	error: 错误信息
func (b *BaziServiceImpl) CreateOneBaziAnalysis(ctx *gin.Context, req *dto.CreateOneBaziAnalysisRequest) (*baziVO.BaziResponse, error) {
	requestID := ctx.GetHeader("X-Request-ID")

	// 计算八字信息
	baziInfo := utils.CalculateBazi(req.BirthTime, req.Calendar)

	analysis, err := b.aiService.AnalyzeBazi(ctx, req.Name, req.Gender, req.BirthTime, req.Calendar, baziInfo)
	if err != nil {
		return nil, fmt.Errorf("八字 AI 分析失败: %w", err)
	}

	record := &bazi.BaziRecord{
		UserID:         0,
		Name:           req.Name,
		Gender:         req.Gender,
		BirthTime:      req.BirthTime,
		Calendar:       req.Calendar,
		YearPillar:     baziInfo["year"],
		MonthPillar:    baziInfo["month"],
		DayPillar:      baziInfo["day"],
		HourPillar:     baziInfo["hour"],
		YearGan:        baziInfo["year_gan"],
		YearZhi:        baziInfo["year_zhi"],
		MonthGan:       baziInfo["month_gan"],
		MonthZhi:       baziInfo["month_zhi"],
		DayGan:         baziInfo["day_gan"],
		DayZhi:         baziInfo["day_zhi"],
		HourGan:        baziInfo["hour_gan"],
		HourZhi:        baziInfo["hour_zhi"],
		YearGanWuXing:  baziInfo["year_gan_wu_xing"],
		MonthGanWuXing: baziInfo["month_gan_wu_xing"],
		DayGanWuXing:   baziInfo["day_gan_wu_xing"],
		HourGanWuXing:  baziInfo["hour_gan_wu_xing"],
		YearNaYin:      baziInfo["year_na_yin"],
		MonthNaYin:     baziInfo["month_na_yin"],
		DayNaYin:       baziInfo["day_na_yin"],
		TimeNaYin:      baziInfo["time_na_yin"],
		YinYang:        baziInfo["yin_yang"],
		WuXing:         baziInfo["wu_xing"],
		Analysis:       analysis,
		RequestID:      requestID,
	}

	if err := mapper.CreateOneBaziRecord(ctx, record); err != nil {
		utils.BizLogger(ctx).Errorf("存储八字记录失败: %v", err)
		return nil, fmt.Errorf("存储八字记录失败: %w", err)
	}

	vo, err := utils.MapModelToVO(record, &baziVO.BaziResponse{})
	if err != nil {
		utils.BizLogger(ctx).Errorf("分析八字结果时映射 VO 失败: %v", err)
		return nil, fmt.Errorf("分析八字结果时映射 VO 失败: %w", err)
	}

	return vo.(*baziVO.BaziResponse), nil
}

// StreamCreateOneBaziAnalysis 流式生成八字分析
// 参数：
//
//	ctx: Gin 上下文
//	req: 请求参数
//
// 返回值：
//
//	error: 错误信息
func (b *BaziServiceImpl) StreamCreateOneBaziAnalysis(ctx *gin.Context, req *dto.CreateOneBaziAnalysisRequest) error {
	requestID := ctx.GetHeader("X-Request-ID")

	baziInfo := utils.CalculateBazi(req.BirthTime, req.Calendar)

	record := &bazi.BaziRecord{
		UserID:         0,
		Name:           req.Name,
		Gender:         req.Gender,
		BirthTime:      req.BirthTime,
		Calendar:       req.Calendar,
		YearPillar:     baziInfo["year"],
		MonthPillar:    baziInfo["month"],
		DayPillar:      baziInfo["day"],
		HourPillar:     baziInfo["hour"],
		YearGan:        baziInfo["year_gan"],
		YearZhi:        baziInfo["year_zhi"],
		MonthGan:       baziInfo["month_gan"],
		MonthZhi:       baziInfo["month_zhi"],
		DayGan:         baziInfo["day_gan"],
		DayZhi:         baziInfo["day_zhi"],
		HourGan:        baziInfo["hour_gan"],
		HourZhi:        baziInfo["hour_zhi"],
		YearGanWuXing:  baziInfo["year_gan_wu_xing"],
		MonthGanWuXing: baziInfo["month_gan_wu_xing"],
		DayGanWuXing:   baziInfo["day_gan_wu_xing"],
		HourGanWuXing:  baziInfo["hour_gan_wu_xing"],
		YearNaYin:      baziInfo["year_na_yin"],
		MonthNaYin:     baziInfo["month_na_yin"],
		DayNaYin:       baziInfo["day_na_yin"],
		TimeNaYin:      baziInfo["time_na_yin"],
		YinYang:        baziInfo["yin_yang"],
		WuXing:         baziInfo["wu_xing"],
		RequestID:      requestID,
	}

	// 发送连接成功消息
	utils.SendSSEEvent(ctx, "", vo.Success(ctx, map[string]string{"status": "connected"}))

	// 保存完整分析结果
	var fullAnalysis string

	// 处理流式响应
	handler := func(chunk *aiVO.StreamChunk) error {
		if chunk.Done {
			record.Analysis = fullAnalysis
			if err := mapper.CreateOneBaziRecord(ctx, record); err != nil {
				utils.BizLogger(ctx).Errorf("存储八字记录失败: %v", err)
				return fmt.Errorf("存储八字记录失败: %w", err)
			}

			// 发送完成事件
			utils.SendSSEEvent(ctx, utils.SSE_EVENT_DONE, vo.Success(ctx, map[string]string{"request_id": requestID}))
			return nil
		}

		utils.SendSSEEvent(ctx, utils.SSE_EVENT_REASONING, vo.Success(ctx, chunk.ReasoningContent))
		fullAnalysis += chunk.Content
		utils.SendSSEEvent(ctx, utils.SSE_EVENT_CONTENT, vo.Success(ctx, fullAnalysis))

		return nil
	}

	if err := b.aiService.StreamAnalyzeBazi(ctx, req.Name, req.Gender, req.BirthTime, req.Calendar, baziInfo, handler); err != nil {
		sysError := bizErr.New(bizErr.SYSTEM_ERROR, err.Error())
		utils.SendSSEEvent(ctx, utils.SSE_EVENT_ERROR, vo.Fail(ctx, err.Error(), sysError))
		return fmt.Errorf("八字分析失败: %s", err.Error())
	}

	return nil
}

// GetOneBaziRecord 获取八字记录
// 参数：
//
//	ctx: 上下文信息
//	req: 请求参数
//
// 返回值：
//
//	*baziVO.BaziResponse: 八字记录视图对象
//	error: 错误信息
func (b *BaziServiceImpl) GetOneBaziRecord(ctx *gin.Context, req *dto.GetOneBaziRecordRequest) (*baziVO.BaziResponse, error) {
	record, err := mapper.GetOneBaziRecordByRequestID(ctx, req.RequestID)
	if err != nil {
		return nil, err
	}

	vo, err := utils.MapModelToVO(record, &baziVO.BaziResponse{})
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取八字记录时映射 VO 失败: %v", err)
		return nil, fmt.Errorf("获取八字记录时映射 VO 失败: %w", err)
	}

	return vo.(*baziVO.BaziResponse), nil
}
