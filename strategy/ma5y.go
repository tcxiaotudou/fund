package strategy

import (
	"encoding/json"
	"fmt"
	"founds/constant"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Ma5y 国证A指 5年均线
func Ma5y() string {
	url := "https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=sz399317&scale=240&ma=no&datalen=1300" // 请求的URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("创建请求失败:", err)
		return ""
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return ""
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	var indexs []constant.Index
	err = json.Unmarshal(responseBody, &indexs)
	if err != nil {
		log.Println("json unmarshal error:", err)
		return ""
	}
	// 计算总和
	var sum = 0.0
	var today = 0.0
	var date = ""
	for _, item := range indexs[len(indexs)-1250:] {
		closeNum, _ := strconv.ParseFloat(item.Close, 64)
		sum += closeNum
		today = closeNum
		date = item.Date
	}
	// 计算均线
	movingAverage := sum / 1250
	diff := (today - movingAverage) / movingAverage * 100
	return fmt.Sprintf("「%s」 %.2f%%", date, diff)
}
