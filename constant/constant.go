package constant

const GUO_ZHENG = "sz399317"

const GUO_ZHAI = "sh511090"

// Index 定义一个结构体，用来存储指数的收盘价和日期
type Index struct {
	Close string `json:"close"` // 收盘价
	Date  string `json:"day"`   // 日期
}

var (
	EtfGroups = map[string]string{
		"中证新能":     "sz399808",
		"中证军工":     "sz399967",
		"中证传媒":     "sz399971",
		"中证银行":     "sz399986",
		"中证医疗":     "sz399989",
		"中证白酒":     "sz399997",
		"中证煤炭":     "sz399998",
		"半导体ETF":   "sh512480",
		"机器人ETF":   "sz159770",
		"智能驾驶ETF":  "sh516520",
		"H股ETF":    "sh510900",
		"芯片ETF":    "sz159995",
		"人工智能ETF":  "sh515070",
		"中证红利ETF":  "sh000015",
		"恒生科技ETF":  "sh513180",
		"港股创新药ETF": "sh513120",
		"电力ETF":    "sh561560",
		"钢铁ETF":    "sh515210",
		"有色金属ETF":  "sh512400",
		"香港证券ETF":  "sh513090",
		"油气ETF":    "sz159697",
		"机床ETF":    "sz159663",
		"通信ETF":    "sh515880",
		"法国ETF":    "sh513080",
		"德国ETF":    "sh513030",
		"标普生物ETF":  "sz161127",
		"纳斯达克ETF":  "sh513300",
		"美国消费ETF":  "sz162415",
		"美国REIT":   "sz160140",
		"日本ETF":    "sh513520",
		"印度ETF":    "sz164824",
		"豆柏ETF":    "sz159985",
		"黄金ETF":    "sz159812",
		"原油ETF":    "sz160723",
	}
)
