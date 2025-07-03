// Package bazi 提供八字分析相关的控制器功能
// 创建者：Done-0
// 创建时间：2023-10-18
package bazi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bizErr "github.com/Done-0/metaphysics/internal/error"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/bazi/dto"
	baziSrv "github.com/Done-0/metaphysics/pkg/serve/service/bazi"
	"github.com/Done-0/metaphysics/pkg/vo"
)

// BaziController 八字控制器
type BaziController struct {
	baziService baziSrv.BaziService
}

// NewBaziController 创建八字控制器
// 参数：
//   - baziService: 八字服务
//
// 返回值：
//   - *BaziController: 八字控制器
func NewBaziController(baziService baziSrv.BaziService) *BaziController {
	return &BaziController{
		baziService: baziService,
	}
}

// AnalyzeOneBazi 分析八字
// @Summary 分析八字
// @Description 根据用户提供的出生日期时间分析八字
// @Tags 八字
// @Accept json
// @Produce json
// @Param request body dto.BaziRequest true "八字分析请求"
// @Success 200 {object} vo.Result{data=baziVo.BaziResponse} "成功"
// @Failure 400 {object} vo.Result "参数错误"
// @Failure 500 {object} vo.Result "服务器内部错误"
// @Router /api/v1/bazi/analyze [post]
func (c *BaziController) AnalyzeOneBazi(ctx *gin.Context) {
	req := new(dto.CreateOneBaziAnalysisRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	validationErrors := utils.Validator(req)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, validationErrors, bizErr.New(bizErr.PARAM_ERROR)))
		return
	}

	response, err := c.baziService.CreateOneBaziAnalysis(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// StreamAnalyzeOneBazi 流式分析八字
// @Summary 流式分析八字
// @Description 使用Server-Sent Events流式返回八字分析结果
// @Tags 八字
// @Accept json
// @Produce text/event-stream
// @Param request body dto.BaziRequest true "八字分析请求"
// @Success 200 {object} string "SSE流"
// @Failure 400 {object} vo.Result "参数错误"
// @Failure 500 {object} vo.Result "服务器内部错误"
// @Router /api/v1/bazi/analyze/stream [post]
func (c *BaziController) StreamAnalyzeOneBazi(ctx *gin.Context) {
	req := new(dto.CreateOneBaziAnalysisRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	validationErrors := utils.Validator(req)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, validationErrors, bizErr.New(bizErr.PARAM_ERROR)))
		return
	}

	// 设置 SSE 响应头
	utils.SetupSSEHeaders(ctx)

	err := c.baziService.StreamCreateOneBaziAnalysis(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}
}

// GetBaziRecord 获取八字记录
// @Summary 获取八字记录
// @Description 根据请求ID获取八字记录
// @Tags 八字
// @Accept json
// @Produce json
// @Param request_id path string true "请求ID"
// @Success 200 {object} vo.Result{data=baziVo.BaziResponse} "成功"
// @Failure 400 {object} vo.Result "参数错误"
// @Failure 404 {object} vo.Result "记录不存在"
// @Failure 500 {object} vo.Result "服务器内部错误"
// @Router /api/v1/bazi/record/{request_id} [get]
func (c *BaziController) GetBaziRecord(ctx *gin.Context) {
	req := new(dto.GetOneBaziRecordRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	validationErrors := utils.Validator(req)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, validationErrors, bizErr.New(bizErr.PARAM_ERROR)))
		return
	}

	response, err := c.baziService.GetOneBaziRecord(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}
