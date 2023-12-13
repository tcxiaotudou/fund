package main

import (
	"fmt"
	"founds/constant"
	"founds/strategy"
	"gopkg.in/gomail.v2"
	"log"
	"sort"
	"time"
)

var (
	rsiList = map[string]*strategy.RsiData{}
	result  = map[string]interface{}{}
)

func main() {
	// 行情数据
	guoZheng14RsiData := strategy.Rsi(constant.GUO_ZHENG, 14)
	result[fmt.Sprintf("14日RSI（%s）", guoZheng14RsiData.Message)] = fmt.Sprintf("%.2f", guoZheng14RsiData.Now)
	guoZheng90RsiData := strategy.Rsi(constant.GUO_ZHENG, 90)
	result["90日RSI（57 点和 70 点卖）"] = fmt.Sprintf("%.2f", guoZheng90RsiData.Now)
	result["股债平衡建议"] = strategy.RsiStockBalance(guoZheng90RsiData.Now)
	result["5年均线"] = strategy.Ma5y()
	result["沪深300风险溢价"] = strategy.Stock300Balance()

	// ETF Rsi
	suggestions := make([]constant.Suggest, 0)
	for name, code := range constant.EtfGroups {
		time.Sleep(5 * time.Second)
		etfRsiData := strategy.Rsi(code, 14)
		if etfRsiData == nil {
			continue
		}
		log.Printf("(%s,%s), %v\n", name, code, etfRsiData)
		// 离最低点还有大于10天的差距，不做处理
		if etfRsiData.NowToLowDays > 10 {
			continue
		}
		// rsi小于30 或者 rsi离最低点小于5天 或者 rsi小于45 && 最低点大于35 && 最高点大于70
		if (etfRsiData.Now <= 30) ||
			(etfRsiData.High >= 70 && etfRsiData.High2NowLow <= 43 && etfRsiData.High2NowLow >= 38) {
			rsiList[name+"("+code+")"] = etfRsiData
			// 纳入购买建议
			suggestion := constant.Suggest{
				CodeName: name + "(" + code + ")",
				Now:      etfRsiData.Now,
				Interval: fmt.Sprintf("(%s, %s, %s, %s)",
					fmt.Sprintf("%.2f", etfRsiData.High),
					fmt.Sprintf("%.2f", etfRsiData.TwoThirds),
					fmt.Sprintf("%.2f", etfRsiData.OneThirds),
					fmt.Sprintf("%.2f", etfRsiData.Low)),
				Remark: etfRsiData.Message,
				Time:   time.Now().Format("2006-01-02 15:04:05"),
			}
			suggestions = append(suggestions, suggestion)
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Now < suggestions[j].Now
	})

	// 推荐基金
	funds := strategy.FundRank()

	SendMail(funds, suggestions, result)
}

/**
- 股债平衡建议: 3股7债
- 14日RSI: 59.17
- 90日RSI（57 点和 70 点卖）: 45.79
- 5年均线: -4.11%

| 编号 | 当前 | 区间 | 备注 |
| ------ | ------ | ------ | ------ |
| 中证银行(sz399986) | 33.62 | (68.86, 56.56, 44.26, 31.96) | 数据81天, 70以上有0天, 65以上有3天, 60以上有4天, 55以上有11天, 当前与最低点之间有4天 |

*/

// SendMail 邮件
func SendMail(funds []*constant.Fund, rsiList []constant.Suggest, result map[string]interface{}) {
	// 创建一个新的邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", "2290262044@qq.com")
	m.SetHeader("To", "2290262044@qq.com")
	m.SetHeader("Subject", fmt.Sprintf("每日行情（%s）", strategy.Date))
	content := "<h4>行情数据：</h4><ul>"
	for key, value := range result {
		content = content + fmt.Sprintf("<li>%s: %s</li>", key, value)
	}
	content += "</ul>"
	risContent := `<h4>场内ETF买入建议:</h4><br/>
		<table border="1" style="border-collapse: collapse;">
		<tr>
			<th>名称</th>
			<th>当前</th>
			<th>区间</th>
			<th>备注</th>
			<th>时间</th>
        </tr>
	`
	for _, rsiData := range rsiList {
		content := fmt.Sprintf(`
		  <tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		  </tr>`, rsiData.CodeName, fmt.Sprintf("%.2f", rsiData.Now), rsiData.Interval, rsiData.Remark, rsiData.Time)
		risContent += content
	}
	risContent += `</table><br/>`
	content += risContent

	// 基金排行榜
	fundContent := `<h4>场外基金推荐:</h4><br/>
		<table border="1" style="border-collapse: collapse;">
		<tr>
			<th>名称</th>
			<th>规模</th>
			<th>近6个月最大回撤</th>
			<th>夏普率</th>
			<th>近3个月收益</th>
			<th>近6个月收益</th>
			<th>近一年收益</th>
			<th>评分</th>
        </tr>
	`
	for _, fund := range funds {
		content := fmt.Sprintf(`
		  <tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
		  </tr>`, fmt.Sprintf("%s(%s)", fund.Name, fund.Code),
			fmt.Sprintf("%.2f亿", fund.Scale),
			fmt.Sprintf("%s%%", fund.Retracement),
			fund.Sharpe,
			fmt.Sprintf("%s%%", fund.Yield3),
			fmt.Sprintf("%s%%", fund.Yield6),
			fmt.Sprintf("%s%%", fund.Yield12),
			fund.Score)
		fundContent += content
	}
	fundContent += `</table><br/>`
	content += fundContent

	// 相关链接
	content += "<h4>相关链接:</h4><ul>"
	content += fmt.Sprintf(`<li><a href="%s" target="_blank">%s</a></li>`, "https://youzhiyouxing.cn/data/market", "有知有行全市场温度")

	content += "</ul>"

	m.SetBody("text/html", content)
	// 创建一个新的SMTP拨号器
	d := gomail.NewDialer("smtp.qq.com", 587, "2290262044", "ehdrbzzctgvoebec")
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
