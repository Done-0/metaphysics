// Package utils 提供八字计算相关功能
// 创建者：Done-0
// 创建时间：2023-10-18
package utils

import (
	"time"

	lunarCalendar "github.com/6tail/lunar-go/calendar"
)

// 日历类型常量
const (
	CALENDAR_LUNAR = "lunar" // 农历
	CALENDAR_SOLAR = "solar" // 公历
)

// CalculateBazi 计算八字
// 参数：
//   - birthTime: 出生时间
//   - calendar: 日历类型 (lunar/solar)
//
// 返回值：
//   - map[string]string: 八字信息
func CalculateBazi(birthTime time.Time, calendar string) map[string]string {
	var lunar *lunarCalendar.Lunar

	// 根据日历类型获取农历对象，默认使用农历
	switch calendar {
	case CALENDAR_SOLAR:
		lunar = lunarCalendar.NewSolarFromDate(birthTime).GetLunar()
	default:
		year, month, day := birthTime.Date()
		hour, minute, second := birthTime.Clock()
		lunar = lunarCalendar.NewLunar(year, int(month), day, hour, minute, second)
	}

	eightChar := lunar.GetEightChar()

	// 四柱\天干\地支
	result := map[string]string{
		"year":  eightChar.GetYear(),
		"month": eightChar.GetMonth(),
		"day":   eightChar.GetDay(),
		"hour":  eightChar.GetTime(),

		"year_gan":  eightChar.GetYearGan(),
		"month_gan": eightChar.GetMonthGan(),
		"day_gan":   eightChar.GetDayGan(),
		"hour_gan":  eightChar.GetTimeGan(),

		"year_zhi":  eightChar.GetYearZhi(),
		"month_zhi": eightChar.GetMonthZhi(),
		"day_zhi":   eightChar.GetDayZhi(),
		"hour_zhi":  eightChar.GetTimeZhi(),
	}

	return result
}
