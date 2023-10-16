package main

import (
	"fmt"
	"testing"
)

func TestRsi(t *testing.T) {
	rsiValue := GetRsi("sh000300")
	fmt.Println(rsiValue)
}
