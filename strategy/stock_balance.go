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

/*
*
90日RSI股债平衡建议
90 日 RSI 大多数时间在43 点(红色水平线)和 57 点(淡绿色水平线)之前盘整，泡沫大牛市会冲到 70 点
*/
func RsiStockBalance(rsi14 float64) string {
	if rsi14 < 30 {
		return "10股0债"
	} else if rsi14 >= 30 && rsi14 < 35 {
		return "9股1债"
	} else if rsi14 >= 35 && rsi14 < 43 {
		return "8股2债"
	} else if rsi14 >= 43 && rsi14 < 47.5 {
		return "7股3债"
	} else if rsi14 >= 47.5 && rsi14 < 52 {
		return "6股4债"
	} else if rsi14 >= 52 && rsi14 < 56.5 {
		return "5股5债"
	} else if rsi14 >= 56.5 && rsi14 < 61 {
		return "4股6债"
	} else if rsi14 >= 61 && rsi14 < 65.5 {
		return "3股7债"
	} else if rsi14 >= 65.5 {
		return "2股8债"
	}
	return ""
}
