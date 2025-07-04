// Package bazi 提供八字分析相关的视图对象
// 创建者：Done-0
// 创建时间：2023-10-18
package bazi

// BaziResponse 八字分析响应
// @Description 八字分析响应
// @Property    Name      string true "姓名"
// @Property    Gender    string true "性别"
// @Property    Calendar  string true "日历类型 (lunar/solar)"
// @Property    YearPillar  string true "年柱（干支）"
// @Property    MonthPillar string true "月柱（干支）"
// @Property    DayPillar   string true "日柱（干支）"
// @Property    HourPillar  string true "时柱（干支）"
// @Property    YearGan     string true "年干"
// @Property    MonthGan    string true "月干"
// @Property    DayGan      string true "日干"
// @Property    HourGan     string true "时干"
// @Property    YearZhi     string true "年支"
// @Property    MonthZhi    string true "月支"
// @Property    DayZhi      string true "日支"
// @Property    HourZhi     string true "时支"
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
}

// BaziListResponse 八字列表响应
// @Description 八字列表响应
// @Property total     int64 true "总条数"
// @Property pageNo    int   true "当前页"
// @Property pageSize  int   true "当前分页记录数"
// @Property list      []BaziResponse true "分页内容"
type BaziListResponse struct {
	Total    int64           `json:"total"`    // 总条数
	PageNo   int             `json:"pageNo"`   // 当前页
	PageSize int             `json:"pageSize"` // 当前分页记录数
	List     []*BaziResponse `json:"list"`     // 分页内容
}
