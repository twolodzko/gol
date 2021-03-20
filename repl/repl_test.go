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
		{"\u001B"},
	}

	for _, test := range testCases {
		reader := bufio.NewReader(strings.NewReader(test.input))
		result, err := Read(reader)

		if err == nil {
			t.Errorf("expected an error, got result: '%s'", result)
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
		{"(word)\n\n\n", "(word)"},
		{"(first\nsecond)", "(first second)"},
		{"(first (second))", "(first (second))"},
		{"(first) (second)", "(first) (second)"},
	}

	for _, tt := range testCases {
		reader := bufio.NewReader(strings.NewReader(tt.input))
		result, err := Read(reader)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("expected: '%s' , got: '%s'", tt.expected, result)
		}
	}
}
