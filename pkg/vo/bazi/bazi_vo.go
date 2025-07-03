// Package bazi 提供八字分析相关的视图对象
// 创建者：Done-0
// 创建时间：2023-10-18
package bazi

// BaziResponse 八字分析响应
type BaziResponse struct {
	// 基本信息
	RequestID string `json:"request_id"` // 请求 ID
	Name      string `json:"name"`       // 姓名
	Gender    string `json:"gender"`     // 性别
	Calendar  string `json:"calendar"`   // 日历类型 (lunar/solar)

	// 四柱干支
	YearPillar  string `json:"year_pillar"`  // 年柱（干支）
	MonthPillar string `json:"month_pillar"` // 月柱（干支）
	DayPillar   string `json:"day_pillar"`   // 日柱（干支）
	HourPillar  string `json:"hour_pillar"`  // 时柱（干支）

	// 天干
	YearGan  string `json:"year_gan"`  // 年干
	MonthGan string `json:"month_gan"` // 月干
	DayGan   string `json:"day_gan"`   // 日干
	HourGan  string `json:"hour_gan"`  // 时干

	// 地支
	YearZhi  string `json:"year_zhi"`  // 年支
	MonthZhi string `json:"month_zhi"` // 月支
	DayZhi   string `json:"day_zhi"`   // 日支
	HourZhi  string `json:"hour_zhi"`  // 时支

	// 命理属性
	YinYang string `json:"yin_yang"` // 阴阳属性
	WuXing  string `json:"wu_xing"`  // 五行属性

	// 分析结果
	Analysis string `json:"analysis"` // AI 分析结果
}
