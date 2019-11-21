package gohjong

import (
	"errors"
	"reflect"
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
	if len(hand) == 0 {
		return nil, errors.New("empty hand")
	}

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
				var tt TileType
				if handstr == "E" || handstr == "S" || handstr == "W" || handstr == "N" {
					tt = WindTile
				} else if handstr == "D" || handstr == "H" || handstr == "T" {
					tt = ThreeDragonTile
				} else {
					return nil, errors.New("unknown hand")
				}
				res = append(res, Tile{name: handstr, tiletype: tt, num: 0})
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

type OutputHand struct {
	DefiniteHand []Tile
	WaitingHand  []Tile
}

func ShowWaiting(hand string) ([]string, error) {
	outputHand, err := CheckWaiting(hand)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0)
	for _, output := range outputHand {
		waitingHand := output.WaitingHand
		if len(waitingHand) == 1 {
			ret = append(ret, waitingHand[0].name)
		} else { // len is 2
			if waitingHand[0].tiletype != SuitTile {
				ret = append(ret, waitingHand[0].name)
			} else {
				n1 := waitingHand[0].num
				n2 := waitingHand[1].num

				if n2-n1 == 1 {
					if n2 == 9 {
						w := "7" + string(waitingHand[0].name[1])
						ret = appendWaiting(ret, w)
					} else if n1 == 1 {
						w := "3" + string(waitingHand[0].name[1])
						ret = appendWaiting(ret, w)
					} else {
						s1 := strconv.Itoa(waitingHand[0].num - 1)
						s2 := strconv.Itoa(waitingHand[1].num + 1)
						w := s1 + string(waitingHand[0].name[1]) + "-" + s2 + string(waitingHand[1].name[1])

						ret = appendWaiting(ret, w)
					}
				} else if n2-n1 == 2 {
					s := strconv.Itoa(waitingHand[0].num + 1)
					ret = appendWaiting(ret, s+string(waitingHand[0].name[1]))
				} else if n2-n1 == 0 {
					ret = appendWaiting(ret, waitingHand[0].name)
				} else {
					return nil, errors.New("unknown hand")
				}
			}
		}
	}

	return ret, nil
}

func appendWaiting(wl []string, w string) []string {
	if !containWaiting(wl, w) {
		wl = append(wl, w)
	}

	return wl
}

func containWaiting(sl []string, s string) bool {
	for _, ss := range sl {
		if ss == s {
			return true
		}
	}

	return false
}

var candidate = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "E", "S", "W", "N", "D", "H", "T"}

// CheckWaiting check waiting tiles
// returns mentsu, machi, and error
func CheckWaiting(hand string) ([]OutputHand, error) {
	handTile, err := ParseHand(hand)
	if err != nil {
		return nil, err
	}

	output := make([]OutputHand, 0)

	// check waiting
	// check toistu kotsu kotsu kotsu kotsu
	for i := 0; i < len(candidate); i++ {
		definite := make([]Tile, 0)
		resthand := make([]Tile, 0)

		resthand, toitsu := checkToitsu(handTile, candidate[i])
		resthand, kotsu1 := checkKotsu(resthand)
		resthand, kotsu2 := checkKotsu(resthand)
		resthand, kotsu3 := checkKotsu(resthand)
		resthand, kotsu4 := checkKotsu(resthand)
		if checkTenpai(resthand) {
			definite = appendTile(definite, kotsu1)
			definite = appendTile(definite, kotsu2)
			definite = appendTile(definite, kotsu3)
			if len(resthand) == 1 {
				definite = appendTile(definite, kotsu4)
			} else {
				definite = appendTile(definite, toitsu)
			}

			o := OutputHand{definite, resthand}
			if !contain(output, o) {
				output = append(output, o)
			}
			definite = nil
		}
	}

	// check toitsu shuntsu kotsu kotsu kotsu
	for i := 0; i < len(candidate); i++ {
		definite := make([]Tile, 0)
		resthand := make([]Tile, 0)

		resthand, toitsu := checkToitsu(handTile, candidate[i])
		for j := 1; j <= 7; j++ {
			resthand, shuntsu := checkShuntsu(resthand, j)
			if shuntsu != nil {
				resthand, kotsu1 := checkKotsu(resthand)
				resthand, kotsu2 := checkKotsu(resthand)
				resthand, kotsu3 := checkKotsu(resthand)
				if checkTenpai(resthand) {
					definite = appendTile(definite, shuntsu)
					definite = appendTile(definite, kotsu1)
					definite = appendTile(definite, kotsu2)
					if len(resthand) == 1 {
						definite = appendTile(definite, kotsu3)
					} else {
						definite = appendTile(definite, toitsu)
					}

					o := OutputHand{definite, resthand}
					if !contain(output, o) {
						output = append(output, o)
					}
					definite = nil
				}
			}
		}
	}

	// check toitsu shuntsu shuntsu kotsu kotsu
	for i := 0; i < len(candidate); i++ {
		definite := make([]Tile, 0)
		resthand := make([]Tile, 0)

		resthand, toitsu := checkToitsu(handTile, candidate[i])
		for j := 1; j <= 7; j++ {
			resthand, shuntsu1 := checkShuntsu(resthand, j)
			if shuntsu1 != nil {
				for k := j; k <= 7; k++ {
					resthand, shuntsu2 := checkShuntsu(resthand, k)
					if shuntsu2 != nil {
						resthand, kotsu1 := checkKotsu(resthand)
						resthand, kotsu2 := checkKotsu(resthand)
						if checkTenpai(resthand) {
							definite = appendTile(definite, shuntsu1)
							definite = appendTile(definite, shuntsu2)
							definite = appendTile(definite, kotsu1)
							if len(resthand) == 1 {
								definite = appendTile(definite, kotsu2)
							} else {
								definite = appendTile(definite, toitsu)
							}

							o := OutputHand{definite, resthand}
							if !contain(output, o) {
								output = append(output, o)
							}
							definite = nil
						}
					}
				}
			}
		}
	}

	// check toitsu shuntsu shuntsu shuntsu kotsu
	for i := 0; i < len(candidate); i++ {
		definite := make([]Tile, 0)
		resthand := make([]Tile, 0)

		resthand, toitsu := checkToitsu(handTile, candidate[i])
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
									definite = appendTile(definite, shuntsu1)
									definite = appendTile(definite, shuntsu2)
									definite = appendTile(definite, shuntsu3)
									if len(resthand) == 1 {
										definite = appendTile(definite, kotsu)
									} else {
										definite = appendTile(definite, toitsu)
									}

									o := OutputHand{definite, resthand}
									if !contain(output, o) {
										output = append(output, o)
									}
									definite = nil
								}
							}
						}
					}
				}
			}
		}
	}

	// check toitsu shuntsu shuntsu shuntsu shuntsu
	for i := 0; i < len(candidate); i++ {
		definite := make([]Tile, 0)
		resthand := make([]Tile, 0)

		resthand, toitsu := checkToitsu(handTile, candidate[i])
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
										definite = appendTile(definite, shuntsu1)
										definite = appendTile(definite, shuntsu2)
										definite = appendTile(definite, shuntsu3)
										if len(resthand) == 1 {
											definite = appendTile(definite, shuntsu4)
										} else {
											definite = appendTile(definite, toitsu)
										}

										o := OutputHand{definite, resthand}
										if !contain(output, o) {
											output = append(output, o)
										}
										definite = nil
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return output, nil
}

func checkTenpai(resthand []Tile) bool {
	if len(resthand) == 1 {
		return true
	}

	if len(resthand) == 2 {
		if resthand[0].tiletype != resthand[1].tiletype {
			return false
		}
		if resthand[0].tiletype == SuitTile {
			if resthand[1].num-resthand[0].num <= 2 {
				return true
			}
		} else {
			if resthand[0].name == resthand[1].name {
				return true
			} else {
				return false
			}
		}
	}

	return false
}

// checkToitsu whether hand has toitsu pairng n like 11
func checkToitsu(hand []Tile, c string) ([]Tile, []Tile) {
	var toitsu []Tile

	n, err := strconv.Atoi(c)

	for i := 0; i < len(hand)-1; i++ {
		h1 := hand[i]
		h2 := hand[i+1]

		if err != nil {
			if h1.name == c && h2.name == c {
				toitsu = []Tile{h1, h2}
				hand = remove(hand, h1)
				hand = remove(hand, h2)
				break
			}
		} else {
			if h1.num == n && h2.num == n {
				toitsu = []Tile{h1, h2}
				hand = remove(hand, h1)
				hand = remove(hand, h2)
				break
			}
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
	if n1.tiletype == SuitTile && n2.tiletype == SuitTile && n3.tiletype == SuitTile {
		if n1.num != -1 && n2.num != -1 && n3.num != -1 {
			out = []Tile{n1, n2, n3}
			hand = remove(hand, n1)
			hand = remove(hand, n2)
			hand = remove(hand, n3)
		}
	}

	return hand, out
}

// checkKotsu whether hand has kotsu like 111
func checkKotsu(hand []Tile) ([]Tile, []Tile) {
	var kotsu []Tile

	for i := 0; i < len(candidate); i++ {
		for j := 0; j < len(hand)-2; j++ {
			h1 := hand[j]
			h2 := hand[j+1]
			h3 := hand[j+2]

			if h1.tiletype == SuitTile && h2.tiletype == SuitTile && h3.tiletype == SuitTile {
				if h1.num == i && h2.num == i && h3.num == i {
					kotsu = []Tile{h1, h2, h3}
					hand = remove(hand, h1)
					hand = remove(hand, h2)
					hand = remove(hand, h3)
					return hand, kotsu
				}
			} else {
				if h1.name == candidate[i] && h2.name == candidate[i] && h3.name == candidate[i] {
					kotsu = []Tile{h1, h2, h3}
					hand = remove(hand, h1)
					hand = remove(hand, h2)
					hand = remove(hand, h3)
					return hand, kotsu
				}
			}
		}
	}

	return hand, kotsu
}

func contain(ol []OutputHand, o OutputHand) bool {
	for _, oh := range ol {
		if reflect.DeepEqual(oh, o) {
			return true
		}
	}

	return false
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
