package strategy

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

// 沪深三百股债平衡
func Stock300Balance() (value string) {
	url := "https://api.jiucaishuo.com/gz/gz/fed"

	// 构建请求数据
	requestData := map[string]interface{}{
		"gu_code": "000300.SH",
		"year":    5,
		"type":    "h5",
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Fatal(err)
	}
	// 发送POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 输出响应结果
	var dataJson map[string]interface{}
	err = json.Unmarshal(responseBody, &dataJson)
	if err != nil || dataJson["message"].(string) != "success" {
		return
	}
	data := dataJson["data"].(map[string]interface{})["new"].(map[string]interface{})
	return strconv.FormatFloat(data["percent"].(float64), 'f', -1, 64) + "%"
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
