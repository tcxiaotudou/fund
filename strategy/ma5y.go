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

// 计算 1250 日均线
func calculateSMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0.0
	}

	sum := 0.0
	for _, price := range prices[len(prices)-period:] {
		sum += price
	}
	return sum / float64(period)
}

// 计算偏离度
func calculateDeviation(currentPrice, sma float64) float64 {
	return ((currentPrice - sma) / sma) * 100
}

func Ma5y() string {
	url := "https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=sz399317&scale=240&ma=no&datalen=1800" // 请求的URL
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

	prices := make([]float64, 0)
	var date = ""
	for _, item := range indexs {
		closeNum, _ := strconv.ParseFloat(item.Close, 64)
		prices = append(prices, closeNum)
		date = item.Date
	}

	period := 1250
	todayClosePrice := prices[len(prices)-1]

	// 计算 1250 日均线
	sma := calculateSMA(prices, period)
	fmt.Printf("1250日均线: %.2f\n", sma)

	// 计算今天收盘价与 1250 日均线的偏离度
	deviation := calculateDeviation(todayClosePrice, sma)
	// fmt.Printf("今天收盘价与 1250 日均线的偏离度: %.2f%%", deviation)

	return fmt.Sprintf("「%s」 %.2f%%", date, deviation)
}
