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
			[]string{"111,222,555,33,(66)", "111,222,555,66,(33)", "123,123,555,66,(12)"},
		},
		{
			"1113335557779",
			[]string{"111,333,555,777,(9)"},
		},
		{
			"1235556668899",
			[]string{"123,555,666,88,(99)", "123,555,666,99,(88)"},
		},
		{
			"1122336667899",
			[]string{"123,123,666,99,(78)", "123,123,789,666,(9)", "123,123,678,66,(99)", "123,123,678,99,(66)"},
		},
		{
			"1112345678999", // Nine gates
			[]string{"123,456,789,11,(99)", "123,456,789,99,(11)", "123,456,999,11,(78)", "123,678,999,11,(45)",
				"234,567,111,99,(89)", "234,567,111,999,(8)", "234,678,111,999,(5)", "345,678,999,11,(12)",
				"345,678,111,999,(2)", "234,789,111,99,(56)", "456,789,111,99,(23)"},
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

func testContain(sl []string, s string) bool {
	for _, ss := range sl {
		if ss == s {
			return true
		}
	}

	return false
}
