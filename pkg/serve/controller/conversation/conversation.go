// Package conversation 提供对话相关的控制器功能
// 创建者：Done-0
// 创建时间：2025-07-03
package conversation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	bizErr "github.com/Done-0/metaphysics/internal/error"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/conversation/dto"
	conversationSrv "github.com/Done-0/metaphysics/pkg/serve/service/conversation"
	"github.com/Done-0/metaphysics/pkg/vo"
	"github.com/Done-0/metaphysics/pkg/vo/conversation"
)

// ConversationController 对话控制器
type ConversationController struct {
	conversationService conversationSrv.ConversationService
}

// NewConversationController 创建对话控制器
// 参数：
//   - conversationService: 对话服务
//
// 返回值：
//   - *ConversationController: 对话控制器
func NewConversationController(conversationService conversationSrv.ConversationService) *ConversationController {
	return &ConversationController{
		conversationService: conversationService,
	}
}

// AnalyzeBazi godoc
// @Summary      分析八字
// @Description  根据用户ID分析八字
// @Tags         对话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  vo.Result{data=interface{}}  "成功"
// @Failure      400  {object}  vo.Result                   "参数错误"
// @Failure      500  {object}  vo.Result                   "服务器内部错误"
// @Router       /api/v1/conversation/bazi/analyze [get]
func (c *ConversationController) AnalyzeBazi(ctx *gin.Context) {
	result, err := c.conversationService.AnalyzeBaziByUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, nil, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, result))
}

// StreamAnalyzeBazi godoc
// @Summary      流式分析八字
// @Description  流式分析用户八字
// @Tags         对话
// @Accept       json
// @Produce      text/event-stream
// @Security     BearerAuth
// @Success      200  {string}  string           "事件流"
// @Failure      400  {object}  vo.Result        "参数错误"
// @Failure      500  {object}  vo.Result        "服务器内部错误"
// @Router       /api/v1/conversation/bazi/analyze/stream [get]
func (c *ConversationController) StreamAnalyzeBazi(ctx *gin.Context) {
	// 设置响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	// 获取请求和响应消息ID
	requestID, responseID, err := c.conversationService.GetMessageIDs(ctx)
	if err != nil {
		// 使用默认值
		requestID = 1
		responseID = 2
		utils.BizLogger(ctx).Errorf("获取消息ID失败: %v", err)
	}

	// 发送ready事件
	readyData := new(conversation.ReadyEventData)
	readyData.RequestMessageID = requestID
	readyData.ResponseMessageID = responseID
	readyJSON, _ := json.Marshal(readyData)
	ctx.SSEvent(string(conversation.EventReady), string(readyJSON))
	ctx.Writer.Flush()

	// 发送update_session事件
	updateSessionData := new(conversation.UpdateSessionData)
	updateSessionData.UpdatedAt = float64(time.Now().Unix()) + float64(time.Now().Nanosecond())/1e9
	updateSessionJSON, _ := json.Marshal(updateSessionData)
	ctx.SSEvent(string(conversation.EventUpdateSession), string(updateSessionJSON))
	ctx.Writer.Flush()

	// 创建初始响应
	initialResponse := new(conversation.ResponseValue)
	initialResponse.Response.MessageID = responseID
	initialResponse.Response.ParentID = requestID
	initialResponse.Response.Model = ""
	initialResponse.Response.Role = "ASSISTANT"
	initialResponse.Response.Content = ""
	initialResponse.Response.ThinkingEnabled = true
	initialResponse.Response.ThinkingContent = nil
	initialResponse.Response.ThinkingElapsedSecs = nil
	initialResponse.Response.BanEdit = false
	initialResponse.Response.BanRegenerate = false
	initialResponse.Response.Status = conversation.StatusWIP
	initialResponse.Response.AccumulatedTokenUsage = 0
	initialResponse.Response.Files = []interface{}{}
	initialResponse.Response.Tips = []interface{}{}
	initialResponse.Response.InsertedAt = float64(time.Now().Unix()) + float64(time.Now().Nanosecond())/1e9
	initialResponse.Response.SearchEnabled = true
	initialResponse.Response.SearchStatus = conversation.StatusInit
	initialResponse.Response.SearchResults = nil

	initialResponseJSON, _ := json.Marshal(initialResponse)
	ctx.SSEvent("", string(initialResponseJSON))
	ctx.Writer.Flush()

	// 发送搜索状态更新
	searchStatusData := new(conversation.StreamData)
	searchStatusData.V = string(conversation.StatusAnswer)
	searchStatusData.P = "response/search_status"
	searchStatusJSON, _ := json.Marshal(searchStatusData)
	ctx.SSEvent("", string(searchStatusJSON))
	ctx.Writer.Flush()

	// 开始思考时间
	thinkingStartTime := time.Now()

	// 思考内容
	thinkingContent := "分析用户八字..."
	thinkingContentData := new(conversation.StreamData)
	thinkingContentData.V = thinkingContent
	thinkingContentData.P = "response/thinking_content"
	thinkingContentJSON, _ := json.Marshal(thinkingContentData)
	ctx.SSEvent("", string(thinkingContentJSON))
	ctx.Writer.Flush()

	// 流式分析八字
	var fullContent string
	var tokenCount int = 0

	err = c.conversationService.StreamAnalyzeBaziByUserID(ctx, func(content string, done bool) error {
		if !done {
			// 累积完整内容
			fullContent += content
			tokenCount += len(content) / 4 // 粗略估算token数量

			// 发送内容更新
			contentData := new(conversation.StreamData)
			contentData.V = content
			contentData.O = conversation.OperationAppend
			if len(fullContent) == len(content) { // 第一次发送
				contentData.P = "response/content"
			}
			contentJSON, _ := json.Marshal(contentData)
			ctx.SSEvent("", string(contentJSON))
			ctx.Writer.Flush()
		} else {
			// 计算思考时间
			thinkingElapsedSecs := int(time.Since(thinkingStartTime).Seconds())
			thinkingElapsedData := new(conversation.StreamData)
			thinkingElapsedData.V = thinkingElapsedSecs
			thinkingElapsedData.P = "response/thinking_elapsed_secs"
			thinkingElapsedData.O = conversation.OperationSet
			thinkingElapsedJSON, _ := json.Marshal(thinkingElapsedData)
			ctx.SSEvent("", string(thinkingElapsedJSON))
			ctx.Writer.Flush()

			// 发送完成状态
			batchValues := []conversation.BatchValue{
				{V: conversation.StatusFinished, P: "status"},
				{V: tokenCount, P: "accumulated_token_usage"},
			}
			batchData := new(conversation.StreamData)
			batchData.V = batchValues
			batchData.P = "response"
			batchData.O = conversation.OperationBatch
			batchJSON, _ := json.Marshal(batchData)
			ctx.SSEvent("", string(batchJSON))
			ctx.Writer.Flush()

			// 发送finish事件
			ctx.SSEvent(string(conversation.EventFinish), "{}")
			ctx.Writer.Flush()

			// 发送update_session事件
			finalUpdateSessionData := new(conversation.UpdateSessionData)
			finalUpdateSessionData.UpdatedAt = float64(time.Now().Unix()) + float64(time.Now().Nanosecond())/1e9
			finalUpdateSessionJSON, _ := json.Marshal(finalUpdateSessionData)
			ctx.SSEvent(string(conversation.EventUpdateSession), string(finalUpdateSessionJSON))
			ctx.Writer.Flush()

			// 发送标题事件
			titleData := new(conversation.TitleEventData)
			titleData.Content = "八字分析结果"
			titleJSON, _ := json.Marshal(titleData)
			ctx.SSEvent(string(conversation.EventTitle), string(titleJSON))
			ctx.Writer.Flush()

			// 发送关闭事件
			closeData := new(conversation.CloseEventData)
			closeData.ClickBehavior = "none"
			closeJSON, _ := json.Marshal(closeData)
			ctx.SSEvent(string(conversation.EventClose), string(closeJSON))
			ctx.Writer.Flush()
		}
		return nil
	})

	if err != nil {
		errorData := new(conversation.StreamData)
		errorData.V = err.Error()
		errorData.P = "error"
		errorJSON, _ := json.Marshal(errorData)
		ctx.SSEvent("", string(errorJSON))
		ctx.Writer.Flush()
	}
}

// ContinueConversation godoc
// @Summary      继续对话
// @Description  继续与AI的对话
// @Tags         对话
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.ContinueConversationRequest  true  "对话请求"
// @Success      200  {object}  vo.Result{data=interface{}}  "成功"
// @Failure      400  {object}  vo.Result                   "参数错误"
// @Failure      500  {object}  vo.Result                   "服务器内部错误"
// @Router       /api/v1/conversation/continue [post]
func (c *ConversationController) ContinueConversation(ctx *gin.Context) {
	req := new(dto.ContinueConversationRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, nil, bizErr.New(bizErr.PARAM_ERROR, "请求参数错误: "+err.Error())))
		return
	}

	validationErrors := utils.Validator(req)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, validationErrors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	result, err := c.conversationService.ContinueConversation(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, vo.Fail(ctx, nil, bizErr.New(bizErr.SYSTEM_ERROR, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, vo.Success(ctx, result))
}

// StreamContinueConversation godoc
// @Summary      流式继续对话
// @Description  流式继续与AI的对话
// @Tags         对话
// @Accept       json
// @Produce      text/event-stream
// @Security     BearerAuth
// @Param        request  body      dto.StreamContinueConversationRequest  true  "对话请求"
// @Success      200  {string}  string           "事件流"
// @Failure      400  {object}  vo.Result        "参数错误"
// @Failure      500  {object}  vo.Result        "服务器内部错误"
// @Router       /api/v1/conversation/continue/stream [post]
func (c *ConversationController) StreamContinueConversation(ctx *gin.Context) {
	req := new(dto.StreamContinueConversationRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, nil, bizErr.New(bizErr.PARAM_ERROR, "请求参数错误: "+err.Error())))
		return
	}

	validationErrors := utils.Validator(req)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, vo.Fail(ctx, validationErrors, bizErr.New(bizErr.PARAM_ERROR, "请求参数校验失败")))
		return
	}

	// 设置响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	// 获取请求和响应消息ID
	requestID, responseID, err := c.conversationService.GetMessageIDs(ctx)
	if err != nil {
		// 使用默认值
		requestID = 1
		responseID = 2
		utils.BizLogger(ctx).Errorf("获取消息ID失败: %v", err)
	}

	// 发送ready事件
	readyData := new(conversation.ReadyEventData)
	readyData.RequestMessageID = requestID
	readyData.ResponseMessageID = responseID
	readyJSON, _ := json.Marshal(readyData)
	ctx.SSEvent(string(conversation.EventReady), string(readyJSON))
	ctx.Writer.Flush()

	// 发送update_session事件
	updateSessionData := new(conversation.UpdateSessionData)
	updateSessionData.UpdatedAt = float64(time.Now().Unix()) + float64(time.Now().Nanosecond())/1e9
	updateSessionJSON, _ := json.Marshal(updateSessionData)
	ctx.SSEvent(string(conversation.EventUpdateSession), string(updateSessionJSON))
	ctx.Writer.Flush()

	// 创建初始响应
	initialResponse := new(conversation.ResponseValue)
	initialResponse.Response.MessageID = responseID
	initialResponse.Response.ParentID = requestID
	initialResponse.Response.Model = ""
	initialResponse.Response.Role = "ASSISTANT"
	initialResponse.Response.Content = ""
	initialResponse.Response.ThinkingEnabled = true
	initialResponse.Response.ThinkingContent = nil
	initialResponse.Response.ThinkingElapsedSecs = nil
	initialResponse.Response.BanEdit = false
	initialResponse.Response.BanRegenerate = false
	initialResponse.Response.Status = conversation.StatusWIP
	initialResponse.Response.AccumulatedTokenUsage = 0
	initialResponse.Response.Files = []interface{}{}
	initialResponse.Response.Tips = []interface{}{}
	initialResponse.Response.InsertedAt = float64(time.Now().Unix()) + float64(time.Now().Nanosecond())/1e9
	initialResponse.Response.SearchEnabled = true
	initialResponse.Response.SearchStatus = conversation.StatusInit
	initialResponse.Response.SearchResults = nil

	initialResponseJSON, _ := json.Marshal(initialResponse)
	ctx.SSEvent("", string(initialResponseJSON))
	ctx.Writer.Flush()

	// 发送搜索状态更新
	searchStatusData := new(conversation.StreamData)
	searchStatusData.V = string(conversation.StatusAnswer)
	searchStatusData.P = "response/search_status"
	searchStatusJSON, _ := json.Marshal(searchStatusData)
	ctx.SSEvent("", string(searchStatusJSON))
	ctx.Writer.Flush()

	// 开始思考时间
	thinkingStartTime := time.Now()

	// 思考内容
	thinkingContent := fmt.Sprintf("思考用户的问题: %s", req.Prompt)
	thinkingContentData := new(conversation.StreamData)
	thinkingContentData.V = thinkingContent
	thinkingContentData.P = "response/thinking_content"
	thinkingContentJSON, _ := json.Marshal(thinkingContentData)
	ctx.SSEvent("", string(thinkingContentJSON))
	ctx.Writer.Flush()

	// 流式继续对话
	var fullContent string
	var tokenCount int = 0

	err = c.conversationService.StreamContinueConversation(ctx, req, func(content string, done bool) error {
		if !done {
			// 累积完整内容
			fullContent += content
			tokenCount += len(content) / 4 // 粗略估算token数量

			// 发送内容更新
			contentData := new(conversation.StreamData)
			contentData.V = content
			contentData.O = conversation.OperationAppend
			if len(fullContent) == len(content) { // 第一次发送
				contentData.P = "response/content"
			}
			contentJSON, _ := json.Marshal(contentData)
			ctx.SSEvent("", string(contentJSON))
			ctx.Writer.Flush()
		} else {
			// 计算思考时间
			thinkingElapsedSecs := int(time.Since(thinkingStartTime).Seconds())
			thinkingElapsedData := new(conversation.StreamData)
			thinkingElapsedData.V = thinkingElapsedSecs
			thinkingElapsedData.P = "response/thinking_elapsed_secs"
			thinkingElapsedData.O = conversation.OperationSet
			thinkingElapsedJSON, _ := json.Marshal(thinkingElapsedData)
			ctx.SSEvent("", string(thinkingElapsedJSON))
			ctx.Writer.Flush()

			// 发送完成状态
			batchValues := []conversation.BatchValue{
				{V: conversation.StatusFinished, P: "status"},
				{V: tokenCount, P: "accumulated_token_usage"},
			}
			batchData := new(conversation.StreamData)
			batchData.V = batchValues
			batchData.P = "response"
			batchData.O = conversation.OperationBatch
			batchJSON, _ := json.Marshal(batchData)
			ctx.SSEvent("", string(batchJSON))
			ctx.Writer.Flush()

			// 发送finish事件
			ctx.SSEvent(string(conversation.EventFinish), "{}")
			ctx.Writer.Flush()

			// 发送update_session事件
			finalUpdateSessionData := new(conversation.UpdateSessionData)
			finalUpdateSessionData.UpdatedAt = float64(time.Now().Unix()) + float64(time.Now().Nanosecond())/1e9
			finalUpdateSessionJSON, _ := json.Marshal(finalUpdateSessionData)
			ctx.SSEvent(string(conversation.EventUpdateSession), string(finalUpdateSessionJSON))
			ctx.Writer.Flush()

			// 发送标题事件
			titleData := new(conversation.TitleEventData)
			titleData.Content = fmt.Sprintf("回复: %s", req.Prompt)
			titleJSON, _ := json.Marshal(titleData)
			ctx.SSEvent(string(conversation.EventTitle), string(titleJSON))
			ctx.Writer.Flush()

			// 发送关闭事件
			closeData := new(conversation.CloseEventData)
			closeData.ClickBehavior = "none"
			closeJSON, _ := json.Marshal(closeData)
			ctx.SSEvent(string(conversation.EventClose), string(closeJSON))
			ctx.Writer.Flush()
		}
		return nil
	})

	if err != nil {
		errorData := new(conversation.StreamData)
		errorData.V = err.Error()
		errorData.P = "error"
		errorJSON, _ := json.Marshal(errorData)
		ctx.SSEvent("", string(errorJSON))
		ctx.Writer.Flush()
	}
}
