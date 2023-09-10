package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	result = map[string]interface{}{}
	date   = ""
)

func main() {
	// sar("000300")
	fear()
	guPercent()
	rsi()
	sendMail()
}

func rsi() {
	rsi := GetRsi()
	result["14日RSI"] = strconv.Itoa(int(rsi))
	rsiInt := int(rsi)
	key := "股债平衡建议"
	if rsiInt < 30 {
		result[key] = "9股1债"
	} else if rsiInt >= 30 && rsiInt < 35 {
		result[key] = "8股2债"
	} else if rsiInt >= 35 && rsiInt < 40 {
		result[key] = "7股3债"
	} else if rsiInt >= 40 && rsiInt < 60 {
		result[key] = "5股5债"
	} else if rsiInt >= 60 && rsiInt < 65 {
		result[key] = "2股8债"
	} else if rsiInt >= 65 {
		result[key] = "1股9债"
	}
}

// 恐贪指数
func fear() {
	var key = "恐贪指数"
	response, err := http.Post("https://api.jiucaishuo.com/v2/kjtl/getbasedata",
		"application/json", nil)
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
	num := dataJson["data"].(map[string]interface{})["num"]
	currentTime := dataJson["data"].(map[string]interface{})["current_time"]
	date = currentTime.(string)
	statusStr := dataJson["data"].(map[string]interface{})["status_str"]
	result[key] = fmt.Sprintf("%s / %s - %v", currentTime, statusStr, num)
}

// 周线SAR
func sar(code string) {
	url := fmt.Sprintf("http://webquoteklinepic.eastmoney.com/GetPic.aspx?nid=1.%s&type=W&unitWidth=-6&ef=EXTENDED_SAR&formula=RSI&AT=1&imageType=KXL", code)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fileName := fmt.Sprintf("%s_sar_image.jpg", code)
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	log.Println("Image downloaded successfully.")
}

// 邮件
func sendMail() {
	// 创建一个新的邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", "2290262044@qq.com")
	m.SetHeader("To", "2290262044@qq.com")
	m.SetHeader("Subject", fmt.Sprintf("每日行情（%s）", date))
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
        <td>40 - 60</td>
        <td>5股-5债</td>
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

	for key, value := range result {
		content = content + fmt.Sprintf("<h2>%s: %s</h2><br/>", key, value)
	}
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
	var key = "沪深300风险溢价"
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
	// 转换为时间
	t := time.Unix((int64(data["date"].(float64)))/1000, 0)
	// 只保留年月日
	date := t.Format("2006-01-02")
	result[key] = date + " / " + strconv.FormatFloat(data["percentile"].(float64), 'f', -1, 64) + "%"
}
