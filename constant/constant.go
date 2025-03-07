package constant

const GUO_ZHENG = "sz399317"

// Index 定义一个结构体，用来存储指数的收盘价和日期
type Index struct {
	Close string `json:"close"` // 收盘价
	Date  string `json:"day"`   // 日期
}

var (
	EtfGroups = map[string]string{
		"沪深三百":     "sh000300",
		"科创50":     "sh000688",
		"中证环保":     "sh000827",
		"中证2000":   "sh563300",
		"中证1000":   "sh000852",
		"中证500":    "sh000905",
		"中证能源":     "sh000928",
		"中证体育":     "sz399804",
		"中证新能":     "sz399808",
		"中证军工":     "sz399967",
		"中证传媒":     "sz399971",
		"中证银行":     "sz399986",
		"中证医疗":     "sz399989",
		"中证白酒":     "sz399997",
		"中证煤炭":     "sz399998",
		"半导体ETF":   "sh512480",
		"中药ETF":    "sz159647",
		"创业板指":     "sz399006",
		"创新药ETF":   "sz159992",
		"智能汽车ETF":  "sh515250",
		"汽车ETF":    "sh516110",
		"科技ETF":    "sh515000",
		"大数据ETF":   "sh515400",
		"机器人ETF":   "sz159770",
		"智能驾驶ETF":  "sh516520",
		"H股ETF":    "sh510900",
		"国证油气":     "sz399439",
		"芯片ETF":    "sz159995",
		"人工智能ETF":  "sh515070",
		"中证红利ETF":  "sh000015",
		"证券公司":     "sz399975",
		"恒生科技ETF":  "sh513180",
		"恒生互联网ETF": "sh513330",
		"中概互联ETF":  "sh513050",
		"港股创新药ETF": "sh513120",
		"教育ETF":    "sh513360",
		"房地产ETF":   "sh512200",
		"绿电50ETF":  "sh561170",
		"基建工程LOF":  "sz165525",
		"电力ETF":    "sh561560",
		"纳斯达克ETF":  "sh513300",
		"标普500ETF": "sh513500",
		"化工ETF":    "sz159870",
		"钢铁ETF":    "sh515210",
		"饮食ETF":    "sz159736",
		"有色金属ETF":  "sh512400",
		"上证50ETF":  "sh000016",
		"香港证券ETF":  "sh513090",
		"法国ETF":    "sh513080",
		"德国ETF":    "sh513030",
		"标普生物ETF":  "sz161127",
		"标普科技ETF":  "sz161128",
		"美国消费ETF":  "sz162415",
		"美国REIT":   "sz160140",
		"中国互联网ETF": "sz164906",
		"日本ETF":    "sh513520",
		"印度ETF":    "sz164824",
		"游戏ETF":    "sh516010",
		"豆柏ETF":    "sz159985",
		"黄金ETF":    "sz159812",
		"油气ETF":    "sz159697",
		"机床ETF":    "sz159663",
		"信创ETF":    "sh562030",
		"通信ETF":    "sh515880",
		"原油ETF":    "sz160723",
		"北证50":     "bj899050",
	}
)
