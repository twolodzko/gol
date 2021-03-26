package lexer

import "testing"

func Test_isWordBoundary(t *testing.T) {
	var testCases = []struct {
		input    rune
		expected bool
	}{
		{' ', true},
		{'\t', true},
		{'\n', true},
		{'(', true},
		{')', true},
		{'a', false},
		{'8', false},
		{'+', false},
	}

	for _, tt := range testCases {
		result := IsWordBoundary(tt.input)
		if result != tt.expected {
			t.Errorf("for %q expected %v, got: %v", tt.input, tt.expected, result)
		}
	}
}
