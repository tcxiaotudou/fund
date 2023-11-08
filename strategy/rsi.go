package strategy

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

type RsiData struct {
	RsiGroup []float64
	Message  string
}

var Date = ""

// RsiGroup https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=sh000300&scale=30&ma=no&datalen=1023
func RsiGroup(code string, dayScale int) ([]float64, string) {
	rsiDataArr := rsiDataArray(code, dayScale)
	message := ""
	if rsiDataArr == nil {
		return []float64{0, 0, 0, 0, 0}, message
	}
	high := 0.0
	avg := 0.0
	low := 100.0
	// 忽略前20个元素
	rsiDataArr = rsiDataArr[20:]
	for _, rsi := range rsiDataArr {
		if rsi == 0 {
			continue
		}
		if rsi > high {
			high = rsi
		}
		if rsi < low {
			low = rsi
		}
	}
	avg = (high - low) / 3.0
	group := []float64{rsiDataArr[len(rsiDataArr)-1], high, high - avg, high - 2*avg, low}

	var dayMaxLow, day70, day65, day60, day55 int
	for _, rsi := range rsiDataArr {
		if rsi >= group[4] && rsi < group[0] {
			dayMaxLow++
		}
		if rsi >= 70 {
			day70++
		}
		if rsi >= 65 {
			day65++
		}
		if rsi >= 60 {
			day60++
		}
		if rsi >= 55 {
			day55++
		}
	}
	message = fmt.Sprintf("数据%d天, "+
		"70以上有%d天, "+
		"65以上有%d天, "+
		"60以上有%d天, "+
		"55以上有%d天, "+
		"当前与最低点之间有%d天", len(rsiDataArr), day70, day65, day60, day55, dayMaxLow)
	return group, message
}

// 获取rsi数组数据
func rsiDataArray(code string, dayScale int) []float64 {
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
		return nil
	}
	jsonStr := fmt.Sprintf("[%s]", string(titleMatches[0][1]))
	var index []Index
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

// Index 定义一个结构体，用来存储指数的收盘价和日期
type Index struct {
	Close string `json:"close"` // 收盘价
	Date  string `json:"day"`   // 日期
}

// 是否为收盘时间
func isCloseTime(timeObj time.Time) bool {
	targetTime := time.Date(timeObj.Year(), timeObj.Month(), timeObj.Day(), 15, 0, 0, 0, timeObj.Location())
	return timeObj.Equal(targetTime)
}

// 根据收盘价和间隔计算RSI
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

// Decimal 四舍五入保留两位小数
func Decimal(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}
