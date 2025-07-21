package strategy

import (
	"encoding/json"
	"fmt"
	"founds/constant"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// MaStrategyData 移动平均线策略数据
type MaStrategyData struct {
	ETFName      string    // ETF名称
	ETFCode      string    // ETF代码
	WeeklyMA60   float64   // 60周均线
	CurrentDaily float64   // 当前日K线收盘价
	DataTime     time.Time // 数据时间
	IsBuySignal  bool      // 是否为买入信号
}

// MaStrategy 移动平均线策略：日线在60周均线上方，日线在60日均线下方时买入
func MaStrategy() []*MaStrategyData {
	var results []*MaStrategyData

	for name, code := range constant.EtfGroups {
		time.Sleep(3 * time.Second) // 避免请求过于频繁

		// 获取60周均线
		weeklyMA60 := calculateMA(code, 1200, 60) // 1200为周K线
		if weeklyMA60 == 0 {
			continue
		}

		time.Sleep(3 * time.Second)

		// 获取60日均线
		dailyMA60 := calculateMA(code, 240, 60) // 240为日K线
		if dailyMA60 == 0 {
			continue
		}

		time.Sleep(3 * time.Second)
		// 获取当前日K线收盘价
		currentDaily := getCurrentPrice(code, 240)
		if currentDaily == 0 {
			continue
		}

		// 判断买入信号：当前周K线在60周均线以上 且 当前日K线在60日均线以下
		isBuySignal := currentDaily > weeklyMA60 && currentDaily < dailyMA60

		data := &MaStrategyData{
			ETFName:      fmt.Sprintf("%s(%s)", name, code),
			ETFCode:      code,
			WeeklyMA60:   weeklyMA60,
			CurrentDaily: currentDaily,
			DataTime:     time.Now(),
			IsBuySignal:  isBuySignal,
		}

		results = append(results, data)

		log.Printf("MA策略 - %s: 日线%.2f/60日均线%.2f, 买入信号:%v",
			name, currentDaily, dailyMA60, isBuySignal)
	}

	return results
}

// calculateMA 计算移动平均线
func calculateMA(code string, scale int, period int) float64 {
	prices := getPricesForMA(code, scale, period+10) // 多取10个数据点确保足够
	if len(prices) < period {
		return 0
	}

	// 计算最近period个周期的平均值
	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}

// getCurrentPrice 获取当前价格（最新收盘价）
func getCurrentPrice(code string, scale int) float64 {
	prices := getPricesForMA(code, scale, 1)
	if len(prices) == 0 {
		return 0
	}
	return prices[len(prices)-1]
}

// getPricesForMA 获取用于移动平均线计算的价格数据
func getPricesForMA(code string, scale int, datalen int) []float64 {
	url := fmt.Sprintf("https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=%s&scale=%d&ma=no&datalen=%d",
		code, scale, datalen)

	response, err := http.Get(url)
	if err != nil {
		log.Printf("获取%s数据失败: %v", code, err)
		return nil
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("读取%s响应失败: %v", code, err)
		return nil
	}

	var index []constant.Index
	err = json.Unmarshal(data, &index)
	if err != nil {
		log.Printf("解析%s数据失败: %v", code, err)
		return nil
	}

	if len(index) == 0 {
		log.Printf("获取%s数据为空", code)
		return nil
	}

	prices := make([]float64, 0)
	for _, data := range index {
		price, err := strconv.ParseFloat(data.Close, 64)
		if err != nil {
			log.Printf("解析价格失败 %s: %v", data.Close, err)
			continue
		}
		prices = append(prices, price)
	}

	return prices
}
