// Package utils 提供各种工具函数
// 创建者：Done-0
// 创建时间：2023-10-18
package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// SSEEventType 定义 SSE 事件类型
type SSEEventType string

// SSE 事件类型常量
const (
	SSE_EVENT_CONTENT   SSEEventType = "content"   // 内容事件
	SSE_EVENT_REASONING SSEEventType = "reasoning" // 推理事件
	SSE_EVENT_DONE      SSEEventType = "done"      // 完成事件
	SSE_EVENT_ERROR     SSEEventType = "error"     // 错误事件
)

// SetupSSEHeaders 设置 SSE 响应头
// 参数：
//
//	ctx: Gin 上下文
func SetupSSEHeaders(ctx *gin.Context) {
	headers := map[string]string{
		"Content-Type":      "text/event-stream",
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"X-Accel-Buffering": "no",
	}

	for key, value := range headers {
		ctx.Header(key, value)
	}
	ctx.Writer.Flush()
}

// SendSSEEvent 发送 SSE 事件
// 参数：
//
//	ctx: Gin 上下文
//	eventType: 事件类型，空字符串表示不带事件类型
//	data: 事件数据
//
// 返回值：
//
//	error: 错误信息
func SendSSEEvent(ctx *gin.Context, eventType SSEEventType, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化 SSE 事件数据失败: %w", err)
	}

	eventString := fmt.Sprintf("data: %s\n\n", jsonData)

	_, err = ctx.Writer.WriteString(eventString)
	if err != nil {
		return fmt.Errorf("写入 SSE 事件失败: %w", err)
	}
	ctx.Writer.Flush()
	return nil
}
