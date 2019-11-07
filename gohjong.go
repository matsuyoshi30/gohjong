package gohjong

import (
	"errors"
	"strconv"
)

// ParseHand parse hand to string slice
func ParseHand(hand string) ([]string, error) {
	res := make([]string, 0)

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
					res = append(res, string(p)+handstr)
				}

				pool = make([]uint8, 0)
			} else { // len(pool) == 0 => must be E, S, W, N, D, H, T
				if handstr != "E" && handstr != "S" && handstr != "W" && handstr != "N" &&
					handstr != "D" && handstr != "H" && handstr != "T" {
					return nil, errors.New("unknown hand")
				}
				res = append(res, handstr)
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
