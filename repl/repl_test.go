package repl

import (
	"bufio"
	"io"
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

func TestRepl(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{"()", "()"},
		{"()\n", "()"},
		{"word", "word"},
		{"(first\t(second))", "(first (second))"},
		{"(first ; ignore this\nsecond);last comment", "(first second)"},
		{"(first\nsecond)", "(first second)"},
		{"(\"first line\nnext line\"\nfoo)", "(\"first line\nnext line\" foo)"},
		// FIXME
		// {`(")")`, `(")")`},
	}

	for _, tt := range testCases {
		reader := bufio.NewReader(strings.NewReader(tt.input))
		result, err := Repl(reader)

		if err != nil && err != io.EOF {
			t.Errorf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(result, tt.expected+"\n") {
			t.Errorf("expected: '%s', got '%s'", tt.expected, result)
		}
	}
}
