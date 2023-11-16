package main

import (
	"fmt"
	"founds/constant"
	"founds/strategy"
	"gopkg.in/gomail.v2"
	"log"
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

	// ETF Rsi
	suggests := make([]constant.Suggest, 0)
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
		// rsi小于35 或者 rsi小于40 && 最低点大于35 或者 rsi离最低点小于5天 或者 rsi处于35-45之间 && 最高点大于70
		if (etfRsiData.Now < 35) ||
			(etfRsiData.Now <= 40 && etfRsiData.Low >= 35) ||
			(etfRsiData.NowToLowDays <= 5) ||
			(etfRsiData.High >= 70 && etfRsiData.Now >= 35 && etfRsiData.Now < 45) {
			rsiList[name+"("+code+")"] = etfRsiData
			// 纳入购买建议
			suggest := constant.Suggest{
				CodeName: name + "(" + code + ")",
				Now:      fmt.Sprintf("%.2f", etfRsiData.Now),
				Interval: fmt.Sprintf("(%s, %s, %s, %s)",
					fmt.Sprintf("%.2f", etfRsiData.High),
					fmt.Sprintf("%.2f", etfRsiData.TwoThirds),
					fmt.Sprintf("%.2f", etfRsiData.OneThirds),
					fmt.Sprintf("%.2f", etfRsiData.Low)),
				Remark: etfRsiData.Message,
				Time:   time.Now().Format("2006-01-02 15:01:05"),
			}
			suggests = append(suggests, suggest)
		}
	}
	SendMail(rsiList, result)
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
func SendMail(rsiList map[string]*strategy.RsiData, result map[string]interface{}) {
	// 创建一个新的邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", "2290262044@qq.com")
	m.SetHeader("To", "2290262044@qq.com")
	m.SetHeader("Subject", fmt.Sprintf("每日行情（%s）", strategy.Date))
	content := "<h4>行情数据：</h4><ul>"
	for key, value := range result {
		content = content + fmt.Sprintf("<li>%s: %s</li>", key, value)
	}
	content += "</ul><br/>"
	risContent := `<h4>买入建议:<h4/>
		<table border="1" style="border-collapse: collapse;">
		<tr>
			<th>名称</th>
			<th>当前</th>
			<th>区间</th>
			<th>备注</th>
			<th>时间</th>
        </tr>
	`
	for name, rsiData := range rsiList {
		content := fmt.Sprintf(`
		  <tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		  </tr>`, name, fmt.Sprintf("%.2f", rsiData.Now), fmt.Sprintf("(%s, %s, %s, %s)",
			fmt.Sprintf("%.2f", rsiData.High),
			fmt.Sprintf("%.2f", rsiData.TwoThirds),
			fmt.Sprintf("%.2f", rsiData.OneThirds),
			fmt.Sprintf("%.2f", rsiData.Low)), rsiData.Message, rsiData.Time.Format("2006-01-02 15:04:05"))
		risContent += content
	}
	risContent += `</table>
  	<br/>`
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
