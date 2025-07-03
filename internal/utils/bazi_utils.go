// Package utils 提供八字计算相关功能
// 创建者：Done-0
// 创建时间：2023-10-18
package utils

import (
	"fmt"
	"strings"
	"time"

	lunarCalendar "github.com/6tail/lunar-go/calendar"
)

// 日历类型常量
const (
	CALENDAR_LUNAR = "lunar" // 农历
	CALENDAR_SOLAR = "solar" // 阳历/公历
)

// 天干阴阳映射
var ganYinYangMap = map[string]string{
	"甲": "阳",
	"乙": "阴",
	"丙": "阳",
	"丁": "阴",
	"戊": "阳",
	"己": "阴",
	"庚": "阳",
	"辛": "阴",
	"壬": "阳",
	"癸": "阴",
}

// 年柱神煞映射
var yearShenShaMap = map[string][]string{
	"子": {"岁破", "劫煞"},
	"丑": {"灾煞", "灾杀"},
	"寅": {"伏位", "天医"},
	"卯": {"绝命", "元武"},
	"辰": {"脱劫", "勾陈"},
	"巳": {"危困", "朱雀"},
	"午": {"大耗", "天刑"},
	"未": {"伏藏", "天姚"},
	"申": {"进神", "白虎"},
	"酉": {"退神", "玄武"},
	"戌": {"阳刃", "天空"},
	"亥": {"阴煞", "地柜"},
}

// 月柱神煞映射
var monthShenShaMap = map[string][]string{
	"子": {"大败", "截路"},
	"丑": {"自刑", "蜕变"},
	"寅": {"三刑", "福星"},
	"卯": {"天罗", "天德"},
	"辰": {"月厌", "化解"},
	"巳": {"地网", "解神"},
	"午": {"阴阳差错", "天哭"},
	"未": {"阴差阳错", "天虚"},
	"申": {"孤辰", "天乙"},
	"酉": {"寡宿", "文曲"},
	"戌": {"命危", "华盖"},
	"亥": {"歧路", "将星"},
}

// 日柱神煞映射
var dayShenShaMap = map[string][]string{
	"子": {"驿马", "天贵"},
	"丑": {"亡神", "天厨"},
	"寅": {"时墓", "文昌"},
	"卯": {"时胎", "金舆"},
	"辰": {"官符", "国印"},
	"巳": {"旬空", "学堂"},
	"午": {"血支", "将军"},
	"未": {"灾煞", "夫星"},
	"申": {"天罗", "魁罡"},
	"酉": {"地网", "羊刃"},
	"戌": {"河魁", "孤鸾"},
	"亥": {"披麻", "寡宿"},
}

// 时柱神煞映射
var hourShenShaMap = map[string][]string{
	"子": {"天梁", "红鸾"},
	"丑": {"文星", "流霞"},
	"寅": {"地空", "天马"},
	"卯": {"天德", "解神"},
	"辰": {"白虎", "天福"},
	"巳": {"天狗", "福德"},
	"午": {"天刑", "恩光"},
	"未": {"天姚", "天寿"},
	"申": {"月德", "三台"},
	"酉": {"红艳", "八座"},
	"戌": {"天喜", "贯索"},
	"亥": {"咸池", "龙德"},
}

// CalculateBazi 计算八字
// 参数：
//   - birthTime: 出生时间
//   - calendar: 日历类型 (lunar/solar)
//
// 返回值：
//   - map[string]string: 八字信息
func CalculateBazi(birthTime time.Time, calendar string) map[string]string {
	// 获取农历对象
	var lunar *lunarCalendar.Lunar
	if calendar == CALENDAR_SOLAR {
		lunar = lunarCalendar.NewSolarFromDate(birthTime).GetLunar()
	} else {
		year, month, day := birthTime.Date()
		hour, minute, second := birthTime.Clock()
		lunar = lunarCalendar.NewLunar(year, int(month), day, hour, minute, second)
	}

	// 获取八字信息
	eightChar := lunar.GetEightChar()

	// 获取日干
	dayGan := eightChar.GetDayGan()

	// 构建结果map
	result := map[string]string{
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

		// 天干五行
		"year_gan_wu_xing":  eightChar.GetYearWuXing(),
		"month_gan_wu_xing": eightChar.GetMonthWuXing(),
		"day_gan_wu_xing":   eightChar.GetDayWuXing(),
		"hour_gan_wu_xing":  eightChar.GetTimeWuXing(),

		// 纳音
		"year_na_yin":  eightChar.GetYearNaYin(),
		"month_na_yin": eightChar.GetMonthNaYin(),
		"day_na_yin":   eightChar.GetDayNaYin(),
		"time_na_yin":  eightChar.GetTimeNaYin(),

		// 命理属性
		"wu_xing":  eightChar.GetDayWuXing(), // 日主五行
		"yin_yang": ganYinYangMap[dayGan],    // 日主阴阳
	}

	return result
}

// FormatBaziString 格式化八字字符串
// 参数：
//   - baziInfo: 八字信息
//
// 返回值：
//   - string: 格式化的八字字符串
func FormatBaziString(baziInfo map[string]string) string {
	// 尝试直接获取完整的四柱
	pillars := []string{"year", "month", "day", "hour"}
	values := make([]string, 0, 4)

	// 先尝试获取完整的干支表示
	allComplete := true
	for _, key := range pillars {
		if value, exists := baziInfo[key]; exists && value != "" {
			values = append(values, value)
		} else {
			allComplete = false
			break
		}
	}

	// 如果获取到完整的四柱，直接返回
	if allComplete && len(values) == 4 {
		return strings.Join(values, " ")
	}

	// 否则尝试从天干地支组合
	values = make([]string, 0, 4)
	for _, pillar := range pillars {
		gan := baziInfo[pillar+"_gan"]
		zhi := baziInfo[pillar+"_zhi"]

		// 如果任何一个为空，返回空字符串
		if gan == "" || zhi == "" {
			return ""
		}

		values = append(values, gan+zhi)
	}

	return strings.Join(values, " ")
}

// GetLunarTimeString 获取农历时间字符串
// 参数：
//   - birthTime: 出生时间
//   - calendar: 日历类型 (lunar/solar)
//
// 返回值：
//   - string: 农历时间字符串
func GetLunarTimeString(birthTime time.Time, calendar string) string {
	// 获取农历对象
	var lunar *lunarCalendar.Lunar
	if calendar == CALENDAR_SOLAR {
		lunar = lunarCalendar.NewSolarFromDate(birthTime).GetLunar()
	} else {
		year, month, day := birthTime.Date()
		hour, minute, second := birthTime.Clock()
		lunar = lunarCalendar.NewLunar(year, int(month), day, hour, minute, second)
	}

	// 直接返回格式化的农历时间字符串
	return fmt.Sprintf("%d年%s月%s %s时",
		lunar.GetYear(),
		lunar.GetMonthInChinese(),
		lunar.GetDayInChinese(),
		lunar.GetTime().GetZhi())
}

// GetSolarTimeString 获取公历时间字符串
// 参数：
//   - birthTime: 出生时间
//   - calendar: 日历类型 (lunar/solar)
//
// 返回值：
//   - string: 公历时间字符串
func GetSolarTimeString(birthTime time.Time, calendar string) string {
	// 获取公历对象
	var solar *lunarCalendar.Solar
	if calendar == CALENDAR_LUNAR {
		// 如果输入是农历，则转换为公历
		year, month, day := birthTime.Date()
		hour, minute, second := birthTime.Clock()
		lunar := lunarCalendar.NewLunar(year, int(month), day, hour, minute, second)
		solar = lunar.GetSolar()
	} else {
		// 如果输入已经是公历，直接使用
		solar = lunarCalendar.NewSolarFromDate(birthTime)
	}

	// 返回格式化的公历时间字符串
	return fmt.Sprintf("%d年%d月%d日 %d时%d分",
		solar.GetYear(),
		solar.GetMonth(),
		solar.GetDay(),
		solar.GetHour(),
		solar.GetMinute())
}

// CalculateQiYun 计算起运时间
// 参数：
//   - birthTime: 出生时间
//   - gender: 性别 ("男"/"女")
//   - baziInfo: 八字信息
//
// 返回值：
//   - string: 起运时间字符串
func CalculateQiYun(birthTime time.Time, gender string, baziInfo map[string]string) string {
	// 获取农历对象
	lunar := lunarCalendar.NewSolarFromDate(birthTime).GetLunar()

	// 获取八字
	eightChar := lunar.GetEightChar()

	// 转换性别参数 (男=1, 女=0)
	genderCode := 0
	if gender == "男" {
		genderCode = 1
	}

	// 使用lunar-go库计算大运
	yun := eightChar.GetYun(genderCode)

	// 获取起运时间
	startSolar := yun.GetStartSolar()

	// 计算年龄
	birthYear := birthTime.Year()
	startYear := startSolar.GetYear()
	age := startYear - birthYear

	// 计算月份差
	birthMonth := int(birthTime.Month())
	startMonth := startSolar.GetMonth()
	months := 0

	if startMonth > birthMonth {
		months = startMonth - birthMonth
	} else if startMonth < birthMonth {
		age--
		months = 12 + startMonth - birthMonth
	}

	// 格式化输出
	return fmt.Sprintf("%d岁%d个月 %d年%d月%d日",
		age, months, startSolar.GetYear(), startSolar.GetMonth(), startSolar.GetDay())
}

// GetShenSha 获取神煞
// 参数：
//   - pillar: 柱（干支组合，如"甲子"）
//   - pillarType: 柱类型（year/month/day/hour）
//
// 返回值：
//   - string: 神煞字符串
func GetShenSha(pillar string, pillarType string) string {
	if pillar == "" || len(pillar) < 2 {
		return ""
	}

	// 获取地支（取干支组合的第二个字）
	zhi := string(pillar[1])

	// 根据柱类型获取神煞
	var shenShaList []string
	switch pillarType {
	case "year":
		shenShaList = yearShenShaMap[zhi]
	case "month":
		shenShaList = monthShenShaMap[zhi]
	case "day":
		shenShaList = dayShenShaMap[zhi]
	case "hour":
		shenShaList = hourShenShaMap[zhi]
	}

	return strings.Join(shenShaList, "、")
}

// GetNaYinWuXing 获取纳音五行
// 参数：
//   - baziInfo: 八字信息
//
// 返回值：
//   - string: 纳音五行字符串
func GetNaYinWuXing(baziInfo map[string]string) string {
	pillars := []string{"year", "month", "day", "hour"}
	naYinParts := make([]string, 0, 4)

	for _, pillar := range pillars {
		naYin := baziInfo[pillar+"_na_yin"]
		ganZhi := baziInfo[pillar]

		if naYin != "" && ganZhi != "" {
			naYinParts = append(naYinParts, fmt.Sprintf("%s(%s)", ganZhi, naYin))
		}
	}

	return strings.Join(naYinParts, "、")
}
