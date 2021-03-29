package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/token"
)

func Test_isWordBoundary(t *testing.T) {
	var testCases = []struct {
		input    rune
		expected bool
	}{
		{' ', true},
		{'\t', true},
		{'\n', true},
		{'(', true},
		{')', true},
		{'a', false},
		{'8', false},
		{'+', false},
	}

	for _, tt := range testCases {
		result := IsWordBoundary(tt.input)
		if result != tt.expected {
			t.Errorf("for %q expected %v, got: %v", tt.input, tt.expected, result)
		}
	}
}

func TestLexer(t *testing.T) {
	var testCases = []struct {
		input    string
		expected []token.Token
	}{
		{
			"42",
			[]token.Token{
				{Literal: "42", Type: token.INT},
			},
		},
		{
			"\n\n  \t()\n 42 3.1415 -45.7e-2 \"Hello World!\" \" \t\n\n\\\"\\)\" foo",
			[]token.Token{
				{Literal: "(", Type: token.LPAREN},
				{Literal: ")", Type: token.RPAREN},
				{Literal: "42", Type: token.INT},
				{Literal: "3.1415", Type: token.FLOAT},
				{Literal: "-45.7e-2", Type: token.FLOAT},
				{Literal: "Hello World!", Type: token.STRING},
				{Literal: " \t\n\n\")", Type: token.STRING},
				{Literal: "foo", Type: token.SYMBOL},
			},
		},
		{
			"3.14\n(foo 42 \"Hello World!\")\n\n",
			[]token.Token{
				{Literal: "3.14", Type: token.FLOAT},
				{Literal: "(", Type: token.LPAREN},
				{Literal: "foo", Type: token.SYMBOL},
				{Literal: "42", Type: token.INT},
				{Literal: "Hello World!", Type: token.STRING},
				{Literal: ")", Type: token.RPAREN},
			},
		},
		{
			"(+ 2 3)",
			[]token.Token{
				{Literal: "(", Type: token.LPAREN},
				{Literal: "+", Type: token.SYMBOL},
				{Literal: "2", Type: token.INT},
				{Literal: "3", Type: token.INT},
				{Literal: ")", Type: token.RPAREN},
			},
		},
		{
			`(print "Hello World!")`,
			[]token.Token{
				{Literal: "(", Type: token.LPAREN},
				{Literal: "print", Type: token.SYMBOL},
				{Literal: "Hello World!", Type: token.STRING},
				{Literal: ")", Type: token.RPAREN},
			},
		},
		{
			"3.1415 1e5 1.3e+5 -.34e-5 .223 +.45 +42 0",
			[]token.Token{
				{Literal: "3.1415", Type: token.FLOAT},
				{Literal: "1e5", Type: token.FLOAT},
				{Literal: "1.3e+5", Type: token.FLOAT},
				{Literal: "-.34e-5", Type: token.FLOAT},
				{Literal: ".223", Type: token.FLOAT},
				{Literal: "+.45", Type: token.FLOAT},
				{Literal: "+42", Type: token.INT},
				{Literal: "0", Type: token.INT},
			},
		},
	}

	for _, tt := range testCases {
		l := NewLexer(strings.NewReader(tt.input))
		result, err := l.Tokenize()

		if err != nil && err != io.EOF {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(tt.expected, result) {
			t.Errorf("expected %v, got: %v", tt.expected, result)
		}
	}
}

func Test_readString(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{`"" ignore me`, ""},
		{`"Hello World!" not this`, "Hello World!"},
		{`"To escape a char use \\" "ignore me"`, `To escape a char use \`},
		{`"Hello \"John\"!"`, "Hello \"John\"!"},
		{`"It\'s alive!"`, "It's alive!"},
	}

	for _, tt := range testCases {
		l := NewLexer(strings.NewReader(tt.input))

		if err := l.NextRune(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result, err := l.readString()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %v (%T)", tt.expected, tt.expected, result, result)
		}
	}
}

func Test_readString_InvalidInput(t *testing.T) {
	var testCases = []string{
		` "Hello World!"`,
		`Hello World!"`,
		`"Hello World!`,
		"\"Hello World!\n",
		"\n\"Hello World!",
		"Hello World!\n\t\"",
	}

	for _, input := range testCases {
		l := NewLexer(strings.NewReader(input))

		if err := l.NextRune(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		result, err := l.readString()

		if err == nil || err == io.EOF {
			t.Errorf("expected an error, got: %q (%v)", result, err)
		}
	}
}
