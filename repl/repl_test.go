package repl

import (
	"strings"
	"testing"

	"github.com/twolodzko/goal/evaluator"
)

func TestRead_InvalidInput(t *testing.T) {
	var testCases = []string{
		")",
		"(",
		"((",
		"))",
		"())",
		"(()",
	}

	env := evaluator.InitBuildin()

	for _, input := range testCases {
		repl := NewREPL(strings.NewReader(input), env)
		result, err := repl.Read()

		if err == nil {
			t.Errorf("for %s expected an error, got '%s'", input, result)
		}
	}
}
