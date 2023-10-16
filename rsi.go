package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
港股：https://stock.xueqiu.com/v5/stock/chart/kline.json?symbol=HKHSTECH&begin=1697337471570&period=day&type=before&count=-284&indicator=kline
*/

func GetRsi(code string) float64 {
	url := fmt.Sprintf("https://quotes.sina.cn/cn/api/jsonp_v2.php/=/CN_MarketDataService.getKLineData?symbol=%s&scale=120&ma=no&datalen=180", code)
	response, err := http.Get(url)
	if err != nil {
		log.Println("http get error:", err)
		return 0
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("io read error:", err)
		return 0
	}
	rex := regexp.MustCompile(`\[([\s\S]*?)]`)
	titleMatches := rex.FindAllSubmatch(data, -1)
	if titleMatches == nil {
		return 0
	}
	jsonStr := fmt.Sprintf("[%s]", string(titleMatches[0][1]))
	var index []Index
	err = json.Unmarshal([]byte(jsonStr), &index)
	if err != nil {
		log.Println("json unmarshal error:", err)
		return 0
	}
	rsiData := make([]float64, 0)
	for i, data := range index {
		if i != len(index)-1 && !strings.Contains(data.Date, "15:00:00") {
			continue
		}
		float, err := strconv.ParseFloat(data.Close, 64)
		if err != nil {
			log.Println("strconv parseFloat error:", err)
			return 0
		}
		rsiData = append(rsiData, float)
	}
	result := caRsi(rsiData, 14)
	return result[len(result)-1]
}

// 定义一个结构体，用来存储指数的收盘价和日期
type Index struct {
	Close string `json:"close"` // 收盘价
	Date  string `json:"day"`   // 日期
}

// 是否为收盘时间
func isTime15(timeObj time.Time) bool {
	targetTime := time.Date(timeObj.Year(), timeObj.Month(), timeObj.Day(), 15, 0, 0, 0, timeObj.Location())
	return timeObj.Equal(targetTime)
}

// 计算RSI
func caRsi(inReal []float64, inTimePeriod int) []float64 {
	outReal := make([]float64, len(inReal))

	if len(inReal) < inTimePeriod {
		return outReal
	}

	if inTimePeriod < 2 {
		return outReal
	}

	// variable declarations
	tempValue1 := 0.0
	tempValue2 := 0.0
	outIdx := inTimePeriod
	today := 0
	prevValue := inReal[today]
	prevGain := 0.0
	prevLoss := 0.0
	today++

	for i := inTimePeriod; i > 0; i-- {
		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
	}

	prevLoss /= float64(inTimePeriod)
	prevGain /= float64(inTimePeriod)

	if today > 0 {

		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++

	} else {

		for today < 0 {
			tempValue1 = inReal[today]
			tempValue2 = tempValue1 - prevValue
			prevValue = tempValue1
			prevLoss *= float64(inTimePeriod - 1)
			prevGain *= float64(inTimePeriod - 1)
			if tempValue2 < 0 {
				prevLoss -= tempValue2
			} else {
				prevGain += tempValue2
			}
			prevLoss /= float64(inTimePeriod)
			prevGain /= float64(inTimePeriod)
			today++
		}
	}

	for today < len(inReal) {

		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		prevLoss *= float64(inTimePeriod - 1)
		prevGain *= float64(inTimePeriod - 1)
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
		prevLoss /= float64(inTimePeriod)
		prevGain /= float64(inTimePeriod)
		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}
