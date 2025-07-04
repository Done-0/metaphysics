// Package impl 提供对话相关的数据访问实现
// 创建者：Done-0
// 创建时间：2025-07-03
package impl

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Done-0/metaphysics/internal/global"
	conversationModel "github.com/Done-0/metaphysics/internal/model/conversation"
	"github.com/Done-0/metaphysics/internal/utils"
	conversationMapper "github.com/Done-0/metaphysics/pkg/serve/mapper/conversation"
)

const (
	// 对话历史记录Redis键前缀
	CONVERSATION_HISTORY_PREFIX = "CONVERSATION:HISTORY:"
	// 对话历史记录过期时间（24小时）
	CONVERSATION_EXPIRE_TIME = 24 * time.Hour
	// 用户最后消息ID Redis键前缀
	LAST_MESSAGE_ID_PREFIX = "USER:LAST_MESSAGE_ID:"
	// 初始消息ID
	INITIAL_MESSAGE_ID = 1
)

// MessageIDs 消息ID结构
type MessageIDs struct {
	RequestID  int `json:"request_id"`  // 请求消息ID
	ResponseID int `json:"response_id"` // 响应消息ID
}

// ConversationMapperImpl 对话数据访问接口实现
type ConversationMapperImpl struct{}

// NewConversationMapper 创建对话数据访问接口实现
// 返回值：
//   - conversationMapper.ConversationMapper: 对话数据访问接口
func NewConversationMapper() conversationMapper.ConversationMapper {
	return &ConversationMapperImpl{}
}

// SaveConversation 保存对话
// 参数：
//   - ctx: 上下文信息
//   - conversation: 对话模型
//
// 返回值：
//   - error: 错误信息
func (m *ConversationMapperImpl) SaveConversation(ctx *gin.Context, conversation *conversationModel.Conversation) error {
	return utils.RunDBTransaction(ctx, func() error {
		db := utils.GetDBFromContext(ctx)
		if err := db.Create(conversation).Error; err != nil {
			return fmt.Errorf("保存对话记录失败: %w", err)
		}
		return nil
	})
}

// GetConversationByID 根据ID获取对话
// 参数：
//   - ctx: 上下文信息
//   - id: 对话ID
//
// 返回值：
//   - *conversationModel.Conversation: 对话模型
//   - error: 错误信息
func (m *ConversationMapperImpl) GetConversationByID(ctx *gin.Context, id int64) (*conversationModel.Conversation, error) {
	var conversation conversationModel.Conversation
	db := utils.GetDBFromContext(ctx)
	err := db.Where("id = ? AND deleted = ?", id, false).First(&conversation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("对话记录不存在")
		}
		return nil, fmt.Errorf("查询对话记录失败: %w", err)
	}
	return &conversation, nil
}

// GetConversationsByUserID 获取用户的所有对话
// 参数：
//   - ctx: 上下文信息
//   - userID: 用户ID
//   - pageNo: 页码
//   - pageSize: 每页数量
//
// 返回值：
//   - []*conversationModel.Conversation: 对话列表
//   - int64: 总记录数
//   - error: 错误信息
func (m *ConversationMapperImpl) GetConversationsByUserID(ctx *gin.Context, userID int64, pageNo, pageSize int) ([]*conversationModel.Conversation, int64, error) {
	var conversations []*conversationModel.Conversation
	var total int64

	db := utils.GetDBFromContext(ctx)
	query := db.Model(&conversationModel.Conversation{}).Where("user_id = ? AND deleted = ?", userID, false)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询对话记录总数失败: %w", err)
	}

	// 查询记录 - 分页查询
	offset := (pageNo - 1) * pageSize
	if err := query.
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&conversations).Error; err != nil {
		return nil, 0, fmt.Errorf("查询对话记录列表失败: %w", err)
	}

	return conversations, total, nil
}

// UpdateConversation 更新对话
// 参数：
//   - ctx: 上下文信息
//   - conversation: 对话模型
//
// 返回值：
//   - error: 错误信息
func (m *ConversationMapperImpl) UpdateConversation(ctx *gin.Context, conversation *conversationModel.Conversation) error {
	return utils.RunDBTransaction(ctx, func() error {
		db := utils.GetDBFromContext(ctx)
		if err := db.Save(conversation).Error; err != nil {
			return fmt.Errorf("更新对话记录失败: %w", err)
		}
		return nil
	})
}

// DeleteConversation 删除对话
// 参数：
//   - ctx: 上下文信息
//   - id: 对话ID
//
// 返回值：
//   - error: 错误信息
func (m *ConversationMapperImpl) DeleteConversation(ctx *gin.Context, id int64) error {
	return utils.RunDBTransaction(ctx, func() error {
		db := utils.GetDBFromContext(ctx)
		if err := db.Model(&conversationModel.Conversation{}).
			Where("id = ?", id).
			Update("deleted", true).Error; err != nil {
			return fmt.Errorf("删除对话记录失败: %w", err)
		}
		return nil
	})
}

// SaveMessage 保存消息
// 参数：
//   - ctx: 上下文信息
//   - message: 消息模型
//
// 返回值：
//   - error: 错误信息
func (m *ConversationMapperImpl) SaveMessage(ctx *gin.Context, message *conversationModel.Message) error {
	return utils.RunDBTransaction(ctx, func() error {
		db := utils.GetDBFromContext(ctx)
		if err := db.Create(message).Error; err != nil {
			return fmt.Errorf("保存消息记录失败: %w", err)
		}
		return nil
	})
}

// GetMessagesByConversationID 获取对话的所有消息
// 参数：
//   - ctx: 上下文信息
//   - conversationID: 对话ID
//
// 返回值：
//   - []*conversationModel.Message: 消息列表
//   - error: 错误信息
func (m *ConversationMapperImpl) GetMessagesByConversationID(ctx *gin.Context, conversationID int64) ([]*conversationModel.Message, error) {
	var messages []*conversationModel.Message
	db := utils.GetDBFromContext(ctx)
	err := db.
		Where("conversation_id = ? AND deleted = ?", conversationID, false).
		Order("gmt_create ASC").
		Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("查询消息记录失败: %w", err)
	}
	return messages, nil
}

// GetLatestConversationHistory 获取最近的对话历史
// 参数：
//   - ctx: 上下文信息
//   - userID: 用户ID
//
// 返回值：
//   - string: 会话ID
//   - string: 对话历史
//   - error: 错误信息
func (m *ConversationMapperImpl) GetLatestConversationHistory(ctx *gin.Context, userID int64) (string, string, error) {
	// 获取用户所有对话ID
	pattern := fmt.Sprintf("%s%d:*", CONVERSATION_HISTORY_PREFIX, userID)
	keys, err := global.RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return "", "", fmt.Errorf("获取对话ID失败: %w", err)
	}

	if len(keys) == 0 {
		return uuid.New().String(), "", nil
	}

	// 获取最近的对话历史
	latestKey := keys[0]
	history, err := global.RedisClient.Get(ctx, latestKey).Result()
	if err != nil {
		return "", "", fmt.Errorf("获取对话历史失败: %w", err)
	}

	// 从键中提取对话ID
	sessionID := latestKey[len(fmt.Sprintf("%s%d:", CONVERSATION_HISTORY_PREFIX, userID)):]

	return sessionID, history, nil
}

// SaveConversationHistory 保存对话历史
// 参数：
//   - ctx: 上下文信息
//   - userID: 用户ID
//   - sessionID: 会话ID
//   - history: 对话历史
//
// 返回值：
//   - error: 错误信息
func (m *ConversationMapperImpl) SaveConversationHistory(ctx *gin.Context, userID int64, sessionID string, history string) error {
	historyKey := fmt.Sprintf("%s%d:%s", CONVERSATION_HISTORY_PREFIX, userID, sessionID)
	if err := global.RedisClient.Set(ctx, historyKey, history, CONVERSATION_EXPIRE_TIME).Err(); err != nil {
		return fmt.Errorf("保存对话历史失败: %w", err)
	}

	return nil
}

// GetNextMessageIDs 获取下一个消息ID
// 参数：
//   - ctx: 上下文信息
//   - userID: 用户ID
//
// 返回值：
//   - int: 请求消息ID
//   - int: 响应消息ID
//   - error: 错误信息
func (m *ConversationMapperImpl) GetNextMessageIDs(ctx *gin.Context, userID int64) (int, int, error) {
	db := utils.GetDBFromContext(ctx)

	key := fmt.Sprintf("%s%d", LAST_MESSAGE_ID_PREFIX, userID)

	// 尝试从数据库获取消息计数器
	var counter conversationModel.MessageCounter
	err := db.Where("user_id = ?", userID).First(&counter).Error

	// 数据库中存在消息计数器
	if err == nil {
		requestID := counter.NextID
		responseID := requestID + 1
		counter.NextID += 2

		// 更新数据库
		if err := db.Save(&counter).Error; err != nil {
			utils.BizLogger(ctx).Errorf("更新消息计数器失败: %v", err)
		}

		// 同步更新Redis
		messageIDs := new(MessageIDs)
		messageIDs.RequestID = requestID
		messageIDs.ResponseID = responseID

		if jsonData, err := json.Marshal(messageIDs); err == nil {
			_ = global.RedisClient.Set(ctx, key, string(jsonData), CONVERSATION_EXPIRE_TIME).Err()
		}

		return requestID, responseID, nil
	}

	// 数据库中不存在记录，尝试从Redis获取
	if err == gorm.ErrRecordNotFound {
		val, err := global.RedisClient.Get(ctx, key).Result()
		messageIDs := new(MessageIDs)

		// Redis中不存在或获取出错
		if err == redis.Nil || val == "" || err != nil {
			if err != nil && err != redis.Nil {
				utils.BizLogger(ctx).Errorf("获取Redis消息ID失败: %v", err)
			}

			// 使用初始值
			messageIDs.RequestID = INITIAL_MESSAGE_ID
			messageIDs.ResponseID = INITIAL_MESSAGE_ID + 1

			// 创建数据库记录
			counter = conversationModel.MessageCounter{
				UserID: userID,
				NextID: INITIAL_MESSAGE_ID + 2,
			}
			if err := db.Create(&counter).Error; err != nil {
				utils.BizLogger(ctx).Errorf("创建消息计数器失败: %v", err)
			}
		} else {
			// Redis中有数据，解析并递增
			if err := json.Unmarshal([]byte(val), messageIDs); err != nil {
				messageIDs.RequestID = INITIAL_MESSAGE_ID
				messageIDs.ResponseID = INITIAL_MESSAGE_ID + 1
			} else {
				messageIDs.RequestID += 2
				messageIDs.ResponseID += 2

				// 同步到数据库
				counter = conversationModel.MessageCounter{
					UserID: userID,
					NextID: messageIDs.RequestID + 2,
				}
				if err := db.Create(&counter).Error; err != nil {
					utils.BizLogger(ctx).Errorf("创建消息计数器失败: %v", err)
				}
			}
		}

		// 更新Redis
		if jsonData, err := json.Marshal(messageIDs); err == nil {
			if err := global.RedisClient.Set(ctx, key, string(jsonData), CONVERSATION_EXPIRE_TIME).Err(); err != nil {
				utils.BizLogger(ctx).Errorf("保存消息ID到Redis失败: %v", err)
			}
		}

		return messageIDs.RequestID, messageIDs.ResponseID, nil
	}

	// 其他数据库错误
	return INITIAL_MESSAGE_ID, INITIAL_MESSAGE_ID + 1, fmt.Errorf("获取消息计数器失败: %w", err)
}
