package parser

import (
	"fmt"
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
		{"123", '1'},
		{"; this is a comment\nb", 'b'},
	}

	for _, tt := range testCases {
		reader, err := NewCodeReader(strings.NewReader(tt.input))
		result := reader.Head

		if result != tt.expected {
			t.Errorf("expected: '%v', got: '%v'", string(tt.expected), string(result))
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestCodeReader_InvalidInput(t *testing.T) {
	input := "\x00abc"
	reader, err := NewCodeReader(strings.NewReader(input))

	if err == nil {
		t.Errorf("expected error, got: %q", reader.Head)
	}
}

func TestReadSequence(t *testing.T) {
	var (
		err error
		cr  *CodeReader
	)

	input := ";; first comment\n(foo ; second comment\nbar;third comment\n42)"

	var testCases = []struct {
		expected rune
		err      error
	}{
		{'(', nil},
		{'f', nil},
		{'o', nil},
		{'o', nil},
		{' ', nil},
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

	cr, err = NewCodeReader(strings.NewReader(input))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var runes []rune

	for i, tt := range testCases {
		result := cr.Head

		fmt.Printf("%q %s\n", result, err)

		if tt.expected != result {
			t.Errorf("at step %d expected %q, got: %q (%v)", i, tt.expected, result, err)
		}

		if err == nil {
			runes = append(runes, result)
		}

		err = cr.NextRune()
	}

	expected := "(foo bar 42)"
	result := string(runes)
	if result != expected {
		t.Errorf("expected '%s' (%d chars), got: '%s' (%d chars)", expected, len(expected), result, len(result))
	}
}

func TestNewCodeReader_EmptyString(t *testing.T) {
	_, err := NewCodeReader(strings.NewReader(""))
	if err != io.EOF {
		t.Errorf("expected EOF error, got: %v", err)
	}
}
