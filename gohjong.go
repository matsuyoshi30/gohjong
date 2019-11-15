package gohjong

import (
	"errors"
	_ "fmt"
	"strconv"
	"strings"
)

type Tile struct {
	name string
}

// ParseHand parse hand to tile slice
func ParseHand(hand string) ([]Tile, error) {
	res := make([]Tile, 0)

	pool := make([]uint8, 0)
	for i := 0; i < len(hand); i++ {
		handstr := string(hand[i])
		_, err := strconv.Atoi(handstr)
		if err != nil { // not num
			if len(pool) != 0 {
				if handstr != "m" && handstr != "p" && handstr != "s" { // len(pool) > 0 => must be m, p, s
					return nil, errors.New("unknown hand")
				}

				for _, p := range pool {
					res = append(res, Tile{name: string(p) + handstr})
				}

				pool = make([]uint8, 0)
			} else { // len(pool) == 0 => must be E, S, W, N, D, H, T
				if handstr != "E" && handstr != "S" && handstr != "W" && handstr != "N" &&
					handstr != "D" && handstr != "H" && handstr != "T" {
					return nil, errors.New("unknown hand")
				}
				res = append(res, Tile{name: handstr})
			}
		} else { // num
			pool = append(pool, hand[i])
		}
	}

	if len(pool) != 0 { // must be empty
		return nil, errors.New("could not parse")
	}

	return res, nil
}

// CheckWaiting check waiting tiles
// returns mentsu, machi, and error
func CheckWaiting(hand string) ([]string, error) {
	output := make([]string, 0)
	// TODO
	// _, err := ParseHand(hand)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// check waiting
	// check toistu kotsu kotsu kotsu kotsu
	for i := 0; i <= 9; i++ {
		resthand := ""

		resthand, toitsu := checkToitsu(hand, i)
		resthand, kotsu1 := checkKotsu(resthand)
		resthand, kotsu2 := checkKotsu(resthand)
		resthand, kotsu3 := checkKotsu(resthand)
		resthand, _ = checkKotsu(resthand)
		if checkTenpai(resthand) {
			out := strings.Join([]string{kotsu1, kotsu2, kotsu3, toitsu, "(" + resthand + ")"}, ",")
			output = append(output, out)
		}
	}

	// check toitsu shuntsu kotsu kotsu kotsu
	for i := 0; i <= 9; i++ {
		resthand := ""

		resthand, toitsu := checkToitsu(hand, i)
		for j := 1; j <= 7; j++ {
			resthand, shuntsu := checkShuntsu(resthand, j)
			if shuntsu != "" {
				resthand, kotsu1 := checkKotsu(resthand)
				resthand, kotsu2 := checkKotsu(resthand)
				resthand, kotsu3 := checkKotsu(resthand)
				if checkTenpai(resthand) {
					out := ""
					if toitsu == "" {
						out = strings.Join([]string{shuntsu, kotsu1, kotsu2, kotsu3, "(" + resthand + ")"}, ",")
					} else {
						out = strings.Join([]string{shuntsu, kotsu1, kotsu2, toitsu, "(" + resthand + ")"}, ",")
					}
					output = append(output, out)
				}
			}
		}
	}

	return output, nil
}

func checkTenpai(resthand string) bool {
	if len(resthand) == 1 {
		return true
	}

	if len(resthand) == 2 {
		r1, _ := strconv.Atoi(string(resthand[0]))
		r2, _ := strconv.Atoi(string(resthand[1]))

		if r2-r1 <= 2 {
			return true
		}
	}

	return false
}

// checkToitsu whether hand has toitsu pairng n like 11
func checkToitsu(hand string, n int) (string, string) {
	tl := []string{"11", "22", "33", "44", "55", "66", "77", "88", "99", "00"}
	toitsu := tl[n]

	handstr := strings.Join(strings.Split(hand, toitsu), "")
	if len(handstr) == len(hand) {
		toitsu = ""
	}

	return handstr, toitsu
}

// checkSHuntsu whether hand has shuntsu like 123
func checkShuntsu(hand string, n int) (string, string) {
	out := ""
	n1 := strings.Index(hand, strconv.Itoa(n))
	n2 := strings.Index(hand, strconv.Itoa(n+1))
	n3 := strings.Index(hand, strconv.Itoa(n+2))

	if n1 >= 0 && n2 >= 0 && n3 >= 0 {
		out = strconv.Itoa(n) + strconv.Itoa(n+1) + strconv.Itoa(n+2)
		hand = hand[0:n1] + hand[n1+1:n2] + hand[n2+1:n3] + hand[n3+1:]
	}

	return hand, out
}

// checkKotsu whether hand has kotsu like 111
func checkKotsu(hand string) (string, string) {
	kl := []string{"111", "222", "333", "444", "555", "666", "777", "888", "999"}

	for _, k := range kl {
		handstr := strings.Join(strings.Split(hand, k), "")
		if len(handstr) != len(hand) {
			return handstr, k
		}
	}

	return hand, ""
}
