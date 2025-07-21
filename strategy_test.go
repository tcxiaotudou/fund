package main

import (
	"fmt"
	"founds/strategy"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestPush(t *testing.T) {
	baseURL := "https://api2.pushdeer.com/message/push"
	pushkey := "PDU26608Tyn0d1mfuJPnqu5qLJrwUQQzzUOS797Zo"
	text := "每日行情"
	desp := `
#### 行情数据：
- 股债平衡建议: 3股7债 
- 14日RSI: 59.17 
- 90日RSI（57 点和 70 点卖）: 45.79 
- 5年均线: -4.11% 
#### 建议买入：
- 名称: 中证银行(sz399986)
- 当前: 33.62
- 区间：[68.86, 56.56, 44.26, 31.96]
- 备注：数据81天, 70以上有0天, 65以上有3天, 60以上有4天, 55以上有11天, 当前与最低点之间有4天
- 时间：2023年11月16日10:35:35
---
- 名称: 中证银行(sz399986)
- 当前: 33.62
- 区间：[68.86, 56.56, 44.26, 31.96]
- 备注：数据81天, 70以上有0天, 65以上有3天, 60以上有4天, 55以上有11天, 当前与最低点之间有4天
- 时间：2023年11月16日10:35:35
---
	`
	// Encode the query parameters
	params := url.Values{}
	params.Set("pushkey", pushkey)
	params.Set("text", text)
	params.Set("desp", desp)
	params.Set("type", "markdown")

	// Create the full request URL with encoded parameters
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	response, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer response.Body.Close()

	fmt.Println("Push notification sent successfully!")
}

func TestStockBalance(t *testing.T) {
	data := strategy.Stock300Balance()
	fmt.Println(data)
}

func TestM5year(t *testing.T) {
	fmt.Println(strategy.Ma5y())
}

func TestFundStrategy(t *testing.T) {
	fundStrategy := strategy.FundStrategy()
	for _, c := range fundStrategy {
		fmt.Println(c.Name)
	}

}

func TestQuantifyFundStrategy(t *testing.T) {
	fundStrategy := strategy.QuantifyFundStrategy()
	for _, c := range fundStrategy {
		fmt.Println(c.Name)
	}

}

func TestEtfPortfolioRsi(t *testing.T) {
	fmt.Println(strategy.EtfPortfolioRsi())
}

func TestFundPortfolioRsi(t *testing.T) {
	fmt.Println(strategy.FundPortfolioRsi())
}

func TestFundRsi(t *testing.T) {
	fmt.Println(strategy.FundRsi("378006", 14))
}

func TestMa60Strategy(t *testing.T) {
	log.Println("开始测试移动平均线策略...")

	// 测试移动平均线策略
	results := strategy.MaStrategy()

	fmt.Printf("移动平均线策略结果（共%d个ETF）:\n", len(results))
	fmt.Println("================================================")
	fmt.Printf("%-20s %-10s %-10s %-10s\n", "ETF名称", "60周均线", "当前日线", "买入信号")
	fmt.Println("------------------------------------------------")

	buyCount := 0
	for _, result := range results {
		buySignal := "否"
		if result.IsBuySignal {
			buySignal = "是"
			buyCount++
		}
		fmt.Printf("%-20s %-10.2f %-10.2f %-10s\n",
			result.ETFName, result.WeeklyMA60, result.CurrentDaily, buySignal)
	}

	fmt.Println("================================================")
	fmt.Printf("总计: %d个ETF，其中%d个有买入信号\n", len(results), buyCount)
}
