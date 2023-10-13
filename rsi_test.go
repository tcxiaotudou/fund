package main

import (
	"fmt"
	"testing"
)

func TestRsi(t *testing.T) {
	rsiValue := GetRsi("sh000688")
	fmt.Println(rsiValue)
}
