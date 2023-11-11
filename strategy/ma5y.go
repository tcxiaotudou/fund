package strategy

import (
	"encoding/json"
	"fmt"
	"founds/constant"
	"founds/utils"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Ma5y 国证A指 5年均线
func Ma5y() string {
	url := "https://quotes.sina.cn/cn/api/json_v2.php/CN_MarketDataService.getKLineData?symbol=sz399317&scale=1200&ma=no&datalen=1950" // 请求的URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("创建请求失败:", err)
		return ""
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return ""
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	var indexs []constant.Index
	err = json.Unmarshal(responseBody, &indexs)
	if err != nil {
		log.Println("json unmarshal error:", err)
		return ""
	}
	lastClose := 0.0
	n := 52                                   // number of trading Days in five years
	sum := 0.0                                // sum of the last n closing prices
	ma5Result := make([]float64, len(indexs)) // result slice
	for i, data := range indexs {
		x, _ := strconv.ParseFloat(data.Close, 64)
		lastClose = x
		if i < n-1 {
			// not enough data to calculate Ma5y, append zero
			sum += x
			ma5Result[i] = 0.0
		} else if i == n-1 {
			// just enough data to calculate the first Ma5y, append sum / n
			sum += x
			ma5Result[i] = sum / float64(n)
		} else {
			// more than enough data to calculate Ma5y, append (sum + x - data[i-n]) / n
			tmp, _ := strconv.ParseFloat(indexs[i-n].Close, 64)
			sum += x - tmp
			ma5Result[i] = sum / float64(n)
		}
	}
	avg := ma5Result[len(ma5Result)-1]
	return fmt.Sprintf("%v", utils.Decimal((lastClose-avg)*100/avg)) + "%"
}
