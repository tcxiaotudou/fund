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

	fmt.Println(index)

	//定义一个切片，存储某个指数在2021年4月5日到2021年4月19日的收盘价数据（假设）
	//index := []Index{
	//	{Close: "5256.06", Date: "1686499200000"},
	//	{Close: "5281.61", Date: "1686585600000"},
	//	{Close: "5290.13", Date: "1686672000000"},
	//	{Close: "5348.01", Date: "1686758400000"},
	//	{Close: "5397.22", Date: "1686844800000"},
	//	{Close: "5383.77", Date: "1687104000000"},
	//	{Close: "5376.56", Date: "1687190400000"},
	//	{Close: "5282.90", Date: "1687276800000"},
	//	{Close: "5194.72", Date: "1687708800000"},
	//	{Close: "5258.24", Date: "1687795200000"},
	//	{Close: "5247.10", Date: "1687881600000"},
	//	{Close: "5252.70", Date: "1687968000000"},
	//	{Close: "5302.15", Date: "1688054400000"},
	//	{Close: "5346.54", Date: "1688313600000"},
	//	{Close: "5362.32", Date: "1688400000000"},
	//	{Close: "5323.45", Date: "1688486400000"},
	//	{Close: "5301.82", Date: "1688572800000"},
	//	{Close: "5274.44", Date: "1688659200000"},
	//	{Close: "5290.28", Date: "1688918400000"},
	//	{Close: "5327.99", Date: "1689004800000"},
	//	{Close: "5277.57", Date: "1689091200000"},
	//	{Close: "5349.18", Date: "1689177600000"},
	//	{Close: "5348.95", Date: "1689264000000"},
	//	{Close: "5322.95", Date: "1689523200000"},
	//	{Close: "5311.78", Date: "1689609600000"},
	//	{Close: "5304.37", Date: "1689696000000"},
	//	{Close: "5253.05", Date: "1689782400000"},
	//	{Close: "5248.87", Date: "1689868800000"},
	//	{Close: "5235.81", Date: "1690128000000"},
	//	{Close: "5351.52", Date: "1690214400000"},
	//	{Close: "5331.92", Date: "1690300800000"},
	//	{Close: "5309.16", Date: "1690387200000"},
	//	{Close: "5394.47", Date: "1690473600000"},
	//	{Close: "5433.57", Date: "1690732800000"},
	//	{Close: "5420.65", Date: "1690819200000"},
	//	{Close: "5396.96", Date: "1690905600000"},
	//	{Close: "5421.90", Date: "1690992000000"},
	//	{Close: "5442.56", Date: "1691078400000"},
	//	{Close: "5406.84", Date: "1691337600000"},
	//	{Close: "5388.45", Date: "1691424000000"},
	//	{Close: "5357.76", Date: "1691510400000"},
	//	{Close: "5368.92", Date: "1691596800000"},
	//	{Close: "5260.44", Date: "1691683200000"},
	//	{Close: "5249.90", Date: "1691942400000"},
	//	{Close: "5223.27", Date: "1692028800000"},
	//	{Close: "5176.86", Date: "1692115200000"},
	//	{Close: "5213.11", Date: "1692201600000"},
	//	{Close: "5134.88", Date: "1692288000000"},
	//	{Close: "5076.33", Date: "1692547200000"},
	//	{Close: "5109.43", Date: "1692633600000"},
	//	{Close: "5018.82", Date: "1692720000000"},
	//	{Close: "5041.88", Date: "1692806400000"},
	//	{Close: "4986.90", Date: "1692892800000"},
	//	{Close: "5038.82", Date: "1693152000000"},
	//	{Close: "5141.84", Date: "1693238400000"},
	//	{Close: "5154.89", Date: "1693324800000"},
	//	{Close: "5124.93", Date: "1693411200000"},
	//	{Close: "5144.28", Date: "1693497600000"},
	//	{Close: "5220.25", Date: "1693756800000"},
	//	{Close: "5186.50", Date: "1693843200000"},
	//	{Close: "5191.63", Date: "1693929600000"},
	//	{Close: "5110.92", Date: "1694016000000"},
	//	{Close: "5105.35", Date: "1694102400000"},
	//}

	rsiData := make([]float64, 0)
	for _, data := range index {
		float, err := strconv.ParseFloat(data.Close, 64)
		if err != nil {
			return 0
		}
		rsiData = append(rsiData, float)
	}

	result := caRsi(rsiData, 14)

	fmt.Println(result)

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
			// outReal[outIdx] =

			// outReal[outIdx] = math.Trunc((100.0*(prevGain/tempValue1)-0.11)*1e2+0.5) * 1e-2
			outReal[outIdx] = float64(int(100.0*(prevGain/tempValue1) - 0.11))
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}
