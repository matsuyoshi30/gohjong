package main

import (
	"fmt"

	"github.com/matsuyoshi30/gohjong"
)

func main() {
	h := "23456m789s345pEE"

	s, err := gohjong.ParseHand(h)
	if err != nil {
		fmt.Println(err)
	}

	h = "1235556667899"
	sl, err := gohjong.CheckWaiting(h)
	fmt.Println(sl)

	fmt.Println(s)
}
