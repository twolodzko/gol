package repl

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/twolodzko/goal/parser"
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
		reader := NewReader(strings.NewReader(tt.input))
		result, err := reader.Read()

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
		// FIXME
		// {"(1 2\n; comment\n 3)", "(1 2 3)"},
		{"(first ; ignore this\nsecond);last comment", "(first second)"},
		{"(first\nsecond)", "(first second)"},
		{"(\"first line\nnext line\"\nfoo)", "(\"first line\nnext line\" foo)"},
		{`(")")`, `(")")`},
		{"(1;)\n2)", "(1 2)"},
		// FIXME
		// {`("\")")`, `("\")")`},
	}

	for _, tt := range testCases {
		reader := bufio.NewReader(strings.NewReader(tt.input))
		result, err := Repl(reader)

		if parser.IsReaderError(err) {
			t.Errorf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(result, tt.expected+"\n") {
			t.Errorf("expected: '%s', got '%s'", tt.expected, result)
		}
	}
}
