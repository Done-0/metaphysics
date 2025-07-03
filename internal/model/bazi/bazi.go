// Package bazi 八字模型，定义了八字计算和存储的相关结构体
// 创建者：Done-0
// 创建时间：2023-10-18
package bazi

import (
	"time"

	"github.com/Done-0/metaphysics/internal/model/base"
)

// BaziRecord 八字记录，存储用户的八字信息和分析结果
type BaziRecord struct {
	base.Base

	// 基本信息
	UserID    int64     `json:"user_id" gorm:"index"`                                                       // 用户 ID
	RequestID string    `json:"request_id" gorm:"size:50;index"`                                            // 请求 ID
	Name      string    `json:"name" gorm:"size:50"`                                                        // 姓名
	Gender    string    `json:"gender" gorm:"size:10;default:male;check:gender IN ('male', 'female')"`      // 性别 (male/female)
	BirthTime time.Time `json:"birth_time"`                                                                 // 出生时间
	Calendar  string    `json:"calendar" gorm:"size:10;default:lunar;check:calendar IN ('lunar', 'solar')"` // 日历类型 (lunar/solar)

	// 四柱干支
	YearPillar  string `json:"year_pillar" gorm:"size:20"`  // 年柱（干支）
	MonthPillar string `json:"month_pillar" gorm:"size:20"` // 月柱（干支）
	DayPillar   string `json:"day_pillar" gorm:"size:20"`   // 日柱（干支）
	HourPillar  string `json:"hour_pillar" gorm:"size:20"`  // 时柱（干支）

	// 天干
	YearGan  string `json:"year_gan" gorm:"size:10"`  // 年干
	MonthGan string `json:"month_gan" gorm:"size:10"` // 月干
	DayGan   string `json:"day_gan" gorm:"size:10"`   // 日干
	HourGan  string `json:"hour_gan" gorm:"size:10"`  // 时干

	// 地支
	YearZhi  string `json:"year_zhi" gorm:"size:10"`  // 年支
	MonthZhi string `json:"month_zhi" gorm:"size:10"` // 月支
	DayZhi   string `json:"day_zhi" gorm:"size:10"`   // 日支
	HourZhi  string `json:"hour_zhi" gorm:"size:10"`  // 时支

	// 天干五行
	YearGanWuXing  string `json:"year_gan_wu_xing" gorm:"size:10"`  // 年干五行
	MonthGanWuXing string `json:"month_gan_wu_xing" gorm:"size:10"` // 月干五行
	DayGanWuXing   string `json:"day_gan_wu_xing" gorm:"size:10"`   // 日干五行
	HourGanWuXing  string `json:"hour_gan_wu_xing" gorm:"size:10"`  // 时干五行

	// 纳音五行
	YearNaYin  string `json:"year_na_yin" gorm:"size:30"`  // 年柱纳音
	MonthNaYin string `json:"month_na_yin" gorm:"size:30"` // 月柱纳音
	DayNaYin   string `json:"day_na_yin" gorm:"size:30"`   // 日柱纳音
	TimeNaYin  string `json:"time_na_yin" gorm:"size:30"`  // 时柱纳音

	// 命理属性
	YinYang string `json:"yin_yang" gorm:"size:10"` // 阴阳属性（基于日干）
	WuXing  string `json:"wu_xing" gorm:"size:10"`  // 五行属性（基于日干）

	// 分析结果
	Analysis string `json:"analysis" gorm:"type:text"` // AI 分析结果
}
