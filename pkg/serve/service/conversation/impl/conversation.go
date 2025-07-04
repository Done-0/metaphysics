// Package impl 提供对话服务层实现
// 创建者：Done-0
// 创建时间：2025-07-03
package impl

import (
	"fmt"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	internalAI "github.com/Done-0/metaphysics/internal/ai"
	"github.com/Done-0/metaphysics/internal/ai/types"
	conversationModel "github.com/Done-0/metaphysics/internal/model/conversation"
	"github.com/Done-0/metaphysics/internal/utils"
	"github.com/Done-0/metaphysics/pkg/serve/controller/conversation/dto"
	baziMapper "github.com/Done-0/metaphysics/pkg/serve/mapper/bazi"
	baziMapperImpl "github.com/Done-0/metaphysics/pkg/serve/mapper/bazi/impl"
	conversationMapper "github.com/Done-0/metaphysics/pkg/serve/mapper/conversation"
	conversationSrv "github.com/Done-0/metaphysics/pkg/serve/service/conversation"
	"github.com/Done-0/metaphysics/pkg/vo/conversation"
)

const (
	INITIAL_MESSAGE_ID = 1 // 初始消息ID
)

// ConversationServiceImpl 对话服务实现
type ConversationServiceImpl struct {
	baziMapper         baziMapper.BaziMapper
	conversationMapper conversationMapper.ConversationMapper
	aiService          types.Service
}

// NewConversationService 创建对话服务实例
func NewConversationService(conversationMapperImpl conversationMapper.ConversationMapper) conversationSrv.ConversationService {
	return &ConversationServiceImpl{
		baziMapper:         baziMapperImpl.NewBaziMapper(),
		conversationMapper: conversationMapperImpl,
		aiService:          internalAI.New(),
	}
}

// AnalyzeBaziByUserID 根据用户ID分析八字
func (s *ConversationServiceImpl) AnalyzeBaziByUserID(ctx *gin.Context) (*conversation.BaziAnalysisResponse, error) {
	// 获取用户ID
	id, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 获取用户八字记录
	records, total, err := s.baziMapper.GetBaziList(ctx, 1, 1)
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取用户八字记录失败: %v", err)
		return nil, fmt.Errorf("获取用户八字记录失败: %w", err)
	}

	if total == 0 || len(records) == 0 {
		utils.BizLogger(ctx).Errorf("未找到用户八字记录")
		return nil, fmt.Errorf("未找到用户八字记录")
	}

	record := records[0]

	// 构建八字信息
	baziInfo := map[string]string{
		"year":      record.YearPillar,
		"month":     record.MonthPillar,
		"day":       record.DayPillar,
		"hour":      record.HourPillar,
		"year_gan":  record.YearGan,
		"month_gan": record.MonthGan,
		"day_gan":   record.DayGan,
		"hour_gan":  record.HourGan,
		"year_zhi":  record.YearZhi,
		"month_zhi": record.MonthZhi,
		"day_zhi":   record.DayZhi,
		"hour_zhi":  record.HourZhi,
	}

	// 调用AI分析八字
	analysisResponse, err := s.aiService.AnalyzeBaziWithReasoning(ctx, record.Name, record.Gender, record.BirthTime, record.Calendar, baziInfo)
	if err != nil {
		utils.BizLogger(ctx).Errorf("AI分析八字失败: %v", err)
		return nil, fmt.Errorf("AI分析八字失败: %w", err)
	}

	// 创建新的会话ID并保存对话记录
	sessionID := uuid.New().String()
	if err := s.conversationMapper.SaveConversationHistory(ctx, id, sessionID, analysisResponse.Analysis); err != nil {
		utils.BizLogger(ctx).Errorf("保存对话历史失败: %v", err)
	}

	// 创建对话记录
	conversationRecord := &conversationModel.Conversation{
		UserID:      id,
		Title:       "八字分析",
		SessionID:   sessionID,
		FirstPrompt: "分析我的八字",
	}
	if err := s.conversationMapper.SaveConversation(ctx, conversationRecord); err != nil {
		utils.BizLogger(ctx).Errorf("保存对话记录失败: %v", err)
	}

	// 设置请求ID并构建响应
	reqID := requestid.Get(ctx)

	// 直接返回从AI获取的分析结果
	analysisResponse.RequestID = reqID
	analysisResponse.UserID = id
	return analysisResponse, nil
}

// StreamAnalyzeBaziByUserID 流式分析用户八字
func (s *ConversationServiceImpl) StreamAnalyzeBaziByUserID(ctx *gin.Context, handler func(content string, done bool) error) error {
	// 获取用户ID
	id, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	// 获取用户八字记录
	records, total, err := s.baziMapper.GetBaziList(ctx, 1, 1)
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取用户八字记录失败: %v", err)
		return fmt.Errorf("获取用户八字记录失败: %w", err)
	}

	if total == 0 || len(records) == 0 {
		utils.BizLogger(ctx).Errorf("未找到用户八字记录")
		return fmt.Errorf("未找到用户八字记录")
	}

	record := records[0]

	// 构建八字信息
	baziInfo := map[string]string{
		"year":      record.YearPillar,
		"month":     record.MonthPillar,
		"day":       record.DayPillar,
		"hour":      record.HourPillar,
		"year_gan":  record.YearGan,
		"month_gan": record.MonthGan,
		"day_gan":   record.DayGan,
		"hour_gan":  record.HourGan,
		"year_zhi":  record.YearZhi,
		"month_zhi": record.MonthZhi,
		"day_zhi":   record.DayZhi,
		"hour_zhi":  record.HourZhi,
	}

	// 获取消息ID
	requestID, responseID, err := s.conversationMapper.GetNextMessageIDs(ctx, id)
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取消息ID失败: %v", err)
		requestID = INITIAL_MESSAGE_ID
		responseID = INITIAL_MESSAGE_ID + 1
	}

	// 创建会话ID
	sessionID := uuid.New().String()

	// 创建对话记录
	conversationRecord := &conversationModel.Conversation{
		UserID:      id,
		Title:       "八字分析",
		SessionID:   sessionID,
		FirstPrompt: "分析我的八字",
	}
	if err := s.conversationMapper.SaveConversation(ctx, conversationRecord); err != nil {
		utils.BizLogger(ctx).Errorf("保存对话记录失败: %v", err)
	}

	// 保存用户消息
	userMessage := &conversationModel.Message{
		ConversationID: conversationRecord.ID,
		UserID:         id,
		SessionID:      sessionID,
		Role:           "USER",
		Content:        "分析我的八字",
		RequestID:      requestID,
		ResponseID:     responseID,
		ParentID:       0,
	}
	if err := s.conversationMapper.SaveMessage(ctx, userMessage); err != nil {
		utils.BizLogger(ctx).Errorf("保存用户消息失败: %v", err)
	}

	// 存储完整的分析结果
	var fullAnalysis string

	// 包装处理函数
	wrappedHandler := types.StreamHandler(func(chunk *conversation.StreamChunk) error {
		// 累积完整内容
		if !chunk.Done {
			fullAnalysis += chunk.Content
		} else {
			// 保存对话历史和AI回复
			if err := s.conversationMapper.SaveConversationHistory(ctx, id, sessionID, fullAnalysis); err != nil {
				utils.BizLogger(ctx).Errorf("保存对话历史失败: %v", err)
			}

			aiMessage := &conversationModel.Message{
				ConversationID: conversationRecord.ID,
				UserID:         id,
				SessionID:      sessionID,
				Role:           "ASSISTANT",
				Content:        fullAnalysis,
				RequestID:      requestID,
				ResponseID:     responseID,
				ParentID:       requestID,
				TokenUsage:     len(fullAnalysis) / 4, // 粗略估算token数量
			}
			if err := s.conversationMapper.SaveMessage(ctx, aiMessage); err != nil {
				utils.BizLogger(ctx).Errorf("保存AI回复失败: %v", err)
			}
		}

		// 转发给原始处理函数
		return handler(chunk.Content, chunk.Done)
	})

	// 调用AI服务进行流式分析
	return s.aiService.StreamAnalyzeBazi(ctx, record.Name, record.Gender, record.BirthTime, record.Calendar, baziInfo, wrappedHandler)
}

// ContinueConversation 继续与AI的对话
func (s *ConversationServiceImpl) ContinueConversation(ctx *gin.Context, req *dto.ContinueConversationRequest) (*conversation.ConversationResponse, error) {
	// 获取用户ID
	id, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 获取最近的对话历史
	sessionID, history, err := s.conversationMapper.GetLatestConversationHistory(ctx, id)
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取对话历史失败: %v", err)
		return nil, fmt.Errorf("获取对话历史失败: %w", err)
	}

	// 获取或创建对话
	var conversationRecord *conversationModel.Conversation
	conversations, total, err := s.conversationMapper.GetConversationsByUserID(ctx, id, 1, 1)
	if err != nil || total == 0 {
		// 创建新对话
		conversationRecord = &conversationModel.Conversation{
			UserID:      id,
			Title:       fmt.Sprintf("对话: %s", req.Prompt),
			SessionID:   sessionID,
			FirstPrompt: req.Prompt,
		}
		if err := s.conversationMapper.SaveConversation(ctx, conversationRecord); err != nil {
			utils.BizLogger(ctx).Errorf("保存对话记录失败: %v", err)
			return nil, fmt.Errorf("保存对话记录失败: %w", err)
		}
	} else {
		conversationRecord = conversations[0]
	}

	// 获取消息ID
	requestID, responseID, err := s.conversationMapper.GetNextMessageIDs(ctx, id)
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取消息ID失败: %v", err)
		requestID = INITIAL_MESSAGE_ID
		responseID = INITIAL_MESSAGE_ID + 1
	}

	// 保存用户消息
	userMessage := &conversationModel.Message{
		ConversationID: conversationRecord.ID,
		UserID:         id,
		SessionID:      sessionID,
		Role:           "USER",
		Content:        req.Prompt,
		RequestID:      requestID,
		ResponseID:     responseID,
		ParentID:       0,
	}
	if err := s.conversationMapper.SaveMessage(ctx, userMessage); err != nil {
		utils.BizLogger(ctx).Errorf("保存用户消息失败: %v", err)
	}

	// 构建对话提示并请求AI回复
	prompt := fmt.Sprintf("以下是之前的对话内容：\n\n%s\n\n用户的新问题是：%s", history, req.Prompt)
	baziInfo := map[string]string{"prompt": prompt}
	analysisResponse, err := s.aiService.AnalyzeBaziWithReasoning(ctx, "", "", time.Time{}, "", baziInfo)
	if err != nil {
		utils.BizLogger(ctx).Errorf("AI回复失败: %v", err)
		return nil, fmt.Errorf("AI回复失败: %w", err)
	}

	// 更新对话历史
	newHistory := fmt.Sprintf("%s\n\n用户：%s\n\nAI：%s", history, req.Prompt, analysisResponse.Analysis)
	if err := s.conversationMapper.SaveConversationHistory(ctx, id, sessionID, newHistory); err != nil {
		utils.BizLogger(ctx).Errorf("更新对话历史失败: %v", err)
	}

	// 保存AI回复
	aiMessage := &conversationModel.Message{
		ConversationID: conversationRecord.ID,
		UserID:         id,
		SessionID:      sessionID,
		Role:           "ASSISTANT",
		Content:        analysisResponse.Analysis,
		RequestID:      requestID,
		ResponseID:     responseID,
		ParentID:       requestID,
		TokenUsage:     len(analysisResponse.Analysis) / 4, // 粗略估算token数量
	}
	if err := s.conversationMapper.SaveMessage(ctx, aiMessage); err != nil {
		utils.BizLogger(ctx).Errorf("保存AI回复失败: %v", err)
	}

	// 构建响应
	reqID := requestid.Get(ctx)
	result := &conversation.ConversationResponse{
		RequestID:      reqID,
		UserID:         id,
		UserMessage:    req.Prompt,
		AIResponse:     analysisResponse.Analysis,
		ConversationID: sessionID,
	}

	return result, nil
}

// StreamContinueConversation 流式继续对话
func (s *ConversationServiceImpl) StreamContinueConversation(ctx *gin.Context, req *dto.StreamContinueConversationRequest, handler func(content string, done bool) error) error {
	// 获取用户ID
	id, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	// 获取最近的对话历史
	sessionID, history, err := s.conversationMapper.GetLatestConversationHistory(ctx, id)
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取对话历史失败: %v", err)
		return fmt.Errorf("获取对话历史失败: %w", err)
	}

	// 获取或创建对话
	var conversationRecord *conversationModel.Conversation
	conversations, total, err := s.conversationMapper.GetConversationsByUserID(ctx, id, 1, 1)
	if err != nil || total == 0 {
		// 创建新对话
		conversationRecord = &conversationModel.Conversation{
			UserID:      id,
			Title:       fmt.Sprintf("对话: %s", req.Prompt),
			SessionID:   sessionID,
			FirstPrompt: req.Prompt,
		}
		if err := s.conversationMapper.SaveConversation(ctx, conversationRecord); err != nil {
			utils.BizLogger(ctx).Errorf("保存对话记录失败: %v", err)
			return fmt.Errorf("保存对话记录失败: %w", err)
		}
	} else {
		conversationRecord = conversations[0]
	}

	// 获取消息ID
	requestID, responseID, err := s.conversationMapper.GetNextMessageIDs(ctx, id)
	if err != nil {
		utils.BizLogger(ctx).Errorf("获取消息ID失败: %v", err)
		requestID = INITIAL_MESSAGE_ID
		responseID = INITIAL_MESSAGE_ID + 1
	}

	// 保存用户消息
	userMessage := &conversationModel.Message{
		ConversationID: conversationRecord.ID,
		UserID:         id,
		SessionID:      sessionID,
		Role:           "USER",
		Content:        req.Prompt,
		RequestID:      requestID,
		ResponseID:     responseID,
		ParentID:       0,
	}
	if err := s.conversationMapper.SaveMessage(ctx, userMessage); err != nil {
		utils.BizLogger(ctx).Errorf("保存用户消息失败: %v", err)
	}

	// 构建对话提示
	prompt := fmt.Sprintf("以下是之前的对话内容：\n\n%s\n\n用户的新问题是：%s", history, req.Prompt)

	// 存储完整的响应
	var fullResponse string

	// 包装处理函数
	wrappedHandler := types.StreamHandler(func(chunk *conversation.StreamChunk) error {
		// 累积完整内容
		if !chunk.Done {
			fullResponse += chunk.Content
		} else {
			// 更新对话历史和保存AI回复
			newHistory := fmt.Sprintf("%s\n\n用户：%s\n\nAI：%s", history, req.Prompt, fullResponse)
			if err := s.conversationMapper.SaveConversationHistory(ctx, id, sessionID, newHistory); err != nil {
				utils.BizLogger(ctx).Errorf("更新对话历史失败: %v", err)
			}

			aiMessage := &conversationModel.Message{
				ConversationID: conversationRecord.ID,
				UserID:         id,
				SessionID:      sessionID,
				Role:           "ASSISTANT",
				Content:        fullResponse,
				RequestID:      requestID,
				ResponseID:     responseID,
				ParentID:       requestID,
				TokenUsage:     len(fullResponse) / 4, // 粗略估算token数量
			}
			if err := s.conversationMapper.SaveMessage(ctx, aiMessage); err != nil {
				utils.BizLogger(ctx).Errorf("保存AI回复失败: %v", err)
			}
		}

		// 转发给原始处理函数
		return handler(chunk.Content, chunk.Done)
	})

	// 调用AI服务进行流式对话
	return s.aiService.StreamAnalyzeBazi(ctx, "", "", time.Time{}, "", map[string]string{"prompt": prompt}, wrappedHandler)
}

// GetMessageIDs 获取当前用户的消息ID
func (s *ConversationServiceImpl) GetMessageIDs(ctx *gin.Context) (int, int, error) {
	id, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return 0, 0, err
	}

	return s.conversationMapper.GetNextMessageIDs(ctx, id)
}
