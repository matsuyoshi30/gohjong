package gohjong

import (
	"errors"
	"strconv"
)

type TileType int

const (
	SuitTile        TileType = iota // 数牌
	WindTile                        // 風牌
	ThreeDragonTile                 // 三元牌
)

type Tile struct {
	name     string
	tiletype TileType
	num      int // if tiletype is suittile, 1 to 9
}

// ParseHand parse hand to tile slice
func ParseHand(hand string) ([]Tile, error) {
	res := make([]Tile, 0)

	pool := make([]int, 0)
	for i := 0; i < len(hand); i++ {
		handstr := string(hand[i])
		_, err := strconv.Atoi(handstr)
		if err != nil { // not num
			if len(pool) != 0 {
				if handstr != "m" && handstr != "p" && handstr != "s" { // len(pool) > 0 => must be m, p, s
					return nil, errors.New("unknown hand")
				}

				for _, p := range pool {
					res = append(res, Tile{name: strconv.Itoa(p) + handstr, tiletype: SuitTile, num: p})
				}

				pool = make([]int, 0)
			} else { // len(pool) == 0 => must be E, S, W, N, D, H, T
				if handstr != "E" && handstr != "S" && handstr != "W" && handstr != "N" &&
					handstr != "D" && handstr != "H" && handstr != "T" {
					return nil, errors.New("unknown hand")
				}
				res = append(res, Tile{name: handstr})
			}
		} else { // num
			i, err := strconv.Atoi(string(hand[i]))
			if err != nil {
				return nil, err
			}
			pool = append(pool, i)
		}
	}

	if len(pool) != 0 { // must be empty
		return nil, errors.New("could not parse")
	}

	return res, nil
}

// CheckWaiting check waiting tiles
// returns mentsu, machi, and error
func CheckWaiting(hand string) ([]Tile, []Tile, error) {
	handTile, err := ParseHand(hand)
	if err != nil {
		return nil, nil, err
	}

	output := make([]Tile, 0)
	waiting := make([]Tile, 0)

	// check waiting
	// check toistu kotsu kotsu kotsu kotsu
	for i := 0; i <= 9; i++ {
		var resthand []Tile

		resthand, toitsu := checkToitsu(handTile, i)
		resthand, kotsu1 := checkKotsu(resthand)
		resthand, kotsu2 := checkKotsu(resthand)
		resthand, kotsu3 := checkKotsu(resthand)
		resthand, kotsu4 := checkKotsu(resthand)
		if checkTenpai(resthand) {
			output = appendTile(output, kotsu1)
			output = appendTile(output, kotsu2)
			output = appendTile(output, kotsu3)
			if len(resthand) == 1 {
				output = appendTile(output, kotsu4)
			} else {
				output = appendTile(output, toitsu)
			}
			waiting = appendTile(waiting, resthand)
		}
	}

	// check toitsu shuntsu kotsu kotsu kotsu
	for i := 0; i <= 9; i++ {
		var resthand []Tile

		resthand, toitsu := checkToitsu(handTile, i)
		for j := 1; j <= 7; j++ {
			resthand, shuntsu := checkShuntsu(resthand, j)
			if shuntsu != nil {
				resthand, kotsu1 := checkKotsu(resthand)
				resthand, kotsu2 := checkKotsu(resthand)
				resthand, kotsu3 := checkKotsu(resthand)
				if checkTenpai(resthand) {
					output = appendTile(output, shuntsu)
					output = appendTile(output, kotsu1)
					output = appendTile(output, kotsu2)
					if len(resthand) == 1 {
						output = appendTile(output, kotsu3)
					} else {
						output = appendTile(output, toitsu)
					}
					waiting = appendTile(waiting, resthand)
				}
			}
		}
	}

	// check toitsu shuntsu shuntsu kotsu kotsu
	for i := 0; i <= 9; i++ {
		var resthand []Tile

		resthand, toitsu := checkToitsu(handTile, i)
		for j := 1; j <= 7; j++ {
			resthand, shuntsu1 := checkShuntsu(resthand, j)
			if shuntsu1 != nil {
				for k := j; k <= 7; k++ {
					resthand, shuntsu2 := checkShuntsu(resthand, k)
					if shuntsu2 != nil {
						resthand, kotsu1 := checkKotsu(resthand)
						resthand, kotsu2 := checkKotsu(resthand)
						if checkTenpai(resthand) {
							output = appendTile(output, shuntsu1)
							output = appendTile(output, shuntsu2)
							output = appendTile(output, kotsu1)
							if len(resthand) == 1 {
								output = appendTile(output, kotsu2)
							} else {
								output = appendTile(output, toitsu)
							}
							waiting = appendTile(waiting, resthand)
						}
					}
				}
			}
		}
	}

	// check toitsu shuntsu shuntsu shuntsu kotsu
	for i := 0; i <= 9; i++ {
		var resthand []Tile

		resthand, toitsu := checkToitsu(handTile, i)
		for j := 1; j <= 7; j++ {
			resthand, shuntsu1 := checkShuntsu(resthand, j)
			if shuntsu1 != nil {
				for k := j; k <= 7; k++ {
					resthand, shuntsu2 := checkShuntsu(resthand, k)
					if shuntsu2 != nil {
						for l := k; l <= 7; l++ {
							resthand, shuntsu3 := checkShuntsu(resthand, l)
							if shuntsu3 != nil {
								resthand, kotsu := checkKotsu(resthand)
								if checkTenpai(resthand) {
									output = appendTile(output, shuntsu1)
									output = appendTile(output, shuntsu2)
									output = appendTile(output, shuntsu3)
									if len(resthand) == 1 {
										output = appendTile(output, kotsu)
									} else {
										output = appendTile(output, toitsu)
									}
									waiting = appendTile(waiting, resthand)
								}
							}
						}
					}
				}
			}
		}
	}

	// check toitsu shuntsu shuntsu shuntsu shuntsu
	for i := 0; i <= 9; i++ {
		var resthand []Tile

		resthand, toitsu := checkToitsu(handTile, i)
		for j := 1; j <= 7; j++ {
			resthand, shuntsu1 := checkShuntsu(resthand, j)
			if shuntsu1 != nil {
				for k := j; k <= 7; k++ {
					resthand, shuntsu2 := checkShuntsu(resthand, k)
					if shuntsu2 != nil {
						for l := k; l <= 7; l++ {
							resthand, shuntsu3 := checkShuntsu(resthand, l)
							if shuntsu3 != nil {
								for m := l; m <= 7; m++ {
									resthand, shuntsu4 := checkShuntsu(resthand, m)
									if checkTenpai(resthand) {
										output = appendTile(output, shuntsu1)
										output = appendTile(output, shuntsu2)
										output = appendTile(output, shuntsu3)
										if len(resthand) == 1 {
											output = appendTile(output, shuntsu4)
										} else {
											output = appendTile(output, toitsu)
										}
										waiting = appendTile(waiting, resthand)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return output, waiting, nil
}

func checkTenpai(resthand []Tile) bool {
	if len(resthand) == 1 {
		return true
	}

	if len(resthand) == 2 {
		if resthand[0].tiletype != SuitTile || resthand[1].tiletype != SuitTile {
			return false
		}

		if resthand[1].num-resthand[0].num < 2 {
			return true
		}
	}

	return false
}

// checkToitsu whether hand has toitsu pairng n like 11
func checkToitsu(hand []Tile, n int) ([]Tile, []Tile) {
	var toitsu []Tile

	for i := 0; i < len(hand)-1; i++ {
		h1 := hand[i]
		h2 := hand[i+1]
		if h1.num == n && h2.num == n {
			toitsu = []Tile{h1, h2}
			hand = remove(hand, h1)
			hand = remove(hand, h2)
			break
		}
	}

	return hand, toitsu
}

// checkSHuntsu whether hand has shuntsu like 123
func checkShuntsu(hand []Tile, n int) ([]Tile, []Tile) {
	n1 := index(hand, n)
	n2 := index(hand, n+1)
	n3 := index(hand, n+2)

	var out []Tile
	if n1.num != -1 && n2.num != -1 && n3.num != -1 {
		out = []Tile{n1, n2, n3}
		hand = remove(hand, n1)
		hand = remove(hand, n2)
		hand = remove(hand, n3)
	}

	return hand, out
}

// checkKotsu whether hand has kotsu like 111
func checkKotsu(hand []Tile) ([]Tile, []Tile) {
	var kotsu []Tile

	for i := 0; i <= 9; i++ {
		for j := 0; j < len(hand)-2; j++ {
			h1 := hand[j]
			h2 := hand[j+1]
			h3 := hand[j+2]
			if h1.num == i && h2.num == i && h3.num == i {
				kotsu = []Tile{h1, h2, h3}
				hand = remove(hand, h1)
				hand = remove(hand, h2)
				hand = remove(hand, h3)
				break
			}
		}
	}

	return hand, kotsu
}

func remove(pl []Tile, p Tile) []Tile {
	ret := make([]Tile, 0)

	i := 0
	for i = 0; i < len(pl); i++ {
		if pl[i] != p {
			ret = append(ret, pl[i])
		} else {
			break
		}
	}

	if i < len(pl) {
		ret = appendTile(ret, pl[i+1:])
	}

	return ret
}

func index(pl []Tile, n int) Tile {
	t := Tile{num: -1}

	for _, p := range pl {
		if p.num == n {
			return p
		}
	}

	return t
}

func appendTile(dest []Tile, source []Tile) []Tile {
	for _, s := range source {
		dest = append(dest, s)
	}

	return dest
}

func contains(ol []Tile, o Tile) bool {
	for _, oo := range ol {
		if oo == o {
			return true
		}
	}
	return false
}
