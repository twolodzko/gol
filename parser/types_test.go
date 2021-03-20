package parser

import "testing"

func TestParseSymbol(t *testing.T) {
	var testCases = []struct {
		input    string
		start    int
		expected Symbol
		pos      int
	}{
		{"foo", 0, Symbol{"foo"}, 2},
		{"(foo bar baz)", 5, Symbol{"bar"}, 7},
	}

	for _, tt := range testCases {
		val, pos, err := ParseSymbol(tt.input, tt.start)

		if val != tt.expected {
			t.Errorf("expected '%v', got '%v'", tt.expected, val)
		}

		if pos != tt.pos {
			t.Errorf("expected position %d, got %d", tt.pos, pos)
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestParseInteger(t *testing.T) {
	var testCases = []struct {
		input    string
		start    int
		expected int
		pos      int
	}{
		{"3", 0, 3, 0},
		{"42", 0, 42, 1},
		{"x 42)", 2, 42, 3},
		{"(list x -42)", 8, -42, 10},
	}

	for _, tt := range testCases {
		val, pos, err := ParseInteger(tt.input, tt.start)

		if val != tt.expected {
			t.Errorf("expected '%v', got '%v'", tt.expected, val)
		}

		if pos != tt.pos {
			t.Errorf("expected position %d, got %d", tt.pos, pos)
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}
