package gohjong

import (
	"reflect"
	"testing"
)

func TestParseHand(t *testing.T) {
	testcase := []struct {
		input    string
		expected []Tile
	}{
		{"", []Tile{}},
		{"111", []Tile{}},
		{"111m", []Tile{
			Tile{name: "1m", tiletype: SuitTile, num: 1},
			Tile{name: "1m", tiletype: SuitTile, num: 1},
			Tile{name: "1m", tiletype: SuitTile, num: 1}}},
		{"m", nil},
	}

	for _, tt := range testcase {
		output, _ := ParseHand(tt.input)
		if len(output) != len(tt.expected) {
			t.Fatalf("Length wrong: expected %v, but got %v", tt.expected, output)
		}

		for idx, o := range output {
			if o != tt.expected[idx] {
				t.Fatalf("Element[%d] wrong: expected %v, but got %v", idx, tt.expected, output)
			}
		}
	}
}

func TestCheckWaiting(t *testing.T) {
	testcase := []struct {
		input    string
		expected []OutputHand
	}{
		{
			"1112224588899m",
			[]OutputHand{
				OutputHand{
					DefiniteHand: []Tile{Tile{"1m", SuitTile, 1}, Tile{"1m", SuitTile, 1}, Tile{"1m", SuitTile, 1},
						Tile{"2m", SuitTile, 2}, Tile{"2m", SuitTile, 2}, Tile{"2m", SuitTile, 2},
						Tile{"8m", SuitTile, 8}, Tile{"8m", SuitTile, 8}, Tile{"8m", SuitTile, 8},
						Tile{"9m", SuitTile, 9}, Tile{"9m", SuitTile, 9}},
					WaitingHand: []Tile{Tile{"4m", SuitTile, 4}, Tile{"5m", SuitTile, 5}},
				},
			},
		},
	}

	for _, tt := range testcase {
		actualOutput, err := CheckWaiting(tt.input)
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}

		if len(tt.expected) != len(actualOutput) {
			t.Errorf("Error: expected length %d, but got %d\n", len(tt.expected), len(actualOutput))
		}

		if len(actualOutput) != 0 { // check slice length before check elements
			if !reflect.DeepEqual(tt.expected, actualOutput) {
				t.Errorf("Error: expected %v, but got %v\n", tt.expected, actualOutput)
			}
		}
	}
}

func testContain(sl []OutputHand, s OutputHand) bool {
	for _, ss := range sl {
		for _, sd := range s.definiteHand {
			if !testContainHand(ss.definiteHand, sd) {
				return false
			}
		}
		for _, sw := range s.waitingHand {
			if !testContainHand(ss.waitingHand, sw) {
				return false
			}
		}
	}

	return true
}

func testContainHand(hl []Tile, h Tile) bool {
	for _, hh := range hl {
		if hh == h {
			return true
		}
	}

	return false
}
