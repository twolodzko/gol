package parser

import (
	"io"
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
		{"bar ", []objects.Object{objects.Symbol{Val: "bar"}}},
		{"foo bar\n", []objects.Object{objects.Symbol{Val: "foo"}, objects.Symbol{Val: "bar"}}},
		{"42", []objects.Object{objects.Int{Val: 42}}},
		{`"Hello World!" `, []objects.Object{objects.String{Val: "Hello World!"}}},
		{"1e-7", []objects.Object{objects.Float{Val: 1e-7}}},
		{" \n\t bar", []objects.Object{objects.Symbol{Val: "bar"}}},
		{"((1 2))", []objects.Object{objects.NewList(objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2}))}},
		{`(foo 42 "bar")`, []objects.Object{objects.NewList(objects.Symbol{Val: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})}},
		{"  \n\t(\nfoo \n\n42\t\"bar\")", []objects.Object{objects.NewList(objects.Symbol{Val: "foo"}, objects.Int{Val: 42}, objects.String{Val: "bar"})}},
		{"(1 2) (3 4)", []objects.Object{objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2}), objects.NewList(objects.Int{Val: 3}, objects.Int{Val: 4})}},
		{"(1 2)\n\n(3\n4)", []objects.Object{objects.NewList(objects.Int{Val: 1}, objects.Int{Val: 2}), objects.NewList(objects.Int{Val: 3}, objects.Int{Val: 4})}},
	}

	for _, tt := range testCases {
		result, err := Parse(strings.NewReader(tt.input))

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !cmp.Equal(result, tt.expected) {
			t.Errorf("expected: %v (%T), got: %s (%T)", tt.expected, tt.expected, result, result)
		}
	}

	result, err := Parse(strings.NewReader("1 2) (3 4)"))

	if err == nil || err == io.EOF {
		t.Errorf("expected and error, got result: %v (error=%v)", result, err)
	}
}

func TestParse_InvalidInput(t *testing.T) {
	var testCases = []string{
		"(",
		")",
		")(",
		"1 2) (3 4)",
		"(1 2) (3 4",
		"(1 2) 3 4)",
		"(1 2) (",
		"(1 2) )",
	}

	for _, input := range testCases {
		result, err := Parse(strings.NewReader(input))

		if err == nil || err == io.EOF {
			t.Errorf("expected and error, got result: %v (error=%v)", result, err)
		}
	}
}
