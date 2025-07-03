// Package bazi 提供八字分析相关的服务层功能
// 创建者：Done-0
// 创建时间：2023-10-18
package bazi

import (
	"github.com/Done-0/metaphysics/pkg/serve/controller/bazi/dto"
	baziVO "github.com/Done-0/metaphysics/pkg/vo/bazi"
	"github.com/gin-gonic/gin"
)

// BaziService 八字服务接口
type BaziService interface {
	// CreateBaziAnalysis 生成八字分析
	// 参数：
	//   - ctx: 上下文信息
	//   - req: 请求参数
	// 返回值：
	//   - *baziVO.BaziResponse: 八字分析结果
	//   - error: 错误信息
	CreateOneBaziAnalysis(ctx *gin.Context, req *dto.CreateOneBaziAnalysisRequest) (*baziVO.BaziResponse, error)

	// GetOneBaziRecord 获取八字记录
	// 参数：
	//   - ctx: 上下文信息
	// 返回值：
	//   - *baziVO.BaziResponse: 八字记录视图对象
	//   - error: 错误信息
	GetOneBaziRecord(ctx *gin.Context, req *dto.GetOneBaziRecordRequest) (*baziVO.BaziResponse, error)

	// StreamCreateOneBaziAnalysis 流式生成八字分析
	// 参数：
	//   - ctx: 上下文信息
	//   - req: 请求参数
	// 返回值：
	//   - error: 错误信息
	StreamCreateOneBaziAnalysis(ctx *gin.Context, req *dto.CreateOneBaziAnalysisRequest) error
}
