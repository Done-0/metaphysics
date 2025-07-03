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
//   - baziInfo: 八字信息
//
// 返回值：
//   - string: 格式化的提示文本
func BuildBaziPrompt(name, gender string, birthTime time.Time, baziInfo map[string]string) string {
	birthTimeStr := birthTime.Format("2006-01-02 15:04:05")

	// 获取八字排盘字符串（年月日时四柱）
	baziStr := utils.FormatBaziString(baziInfo)

	// 提取天干
	yearGan := baziInfo["year_gan"]
	monthGan := baziInfo["month_gan"]
	dayGan := baziInfo["day_gan"]
	hourGan := baziInfo["hour_gan"]

	// 提取地支
	yearZhi := baziInfo["year_zhi"]
	monthZhi := baziInfo["month_zhi"]
	dayZhi := baziInfo["day_zhi"]
	hourZhi := baziInfo["hour_zhi"]

	// 提取命理属性
	yinYang := baziInfo["yin_yang"]
	wuXing := baziInfo["wu_xing"]

	return fmt.Sprintf(BAZI_ANALYSIS_PROMPT,
		name, gender, birthTimeStr, baziStr,
		yearGan, yearZhi, monthGan, monthZhi,
		dayGan, dayZhi, hourGan, hourZhi,
		yinYang, wuXing)
}
