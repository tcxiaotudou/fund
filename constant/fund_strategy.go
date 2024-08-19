package constant

type FundStrategy struct {
	Name            string `json:"name"`            // 基金名称
	Code            string `json:"code"`            // 基金代码
	PersonName      string `json:"PersonName"`      // 基金经理
	PersonYear      string `json:"personYear"`      // 经理管理时长
	Gm              string `json:"gm"`              // 规模
	YearTodayIncome string `json:"yearTodayIncome"` // 今年以来收益率

	Year5Income string `json:"year5Income"` // 5年年化收益率

	HcMaxYear5 string `json:"hcMax"` // 近 5 年最大回撤
	HcCurYear5 string `json:"hcCur"` // 近 5 年当前回撤

}
