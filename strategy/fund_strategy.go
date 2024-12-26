package strategy

import (
	"bytes"
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
	"161611": "融通内需驱动混合A",
	"519702": "交银趋势混合A",
	"260112": "景顺长城能源基建混合A",
	"006624": "中泰玉衡价值优选混合A",
	"121010": "国投瑞银瑞源灵活配置混合A",
	"004475": "华泰柏瑞富利混合A",
	"090007": "大成策略回报混合A",
	"004814": "中欧红利优享混合A",
}

func FundStrategy() []*constant.FundStrategy {
	url := "https://api.jiucaishuo.com/v2/fundchoose/result2"
	method := "POST"

	payload := []byte(`{
        "condition_id": "2199957"
    }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
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
		log.Println(err)
		return nil
	}
	response := string(responseBody)

	if err != nil || gjson.Get(response, "code").Int() != 0 {
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
		item.Year5Sharpe = year5Sharpe
		item.Gm = fundInfo[1].Map()["val"].String()
		item.YearTodayIncome = fundInfo[9].Map()["val"].String()
		item = setRate(item)
		if item == nil {
			continue
		}
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Year5Sharpe < result[j].Year5Sharpe
	})

	list := make([]*constant.FundStrategy, 0)

	// 去重
	for _, fund := range result {
		_, ok := cache[fund.PersonName]
		if ok {
			continue
		}
		cache[fund.PersonName] = "1"
		_, ok = existFund[fund.Code]
		if !ok {
			fund.Name = "**" + fund.Name
		}
		list = append(list, fund)
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

func setHc(strategy *constant.FundStrategy) *constant.FundStrategy {
	url := "https://api.jiucaishuo.com/fundetail/details/showhc"
	//currentTime := time.Now()
	//milliseconds := currentTime.UnixNano() / int64(time.Millisecond)

	request := fmt.Sprintf(`{
        "date": 60,
    "ben": "",
    "search_code": "",
    "is_jl": "",
    "s_time": "",
    "e_time": "",
    "fund_code": "%s",
    "type": "h5",
    "version": "2.5.6",
    "authtoken": "FRfWjc2EFWmDZ55cIm4xat7RaBA3Tl1p",
    "ss": "",
    "act_time": %d,
    "tirgkjfs": "b9",
    "abiokytke": "c0",
    "u54rg5d": "2c",
    "kf54ge7": "d",
    "tiklsktr4": "9",
    "lksytkjh": "aa83",
    "sbnoywr": "4e",
    "bgd7h8tyu54": "e8",
    "y654b5fs3tr": "9",
    "bioduytlw": "5",
    "bd4uy742": "5",
    "h67456y": "9aa",
    "bvytikwqjk": "e8",
    "ngd4uy551": "aa",
    "bgiuytkw": "44",
    "nd354uy4752": "a",
    "ghtoiutkmlg": "9de",
    "bd24y6421f": "e9",
    "tbvdiuytk": "9",
    "ibvytiqjek": "5f",
    "jnhf8u5231": "44",
    "fjlkatj": "2c7",
    "hy5641d321t": "95",
    "iogojti": "9",
    "ngd4yut78": "de",
    "nkjhrew": "5",
    "yt447e13f": "7",
    "n3bf4uj7y7": "a",
    "nbf4uj7y432": "c0",
    "yi854tew": "1a",
    "h13ey474": "1ad",
    "quikgdky": "3c"
}`, strategy.Code, 1723991849512)

	payload := []byte(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return strategy
	}

	req.Header.Set("priority", "u=1, i")
	req.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "api.jiucaishuo.com")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return strategy
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	response := string(responseBody)
	if err != nil || gjson.Get(response, "code").Int() != 0 {
		return strategy
	}

	maxHc := gjson.Get(response, "data.max_hc").String()
	curHc := gjson.Get(response, "data.cur_hc").String()
	strategy.HcCurYear5 = curHc
	strategy.HcMaxYear5 = maxHc
	return strategy
}

func setRate(strategy *constant.FundStrategy) *constant.FundStrategy {
	url := "https://danjuanfunds.com/djapi/fund/" + strategy.Code

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return strategy
	}

	req.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "danjuanfunds.com")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return strategy
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	response := string(responseBody)
	if err != nil || gjson.Get(response, "result_code").Int() != 0 {
		log.Println(response)
		return strategy
	}

	if gjson.Get(response, "data.declare_status").String() == "0" {
		log.Println(response)
		return nil
	}

	baseDataArr := gjson.Get(response, "data.fir_header_base_data").Array()

	for _, baseData := range baseDataArr {
		if baseData.Map()["data_name"].String() == "年化收益（近5年）" {
			strategy.Year5Income = baseData.Map()["data_value_str"].String()
			strategy.Year5IncomeNumber = baseData.Map()["data_value_number"].Num
		}
	}

	if strategy.Year5IncomeNumber < 10 {
		return nil
	}

	return strategy
}
