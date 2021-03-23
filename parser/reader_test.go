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
		result, err := reader.ReadRune()

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

	reader := NewCodeReader(strings.NewReader(input))
	result := []rune{}

	for {
		r, err := reader.ReadRune()

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

func TestPeekRune(t *testing.T) {
	input := "你好，世界"
	reader := NewCodeReader(strings.NewReader(input))

	for i := 0; i <= 3; i++ {
		r, err := reader.PeekRune()

		if r != '你' {
			t.Errorf("unexpected result: '%v'", string(r))
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestPeekRune_EmptyString(t *testing.T) {
	reader := NewCodeReader(strings.NewReader(""))
	_, err := reader.PeekRune()
	if err != io.EOF {
		t.Errorf("expected EOF error, got: %v", err)
	}
}
