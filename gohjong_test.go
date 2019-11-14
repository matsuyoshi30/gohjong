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
		{"111m", []Tile{Tile{"1m"}, Tile{"1m"}, Tile{"1m"}}},
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
		expectedWaiting []string
	}{
		{
			"1112224588899",
			[]string{"111,222,888,99,(45)"},
		},
		{
			"1112223355566",
			[]string{"111,222,555,33,(66)", "111,222,555,66,(33)"},
		},
	}

	for _, tt := range testcase {
		actualWaiting, err := CheckWaiting(tt.input)
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}

		if len(tt.expectedWaiting) != len(actualWaiting) {
			t.Errorf("Error: expected length %d, but got %d\n", len(tt.expectedWaiting), len(actualWaiting))
		}

		for idx, e := range tt.expectedWaiting {
			if e != actualWaiting[idx] {
				t.Errorf("Error: expected %v, but got %v\n", e, actualWaiting[idx])
			}
		}
	}
}
