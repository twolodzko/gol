package parser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twolodzko/goal/objects"
)

func TestParse(t *testing.T) {
	var testCases = []struct {
		input    string
		expected []objects.Object
	}{
		{" ", nil},
		{"\n", nil},
		{"bar ", []objects.Object{objects.Symbol{Name: "bar"}}},
		{"foo bar\n", []objects.Object{objects.Symbol{Name: "foo"}, objects.Symbol{Name: "bar"}}},
		{"42", []objects.Object{objects.Int{Val: 42}}},
		{`"Hello World!" `, []objects.Object{objects.String{Val: "Hello World!"}}},
		{"42", []objects.Object{objects.Int{Val: 42}}},
		{" \n\t bar", []objects.Object{objects.Symbol{Name: "bar"}}},
		{"((1 2))", []objects.Object{objects.NewList(objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2}))}},
		{`(foo 42 "bar")`, []objects.Object{objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})}},
		{"  \n\t(\nfoo \n\n42\t\"bar\")", []objects.Object{objects.NewList(objects.Symbol{Name: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})}},
	}

	for _, tt := range testCases {
		lexer := NewLexer(strings.NewReader(tt.input))
		tokens, _ := lexer.Tokenize()

		parser := NewParser(tokens)

		result, err := parser.Parse()

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}
}
