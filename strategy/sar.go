package strategy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// 周线SAR
func sar(code string) {
	url := fmt.Sprintf("http://webquoteklinepic.eastmoney.com/GetPic.aspx?nid=1.%s&type=W&unitWidth=-6&ef=EXTENDED_SAR&formula=RSI&AT=1&imageType=KXL", code)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fileName := fmt.Sprintf("%s_sar_image.jpg", code)
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	log.Println("Image downloaded successfully.")
}
