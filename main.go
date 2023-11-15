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
	guoZheng14RsiData := strategy.Rsi(constant.GUO_ZHENG, 14)
	result["股债平衡建议"] = strategy.RsiStockBalance(guoZheng14RsiData.Now)
	result[fmt.Sprintf("14日RSI（%s）", guoZheng14RsiData.Message)] = fmt.Sprintf("%.2f", guoZheng14RsiData.Now)

	guoZheng90RsiData := strategy.Rsi(constant.GUO_ZHENG, 90)
	result["90日RSI（57 点和 70 点卖）"] = fmt.Sprintf("%.2f", guoZheng90RsiData.Now)

	// ETF Rsi
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
		}
	}
	result["5年均线"] = strategy.Ma5y()
	SendMail(rsiList, result)
}

// SendMail 邮件
func SendMail(rsiList map[string]*strategy.RsiData, result map[string]interface{}) {
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
	for name, rsiData := range rsiList {
		content := fmt.Sprintf(`
		  <tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		  </tr>`, name, fmt.Sprintf("%.2f", rsiData.Now), fmt.Sprintf("(%s, %s, %s, %s)",
			fmt.Sprintf("%.2f", rsiData.High),
			fmt.Sprintf("%.2f", rsiData.TwoThirds),
			fmt.Sprintf("%.2f", rsiData.OneThirds),
			fmt.Sprintf("%.2f", rsiData.Low)), rsiData.Message)
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
