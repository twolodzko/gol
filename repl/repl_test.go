package repl

import (
	"strings"
	"testing"
)

func TestInvalidInput(t *testing.T) {
	var testCases = []string{
		")",
		"(",
		"((",
		"))",
		"())",
		"(()",
	}

	for _, input := range testCases {
		result, err := Read(strings.NewReader(input))

		if err == nil {
			t.Errorf("for %s expected an error, got '%s'", input, result)
		}
	}
}
