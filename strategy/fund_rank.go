package strategy

import (
	"bytes"
	"fmt"
	"founds/constant"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

/**
1、过去1个月，过去3个月，过去6个月，过去1年 收益率前20
2、排除定开，排除规模大于50亿
3、进行PK，按得分取前5
*/

func FundRank() []*constant.Fund {
	fundMap := make(map[string]*constant.Fund)
	list := make([]*constant.Fund, 0)
	fundRank1 := rank(1)
	fundRank3 := rank(3)
	fundRank6 := rank(6)
	list = append(list, fundRank1...)
	list = append(list, fundRank3...)
	list = append(list, fundRank6...)
	//list = append(list, fundRank12...)
	for _, fund := range list {
		if fundMap[fund.Code] == nil {
			fundMap[fund.Code] = fund
		}
	}
	list = make([]*constant.Fund, 0)
	for _, fund := range fundMap {
		list = append(list, fund)
	}
	pages := make([]*constant.Fund, 0)
	result := make([]*constant.Fund, 0)
	for _, fund := range list {
		pages = append(pages, fund)
		if len(pages) == 5 {
			pages = pk(pages)
			pages = productInfo(pages)
			pages = sharpeAndRetracement(pages)
			pages = filter(pages)
			pages = performence(pages)
			result = append(result, pages...)
			pages = make([]*constant.Fund, 0)
		}
	}
	result = topN(result, 5)
	return result
}

func rank(monthScale int) []*constant.Fund {
	result := make([]*constant.Fund, 0)
	url := "https://api.jiucaishuo.com/v2/fund-lists/fundrank"
	payload := []byte(fmt.Sprintf(`{
    "page": 1,
    "cate_id": "1",
    "time": %d,
    "type_id": "1",
    "rank_type": "details",
    "orderby": "desc",
    "data_source": "xichou",
    "type": "h5"
	}`, monthScale))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return result
	}
	req.Header.Set("authority", "api.jiucaishuo.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "api.jiucaishuo.com")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return result
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result
	}
	response := string(responseBody)

	if err != nil || gjson.Get(response, "code").Int() != 0 {
		return result
	}

	fundInfoList := gjson.Get(response, "data.list").Array()
	for _, fundInfoItem := range fundInfoList {
		fundInfo := fundInfoItem.Map()
		name := fundInfo["name"].String()
		code := fundInfo["code"].String()
		fund := &constant.Fund{
			Name:   name,
			Code:   code,
			Reason: fmt.Sprintf("近%d个月收益", monthScale),
		}
		result = append(result, fund)
	}
	return result
}

func filter(funds []*constant.Fund) []*constant.Fund {
	result := make([]*constant.Fund, 0)
	for _, fund := range funds {
		if strings.Contains(fund.Name, "定开") {
			continue
		}
		if strings.Contains(fund.Name, "一年") {
			continue
		}
		if strings.Contains(fund.Name, "三年") {
			continue
		}
		if strings.Contains(fund.Name, "持有期") {
			continue
		}
		if fund.Scale > 50 {
			continue
		}
		retracement, err := strconv.ParseFloat(fund.Retracement, 64)
		if err == nil && retracement < -15 {
			continue
		}

		sharpe, err := strconv.ParseFloat(fund.Sharpe, 64)
		if err == nil && sharpe < 0 {
			continue
		}

		result = append(result, fund)
	}
	return result
}

func pk(funds []*constant.Fund) []*constant.Fund {
	var codes []string
	for _, item := range funds {
		codes = append(codes, item.Code)
	}
	requestCodes := strings.Join(codes, ",")

	url := "https://api.jiucaishuo.com/v2/fundpk/basic-info"
	payload := []byte(fmt.Sprintf(`{
    "fund_code": "%s",
    "type": "h5",
    "version": "2.4.8"
	}`, requestCodes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return funds
	}

	req.Header.Set("authority", "api.jiucaishuo.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "api.jiucaishuo.com")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return funds
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return funds
	}
	response := string(responseBody)

	if err != nil || gjson.Get(response, "code").Int() != 0 {
		return funds
	}

	pkshareData := gjson.Get(response, "data.pkshare_data").Array()
	for _, data := range pkshareData {
		pkInfo := data.Map()
		fundCode := pkInfo["fund_code"].String()
		score := pkInfo["sorce"].Int()
		for _, fund := range funds {
			if fund.Code == fundCode {
				fund.Score = score
			}
		}
	}
	return funds
}

func productInfo(funds []*constant.Fund) []*constant.Fund {
	var codes []string
	for _, item := range funds {
		codes = append(codes, item.Code)
	}
	requestCodes := strings.Join(codes, ",")

	url := "https://api.jiucaishuo.com/v2/fundpk/fundproduct"
	payload := []byte(fmt.Sprintf(`{
    "fund_code": "%s",
    "type": "h5",
    "version": "2.4.8"
	}`, requestCodes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return funds
	}

	req.Header.Set("authority", "api.jiucaishuo.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "api.jiucaishuo.com")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return funds
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return funds
	}
	response := string(responseBody)

	if err != nil || gjson.Get(response, "code").Int() != 0 {
		return funds
	}

	// 创建正则表达式模式
	pattern := `[0-9.]+`
	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	pkshareData := gjson.Get(response, "data.list").Array()
	for _, data := range pkshareData {
		pkInfo := data.Map()
		fundCreateTime := pkInfo["fund_create_time"].String()
		style := pkInfo["fg"].String()
		fundCode := pkInfo["fund_code"].String()
		scale := pkInfo["jjgm"].String()
		for _, fund := range funds {
			if fund.Code == fundCode {
				fund.CreateYears = fundCreateTime
				fund.Style = style
				match := re.FindString(scale)
				num, err := strconv.ParseFloat(match, 64)
				if err != nil {
					fmt.Println("Error parsing float:", err)
					continue
				}
				fund.Scale = num
			}
		}
	}
	return funds
}

/*
*

	设置近一年的夏普比率和最大回撤
*/
func sharpeAndRetracement(funds []*constant.Fund) []*constant.Fund {
	var codes []string
	for _, item := range funds {
		codes = append(codes, item.Code)
	}
	requestCodes := strings.Join(codes, ",")

	url := "https://api.jiucaishuo.com/v2/fundpk/fundrank"
	payload := []byte(fmt.Sprintf(`{
    "fund_code": "%s",
    "type": "h5",
 	"fund_day": 6,
    "version": "2.4.8"
	}`, requestCodes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return funds
	}

	req.Header.Set("authority", "api.jiucaishuo.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "api.jiucaishuo.com")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return funds
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return funds
	}
	response := string(responseBody)

	if err != nil || gjson.Get(response, "code").Int() != 0 {
		return funds
	}

	list := gjson.Get(response, "data.list").Array()
	for _, data := range list {
		pkInfo := data.Map()
		fundCode := pkInfo["fund_code"].String()

		for _, fund := range funds {
			if fund.Code == fundCode {
				sharpeInfo := pkInfo["list"].Array()
				for _, sharpeItem := range sharpeInfo {
					sharpeMap := sharpeItem.Map()
					name := sharpeMap["name"]
					if name.String() == "回撤率" {
						fund.Retracement = sharpeMap["num"].String()
					}
					if name.String() == "夏普率" {
						fund.Sharpe = sharpeMap["num"].String()
					}
				}
			}
		}
	}
	return funds
}

/*
*
业绩
*/
func performence(funds []*constant.Fund) []*constant.Fund {
	var codes []string
	for _, item := range funds {
		codes = append(codes, item.Code)
	}
	requestCodes := strings.Join(codes, ",")

	url := "https://api.jiucaishuo.com/v2/fundpk/performance"
	payload := []byte(fmt.Sprintf(`{
    "fund_code": "%s",
    "type": "h5",
 	"sign": 1,
    "version": "2.4.8"
	}`, requestCodes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return funds
	}

	req.Header.Set("authority", "api.jiucaishuo.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "api.jiucaishuo.com")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return funds
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return funds
	}
	response := string(responseBody)

	if err != nil || gjson.Get(response, "code").Int() != 0 {
		return funds
	}

	list := gjson.Get(response, "data.result_list.series").Array()
	for _, data := range list {
		pkInfo := data.Map()
		fundCode := pkInfo["fund_code"].String()

		for _, fund := range funds {
			if fund.Code == fundCode {
				sharpeInfo := pkInfo["list"].Array()
				for _, sharpeItem := range sharpeInfo {
					sharpeMap := sharpeItem.Map()
					name := sharpeMap["qx"]
					if name.String() == "近3月" {
						fund.Yield3 = sharpeMap["num"].String()
					}
					if name.String() == "近6月" {
						fund.Yield6 = sharpeMap["num"].String()
					}
					if name.String() == "近1年" {
						fund.Yield12 = sharpeMap["num"].String()
					}
				}
			}
		}
	}
	return funds
}

func topN(funds []*constant.Fund, topN int) []*constant.Fund {
	// 自定义排序规则
	sort.Slice(funds, func(i, j int) bool {
		// 先按 Score 从大到小排序
		if funds[i].Score > funds[j].Score {
			return true
		} else if funds[i].Score < funds[j].Score {
			return false
		}
		// 当 Score 相等时，按 Sharpe 从大到小排序
		return funds[i].Sharpe > funds[j].Sharpe
	})
	if len(funds) > topN {
		return funds[0:topN]
	}
	return funds
}
