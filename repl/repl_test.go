package repl

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestInvalidInput(t *testing.T) {
	var testCases = []struct {
		input string
	}{
		{")"},
		{"("},
		{"(("},
		{"))"},
		{"())"},
		{"(()"},
	}

	for _, tt := range testCases {
		reader := bufio.NewReader(strings.NewReader(tt.input))
		result, err := Read(reader)

		if err == nil {
			t.Errorf("for %s expected an error, got '%s'", tt.input, result)
		}
	}
}

func TestRead(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{"()", "()"},
		{"word", "word"},
		{"(first\t(second))", "(first\t(second))"},
		// {"(first ; ignore this\nsecond);last comment", "(first second)"},
		// {"(first\nsecond)", "(first second)"},
		// {"(\")\")", "(\")\")"},
		// {"(\"first line\nnext line\"\nfoo)", "(\"first line\nnext line\" foo)"},
	}

	for _, tt := range testCases {
		reader := bufio.NewReader(strings.NewReader(tt.input))
		result, err := Read(reader)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("expected: '%s', got '%s'", tt.expected, result)
		}
	}
}
