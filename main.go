package main

import (
	"encoding/json"
	"fmt"
	"founds/strategy"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var (
	rsiList   = map[string][]float64{}
	result    = map[string]interface{}{}
	rsiSource = map[string]string{
		"沪深三百":     "sh000300",
		"科创50":     "sh000688",
		"中证环保":     "sh000827",
		"中证1000":   "sh000852",
		"中证100":    "sh000903",
		"中证500":    "sh000905",
		"中证800":    "sh000906",
		"中证能源":     "sh000928",
		"中证消费":     "sh000932",
		"中证信息":     "sh000935",
		"中证体育":     "sz399804",
		"中证新能":     "sz399808",
		"中证国安":     "sz399813",
		"中证军工":     "sz399967",
		"中证传媒":     "sz399971",
		"中证国防":     "sz399973",
		"中证银行":     "sz399986",
		"中证酒":      "sz399987",
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
		"中证红利ETF":  "sh515080",
		"证券公司":     "sz399975",
		"恒生科技ETF":  "sh513180",
		"恒生互联网ETF": "sh513330",
		"中概互联ETF":  "sh510900",
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
		"上证50ETF":  "sh510050",
		"新能源ETF":   "sh516160",
		"科创成长ETF":  "sh588110",
		"香港证券ETF":  "sh513090",
		"法国ETF":    "sh513080",
		"德国ETF":    "sh513030",
		"标普生物ETF":  "sz161127",
		"日本ETF":    "sh513520",
		"游戏ETF":    "sh516010",
		"豆柏ETF":    "sz159985",
		"黄金ETF":    "sz159812",
		"油气ETF":    "sz159697",
		"创业板成长ETF": "sz159967",
		"机床ETF":    "sz159663",
		"信创ETF":    "sh562030",
		"通信ETF":    "sh515880",
	}
)

func init() {
	_, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return
	}
}

func main() {
	// guPercent()
	rsi()
	Ma5y()
	sendMail()
}

func rsi() {
	guozheng14Rsi := strategy.RsiGroup("sz399317", 14)[0]
	result["14日RSI（60点,65点,70点卖）"] = strconv.Itoa(int(guozheng14Rsi))
	guozhengRsiInt := int(guozheng14Rsi)
	key := "股债平衡建议"
	if guozhengRsiInt < 30 {
		result[key] = "9股1债"
	} else if guozhengRsiInt >= 30 && guozhengRsiInt < 35 {
		result[key] = "8股2债"
	} else if guozhengRsiInt >= 35 && guozhengRsiInt < 40 {
		result[key] = "7股3债"
	} else if guozhengRsiInt >= 40 && guozhengRsiInt < 50 {
		result[key] = "5股5债"
	} else if guozhengRsiInt >= 50 && guozhengRsiInt < 55 {
		result[key] = "4股6债"
	} else if guozhengRsiInt >= 55 && guozhengRsiInt < 60 {
		result[key] = "3股7债"
	} else if guozhengRsiInt >= 60 && guozhengRsiInt < 65 {
		result[key] = "2股8债"
	} else if guozhengRsiInt >= 65 {
		result[key] = "1股9债"
	}
	guozheng90Rsi := strategy.RsiGroup("sz399317", 90)[0]
	result["90日RSI（57 点和 70 点卖）"] = strconv.Itoa(int(guozheng90Rsi))

	for name, code := range rsiSource {
		rsiList[name+"("+code+")"] = strategy.RsiGroup(code, 14)
	}
}

// 邮件
func sendMail() {
	// 创建一个新的邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", "2290262044@qq.com")
	m.SetHeader("To", "2290262044@qq.com")
	m.SetHeader("Subject", fmt.Sprintf("每日行情（%s）", strategy.Date))
	content := `<div>
    <table border="1">
      <tr>
        <th>14日RSI</th>
        <th>股债比</th>
      </tr>
      <tr>
        <td>...</td>
        <td>9股-1债</td>
      </tr>
      <tr>
        <td>30 - 35</td>
        <td>8股-2债</td>
      </tr>
      <tr>
        <td>35 - 40</td>
        <td>7股-3债</td>
      </tr>
      <tr>
        <td>40 - 50</td>
        <td>5股-5债</td>
      </tr>
      <tr>
        <td>50 - 55</td>
        <td>4股-6债</td>
      </tr>
      <tr>
        <td>55 - 60</td>
        <td>3股-7债</td>
      </tr>
      <tr>
        <td>60 - 65</td>
        <td>2股-8债</td>
      </tr>
      <tr>
        <td>...</td>
        <td>1股-9债</td>
      </tr>
    </table>
  </div><br/>`

	risContent := `各行业RSI:<br/><div>
		<table border="1">
	`
	for name, rsiGroup := range rsiList {
		rsiValue := int(rsiGroup[0])
		if rsiValue >= 35 {
			continue
		}
		content := fmt.Sprintf(`
      <tr>
        <td>%s</td>
        <td>%d</td>
		<td>%s</td>
      </tr>`, name, rsiValue, fmt.Sprintf("(%s, %s, %s, %s)",
			fmt.Sprintf("%.2f", rsiGroup[1]),
			fmt.Sprintf("%.2f", rsiGroup[2]),
			fmt.Sprintf("%.2f", rsiGroup[3]),
			fmt.Sprintf("%.2f", rsiGroup[4])))
		risContent += content
	}
	risContent += `</table>
  	</div><br/>`

	for key, value := range result {
		content = content + fmt.Sprintf("<h2>%s: %s</h2>", key, value)
	}

	content += risContent

	content += `相关链接:<br/>`
	content += fmt.Sprintf(`<a href="%s" target="_blank">%s</a><br/>`, "https://youzhiyouxing.cn/data/market", "有知有行全市场温度")

	m.SetBody("text/html", content)
	// 创建一个新的SMTP拨号器
	d := gomail.NewDialer("smtp.qq.com", 587, "2290262044", "ehdrbzzctgvoebec")
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// 沪深300风险溢价
func guPercent() {
	var key = "股债百分位"
	url := "http://f.gushiyaowan.cn/v1/portfolio/stockBondYRDiff/list?indexCode=000300&bondCode=CN10YR&month=0&startDate=&endDate=" // 请求的URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	// 设置自定义请求头
	req.Header.Set("accessToken", "1e2e3c4cb0114a1797b276f07cc2b09e")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	var dataJson map[string]interface{}
	err = json.Unmarshal(responseBody, &dataJson)
	if err != nil {
		return
	}
	lists := dataJson["data"].(map[string]interface{})["list"].([]interface{})
	todayData := lists[len(lists)-1]
	data := todayData.(map[string]interface{})
	result[key] = strconv.FormatFloat(data["percentile"].(float64), 'f', -1, 64) + "%"
}

func sortByValue(m map[string]int) []string {
	// 将 map 数据转换为切片
	var pairs []struct {
		Key   string
		Value int
	}
	for k, v := range m {
		pairs = append(pairs, struct {
			Key   string
			Value int
		}{k, v})
	}

	// 自定义排序函数
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value < pairs[j].Value
	})

	// 提取排序后的键值并返回
	var sortedKeys []string
	for _, pair := range pairs {
		sortedKeys = append(sortedKeys, pair.Key)
	}

	return sortedKeys
}

// Ma5y 5年均线
func Ma5y() {
	url := "https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=sz399317&scale=1200&ma=no&datalen=1950" // 请求的URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	var index []strategy.Index
	err = json.Unmarshal(responseBody, &index)
	if err != nil {
		log.Println("json unmarshal error:", err)
	}
	lastClose := 0.0
	n := 52                                  // number of trading days in five years
	sum := 0.0                               // sum of the last n closing prices
	ma5Result := make([]float64, len(index)) // result slice
	for i, data := range index {
		x, _ := strconv.ParseFloat(data.Close, 64)
		lastClose = x
		if i < n-1 {
			// not enough data to calculate Ma5y, append zero
			sum += x
			ma5Result[i] = 0.0
		} else if i == n-1 {
			// just enough data to calculate the first Ma5y, append sum / n
			sum += x
			ma5Result[i] = sum / float64(n)
		} else {
			// more than enough data to calculate Ma5y, append (sum + x - data[i-n]) / n
			tmp, _ := strconv.ParseFloat(index[i-n].Close, 64)
			sum += x - tmp
			ma5Result[i] = sum / float64(n)
		}
	}
	avg := ma5Result[len(ma5Result)-1]
	result["5年均线"] = fmt.Sprintf("%v", strategy.Decimal((lastClose-avg)*100/avg)) + "%"
}
