// Package prompt 提供 AI 提示模板和格式化工具
// 创建者：Done-0
// 创建时间：2024-06-10
package prompt

import (
	"fmt"
	"time"

	"github.com/Done-0/metaphysics/internal/utils"
)

// BuildBaziPrompt 构建八字分析提示
// 参数：
//   - name: 姓名
//   - gender: 性别
//   - birthTime: 出生时间
//   - calendar: 日历类型 (solar/lunar)
//   - baziInfo: 八字信息
//
// 返回值：
//   - string: 格式化的提示文本
func BuildBaziPrompt(name, gender string, birthTime time.Time, calendar string, baziInfo map[string]string) string {
	var calendarType string
	timeStr := birthTime.Format("2006-01-02 15:04:05")

	// 根据日历类型设置显示文本，默认使用农历
	switch calendar {
	case utils.CALENDAR_SOLAR:
		calendarType = "公历"
	default:
		calendarType = "农历"
	}

	return fmt.Sprintf(BAZI_ANALYSIS_PROMPT,
		name, gender, timeStr, calendarType,
		baziInfo["year"], baziInfo["month"], baziInfo["day"], baziInfo["hour"])
}
