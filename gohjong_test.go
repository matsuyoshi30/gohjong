package gohjong

import (
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
		input           string
		expectedWaiting []Tile
	}{
		{
			"1112224588899m",
			[]Tile{Tile{"4m", SuitTile, 4}, Tile{"5m", SuitTile, 5}},
		},
	}

	for _, tt := range testcase {
		_, actualWaiting, err := CheckWaiting(tt.input)
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}

		if len(tt.expectedWaiting) != len(actualWaiting) {
			t.Errorf("Error: expected length %d, but got %d\n", len(tt.expectedWaiting), len(actualWaiting))
		}

		if len(tt.expectedWaiting) != 0 { // check slice length before check elements
			for _, e := range tt.expectedWaiting {
				if !testContain(actualWaiting, e) {
					t.Errorf("Error: expected %v, but does not\n", e)
				}
			}

			for _, a := range actualWaiting {
				if !testContain(tt.expectedWaiting, a) {
					t.Errorf("Error: unexpected %v\n", a)
				}
			}

		}
	}
}

func testContain(sl []Tile, s Tile) bool {
	for _, ss := range sl {
		if ss == s {
			return true
		}
	}

	return false
}
