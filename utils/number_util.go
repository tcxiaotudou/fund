package utils

import (
	"fmt"
	"strconv"
)

// Decimal 四舍五入保留两位小数
func Decimal(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}
