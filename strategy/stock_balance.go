package strategy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// 沪深三百股债平衡
func stock300Balance() (key, value string) {
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
	return "股债百分位", strconv.FormatFloat(data["percentile"].(float64), 'f', -1, 64) + "%"
}

// RsiStockBalance 14日RSI股债平衡建议
func RsiStockBalance(rsi14 float64) string {
	if rsi14 < 30 {
		return "9股1债"
	} else if rsi14 >= 30 && rsi14 < 35 {
		return "8股2债"
	} else if rsi14 >= 35 && rsi14 < 40 {
		return "7股3债"
	} else if rsi14 >= 40 && rsi14 < 50 {
		return "5股5债"
	} else if rsi14 >= 50 && rsi14 < 55 {
		return "4股6债"
	} else if rsi14 >= 55 && rsi14 < 60 {
		return "3股7债"
	} else if rsi14 >= 60 && rsi14 < 65 {
		return "2股8债"
	} else if rsi14 >= 65 {
		return "1股9债"
	}
	return ""
}
