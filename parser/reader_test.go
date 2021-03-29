package parser

import (
	"io"
	"strings"
	"testing"
)

func TestCodeReader(t *testing.T) {
	var (
		err error
		cr  *CodeReader
	)

	input := ";; first comment\n(foo ; second comment\nbar;third comment\n42)"

	var testCases = []struct {
		expected rune
		err      error
	}{
		{' ', nil}, // comment
		{'(', nil},
		{'f', nil},
		{'o', nil},
		{'o', nil},
		{' ', nil},
		{' ', nil}, // comment
		{'b', nil},
		{'a', nil},
		{'r', nil},
		{' ', nil},
		{'4', nil},
		{'2', nil},
		{')', nil},
		{'\x00', io.EOF},
		{'\x00', io.EOF},
	}

	cr = NewCodeReader(strings.NewReader(input))

	for i, tt := range testCases {
		err = cr.NextRune()

		if err != nil && err != io.EOF {
			break
		}

		result := cr.Head

		if tt.expected != result {
			t.Errorf("at step %d expected %q, got: %q (%v)", i, tt.expected, result, err)
		}
	}
}
