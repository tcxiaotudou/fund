package strategy

import (
	"encoding/json"
	"fmt"
	"founds/constant"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RsiData RSI数据
type RsiData struct {
	Days      int     // RSI数据的天数
	Now       float64 // 当前RSI
	High      float64 // RSI最高点
	TwoThirds float64 // RSI平均次高点
	OneThirds float64 // RSI平均次低点
	Low       float64 // RSI最低点

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
	high := 0.0
	avg := 0.0
	low := 100.0
	High2NowLow := 100.0
	// 忽略前20个元素
	rsiArr = rsiArr[dayScale:]
	if len(rsiArr) < dayScale {
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

		if rsi < low {
			low = rsi
		}
	}

	avg = (high - low) / 3.0
	rsiData := &RsiData{
		Days:      len(rsiArr),
		Now:       rsiArr[len(rsiArr)-1],
		High:      high,
		TwoThirds: high - avg,
		OneThirds: high - 2*avg,
		Low:       low,
		Message:   message,
		Time:      time.Now(),
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
	rsiData.High2NowLow = High2NowLow
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
	prices := getPrices(code, dayScale)
	rsi := calculateRSI(prices, dayScale)
	return rsi
}

func getPrices(code string, dayScale int) []float64 {
	defaultDay := 201
	if dayScale > defaultDay/3 {
		defaultDay = dayScale * 11
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
	prices := make([]float64, 0)
	for i, data := range index {
		if i != len(index)-1 && !strings.Contains(data.Date, "15:00:00") {
			continue
		}
		float, err := strconv.ParseFloat(data.Close, 64)
		if err != nil {
			log.Println("strconv parseFloat error:", err)
			return nil
		}
		prices = append(prices, float)
		Date = data.Date
	}
	return prices
}

// 计算 RSI
func calculateRSI(prices []float64, period int) []float64 {
	if prices == nil || len(prices) < period {
		return []float64{}
	}

	var gains, losses, rsi []float64

	// 计算每日涨跌幅
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, -change)
		}
	}

	// 填充前期空值
	for i := 0; i < period-1; i++ {
		rsi = append(rsi, 0)
	}

	// 计算第一个RSI值
	avgGain := sumSlice(gains[:period]) / float64(period)
	avgLoss := sumSlice(losses[:period]) / float64(period)
	rs := avgGain / (avgLoss + math.SmallestNonzeroFloat64)
	rsi = append(rsi, toFixed(100-(100/(1+rs)), 2))

	// 计算后续的RSI值
	for i := period; i < len(gains); i++ {
		avgGain = ((avgGain * float64(period-1)) + gains[i]) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + losses[i]) / float64(period)
		rs = avgGain / (avgLoss + math.SmallestNonzeroFloat64)
		rsi = append(rsi, toFixed(100-(100/(1+rs)), 2))
	}

	return rsi
}

// 计算切片的总和
func sumSlice(slice []float64) float64 {
	sum := 0.0
	for _, value := range slice {
		sum += value
	}
	return sum
}

// 保留指定小数位数
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}
