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
		reader, err := NewCodeReader(r)
		result := reader.Head

		if result != tt.expected {
			t.Errorf("expected: '%v', got: '%v'", string(tt.expected), string(result))
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestReadSequence(t *testing.T) {
	input := ";; first comment\n(foo ; second comment\nbar)"
	expected := "(foo bar)"

	reader, err := NewCodeReader(strings.NewReader(input))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	result := []rune{}

	for {
		result = append(result, reader.Head)
		err := reader.NextRune()

		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	if string(result) != expected {
		t.Errorf("expected: '%s', got: '%s'", expected, string(result))
	}
}

func TestPeekRune_EmptyString(t *testing.T) {
	_, err := NewCodeReader(strings.NewReader(""))
	if err != io.EOF {
		t.Errorf("expected EOF error, got: %v", err)
	}
}
