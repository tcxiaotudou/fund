package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

func GetRsi() float64 {
	url := "https://quotes.sina.cn/cn/api/jsonp_v2.php/var%20_sh600036_240_1577432551767=/CN_MarketDataService.getKLineData?symbol=sz399317&scale=240&ma=no&datalen=63"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
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
		panic(err)
		return 0
	}
	rsiData := make([]float64, 0)
	for _, data := range index {
		float, err := strconv.ParseFloat(data.Close, 64)
		if err != nil {
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

// --------
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
			// outReal[outIdx] = math.Trunc((100.0*(prevGain/tempValue1)-0.11)*1e2+0.5) * 1e-2
			outReal[outIdx] = float64(int(100.0*(prevGain/tempValue1) - 0.11))
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}
