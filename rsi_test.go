package main

import (
	"fmt"
	"testing"
)

func TestRsi(t *testing.T) {
	date, rsiValue := GetRsi("sh000300")
	fmt.Println(date, rsiValue)
}
