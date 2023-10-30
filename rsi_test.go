package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func TestRsi(t *testing.T) {
	date, rsiValue := GetRsi("sh000300", 14)
	fmt.Println(date, rsiValue)
}

// 5年均线
func Test5yearAVG(t *testing.T) {
	var key = "股债百分位"
	url := "https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=sz399317&scale=1200&ma=no&datalen=250" // 请求的URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	var index []Index
	err = json.Unmarshal(responseBody, &index)
	if err != nil {
		log.Println("json unmarshal error:", err)
	}

	sum := 0.0
	for _, item := range index {
		closeData, _ := strconv.ParseFloat(item.Close, 64)
		sum += closeData
	}

	log.Println("closeData: ", key)
}
