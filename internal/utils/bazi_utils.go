// Package utils 提供八字计算相关功能
// 创建者：Done-0
// 创建时间：2023-10-18
package utils

import (
	"strings"
	"time"

	lunarCalendar "github.com/6tail/lunar-go/calendar"
)

// 日历类型常量
const (
	CALENDAR_LUNAR = "lunar" // 农历
	CALENDAR_SOLAR = "solar" // 阳历/公历
)

// 天干阴阳映射表：奇数为阳，偶数为阴
var ganYinYangMap = map[string]string{
	"甲": "阳", "乙": "阴",
	"丙": "阳", "丁": "阴",
	"戊": "阳", "己": "阴",
	"庚": "阳", "辛": "阴",
	"壬": "阳", "癸": "阴",
}

// 天干五行映射表
var ganWuXingMap = map[string]string{
	"甲": "木", "乙": "木",
	"丙": "火", "丁": "火",
	"戊": "土", "己": "土",
	"庚": "金", "辛": "金",
	"壬": "水", "癸": "水",
}

// CalculateBazi 计算八字
// 参数：
//   - birthTime: 出生时间
//   - calendar: 日历类型 (lunar/solar)
//
// 返回值：
//   - map[string]string: 八字信息
func CalculateBazi(birthTime time.Time, calendar string) map[string]string {
	var lunar *lunarCalendar.Lunar
	if calendar == CALENDAR_SOLAR {
		lunar = lunarCalendar.NewSolarFromDate(birthTime).GetLunar()
	} else {
		// 默认使用农历
		year, month, day := birthTime.Date()
		hour, minute, second := birthTime.Clock()
		lunar = lunarCalendar.NewLunar(year, int(month), day, hour, minute, second)
	}

	// 获取八字信息
	eightChar := lunar.GetEightChar()

	// 获取日柱天干
	dayGan := eightChar.GetDayGan()

	// 根据天干确定阴阳属性
	yinyang := ganYinYangMap[dayGan]

	// 根据天干确定五行属性
	wuxing := ganWuXingMap[dayGan]

	// 返回八字信息
	return map[string]string{
		// 四柱完整表示
		"year":  eightChar.GetYear(),
		"month": eightChar.GetMonth(),
		"day":   eightChar.GetDay(),
		"hour":  eightChar.GetTime(),

		// 天干
		"year_gan":  eightChar.GetYearGan(),
		"month_gan": eightChar.GetMonthGan(),
		"day_gan":   eightChar.GetDayGan(),
		"hour_gan":  eightChar.GetTimeGan(),

		// 地支
		"year_zhi":  eightChar.GetYearZhi(),
		"month_zhi": eightChar.GetMonthZhi(),
		"day_zhi":   eightChar.GetDayZhi(),
		"hour_zhi":  eightChar.GetTimeZhi(),

		// 命理属性
		"yin_yang": yinyang,
		"wu_xing":  wuxing,
	}
}

// FormatBaziString 格式化八字字符串
// 参数：
//   - baziInfo: 八字信息
//
// 返回值：
//   - string: 格式化的八字字符串
func FormatBaziString(baziInfo map[string]string) string {
	keys := []string{"year", "month", "day", "hour"}
	values := make([]string, 0, 4)

	for _, key := range keys {
		if value, exists := baziInfo[key]; exists && value != "" {
			values = append(values, value)
		}
	}

	// 如果有完整的四柱，直接返回
	if len(values) == 4 {
		return strings.Join(values, " ")
	}

	// 否则尝试从天干地支组合
	values = make([]string, 0, 4)
	pairs := []struct{ gan, zhi string }{
		{"year_gan", "year_zhi"},
		{"month_gan", "month_zhi"},
		{"day_gan", "day_zhi"},
		{"hour_gan", "hour_zhi"},
	}

	for _, pair := range pairs {
		gan, hasGan := baziInfo[pair.gan]
		zhi, hasZhi := baziInfo[pair.zhi]

		if !hasGan || !hasZhi || gan == "" || zhi == "" {
			return "" // 任何一柱不完整，返回空字符串
		}

		values = append(values, gan+zhi)
	}

	return strings.Join(values, " ")
}
