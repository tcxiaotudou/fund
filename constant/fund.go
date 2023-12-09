package constant

type Fund struct {
	Name        string  `json:"name"`         // 基金名称
	Code        string  `json:"code"`         // 基金代码
	Score       int64   `json:"score"`        // PK分数
	Style       string  `json:"style"`        // 风格
	Scale       float64 `json:"scale"`        // 规模
	CreateYears string  `json:"create_years"` // 创建年限
	Retracement string  `json:"retracement"`  // 近一年最大回撤
	Sharpe      string  `json:"sharpe"`       // 近一年夏普比率
	Yield3      string  `json:"yield3"`       // 近3个月收益率
	Yield6      string  `json:"yield6"`       // 近6个月收益率
	Yield12     string  `json:"yield12"`      // 近一年收益率
	Reason      string  `json:"reason"`       // 入选理由
}
