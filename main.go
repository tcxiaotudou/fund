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
	rsiList   = map[string]int{}
	result    = map[string]interface{}{}
	date      = ""
	rsiSource = map[string]string{
		"科创50":   "sh000688",
		"中证环保":   "sh000827",
		"中证1000": "sh000852",
		"中证100":  "sh000903",
		"中证500":  "sh000905",
		"中证800":  "sh000906",
		"中证能源":   "sh000928",
		"中证消费":   "sh000932",
		"中证信息":   "sh000935",
		"中证体育":   "sz399804",
		"中证新能":   "sz399808",
		"中证国安":   "sz399813",
		"中证军工":   "sz399967",
		"中证传媒":   "sz399971",
		"中证国防":   "sz399973",
		"中证银行":   "sz399986",
		"中证酒":    "sz399987",
		"中证医疗":   "sz399989",
		"中证白酒":   "sz399997",
		"中证煤炭":   "sz399998",
	}
)

func init() {
	_, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return
	}
}

func main() {
	// sar("000300")
	fear()
	// guPercent()
	rsi()
	sendMail()
}

func rsi() {
	guozhengRsi := GetRsi("sz399317")
	result["14日RSI"] = strconv.Itoa(int(guozhengRsi))
	guozhengRsiInt := int(guozhengRsi)
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

	for name, code := range rsiSource {
		rsi := GetRsi(code)
		rsiList[name+"("+code+")"] = int(rsi)
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
	result[key] = fmt.Sprintf("%s - %v", statusStr, num)
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

	risContent := `各行业RSI: <br/><div>
		<table border="1">
	`

	for name, rsiValue := range rsiList {
		if rsiValue > 40 {
			continue
		}
		content := fmt.Sprintf(`
      <tr>
        <td>%s</td>
        <td>%s</td>
      </tr>`, name, rsiValue)
		risContent += content
	}
	risContent += `</table>
  	</div><br/>`

	for key, value := range result {
		content = content + fmt.Sprintf("<h2>%s: %s</h2><br/>", key, value)
	}

	content += risContent

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
