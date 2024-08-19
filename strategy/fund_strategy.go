package strategy

import (
	"bytes"
	"fmt"
	"founds/constant"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
)

func FundStrategy() []*constant.FundStrategy {
	url := "https://api.jiucaishuo.com/v2/fundchoose/result2"
	method := "POST"

	payload := []byte(`{
        "condition_id": "2196916"
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
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	response := string(responseBody)

	if err != nil || gjson.Get(response, "code").Int() != 0 {
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
		item.PersonName = fundInfo[11].Map()["val"].String()
		item.PersonYear = fundInfo[4].Map()["val"].String()
		item.Gm = fundInfo[1].Map()["val"].String()
		item.YearTodayIncome = fundInfo[10].Map()["val"].String()
		result = append(result, item)
	}

	list := make([]*constant.FundStrategy, 0)

	// 回撤
	for _, fund := range result {
		// setHc(fund)
		item := setRate(fund)
		if item != nil {
			_, ok := cache[item.PersonName]
			if ok {
				continue
			}
			cache[item.PersonName] = "1"
			list = append(list, fund)
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
		return strategy
	}

	if gjson.Get(response, "data.declare_status").String() == "0" {
		return nil
	}

	baseDataArr := gjson.Get(response, "data.fir_header_base_data").Array()

	for _, baseData := range baseDataArr {
		if baseData.Map()["data_name"].String() == "年化收益（近5年）" {
			strategy.Year5Income = baseData.Map()["data_value_str"].String()
		}
	}

	return strategy
}
