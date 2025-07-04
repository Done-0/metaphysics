// Package bazi 提供八字相关的控制器功能
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

// CalculateBazi 计算八字
// @Summary 计算八字
// @Description 根据用户提供的出生日期时间计算八字
// @Tags 八字
// @Accept json
// @Produce json
// @Param request body dto.BaziRequest true "八字计算请求"
// @Success 200 {object} vo.Result{data=baziVo.BaziResponse} "成功"
// @Failure 400 {object} vo.Result "参数错误"
// @Failure 500 {object} vo.Result "服务器内部错误"
// @Router /api/v1/bazi/calculate [post]
func (c *BaziController) CalculateOneBazi(ctx *gin.Context) {
	req := new(dto.CalculateBaziRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	validationErrors := utils.Validator(req)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, validationErrors, bizErr.New(bizErr.PARAM_ERROR)))
		return
	}

	response, err := c.baziService.CalculateOneBazi(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// GetBaziRecord 获取八字记录
// @Summary 获取八字记录
// @Description 根据ID获取八字记录
// @Tags 八字
// @Accept json
// @Produce json
// @Param id path int64 true "八字记录ID"
// @Success 200 {object} vo.Result{data=baziVo.BaziResponse} "成功"
// @Failure 400 {object} vo.Result "参数错误"
// @Failure 404 {object} vo.Result "记录不存在"
// @Failure 500 {object} vo.Result "服务器内部错误"
// @Router /api/v1/bazi/record/{id} [get]
func (c *BaziController) GetOneBazi(ctx *gin.Context) {
	req := new(dto.GetOneBaziRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	validationErrors := utils.Validator(req)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, validationErrors, bizErr.New(bizErr.PARAM_ERROR)))
		return
	}

	response, err := c.baziService.GetOneBazi(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}

// GetBaziRecordList 获取八字记录列表
// @Summary 获取八字记录列表
// @Description 分页获取八字记录列表
// @Tags 八字
// @Accept json
// @Produce json
// @Param page_no query int false "页码，默认为1"
// @Param page_size query int false "每页记录数，默认为10，最大为100"
// @Success 200 {object} vo.Result{data=baziVo.BaziRecordListResponse} "成功"
// @Failure 400 {object} vo.Result "参数错误"
// @Failure 500 {object} vo.Result "服务器内部错误"
// @Router /api/v1/bazi/records [get]
func (c *BaziController) GetBaziList(ctx *gin.Context) {
	req := new(dto.GetBaziListRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, err, bizErr.New(bizErr.PARAM_ERROR, err.Error())))
		return
	}

	validationErrors := utils.Validator(req)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, validationErrors, bizErr.New(bizErr.PARAM_ERROR)))
		return
	}

	response, err := c.baziService.GetBaziList(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, err, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, response))
}
