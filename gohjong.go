package gohjong

import (
	"errors"
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
func CheckWaiting(hand string) (string, string, error) {
	// TODO
	// _, err := ParseHand(hand)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// check waiting
	for i := 0; i < 9; i++ {
		resthand := ""

		// check toistu kotsu kotsu kotsu kotsu
		resthand, toitsu := checkToitsu(hand, i)
		resthand, kotsu1 := checkKotsu(resthand)
		resthand, kotsu2 := checkKotsu(resthand)
		resthand, kotsu3 := checkKotsu(resthand)
		resthand, _ = checkKotsu(resthand)
		if checkTenpai(resthand) {
			out := strings.Join([]string{kotsu1, kotsu2, kotsu3, toitsu}, ",")
			return out, resthand, nil
		}
	}

	return "", "", errors.New("not tenpai")
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
	tl := []string{"11", "22", "33", "44", "55", "66", "77", "88", "99"}
	toitsu := tl[n]

	handstr := strings.Join(strings.Split(hand, toitsu), "")
	if len(handstr) == len(hand) {
		toitsu = ""
	}

	return handstr, toitsu
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
