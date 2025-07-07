package strategy

import (
	"fmt"
	"time"
)

var (
	EtfGroups = map[string]int64{
		"sz161119": 15, // 中债ETF
		"sh512890": 35, // 红利低波ETF
		"sh513300": 35, // 纳斯达克ETF
		"sz159985": 15, // 豆柏ETF
	}
)

func EtfPortfolioRsi() string {
	// Initialize a slice to store the weighted prices for each day
	var dailyWeightedPrices []float64

	// Iterate over each ETF in the group
	for code, weight := range EtfGroups {
		prices := getPrices(code, 14)
		if dailyWeightedPrices == nil {
			dailyWeightedPrices = make([]float64, len(prices))
		}
		// Accumulate the weighted prices for each day
		for i := 0; i < len(prices); i++ {
			dailyWeightedPrices[i] += prices[i] * float64(weight)
		}
		time.Sleep(4 * time.Second)
	}
	rsi := calculateRSI(dailyWeightedPrices, 14)
	return fmt.Sprintf("%.2f", rsi[len(rsi)-1])
}
