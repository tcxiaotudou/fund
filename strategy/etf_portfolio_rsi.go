package strategy

import (
	"fmt"
)

var (
	EtfGroups = map[string]int64{
		"sz159696": 50, // 纳斯达克ETF
		"sz159307": 50, // 红利低波100ETF
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
			dailyWeightedPrices[i] += prices[i] * 100 * float64(weight)
		}
	}
	rsi := calculateRSI(dailyWeightedPrices, 14)
	return fmt.Sprintf("%.2f", rsi[len(rsi)-1])
}
