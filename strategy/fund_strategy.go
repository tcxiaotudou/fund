package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"founds/constant"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// https://data.howbuy.com/cgi/fund/v800z/zjzhchartdthc.json?zhid=67888190128&range=5N

var existFund = map[string]string{
	"004475": "华泰柏瑞富利混合A",
	"006624": "中泰玉衡价值优选混合A",
	"210002": "金鹰红利价值混合A",
	"008271": "大成优势企业混合A",
	"004814": "中欧红利优享混合A",
	"090013": "大成竞争优势混合A",
	"001564": "东方红京东大数据混合A",
	"005576": "华泰柏瑞新金融地产混合A",
}

var exclude = map[string]string{
	"001247": "华泰柏瑞新利混合A",
	"004685": "金元顺安元启灵活配置混合",
}

var overseasFund = map[string]string{
	"539001": "建信纳斯达克100指数(QDII)A人民币",
	// "050025": "博时标普500ETF联接A",
}

func FundStrategy() []*constant.FundStrategy {
	log.Println("获取精选策略开始...")
	url := "https://api.jiucaishuo.com/v2/fundchoose/result2"
	method := "POST"

	payload := []byte(`{
        "condition_id": "2344760"
    }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.jiucaishuo.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("获取精选策略失败: %s\n", responseBody)
		log.Println(err)
		return nil
	}
	response := string(responseBody)

	log.Printf("获取精选策略结束: %s \n", response)

	if gjson.Get(response, "code").Int() != 0 {
		log.Println(response)
		return nil
	}

	fundList := gjson.Get(response, "data.position_table_data").Array()

	result := make([]*constant.FundStrategy, 0)

	cache := make(map[string]string)

	for _, data := range fundList {
		fundData := data.Map()
		item := &constant.FundStrategy{}
		fundName := fundData["name"].String()
		fundCode := fundData["code"].String()
		item.Name = fundName
		item.Code = fundCode
		fundInfo := fundData["list"].Array()
		item.PersonName = fundInfo[10].Map()["val"].String()
		item.PersonYear = fundInfo[3].Map()["val"].String()
		year5Sharpe, _ := strconv.Atoi(strings.Split(fundInfo[5].Map()["val"].String(), "/")[0])
		year5Calmar, _ := strconv.Atoi(strings.Split(fundInfo[6].Map()["val"].String(), "/")[0])
		item.Year5Sharpe = year5Sharpe
		item.Year5Calmar = year5Calmar
		item.Gm = fundInfo[1].Map()["val"].String()
		item.YearTodayIncome = fundInfo[9].Map()["val"].String()
		item = setFundRate(item)
		if item == nil {
			continue
		}
		result = append(result, item)
	}

	// 按近 5 年收益率排名
	sort.Slice(result, func(i, j int) bool {
		return result[i].Year5IncomeNumber > result[j].Year5IncomeNumber
	})

	list := make([]*constant.FundStrategy, 0)

	size := 10
	// 去重
	for _, fund := range result {

		// 排除
		if exclude[fund.Code] != "" {
			continue
		}

		_, ok := cache[fund.PersonName]
		if ok {
			continue
		}
		cache[fund.PersonName] = "1"
		_, ok = existFund[fund.Code]
		if !ok {
			fund.Name = "**" + fund.Name
		} else {
			fund.Name = existFund[fund.Code]
		}
		if len(list) < size {
			list = append(list, fund)
		}
	}

	for existCode, existName := range existFund {
		isDelete := true
		for _, strategy := range list {
			if existCode == strategy.Code {
				isDelete = false
			}
		}
		if isDelete {
			deleteFund := &constant.FundStrategy{Code: existCode, Name: "xx" + existName}
			list = append(list, deleteFund)
		}
	}

	return list
}

func setFundRate(strategy *constant.FundStrategy) *constant.FundStrategy {
	log.Printf("获取韭圈收益率开始: %s \n", strategy.Code)
	url := "https://api.jiucaishuo.com/fundetail/details/fundinfo"
	method := "POST"
	payload := []byte(fmt.Sprintf(`{
					"fund_code": "%s",
					"type": "h5",
					"version": "2.5.6"
					}`, strategy.Code))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.jiucaishuo.com")
	req.Header.Add("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	responseBody, _ := ioutil.ReadAll(res.Body)

	response := string(responseBody)
	log.Printf("获取韭圈收益率结束: %s \n", response)
	if err != nil || gjson.Get(response, "code").Int() != 0 {
		log.Println(response)
		return strategy
	}

	if strings.Contains(gjson.Get(response, "data.sx").String(), "暂停") {
		log.Println(response)
		return nil
	}

	// ----------
	url = "https://api.jiucaishuo.com/fundetail/details/earn-line"
	payload = []byte(fmt.Sprintf(`{
					"fund_code": "%s",
					"type": "h5",
					"date": 60,
					"version": "2.5.6"
					}`, strategy.Code))

	req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.jiucaishuo.com")
	req.Header.Add("Connection", "keep-alive")
	res, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	responseBody, _ = ioutil.ReadAll(res.Body)
	response = string(responseBody)
	year5IncomeNumber := gjson.Get(response, "data.year_income").Float()
	//
	//if year5IncomeNumber < 10 {
	//	return nil
	//}

	strategy.Year5IncomeNumber = year5IncomeNumber
	strategy.Year5Income = fmt.Sprintf("%.2f%%", year5IncomeNumber)

	return strategy
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		CateList []interface{} `json:"cate_list"`
		Category string        `json:"category"`
		List     [][]struct {
			Name  string `json:"name"`
			Color string `json:"color,omitempty"`
		} `json:"list"`
		TwoCategory      []interface{} `json:"two_category"`
		TwoCategoryIndex int           `json:"two_category_index"`
	} `json:"data"`
}

func GetFundData(fundCode string) ([]float64, error) {
	url := "https://apiv2.jiucaishuo.com/funddetail/changepercent/achieve"
	method := "POST"

	payload := map[string]interface{}{
		"fund_code": fundCode,
		"tags_id":   4,
		"limit":     200,
		"type":      "h5",
		"version":   "2.5.6",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("priority", "u=1, i")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "apiv2.jiucaishuo.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	if len(response.Data.List) < 2 {
		return nil, fmt.Errorf("the list does not contain enough elements")
	}

	var floatValues []float64
	for _, item := range response.Data.List[1] {
		value, err := strconv.ParseFloat(item.Name, 64)
		if err == nil {
			floatValues = append(floatValues, value)
		}
	}

	// floatValues反转
	for i, j := 0, len(floatValues)-1; i < j; i, j = i+1, j-1 {
		floatValues[i], floatValues[j] = floatValues[j], floatValues[i]
	}

	gsjz, err := FetchGSJZ(fundCode)
	if err == nil {
		floatValues = append(floatValues, gsjz)
	}

	return floatValues, nil
}

func FundPortfolioRsi() string {
	// Initialize a slice to store the weighted prices for each day
	var dailyWeightedPrices []float64

	// Iterate over each ETF in the group
	for code, _ := range existFund {
		prices, _ := GetFundData(code)
		if dailyWeightedPrices == nil {
			dailyWeightedPrices = make([]float64, len(prices))
		}
		// Accumulate the weighted prices for each day
		for i := 0; i < len(prices); i++ {
			dailyWeightedPrices[i] += prices[i] * 6.25
		}
	}
	// 海外
	for code, _ := range overseasFund {
		prices, _ := GetFundData(code)
		if dailyWeightedPrices == nil {
			dailyWeightedPrices = make([]float64, len(prices))
		}
		// Accumulate the weighted prices for each day
		for i := 0; i < len(prices); i++ {
			dailyWeightedPrices[i] += prices[i] * float64(50)
		}
	}
	rsi := calculateRSI(dailyWeightedPrices, 14)
	return fmt.Sprintf("%.2f", rsi[len(rsi)-1])
}

func FundRsi(code string, period int) float64 {
	prices, _ := GetFundData(code)
	for i := 0; i < len(prices); i++ {
		prices[i] = prices[i] * 100
	}
	rsi := calculateRSI(prices, period)
	return rsi[len(rsi)-1]
}

// FetchGSJZ retrieves the "gsjz" field from the API response as a float64.
func FetchGSJZ(fundCode string) (float64, error) {
	url := "https://api.jiucaishuo.com/fundetail/ttm/info"
	method := "POST"

	// Construct the request payload
	payload := []byte(fmt.Sprintf(`{
        "fund_code": "%s",
        "type": "h5",
        "version": "2.5.6"
    }`, fundCode))

	// Create the HTTP request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.jiucaishuo.com")
	req.Header.Add("Connection", "keep-alive")

	// Execute the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse JSON response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Extract data.gsjz and convert to float64
	if data, ok := response["data"].(map[string]interface{}); ok {
		if gsjz, ok := data["gsjz"].(string); ok {
			value, err := strconv.ParseFloat(gsjz, 64)
			if err != nil {
				return 0, fmt.Errorf("error converting gsjz to float64: %v", err)
			}
			return value, nil
		}
		return 0, fmt.Errorf("gsjz field not found or not a string")
	}

	return 0, fmt.Errorf("data field not found or not an object")
}

// 量化基金策略

var existQuantifyFund = map[string]string{
	"015880": "中欧小盘成长混合A",
	"005437": "易方达易百智能量化策略A",
	"005616": "东方量化成长灵活配置混合A",
	"006267": "诺德量化核心A",
	"001990": "中欧数据挖掘多因子混合A",
	"011868": "中信建投远见回报混合A",
	"000006": "西部利得量化成长混合A",
	"014201": "天弘中证1000指数增强A",
	"005457": "景顺长城量化小盘股票A",
	"015453": "中欧中证500指数增强A",
}

func QuantifyFundStrategy() []*constant.FundStrategy {
	log.Println("获取量化策略开始...")
	url := "https://api.jiucaishuo.com/v2/fundchoose/result2"
	method := "POST"

	payload := []byte(`{
        "condition_id": "2368605"
    }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.jiucaishuo.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("获取量化策略失败: %s\n", responseBody)
		log.Println(err)
		return nil
	}
	response := string(responseBody)

	log.Printf("获取量化策略结束: %s \n", response)

	if gjson.Get(response, "code").Int() != 0 {
		log.Println(response)
		return nil
	}

	fundList := gjson.Get(response, "data.position_table_data").Array()

	result := make([]*constant.FundStrategy, 0)

	cache := make(map[string]string)

	for _, data := range fundList {
		fundData := data.Map()
		item := &constant.FundStrategy{}
		fundName := fundData["name"].String()
		fundCode := fundData["code"].String()
		item.Name = fundName
		item.Code = fundCode
		fundInfo := fundData["list"].Array()
		item.PersonName = fundInfo[8].Map()["val"].String()
		item.PersonYear = fundInfo[4].Map()["val"].String()
		item.Gm = fundInfo[2].Map()["val"].String()
		item.YearTodayIncome = strings.Split(fundInfo[7].Map()["val"].String(), "%")[0]

		yearTodayIncomeNumber, _ := strconv.ParseFloat(item.YearTodayIncome, 64)
		item.YearTodayIncomeNumber = yearTodayIncomeNumber

		result = append(result, item)
	}

	// 按今年收益率排名
	sort.Slice(result, func(i, j int) bool {
		return result[i].YearTodayIncomeNumber > result[j].YearTodayIncomeNumber
	})

	list := make([]*constant.FundStrategy, 0)

	size := 12
	// 去重
	for _, fund := range result {

		// 排除
		if exclude[fund.Code] != "" {
			continue
		}

		_, ok := cache[fund.PersonName]
		if ok {
			continue
		}
		cache[fund.PersonName] = "1"
		_, ok = existQuantifyFund[fund.Code]
		if !ok {
			fund.Name = "**" + fund.Name
		} else {
			fund.Name = existQuantifyFund[fund.Code]
		}
		if len(list) < size {
			list = append(list, fund)
		}
	}

	for existCode, existName := range existQuantifyFund {
		isDelete := true
		for _, strategy := range list {
			if existCode == strategy.Code {
				isDelete = false
			}
		}
		if isDelete {
			deleteFund := &constant.FundStrategy{Code: existCode, Name: "xx" + existName}
			list = append(list, deleteFund)
		}
	}

	return list
}

func QuantifyFundPortfolioRsi() string {
	// Initialize a slice to store the weighted prices for each day
	var dailyWeightedPrices []float64

	// Iterate over each ETF in the group
	for code, _ := range existQuantifyFund {
		prices, _ := GetFundData(code)
		if dailyWeightedPrices == nil {
			dailyWeightedPrices = make([]float64, len(prices))
		}
		// Accumulate the weighted prices for each day
		for i := 0; i < len(prices); i++ {
			dailyWeightedPrices[i] += prices[i] * 100
		}
	}
	rsi := calculateRSI(dailyWeightedPrices, 14)
	return fmt.Sprintf("%.2f", rsi[len(rsi)-1])
}
