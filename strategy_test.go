package main

import (
	"fmt"
	"founds/strategy"
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

func TestFundRank(t *testing.T) {
	list := strategy.FundRank()
	fmt.Println(list)
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
