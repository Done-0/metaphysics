// Package conversation 提供对话相关的视图对象
// 创建者：Done-0
// 创建时间：2025-07-03
package conversation

// EventType 事件类型
type EventType string

const (
	EventReady         EventType = "ready"          // 准备就绪
	EventUpdateSession EventType = "update_session" // 更新会话
	EventFinish        EventType = "finish"         // 完成
	EventTitle         EventType = "title"          // 标题
	EventClose         EventType = "close"          // 关闭
)

// Status 状态
type Status string

const (
	StatusInit      Status = "init"      // 初始状态
	StatusWIP       Status = "wip"       // 处理中
	StatusAnswer    Status = "answer"    // 回答中
	StatusFinished  Status = "finished"  // 已完成
	StatusError     Status = "error"     // 错误
	StatusRejection Status = "rejection" // 拒绝
)

// OperationType 操作类型
type OperationType string

const (
	OperationAppend OperationType = "append" // 追加操作
	OperationSet    OperationType = "set"    // 设置操作
	OperationBatch  OperationType = "batch"  // 批量操作
)

// ReadyEventData 准备就绪事件数据
// @Description 准备就绪事件数据
// @Property RequestMessageID int true "请求消息ID"
// @Property ResponseMessageID int true "响应消息ID"
type ReadyEventData struct {
	RequestMessageID  int `json:"request_message_id"`  // 请求消息ID
	ResponseMessageID int `json:"response_message_id"` // 响应消息ID
}

// UpdateSessionData 更新会话事件数据
// @Description 更新会话事件数据
// @Property UpdatedAt float64 true "更新时间"
type UpdateSessionData struct {
	UpdatedAt float64 `json:"updated_at"` // 更新时间
}

// TitleEventData 标题事件数据
// @Description 标题事件数据
// @Property Content string true "标题内容"
type TitleEventData struct {
	Content string `json:"content"` // 标题内容
}

// CloseEventData 关闭事件数据
// @Description 关闭事件数据
// @Property ClickBehavior string true "点击行为"
type CloseEventData struct {
	ClickBehavior string `json:"click_behavior"` // 点击行为
}

// StreamChunk 流式响应数据块
// @Description 流式响应数据块
// @Property Content string true "内容"
// @Property Done bool true "是否完成"
type StreamChunk struct {
	Content string `json:"content"` // 内容
	Done    bool   `json:"done"`    // 是否完成
}

// StreamData 流式数据
// @Description 流式数据
// @Property V any true "值"
// @Property P string true "路径"
// @Property O OperationType true "操作类型"
type StreamData struct {
	V any           `json:"v"`           // 值
	P string        `json:"p,omitempty"` // 路径
	O OperationType `json:"o,omitempty"` // 操作类型
}

// BatchValue 批量值
// @Description 批量值
// @Property V any true "值"
// @Property P string true "路径"
type BatchValue struct {
	V any    `json:"v"` // 值
	P string `json:"p"` // 路径
}

// ResponseDetail 响应详情
// @Description 响应详情
// @Property MessageID int true "消息ID"
// @Property ParentID int true "父消息ID"
// @Property Model string true "模型"
// @Property Role string true "角色"
// @Property Content string true "内容"
// @Property ThinkingEnabled bool true "是否启用思考"
// @Property ThinkingContent interface{} true "思考内容"
// @Property ThinkingElapsedSecs interface{} true "思考耗时"
// @Property BanEdit bool true "是否禁止编辑"
// @Property BanRegenerate bool true "是否禁止重新生成"
// @Property Status Status true "状态"
// @Property AccumulatedTokenUsage int true "累计Token使用量"
// @Property Files []interface{} true "文件"
// @Property Tips []interface{} true "提示"
// @Property InsertedAt float64 true "插入时间"
// @Property SearchEnabled bool true "是否启用搜索"
// @Property SearchStatus Status true "搜索状态"
// @Property SearchResults interface{} true "搜索结果"
type ResponseDetail struct {
	MessageID             int           `json:"message_id"`              // 消息ID
	ParentID              int           `json:"parent_id"`               // 父消息ID
	Model                 string        `json:"model"`                   // 模型
	Role                  string        `json:"role"`                    // 角色
	Content               string        `json:"content"`                 // 内容
	ThinkingEnabled       bool          `json:"thinking_enabled"`        // 是否启用思考
	ThinkingContent       interface{}   `json:"thinking_content"`        // 思考内容
	ThinkingElapsedSecs   interface{}   `json:"thinking_elapsed_secs"`   // 思考耗时
	BanEdit               bool          `json:"ban_edit"`                // 是否禁止编辑
	BanRegenerate         bool          `json:"ban_regenerate"`          // 是否禁止重新生成
	Status                Status        `json:"status"`                  // 状态
	AccumulatedTokenUsage int           `json:"accumulated_token_usage"` // 累计Token使用量
	Files                 []interface{} `json:"files"`                   // 文件
	Tips                  []interface{} `json:"tips"`                    // 提示
	InsertedAt            float64       `json:"inserted_at"`             // 插入时间
	SearchEnabled         bool          `json:"search_enabled"`          // 是否启用搜索
	SearchStatus          Status        `json:"search_status"`           // 搜索状态
	SearchResults         interface{}   `json:"search_results"`          // 搜索结果
}

// ResponseValue 响应值
// @Description 响应值
// @Property Response ResponseDetail true "响应详情"
type ResponseValue struct {
	Response ResponseDetail `json:"response"` // 响应详情
}

// BaziAnalysisResponse 八字分析响应
// @Description 八字分析响应
// @Property RequestID string true "请求ID"
// @Property UserID int64 true "用户ID"
// @Property Name string true "姓名"
// @Property Gender string true "性别"
// @Property YearPillar string true "年柱"
// @Property MonthPillar string true "月柱"
// @Property DayPillar string true "日柱"
// @Property HourPillar string true "时柱"
// @Property Analysis string true "分析结果"
type BaziAnalysisResponse struct {
	RequestID   string `json:"request_id"`   // 请求ID
	UserID      int64  `json:"user_id"`      // 用户ID
	Name        string `json:"name"`         // 姓名
	Gender      string `json:"gender"`       // 性别
	YearPillar  string `json:"year_pillar"`  // 年柱
	MonthPillar string `json:"month_pillar"` // 月柱
	DayPillar   string `json:"day_pillar"`   // 日柱
	HourPillar  string `json:"hour_pillar"`  // 时柱
	Analysis    string `json:"analysis"`     // 分析结果
}

// ConversationResponse 对话响应
// @Description 对话响应
// @Property RequestID string true "请求ID"
// @Property UserID int64 true "用户ID"
// @Property UserMessage string true "用户消息"
// @Property AIResponse string true "AI响应"
// @Property ConversationID string true "对话ID"
type ConversationResponse struct {
	RequestID      string `json:"request_id"`      // 请求ID
	UserID         int64  `json:"user_id"`         // 用户ID
	UserMessage    string `json:"user_message"`    // 用户消息
	AIResponse     string `json:"ai_response"`     // AI响应
	ConversationID string `json:"conversation_id"` // 对话ID
}
