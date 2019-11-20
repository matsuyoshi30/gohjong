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
	fmt.Println(s)

	h = "1235556667899m"
	sl, err := gohjong.CheckWaiting(h)
	for _, s := range sl {
		fmt.Println("Definite", s.DefiniteHand)
		fmt.Println("Waiting", s.WaitingHand)
	}

	h = "1112224577799m"
	sw, err := gohjong.ShowWaiting(h)
	for _, s := range sw {
		fmt.Println(h, "is waiting", s)
	}
}
