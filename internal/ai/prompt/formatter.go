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
	var timeStr, calendarType string

	// 根据日历类型处理时间
	if calendar == "lunar" {
		// 如果输入是农历，直接使用农历时间
		timeStr = birthTime.Format("2006-01-02 15:04:05")
		calendarType = "农历"
	} else {
		// 如果输入是公历，直接使用公历时间
		timeStr = birthTime.Format("2006-01-02 15:04:05")
		calendarType = "公历"
	}

	// 获取八字排盘字符串
	baziStr := utils.FormatBaziString(baziInfo)

	// 计算起运时间
	qiYunTime := utils.CalculateQiYun(birthTime, gender, baziInfo)

	// 计算神煞信息
	yearShenSha := utils.GetShenSha(baziInfo["year"], "year")
	monthShenSha := utils.GetShenSha(baziInfo["month"], "month")
	dayShenSha := utils.GetShenSha(baziInfo["day"], "day")
	hourShenSha := utils.GetShenSha(baziInfo["hour"], "hour")

	// 计算纳音五行
	naYinWuXing := utils.GetNaYinWuXing(baziInfo)

	// 返回格式化的提示文本
	return fmt.Sprintf(BAZI_ANALYSIS_PROMPT,
		name, gender, timeStr, calendarType, baziStr,
		qiYunTime, yearShenSha, monthShenSha, dayShenSha,
		hourShenSha, naYinWuXing)
}
