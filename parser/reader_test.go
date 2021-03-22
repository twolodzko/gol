package parser

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestCodeReader(t *testing.T) {
	var testCases = []struct {
		input    string
		expected rune
	}{
		{"ax", 'a'},
		{"; this is a comment\nb", 'b'},
	}

	for _, tt := range testCases {
		r := bufio.NewReader(strings.NewReader(tt.input))
		reader := NewCodeReader(r)
		result, _, err := reader.ReadRune()

		if result != tt.expected {
			t.Errorf("expected: '%v', got: '%v'", string(tt.expected), string(result))
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestReadSequence(t *testing.T) {
	input := ";; first comment\n((lambda (x y) ; second comment\n\t\u001B(+ x y))\n;; third comment\n42 13.6)"
	expected := "((lambda (x y) (+ x y)) 42 13.6)"

	r := bufio.NewReader(strings.NewReader(input))
	reader := NewCodeReader(r)
	result := []rune{}

	for {
		r, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result = append(result, r)
	}

	if string(result) != expected {
		t.Errorf("expected: '%s', got: '%s'", expected, string(result))
	}
}
