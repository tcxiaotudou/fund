package strategy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Fear 韭圈恐贪指数
func Fear() (key, value string) {
	response, err := http.Post("https://api.jiucaishuo.com/v2/kjtl/getbasedata",
		"application/json", nil)
	if err != nil {
		return
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	var dataJson map[string]interface{}
	err = json.Unmarshal(responseBody, &dataJson)
	if err != nil {
		return
	}
	num := dataJson["data"].(map[string]interface{})["num"]
	statusStr := dataJson["data"].(map[string]interface{})["status_str"]
	return "恐贪指数", fmt.Sprintf("%s - %v", statusStr, num)
}
