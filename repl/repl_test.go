package repl

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestInvalidBraces(t *testing.T) {
	input := ")("
	reader := bufio.NewReader(strings.NewReader(input))
	result, err := Read(reader)

	if err == nil {
		t.Errorf("Expected an error, got result: '%s'", result)
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

	for _, test := range testCases {
		reader := bufio.NewReader(strings.NewReader(test.input))
		result, err := Read(reader)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Experced: '%s' , got: '%s'", test.expected, result)
		}
	}
}
