package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	query := bytes.NewReader([]byte(`{"category_type":"cz","pe_category":"fed","region":"","year":10,"type":"pc","version":"2.2.7","authtoken":"","act_time":1694067833457,"tirgkjfs":"53","abiokytke":"77","u54rg5d":"81","kf54ge7":"0","tiklsktr4":"3","lksytkjh":"42c7","sbnoywr":"44","bgd7h8tyu54":"cc","y654b5fs3tr":"7","bioduytlw":"8","bd4uy742":"e","h67456y":"042","bvytikwqjk":"cc","ngd4uy551":"42","bgiuytkw":"24","nd354uy4752":"f","ghtoiutkmlg":"7b3","bd24y6421f":"42","tbvdiuytk":"0","ibvytiqjek":"da","jnhf8u5231":"24","fjlkatj":"819","hy5641d321t":"2e","iogojti":"2","ngd4yut78":"b3","nkjhrew":"e","yt447e13f":"b","n3bf4uj7y7":"2","nbf4uj7y432":"77","yi854tew":"5f","h13ey474":"5f0","quikgdky":"a7"}`))

	response, err := http.Post("https://api.jiucaishuo.com/v2/guzhi/fedshowbasedata",
		"application/json", query)
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
	lists := dataJson["data"].(map[string]interface{})["lists"].([]interface{})
	for _, list := range lists {
		data := list.(map[string]interface{})
		if data["gu_code"] == "000300.SH" {
			result[key] = data["gu_date"].(string) + " / " + strconv.Itoa(int(100.0-(data["gu_pe"].(float64)))) + "%"
		}
	}
}
