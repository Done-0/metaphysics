// Package bazi 提供八字相关的服务层功能
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
	// CalculateOneBazi 计算八字
	// 参数：
	//   - ctx: 上下文信息
	//   - req: 请求参数
	// 返回值：
	//   - *baziVO.BaziResponse: 八字计算结果
	//   - error: 错误信息
	CalculateOneBazi(ctx *gin.Context, req *dto.CalculateBaziRequest) (*baziVO.BaziResponse, error)

	// GetOneBazi 获取八字
	// 参数：
	//   - ctx: 上下文信息
	// 返回值：
	//   - *baziVO.BaziResponse: 八字视图对象
	//   - error: 错误信息
	GetOneBazi(ctx *gin.Context, req *dto.GetOneBaziRequest) (*baziVO.BaziResponse, error)

	// GetBaziList 获取八字列表
	// 参数：
	//   - ctx: 上下文信息
	//   - req: 请求参数
	// 返回值：
	//   - *baziVO.BaziListResponse: 八字列表视图对象
	//   - error: 错误信息
	GetBaziList(ctx *gin.Context, req *dto.GetBaziListRequest) (*baziVO.BaziListResponse, error)
}
