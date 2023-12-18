package strategy

import (
	"encoding/json"
	"fmt"
	"founds/constant"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RsiData RSI数据
type RsiData struct {
	Days       int     // RSI数据的天数
	Now        float64 // 当前RSI
	LatestHigh float64 // RSI最近一次最高点
	High       float64 // RSI最高点
	TwoThirds  float64 // RSI平均次高点
	OneThirds  float64 // RSI平均次低点
	Low        float64 // RSI最低点

	High2NowLow float64 // RSI最高点 到 当前位置，中间出现的最低点

	NowToLowDays int // 当前RSI距离最低点的天数

	Rsi70Days int // RSI大于等于70的天数
	Rsi65Days int // RSI大于等于65的天数
	Rsi60Days int // RSI大于等于60的天数
	Rsi55Days int // RSI大于等于55的天数

	Message string    // 额外提示信息
	Time    time.Time // 时间
}

var Date = ""

// Rsi https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=sh000300&scale=30&ma=no&datalen=1023
func Rsi(code string, dayScale int) *RsiData {
	rsiArr := rsiArray(code, dayScale)
	message := ""
	if rsiArr == nil {
		return nil
	}
	latestHigh := 0.0
	high := 0.0
	avg := 0.0
	low := 100.0
	High2NowLow := 100.0
	// 忽略前20个元素
	rsiArr = rsiArr[20:]
	if len(rsiArr) < 10 {
		return nil
	}
	var highIndex int
	for index, rsi := range rsiArr {
		if rsi == 0 {
			continue
		}
		if rsi > high {
			high = rsi
			highIndex = index
		}

		if (index > 0 && index < len(rsiArr)-1) && rsi >= rsiArr[index-1] && rsi >= rsiArr[index+1] {
			latestHigh = rsi
		}
		if rsi < low {
			low = rsi
		}
	}

	if rsiArr[0] >= latestHigh {
		latestHigh = rsiArr[0]
	}

	if rsiArr[len(rsiArr)-1] >= latestHigh {
		latestHigh = rsiArr[len(rsiArr)-1]
	}

	avg = (high - low) / 3.0
	rsiData := &RsiData{
		Days:       len(rsiArr),
		Now:        rsiArr[len(rsiArr)-1],
		High:       high,
		LatestHigh: latestHigh,
		TwoThirds:  high - avg,
		OneThirds:  high - 2*avg,
		Low:        low,
		Message:    message,
		Time:       time.Now(),
	}
	for index, rsi := range rsiArr {
		if rsi >= rsiData.Low && rsi < rsiData.Now {
			rsiData.NowToLowDays++
		}
		if rsi >= 70 {
			rsiData.Rsi70Days++
		}
		if rsi >= 65 {
			rsiData.Rsi65Days++
		}
		if rsi >= 60 {
			rsiData.Rsi60Days++
		}
		if rsi >= 55 {
			rsiData.Rsi55Days++
		}
		if index > highIndex {
			if rsi < High2NowLow {
				High2NowLow = rsi
			}
		}
	}
	rsiData.Message = fmt.Sprintf("数据%d天, "+
		"70以上有%d天, "+
		"65以上有%d天, "+
		"60以上有%d天, "+
		"55以上有%d天, "+
		"当前与最低点之间有%d天", rsiData.Days, rsiData.Rsi70Days, rsiData.Rsi65Days, rsiData.Rsi60Days, rsiData.Rsi55Days, rsiData.NowToLowDays)
	return rsiData
}

// rsiArray 获取rsi数组数据
func rsiArray(code string, dayScale int) []float64 {
	defaultDay := 201
	if dayScale > defaultDay/2 {
		defaultDay = dayScale * 4
	}
	url := fmt.Sprintf("https://quotes.sina.cn/cn/api/jsonp_v2.php/=/CN_MarketDataService.getKLineData?symbol=%s&scale=120&ma=no&datalen=%d", code, defaultDay)
	response, err := http.Get(url)
	if err != nil {
		log.Println("http get error:", err)
		return nil
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("io read error:", err)
		return nil
	}
	rex := regexp.MustCompile(`\[([\s\S]*?)]`)
	titleMatches := rex.FindAllSubmatch(data, -1)
	if titleMatches == nil {
		log.Println("regexp error:", err)
		return nil
	}
	jsonStr := fmt.Sprintf("[%s]", string(titleMatches[0][1]))
	var index []constant.Index
	err = json.Unmarshal([]byte(jsonStr), &index)
	if err != nil {
		log.Println("json unmarshal error:", err)
		return nil
	}
	rsiData := make([]float64, 0)
	for i, data := range index {
		if i != len(index)-1 && !strings.Contains(data.Date, "15:00:00") {
			continue
		}
		float, err := strconv.ParseFloat(data.Close, 64)
		if err != nil {
			log.Println("strconv parseFloat error:", err)
			return nil
		}
		rsiData = append(rsiData, float)
		Date = data.Date
	}
	rsi := calRsi(rsiData, dayScale)
	return rsi
}

// calRsi 根据收盘价和间隔计算RSI
func calRsi(inReal []float64, inTimePeriod int) []float64 {
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
